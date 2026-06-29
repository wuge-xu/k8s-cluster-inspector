package checker

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func CheckEvents(events []corev1.Event, r *model.Report) {
	r.Events.Total = len(events)

	for _, event := range events {
		if event.Type == corev1.EventTypeWarning {
			r.Events.Warning++

			if len(r.AbnormalEvents) < 10 {
				r.AbnormalEvents = append(r.AbnormalEvents, model.AbnormalEvent{
					Namespace:      event.Namespace,
					Name:           event.Name,
					InvolvedObject: event.InvolvedObject.Kind + "/" + event.InvolvedObject.Name,
					Reason:         event.Reason,
					Message:        event.Message,
					Count:          event.Count,
					Suggestion:     eventSuggestion(event.Reason),
				})
			}
		} else {
			r.Events.Normal++
		}
	}
}

func eventSuggestion(reason string) string {
	switch reason {
	case "FailedScheduling":
		return "Check node resources, taints, tolerations, affinity rules, and PVC binding."
	case "FailedMount":
		return "Check volume configuration, PVC status, StorageClass, and mount permissions."
	case "BackOff", "CrashLoopBackOff":
		return "Check container logs, startup command, probes, and application runtime errors."
	case "FailedPull", "ErrImagePull", "ImagePullBackOff":
		return "Check image name, tag, registry connectivity, and imagePullSecrets."
	default:
		return "Check kubectl describe output and related workload events."
	}
}
