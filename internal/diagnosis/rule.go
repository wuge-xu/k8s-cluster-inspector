package diagnosis

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

type ClusterData struct {
	Pods        []corev1.Pod
	Nodes       []corev1.Node
	Deployments []appsv1.Deployment
	PVCs        []corev1.PersistentVolumeClaim
}

type Rule interface {
	ID() string
	Name() string
	Severity() model.Severity
	Enabled() bool
	Evaluate(data ClusterData) []model.Diagnosis
}
