package rules

import (
	"fmt"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

type PodCrashLoopRule struct {
	enabled bool
}

func NewPodCrashLoopRule(enabled bool) *PodCrashLoopRule {
	return &PodCrashLoopRule{enabled: enabled}
}

func (r *PodCrashLoopRule) ID() string {
	return "pod-crash-loop-back-off"
}

func (r *PodCrashLoopRule) Name() string {
	return "Pod CrashLoopBackOff"
}

func (r *PodCrashLoopRule) Severity() model.Severity {
	return model.SeverityCritical
}

func (r *PodCrashLoopRule) Enabled() bool {
	return r.enabled
}

func (r *PodCrashLoopRule) Evaluate(data diagnosis.ClusterData) []model.Diagnosis {
	results := make([]model.Diagnosis, 0)

	for _, pod := range data.Pods {
		for _, status := range pod.Status.ContainerStatuses {
			waiting := status.State.Waiting

			if waiting == nil || waiting.Reason != "CrashLoopBackOff" {
				continue
			}

			namespace := pod.Namespace
			if namespace == "" {
				namespace = "default"
			}

			results = append(results, model.Diagnosis{
				RuleID:        r.ID(),
				RuleName:      r.Name(),
				Severity:      r.Severity(),
				ResourceKind:  "Pod",
				Namespace:     namespace,
				ResourceName:  pod.Name,
				ContainerName: status.Name,
				Message: fmt.Sprintf(
					"Container %s is in CrashLoopBackOff and has restarted %d times.",
					status.Name,
					status.RestartCount,
				),
				PossibleCauses: []string{
					"Application startup failure.",
					"Incorrect command, arguments, or environment variables.",
					"ConfigMap or Secret configuration error.",
					"Dependent service is unavailable.",
					"Incorrect health probe configuration.",
				},
				Suggestions: []string{
					"Inspect Pod events.",
					"Check current container logs.",
					"Check logs from the previous container instance.",
					"Verify workload configuration and dependent services.",
				},
				Commands: []string{
					fmt.Sprintf(
						"kubectl describe pod %s -n %s",
						pod.Name,
						namespace,
					),
					fmt.Sprintf(
						"kubectl logs %s -n %s -c %s",
						pod.Name,
						namespace,
						status.Name,
					),
					fmt.Sprintf(
						"kubectl logs %s -n %s -c %s --previous",
						pod.Name,
						namespace,
						status.Name,
					),
				},
			})
		}
	}

	return results
}
