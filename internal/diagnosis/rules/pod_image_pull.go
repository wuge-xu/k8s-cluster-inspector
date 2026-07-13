package rules

import (
	"fmt"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

type PodImagePullRule struct {
	enabled bool
}

func NewPodImagePullRule(enabled bool) *PodImagePullRule {
	return &PodImagePullRule{enabled: enabled}
}

func (r *PodImagePullRule) ID() string {
	return "pod-image-pull-failure"
}

func (r *PodImagePullRule) Name() string {
	return "Pod Image Pull Failure"
}

func (r *PodImagePullRule) Severity() model.Severity {
	return model.SeverityCritical
}

func (r *PodImagePullRule) Enabled() bool {
	return r.enabled
}

func (r *PodImagePullRule) Evaluate(
	data diagnosis.ClusterData,
) []model.Diagnosis {
	results := make([]model.Diagnosis, 0)

	for _, pod := range data.Pods {
		for _, status := range pod.Status.ContainerStatuses {
			waiting := status.State.Waiting
			if waiting == nil {
				continue
			}

			if waiting.Reason != "ImagePullBackOff" &&
				waiting.Reason != "ErrImagePull" {
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
					"Container %s cannot pull its image: %s.",
					status.Name,
					waiting.Reason,
				),
				PossibleCauses: []string{
					"Image name or tag does not exist.",
					"Container registry is unavailable.",
					"ImagePullSecret is missing or invalid.",
					"Registry authentication failed.",
					"Node cannot connect to the registry.",
				},
				Suggestions: []string{
					"Check the image name and tag.",
					"Inspect Pod events for registry errors.",
					"Verify imagePullSecrets.",
					"Check registry connectivity from the node.",
				},
				Commands: []string{
					fmt.Sprintf(
						"kubectl describe pod %s -n %s",
						pod.Name,
						namespace,
					),
					fmt.Sprintf(
						"kubectl get pod %s -n %s -o jsonpath='{.spec.containers[*].image}'",
						pod.Name,
						namespace,
					),
					fmt.Sprintf(
						"kubectl get events -n %s --field-selector involvedObject.name=%s",
						namespace,
						pod.Name,
					),
				},
			})
		}
	}

	return results
}
