package checker

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func CheckPods(pods []corev1.Pod, r *model.Report) {
	r.Pods.Total = len(pods)

	for _, pod := range pods {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			r.Pods.Running++
		case corev1.PodPending:
			r.Pods.Pending++
			addAbnormalPod(r, pod, "Pending", 0)
		case corev1.PodSucceeded:
			r.Pods.Succeeded++
		case corev1.PodFailed:
			r.Pods.Failed++
			addAbnormalPod(r, pod, "Failed", 0)
		}

		for _, cs := range pod.Status.ContainerStatuses {
			if cs.RestartCount >= 3 {
				r.Pods.HighRestartPods++
				addAbnormalPod(r, pod, "HighRestart", cs.RestartCount)
			}

			if cs.State.Waiting != nil {
				reason := cs.State.Waiting.Reason
				if reason == "CrashLoopBackOff" {
					r.Pods.CrashLoopBackOff++
					addAbnormalPod(r, pod, reason, cs.RestartCount)
				}
				if reason == "ImagePullBackOff" || reason == "ErrImagePull" {
					r.Pods.ImagePullBackOff++
					addAbnormalPod(r, pod, reason, cs.RestartCount)
				}
			}
		}
	}
}

func addAbnormalPod(r *model.Report, pod corev1.Pod, reason string, restart int32) {
	key := pod.Namespace + "/" + pod.Name + "/" + reason

	for _, existing := range r.AbnormalPods {
		existingKey := existing.Namespace + "/" + existing.Name + "/" + existing.Reason
		if existingKey == key {
			return
		}
	}

	r.AbnormalPods = append(r.AbnormalPods, model.AbnormalPod{
		Namespace:    pod.Namespace,
		Name:         pod.Name,
		Phase:        string(pod.Status.Phase),
		Reason:       reason,
		RestartCount: restart,
		Suggestion:   Suggestion(reason),
	})
}

func Suggestion(reason string) string {
	switch reason {
	case "Pending":
		return "Check node resources, PVC binding, image pulling, taints, tolerations, and scheduling events."
	case "Failed":
		return "Check pod logs, exit code, and recent Kubernetes events."
	case "CrashLoopBackOff":
		return "Check container logs, command args, environment variables, config files, probes, and resource limits."
	case "ImagePullBackOff", "ErrImagePull":
		return "Check image name, tag, registry connectivity, and imagePullSecrets."
	case "HighRestart":
		return "Check restart reason, application logs, liveness probe, memory limit, and runtime errors."
	default:
		return "Run kubectl describe pod and kubectl logs for further diagnosis."
	}
}
