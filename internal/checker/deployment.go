package checker

import (
	appsv1 "k8s.io/api/apps/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func CheckDeployments(deployments []appsv1.Deployment, r *model.Report) {
	r.Deployments.Total = len(deployments)

	for _, d := range deployments {
		desired := int32(1)
		if d.Spec.Replicas != nil {
			desired = *d.Spec.Replicas
		}

		if d.Status.AvailableReplicas >= desired {
			r.Deployments.Available++
		} else {
			r.Deployments.Unavailable++
			r.AbnormalDeployments = append(r.AbnormalDeployments, model.AbnormalDeployment{
				Namespace:         d.Namespace,
				Name:              d.Name,
				DesiredReplicas:   desired,
				AvailableReplicas: d.Status.AvailableReplicas,
				Reason:            "Unavailable",
				Suggestion:        "Check deployment events, pod status, image, probes, resource limits, and rollout history.",
			})
		}
	}
}
