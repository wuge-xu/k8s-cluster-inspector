package checker

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func CheckNamespaces(
	namespaces []corev1.Namespace,
	pods []corev1.Pod,
	deployments []appsv1.Deployment,
	pvcs []corev1.PersistentVolumeClaim,
	r *model.Report,
) {
	r.Namespaces.Total = len(namespaces)

	items := make(map[string]*model.NamespaceSummary)

	for _, ns := range namespaces {
		items[ns.Name] = &model.NamespaceSummary{Name: ns.Name}
	}

	for _, pod := range pods {
		item := ensureNamespaceSummary(items, pod.Namespace)

		item.PodsTotal++

		switch pod.Status.Phase {
		case corev1.PodRunning:
			item.PodsRunning++
		case corev1.PodPending:
			item.PodsPending++
		case corev1.PodFailed:
			item.PodsFailed++
		}

		for _, cs := range pod.Status.ContainerStatuses {
			if cs.RestartCount >= 3 {
				item.PodsHighRestart++
				break
			}
		}
	}

	for _, deploy := range deployments {
		item := ensureNamespaceSummary(items, deploy.Namespace)

		item.DeploymentsTotal++

		desired := int32(1)
		if deploy.Spec.Replicas != nil {
			desired = *deploy.Spec.Replicas
		}

		if deploy.Status.AvailableReplicas < desired {
			item.DeploymentsUnavailable++
		}
	}

	for _, pvc := range pvcs {
		item := ensureNamespaceSummary(items, pvc.Namespace)

		item.PVCsTotal++

		switch pvc.Status.Phase {
		case corev1.ClaimPending:
			item.PVCsPending++
		case corev1.ClaimLost:
			item.PVCsLost++
		}
	}

	for _, item := range items {
		r.Namespaces.Items = append(r.Namespaces.Items, *item)
	}
}

func ensureNamespaceSummary(items map[string]*model.NamespaceSummary, name string) *model.NamespaceSummary {
	if item, ok := items[name]; ok {
		return item
	}

	items[name] = &model.NamespaceSummary{Name: name}
	return items[name]
}
