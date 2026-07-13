package rules

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func TestPodImagePullRuleDetectsFailure(t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bad-image-pod",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "api",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason: "ImagePullBackOff",
						},
					},
				},
			},
		},
	}

	rule := NewPodImagePullRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Pods: []corev1.Pod{pod},
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 diagnosis, got %d", len(results))
	}

	result := results[0]

	if result.RuleID != "pod-image-pull-failure" {
		t.Fatalf("unexpected rule ID: %s", result.RuleID)
	}

	if result.Severity != model.SeverityCritical {
		t.Fatalf("expected critical, got %s", result.Severity)
	}

	if result.ResourceName != "bad-image-pod" {
		t.Fatalf(
			"expected bad-image-pod, got %s",
			result.ResourceName,
		)
	}

	if len(result.Commands) != 3 {
		t.Fatalf(
			"expected 3 commands, got %d",
			len(result.Commands),
		)
	}
}

func TestPodImagePullRuleDetectsErrImagePull(t *testing.T) {
	pod := corev1.Pod{
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "api",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason: "ErrImagePull",
						},
					},
				},
			},
		},
	}

	rule := NewPodImagePullRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Pods: []corev1.Pod{pod},
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 diagnosis, got %d", len(results))
	}
}

func TestPodImagePullRuleIgnoresHealthyPod(t *testing.T) {
	pod := corev1.Pod{
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "api",
					State: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{},
					},
				},
			},
		},
	}

	rule := NewPodImagePullRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Pods: []corev1.Pod{pod},
	})

	if len(results) != 0 {
		t.Fatalf("expected no diagnosis, got %d", len(results))
	}
}
