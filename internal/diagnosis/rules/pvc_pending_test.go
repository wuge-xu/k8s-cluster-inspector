package rules

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func stringPointer(value string) *string {
	return &value
}

func TestPVCPendingRuleDetectsFailure(t *testing.T) {
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "database-data",
			Namespace: "default",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: stringPointer("fast-storage"),
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("10Gi"),
				},
			},
		},
		Status: corev1.PersistentVolumeClaimStatus{
			Phase: corev1.ClaimPending,
		},
	}

	rule := NewPVCPendingRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		PVCs: []corev1.PersistentVolumeClaim{pvc},
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 diagnosis, got %d", len(results))
	}

	result := results[0]

	if result.RuleID != "pvc-pending" {
		t.Fatalf("unexpected rule ID: %s", result.RuleID)
	}

	if result.Severity != model.SeverityWarning {
		t.Fatalf("expected warning, got %s", result.Severity)
	}

	if result.ResourceKind != "PersistentVolumeClaim" {
		t.Fatalf(
			"expected PersistentVolumeClaim, got %s",
			result.ResourceKind,
		)
	}

	if result.ResourceName != "database-data" {
		t.Fatalf(
			"expected database-data, got %s",
			result.ResourceName,
		)
	}

	if len(result.Commands) != 5 {
		t.Fatalf(
			"expected 5 commands, got %d",
			len(result.Commands),
		)
	}
}

func TestPVCPendingRuleIgnoresBoundPVC(t *testing.T) {
	pvc := corev1.PersistentVolumeClaim{
		Status: corev1.PersistentVolumeClaimStatus{
			Phase: corev1.ClaimBound,
		},
	}

	rule := NewPVCPendingRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		PVCs: []corev1.PersistentVolumeClaim{pvc},
	})

	if len(results) != 0 {
		t.Fatalf("expected no diagnosis, got %d", len(results))
	}
}

func TestPVCPendingRuleCanBeDisabled(t *testing.T) {
	rule := NewPVCPendingRule(false)
	engine := diagnosis.NewEngine(rule)

	results := engine.Run(diagnosis.ClusterData{})

	if len(results) != 0 {
		t.Fatalf("expected no results from disabled rule")
	}
}
