# Metrics

When running with metrics enabled:

    go run ./cmd/inspector --metrics

The tool exposes metrics at:

    http://localhost:9090/metrics

## Example Metrics

    k8s_cluster_health_score 77
    k8s_nodes_total 1
    k8s_nodes_ready 1
    k8s_nodes_not_ready 0
    k8s_pods_total 15
    k8s_pods_running 12
    k8s_pods_pending 1
    k8s_pods_failed 0
    k8s_pods_high_restart 6
    k8s_deployments_total 8
    k8s_deployments_available 8
    k8s_deployments_unavailable 0
    k8s_services_total 19
