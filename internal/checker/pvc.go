package checker

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func CheckPVCs(pvcs []corev1.PersistentVolumeClaim, r *model.Report) {
	r.PVCs.Total = len(pvcs)

	for _, pvc := range pvcs {
		switch pvc.Status.Phase {
		case corev1.ClaimBound:
			r.PVCs.Bound++
		case corev1.ClaimPending:
			r.PVCs.Pending++
			addAbnormalPVC(r, pvc, "Check StorageClass, PersistentVolume availability, access modes, and PVC events.")
		case corev1.ClaimLost:
			r.PVCs.Lost++
			addAbnormalPVC(r, pvc, "PVC is lost. Check the related PersistentVolume and storage backend status.")
		}
	}
}

func addAbnormalPVC(r *model.Report, pvc corev1.PersistentVolumeClaim, suggestion string) {
	r.AbnormalPVCs = append(r.AbnormalPVCs, model.AbnormalPVC{
		Namespace:  pvc.Namespace,
		Name:       pvc.Name,
		Phase:      string(pvc.Status.Phase),
		VolumeName: pvc.Spec.VolumeName,
		Suggestion: suggestion,
	})
}
