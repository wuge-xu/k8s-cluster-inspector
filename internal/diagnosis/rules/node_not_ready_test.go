package rules

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func TestNodeNotReadyRuleDetectsFailure(t *testing.T) {
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "worker-01",
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{
					Type:    corev1.NodeReady,
					Status:  corev1.ConditionFalse,
					Reason:  "KubeletNotReady",
					Message: "container runtime is down",
				},
			},
		},
	}

	rule := NewNodeNotReadyRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Nodes: []corev1.Node{node},
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 diagnosis, got %d", len(results))
	}

	result := results[0]

	if result.RuleID != "node-not-ready" {
		t.Fatalf("unexpected rule ID: %s", result.RuleID)
	}

	if result.Severity != model.SeverityCritical {
		t.Fatalf("expected critical, got %s", result.Severity)
	}

	if result.ResourceKind != "Node" {
		t.Fatalf("expected Node, got %s", result.ResourceKind)
	}

	if result.ResourceName != "worker-01" {
		t.Fatalf("expected worker-01, got %s", result.ResourceName)
	}

	if len(result.Commands) != 4 {
		t.Fatalf("expected 4 commands, got %d", len(result.Commands))
	}
}

func TestNodeNotReadyRuleIgnoresReadyNode(t *testing.T) {
	node := corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "worker-01",
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeReady,
					Status: corev1.ConditionTrue,
					Reason: "KubeletReady",
				},
			},
		},
	}

	rule := NewNodeNotReadyRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Nodes: []corev1.Node{node},
	})

	if len(results) != 0 {
		t.Fatalf("expected no diagnosis, got %d", len(results))
	}
}

func TestNodeNotReadyRuleCanBeDisabled(t *testing.T) {
	rule := NewNodeNotReadyRule(false)
	engine := diagnosis.NewEngine(rule)

	results := engine.Run(diagnosis.ClusterData{})

	if len(results) != 0 {
		t.Fatalf("expected no results from disabled rule")
	}
}
