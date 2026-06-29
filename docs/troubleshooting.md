# Troubleshooting Notes

This project is designed for real Kubernetes troubleshooting scenarios.

## Common Issues Detected

### Pod Pending

Possible causes:

    insufficient node resources
    PVC not bound
    image pulling problem
    node selector or affinity mismatch
    taints and tolerations mismatch

Useful commands:

    kubectl describe pod <pod-name> -n <namespace>
    kubectl get events -n <namespace> --sort-by=.lastTimestamp

### CrashLoopBackOff

Possible causes:

    application startup failure
    wrong command or args
    missing environment variables
    failed readiness or liveness probe
    memory limit too small

Useful commands:

    kubectl logs <pod-name> -n <namespace>
    kubectl logs <pod-name> -n <namespace> --previous
    kubectl describe pod <pod-name> -n <namespace>

### ImagePullBackOff

Possible causes:

    wrong image name or tag
    private image without imagePullSecret
    registry network problem
    image does not exist

Useful commands:

    kubectl describe pod <pod-name> -n <namespace>
    kubectl get secret -n <namespace>

### PVC Pending

Possible causes:

    no available PersistentVolume
    StorageClass does not exist
    access mode mismatch
    storage provisioner failure

Useful commands:

    kubectl get pvc -A
    kubectl describe pvc <pvc-name> -n <namespace>
    kubectl get storageclass
