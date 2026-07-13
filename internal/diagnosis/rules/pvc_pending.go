package rules

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

type PVCPendingRule struct {
	enabled bool
}

func NewPVCPendingRule(enabled bool) *PVCPendingRule {
	return &PVCPendingRule{enabled: enabled}
}

func (r *PVCPendingRule) ID() string {
	return "pvc-pending"
}

func (r *PVCPendingRule) Name() string {
	return "PVC Pending"
}

func (r *PVCPendingRule) Severity() model.Severity {
	return model.SeverityWarning
}

func (r *PVCPendingRule) Enabled() bool {
	return r.enabled
}

func (r *PVCPendingRule) Evaluate(
	data diagnosis.ClusterData,
) []model.Diagnosis {
	results := make([]model.Diagnosis, 0)

	for _, pvc := range data.PVCs {
		if pvc.Status.Phase != corev1.ClaimPending {
			continue
		}

		namespace := pvc.Namespace
		if namespace == "" {
			namespace = "default"
		}

		storageClass := "<default>"
		if pvc.Spec.StorageClassName != nil &&
			*pvc.Spec.StorageClassName != "" {
			storageClass = *pvc.Spec.StorageClassName
		}

		requestedStorage := pvc.Spec.Resources.Requests.Storage()
		requestedSize := "unknown"
		if requestedStorage != nil {
			requestedSize = requestedStorage.String()
		}

		results = append(results, model.Diagnosis{
			RuleID:       r.ID(),
			RuleName:     r.Name(),
			Severity:     r.Severity(),
			ResourceKind: "PersistentVolumeClaim",
			Namespace:    namespace,
			ResourceName: pvc.Name,
			Message: fmt.Sprintf(
				"PVC is Pending: storageClass=%s requestedSize=%s.",
				storageClass,
				requestedSize,
			),
			PossibleCauses: []string{
				"No matching PersistentVolume is available.",
				"StorageClass does not exist or is configured incorrectly.",
				"Dynamic volume provisioning failed.",
				"Requested storage size or access mode cannot be satisfied.",
				"Cloud or local storage provisioner is unavailable.",
			},
			Suggestions: []string{
				"Inspect PVC events and status.",
				"Verify the configured StorageClass.",
				"Check available PersistentVolumes.",
				"Check the storage provisioner Pods and logs.",
				"Verify requested access modes and storage size.",
			},
			Commands: []string{
				fmt.Sprintf(
					"kubectl describe pvc %s -n %s",
					pvc.Name,
					namespace,
				),
				fmt.Sprintf(
					"kubectl get pvc %s -n %s -o yaml",
					pvc.Name,
					namespace,
				),
				"kubectl get storageclass",
				"kubectl get persistentvolume",
				fmt.Sprintf(
					"kubectl get events -n %s --field-selector involvedObject.kind=PersistentVolumeClaim,involvedObject.name=%s --sort-by=.lastTimestamp",
					namespace,
					pvc.Name,
				),
			},
		})
	}

	return results
}
