package rules

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

type DeploymentUnavailableRule struct {
	enabled bool
}

func NewDeploymentUnavailableRule(enabled bool) *DeploymentUnavailableRule {
	return &DeploymentUnavailableRule{enabled: enabled}
}

func (r *DeploymentUnavailableRule) ID() string {
	return "deployment-insufficient-replicas"
}

func (r *DeploymentUnavailableRule) Name() string {
	return "Deployment Insufficient Replicas"
}

func (r *DeploymentUnavailableRule) Severity() model.Severity {
	return model.SeverityWarning
}

func (r *DeploymentUnavailableRule) Enabled() bool {
	return r.enabled
}

func (r *DeploymentUnavailableRule) Evaluate(
	data diagnosis.ClusterData,
) []model.Diagnosis {
	results := make([]model.Diagnosis, 0)

	for _, deployment := range data.Deployments {
		desired := int32(1)
		if deployment.Spec.Replicas != nil {
			desired = *deployment.Spec.Replicas
		}

		if desired == 0 {
			continue
		}

		if deployment.Status.AvailableReplicas >= desired {
			continue
		}

		namespace := deployment.Namespace
		if namespace == "" {
			namespace = "default"
		}

		selector := metav1.FormatLabelSelector(
			deployment.Spec.Selector,
		)

		results = append(results, model.Diagnosis{
			RuleID:       r.ID(),
			RuleName:     r.Name(),
			Severity:     r.Severity(),
			ResourceKind: "Deployment",
			Namespace:    namespace,
			ResourceName: deployment.Name,
			Message: fmt.Sprintf(
				"Deployment expects %d replicas, but only %d are available.",
				desired,
				deployment.Status.AvailableReplicas,
			),
			PossibleCauses: []string{
				"Pods cannot be scheduled.",
				"Containers are failing to start.",
				"Readiness probes are failing.",
				"Image pulling or volume mounting failed.",
				"Rolling update is blocked or incomplete.",
			},
			Suggestions: []string{
				"Inspect Deployment conditions and events.",
				"Check rollout status.",
				"Inspect Pods managed by the Deployment.",
				"Check failed scheduling, image pull, probe, and volume errors.",
			},
			Commands: []string{
				fmt.Sprintf(
					"kubectl describe deployment %s -n %s",
					deployment.Name,
					namespace,
				),
				fmt.Sprintf(
					"kubectl rollout status deployment/%s -n %s",
					deployment.Name,
					namespace,
				),
				fmt.Sprintf(
					"kubectl get pods -n %s -l '%s' -o wide",
					namespace,
					selector,
				),
				fmt.Sprintf(
					"kubectl get events -n %s --sort-by=.lastTimestamp",
					namespace,
				),
			},
		})
	}

	return results
}
