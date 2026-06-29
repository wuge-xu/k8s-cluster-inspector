# Kubernetes Cluster Inspector

[![CI](https://github.com/wuge-xu/k8s-cluster-inspector/actions/workflows/ci.yml/badge.svg)](https://github.com/wuge-xu/k8s-cluster-inspector/actions/workflows/ci.yml)

A lightweight Kubernetes cluster inspection tool for SRE and platform engineering scenarios.

This project is built with Go and Kubernetes client-go. It inspects Kubernetes cluster resources, calculates a health score, diagnoses abnormal workloads, exports JSON and HTML reports, and exposes Prometheus-style metrics.

## Features

- Kubernetes cluster health score
- Node health inspection
- Pod status inspection
- Abnormal Pod diagnosis with troubleshooting suggestions
- Deployment availability inspection
- Service type summary
- PVC status inspection
- Kubernetes Warning Event analysis
- JSON report export
- HTML report export
- Prometheus-style metrics endpoint
- Docker build support
- GitHub Actions CI

## Tech Stack

- Go
- Kubernetes client-go
- Kubernetes API
- Prometheus-style metrics
- Docker
- GitHub Actions

## Project Structure

    cmd/inspector
    internal/client
    internal/checker
    internal/model
    internal/report
    internal/metrics
    docs
    examples

## Quick Start

Run inspector:

    go run ./cmd/inspector

The tool will generate:

    report.json
    report.html

Open HTML report in WSL:

    explorer.exe report.html

Run with metrics endpoint:

    go run ./cmd/inspector --metrics

Check metrics:

    curl http://localhost:9090/metrics

## Build

    make build

Run binary:

    ./bin/k8s-cluster-inspector

## Docker

Build image:

    make docker-build

Run container:

    make docker-run

Note: Docker runtime needs access to kubeconfig and Kubernetes API.

## Example Output

See:

    examples/sample-output.txt

Example summary:

    Health Score: 77 / 100

    Nodes: total=1 ready=1 notReady=0
    Pods: total=15 running=12 pending=1 failed=0 restart>=3=6
    Deployments: total=8 available=8 unavailable=0
    Services: total=19 clusterIP=18 nodePort=0 loadBalancer=1
    PVCs: total=0 bound=0 pending=0 lost=0
    Events: total=120 normal=80 warning=40

## Metrics

Example metrics:

    k8s_cluster_health_score 77
    k8s_nodes_total 1
    k8s_nodes_ready 1
    k8s_pods_total 15
    k8s_pods_running 12
    k8s_pods_pending 1
    k8s_deployments_total 8
    k8s_pvcs_total 0
    k8s_events_warning 40

More details:

    docs/metrics.md

## Architecture

See:

    docs/architecture.md

## Troubleshooting

See:

    docs/troubleshooting.md

## Why This Project

In SRE and platform engineering work, engineers often need to quickly understand whether a Kubernetes cluster is healthy, which workloads are abnormal, and what the likely troubleshooting direction is.

This project focuses on that scenario. It collects cluster runtime information from the Kubernetes API, summarizes key health indicators, detects abnormal resources, and provides basic diagnosis suggestions.


## Testing

Run unit tests:

    go test ./...

Current unit tests cover:

    health score calculation
    PVC status inspection
    Kubernetes Warning Event aggregation
    Namespace summary calculation

## Roadmap

- Node pressure inspection
- Namespace summary
- Kubernetes Event aggregation by reason
- Grafana dashboard
- Native prometheus/client_golang support
- Multi-cluster support
- More unit tests
