package rules

import (
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func TestPodCrashLoopRuleDetectsFailure(t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-api",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:         "api",
					RestartCount: 5,
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason: "CrashLoopBackOff",
						},
					},
				},
			},
		},
	}

	rule := NewPodCrashLoopRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Pods: []corev1.Pod{pod},
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 diagnosis, got %d", len(results))
	}

	result := results[0]

	if result.Severity != model.SeverityCritical {
		t.Fatalf("expected critical, got %s", result.Severity)
	}

	if result.ResourceName != "test-api" {
		t.Fatalf("expected test-api, got %s", result.ResourceName)
	}

	if len(result.Commands) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(result.Commands))
	}

	if !strings.Contains(result.Commands[2], "--previous") {
		t.Fatalf("expected previous logs command")
	}
}

func TestPodCrashLoopRuleIgnoresHealthyPod(t *testing.T) {
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

	rule := NewPodCrashLoopRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Pods: []corev1.Pod{pod},
	})

	if len(results) != 0 {
		t.Fatalf("expected no diagnosis, got %d", len(results))
	}
}

func TestPodCrashLoopRuleDisabled(t *testing.T) {
	rule := NewPodCrashLoopRule(false)

	engine := diagnosis.NewEngine(rule)

	results := engine.Run(diagnosis.ClusterData{})

	if len(results) != 0 {
		t.Fatalf("expected disabled rule to return no results")
	}
}
