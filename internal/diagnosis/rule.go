package diagnosis

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

type ClusterData struct {
	Pods []corev1.Pod
}

type Rule interface {
	ID() string
	Name() string
	Severity() model.Severity
	Enabled() bool
	Evaluate(data ClusterData) []model.Diagnosis
}
