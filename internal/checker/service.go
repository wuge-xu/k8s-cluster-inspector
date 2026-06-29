package checker

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func CheckServices(services []corev1.Service, r *model.Report) {
	r.Services.Total = len(services)

	for _, svc := range services {
		switch svc.Spec.Type {
		case corev1.ServiceTypeClusterIP:
			r.Services.ClusterIP++
		case corev1.ServiceTypeNodePort:
			r.Services.NodePort++
		case corev1.ServiceTypeLoadBalancer:
			r.Services.LoadBalancer++
		}
	}
}
