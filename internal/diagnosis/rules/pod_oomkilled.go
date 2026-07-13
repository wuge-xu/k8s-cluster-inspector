package rules

import (
	"fmt"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

type PodOOMKilledRule struct {
	enabled bool
}

func NewPodOOMKilledRule(enabled bool) *PodOOMKilledRule {
	return &PodOOMKilledRule{enabled: enabled}
}

func (r *PodOOMKilledRule) ID() string {
	return "pod-oom-killed"
}

func (r *PodOOMKilledRule) Name() string {
	return "Pod OOMKilled"
}

func (r *PodOOMKilledRule) Severity() model.Severity {
	return model.SeverityCritical
}

func (r *PodOOMKilledRule) Enabled() bool {
	return r.enabled
}

func (r *PodOOMKilledRule) Evaluate(
	data diagnosis.ClusterData,
) []model.Diagnosis {
	results := make([]model.Diagnosis, 0)

	for _, pod := range data.Pods {
		for _, status := range pod.Status.ContainerStatuses {
			terminated := status.LastTerminationState.Terminated

			if terminated == nil || terminated.Reason != "OOMKilled" {
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
					"Container %s was terminated because it exceeded its memory limit.",
					status.Name,
				),
				PossibleCauses: []string{
					"Container memory limit is too low.",
					"Application memory usage increased unexpectedly.",
					"Application contains a memory leak.",
					"Large requests or workloads caused a memory spike.",
				},
				Suggestions: []string{
					"Check the configured memory requests and limits.",
					"Inspect previous container logs.",
					"Check container memory metrics before the restart.",
					"Increase the memory limit only after identifying expected usage.",
				},
				Commands: []string{
					fmt.Sprintf(
						"kubectl describe pod %s -n %s",
						pod.Name,
						namespace,
					),
					fmt.Sprintf(
						"kubectl logs %s -n %s -c %s --previous",
						pod.Name,
						namespace,
						status.Name,
					),
					fmt.Sprintf(
						"kubectl get pod %s -n %s -o jsonpath='{.spec.containers[*].resources}'",
						pod.Name,
						namespace,
					),
					fmt.Sprintf(
						"kubectl top pod %s -n %s --containers",
						pod.Name,
						namespace,
					),
				},
			})
		}
	}

	return results
}
