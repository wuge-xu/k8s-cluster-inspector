package checker

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func TestCheckPVCs(t *testing.T) {
	pvcs := []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bound-pvc",
				Namespace: "default",
			},
			Status: corev1.PersistentVolumeClaimStatus{
				Phase: corev1.ClaimBound,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pending-pvc",
				Namespace: "default",
			},
			Status: corev1.PersistentVolumeClaimStatus{
				Phase: corev1.ClaimPending,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "lost-pvc",
				Namespace: "monitoring",
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				VolumeName: "pv-lost",
			},
			Status: corev1.PersistentVolumeClaimStatus{
				Phase: corev1.ClaimLost,
			},
		},
	}

	var r model.Report
	CheckPVCs(pvcs, &r)

	if r.PVCs.Total != 3 {
		t.Fatalf("expected total pvcs 3, got %d", r.PVCs.Total)
	}

	if r.PVCs.Bound != 1 {
		t.Fatalf("expected bound pvcs 1, got %d", r.PVCs.Bound)
	}

	if r.PVCs.Pending != 1 {
		t.Fatalf("expected pending pvcs 1, got %d", r.PVCs.Pending)
	}

	if r.PVCs.Lost != 1 {
		t.Fatalf("expected lost pvcs 1, got %d", r.PVCs.Lost)
	}

	if len(r.AbnormalPVCs) != 2 {
		t.Fatalf("expected 2 abnormal pvcs, got %d", len(r.AbnormalPVCs))
	}
}
