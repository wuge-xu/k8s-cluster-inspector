# Kubernetes Cluster Inspector

A lightweight Kubernetes cluster inspection tool for SRE and platform engineering scenarios.

## Features

- Kubernetes cluster health score
- Node, Pod, Deployment and Service inspection
- Abnormal Pod diagnosis with troubleshooting suggestions
- Deployment availability inspection
- JSON report export
- Prometheus metrics endpoint

## Quick Start

Run inspector:

    go run ./cmd/inspector

Run with Prometheus metrics endpoint:

    go run ./cmd/inspector --metrics

Check metrics:

    curl http://localhost:9090/metrics

## Build

    make build

Run binary:

    ./bin/k8s-cluster-inspector

## Example Output

    Kubernetes Cluster Inspector
    Health Score: 77 / 100
    Nodes: total=1 ready=1 notReady=0
    Pods: total=15 running=12 pending=1 failed=0 restart>=3=6
    Deployments: total=8 available=8 unavailable=0
    Services: total=19 clusterIP=18 nodePort=0 loadBalancer=1

## Metrics

Example metrics:

    k8s_cluster_health_score 77
    k8s_nodes_total 1
    k8s_nodes_ready 1
    k8s_pods_total 15
    k8s_pods_running 12
    k8s_pods_pending 1
    k8s_deployments_total 8
    k8s_services_total 19

## Roadmap

- PVC inspection
- Kubernetes Event analysis
- HTML report
- Docker support
- GitHub Actions
- Grafana dashboard
