# Architecture

Kubernetes Cluster Inspector is designed as a lightweight SRE-oriented inspection tool.

## Components

    cmd/inspector
        CLI entrypoint.

    internal/client
        Builds Kubernetes clientset from local kubeconfig.

    internal/checker
        Contains inspection logic for Nodes, Pods, Deployments and Services.

    internal/model
        Defines report structures and diagnosis data models.

    internal/report
        Prints console report and exports JSON report.

    internal/metrics
        Exposes Prometheus-style metrics endpoint.

## Workflow

    Load kubeconfig
        |
        v
    Create Kubernetes client
        |
        v
    Fetch cluster resources
        |
        v
    Run checkers
        |
        v
    Generate health score
        |
        v
    Print console report / export JSON / expose metrics
