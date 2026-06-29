package metrics

import (
	"fmt"
	"net/http"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func Start(r model.Report, addr string) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "k8s_cluster_health_score %d\n", r.Score)
		fmt.Fprintf(w, "k8s_nodes_total %d\n", r.Nodes.Total)
		fmt.Fprintf(w, "k8s_nodes_ready %d\n", r.Nodes.Ready)
		fmt.Fprintf(w, "k8s_nodes_not_ready %d\n", r.Nodes.NotReady)
		fmt.Fprintf(w, "k8s_pods_total %d\n", r.Pods.Total)
		fmt.Fprintf(w, "k8s_pods_running %d\n", r.Pods.Running)
		fmt.Fprintf(w, "k8s_pods_pending %d\n", r.Pods.Pending)
		fmt.Fprintf(w, "k8s_pods_failed %d\n", r.Pods.Failed)
		fmt.Fprintf(w, "k8s_pods_high_restart %d\n", r.Pods.HighRestartPods)
		fmt.Fprintf(w, "k8s_deployments_total %d\n", r.Deployments.Total)
		fmt.Fprintf(w, "k8s_deployments_available %d\n", r.Deployments.Available)
		fmt.Fprintf(w, "k8s_deployments_unavailable %d\n", r.Deployments.Unavailable)
		fmt.Fprintf(w, "k8s_services_total %d\n", r.Services.Total)
	})

	fmt.Println("Metrics server:", addr)
	_ = http.ListenAndServe(addr, nil)
}
