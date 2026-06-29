package checker

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func CheckNodes(nodes []corev1.Node, r *model.Report) {
	r.Nodes.Total = len(nodes)

	for _, node := range nodes {
		ready := false
		for _, c := range node.Status.Conditions {
			if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
				ready = true
				break
			}
		}

		if ready {
			r.Nodes.Ready++
		} else {
			r.Nodes.NotReady++
		}
	}
}
