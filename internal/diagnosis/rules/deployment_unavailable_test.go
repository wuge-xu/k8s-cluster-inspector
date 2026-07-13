package rules

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func int32Pointer(value int32) *int32 {
	return &value
}

func TestDeploymentUnavailableRuleDetectsFailure(t *testing.T) {
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "api-server",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Pointer(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "api-server",
				},
			},
		},
		Status: appsv1.DeploymentStatus{
			ReadyReplicas:     1,
			AvailableReplicas: 1,
			UpdatedReplicas:   2,
		},
	}

	rule := NewDeploymentUnavailableRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Deployments: []appsv1.Deployment{deployment},
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 diagnosis, got %d", len(results))
	}

	result := results[0]

	if result.RuleID != "deployment-insufficient-replicas" {
		t.Fatalf("unexpected rule ID: %s", result.RuleID)
	}

	if result.Severity != model.SeverityWarning {
		t.Fatalf("expected warning, got %s", result.Severity)
	}

	if result.ResourceKind != "Deployment" {
		t.Fatalf(
			"expected Deployment, got %s",
			result.ResourceKind,
		)
	}

	if result.ResourceName != "api-server" {
		t.Fatalf(
			"expected api-server, got %s",
			result.ResourceName,
		)
	}

	if len(result.Commands) != 4 {
		t.Fatalf(
			"expected 4 commands, got %d",
			len(result.Commands),
		)
	}
}

func TestDeploymentUnavailableRuleIgnoresHealthyDeployment(
	t *testing.T,
) {
	deployment := appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Pointer(2),
		},
		Status: appsv1.DeploymentStatus{
			AvailableReplicas: 2,
		},
	}

	rule := NewDeploymentUnavailableRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Deployments: []appsv1.Deployment{deployment},
	})

	if len(results) != 0 {
		t.Fatalf("expected no diagnosis, got %d", len(results))
	}
}

func TestDeploymentUnavailableRuleIgnoresScaledToZero(
	t *testing.T,
) {
	deployment := appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Pointer(0),
		},
	}

	rule := NewDeploymentUnavailableRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Deployments: []appsv1.Deployment{deployment},
	})

	if len(results) != 0 {
		t.Fatalf("expected no diagnosis, got %d", len(results))
	}
}
