package rules

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

type NodeNotReadyRule struct {
	enabled bool
}

func NewNodeNotReadyRule(enabled bool) *NodeNotReadyRule {
	return &NodeNotReadyRule{enabled: enabled}
}

func (r *NodeNotReadyRule) ID() string {
	return "node-not-ready"
}

func (r *NodeNotReadyRule) Name() string {
	return "Node NotReady"
}

func (r *NodeNotReadyRule) Severity() model.Severity {
	return model.SeverityCritical
}

func (r *NodeNotReadyRule) Enabled() bool {
	return r.enabled
}

func (r *NodeNotReadyRule) Evaluate(
	data diagnosis.ClusterData,
) []model.Diagnosis {
	results := make([]model.Diagnosis, 0)

	for _, node := range data.Nodes {
		for _, condition := range node.Status.Conditions {
			if condition.Type != corev1.NodeReady {
				continue
			}

			if condition.Status == corev1.ConditionTrue {
				break
			}

			reason := condition.Reason
			if reason == "" {
				reason = "Unknown"
			}

			results = append(results, model.Diagnosis{
				RuleID:       r.ID(),
				RuleName:     r.Name(),
				Severity:     r.Severity(),
				ResourceKind: "Node",
				ResourceName: node.Name,
				Message: fmt.Sprintf(
					"Node %s is not ready: status=%s reason=%s.",
					node.Name,
					condition.Status,
					reason,
				),
				PossibleCauses: []string{
					"Kubelet is stopped or unhealthy.",
					"Node cannot communicate with the Kubernetes control plane.",
					"Container runtime is unavailable.",
					"Node has disk, memory, PID, or network pressure.",
					"Underlying virtual machine or host is unavailable.",
				},
				Suggestions: []string{
					"Inspect Node conditions and recent events.",
					"Check kubelet and container runtime status.",
					"Check node resource pressure.",
					"Inspect workloads currently assigned to the node.",
				},
				Commands: []string{
					fmt.Sprintf(
						"kubectl describe node %s",
						node.Name,
					),
					fmt.Sprintf(
						"kubectl get node %s -o wide",
						node.Name,
					),
					fmt.Sprintf(
						"kubectl get pods -A --field-selector spec.nodeName=%s",
						node.Name,
					),
					fmt.Sprintf(
						"kubectl top node %s",
						node.Name,
					),
				},
			})

			break
		}
	}

	return results
}
