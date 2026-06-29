package checker

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func TestCheckNamespaces(t *testing.T) {
	replicas := int32(2)

	namespaces := []corev1.Namespace{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "monitoring",
			},
		},
	}

	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "running-pod",
				Namespace: "default",
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pending-pod",
				Namespace: "default",
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodPending,
				ContainerStatuses: []corev1.ContainerStatus{
					{
						RestartCount: 5,
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "failed-pod",
				Namespace: "monitoring",
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodFailed,
			},
		},
	}

	deployments := []appsv1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "web",
				Namespace: "default",
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &replicas,
			},
			Status: appsv1.DeploymentStatus{
				AvailableReplicas: 1,
			},
		},
	}

	pvcs := []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "data",
				Namespace: "default",
			},
			Status: corev1.PersistentVolumeClaimStatus{
				Phase: corev1.ClaimPending,
			},
		},
	}

	var r model.Report
	CheckNamespaces(namespaces, pods, deployments, pvcs, &r)

	if r.Namespaces.Total != 2 {
		t.Fatalf("expected namespaces total 2, got %d", r.Namespaces.Total)
	}

	defaultNS := findNamespaceSummary(t, r.Namespaces.Items, "default")
	monitoringNS := findNamespaceSummary(t, r.Namespaces.Items, "monitoring")

	if defaultNS.PodsTotal != 2 {
		t.Fatalf("expected default pods total 2, got %d", defaultNS.PodsTotal)
	}

	if defaultNS.PodsPending != 1 {
		t.Fatalf("expected default pending pods 1, got %d", defaultNS.PodsPending)
	}

	if defaultNS.PodsHighRestart != 1 {
		t.Fatalf("expected default high restart pods 1, got %d", defaultNS.PodsHighRestart)
	}

	if defaultNS.DeploymentsUnavailable != 1 {
		t.Fatalf("expected default unavailable deployments 1, got %d", defaultNS.DeploymentsUnavailable)
	}

	if defaultNS.PVCsPending != 1 {
		t.Fatalf("expected default pending pvcs 1, got %d", defaultNS.PVCsPending)
	}

	if monitoringNS.PodsFailed != 1 {
		t.Fatalf("expected monitoring failed pods 1, got %d", monitoringNS.PodsFailed)
	}
}

func findNamespaceSummary(t *testing.T, items []model.NamespaceSummary, name string) model.NamespaceSummary {
	t.Helper()

	for _, item := range items {
		if item.Name == name {
			return item
		}
	}

	t.Fatalf("namespace summary %s not found", name)
	return model.NamespaceSummary{}
}
