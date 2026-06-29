package checker

import "github.com/boserwuge/k8s-cluster-inspector/internal/model"

func CalculateScore(r model.Report) int {
	score := 100
	score -= r.Nodes.NotReady * 20
	score -= r.Pods.Pending * 5
	score -= r.Pods.Failed * 10
	score -= r.Pods.CrashLoopBackOff * 10
	score -= r.Pods.ImagePullBackOff * 8
	score -= r.Pods.HighRestartPods * 3
	score -= r.Deployments.Unavailable * 8

	if score < 0 {
		return 0
	}
	return score
}
