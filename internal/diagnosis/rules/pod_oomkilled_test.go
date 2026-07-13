package rules

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/diagnosis"
	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func TestPodOOMKilledRuleDetectsFailure(t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "memory-worker",
			Namespace: "default",
		},
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name:         "worker",
					RestartCount: 2,
					LastTerminationState: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							Reason:   "OOMKilled",
							ExitCode: 137,
						},
					},
				},
			},
		},
	}

	rule := NewPodOOMKilledRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Pods: []corev1.Pod{pod},
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 diagnosis, got %d", len(results))
	}

	result := results[0]

	if result.RuleID != "pod-oom-killed" {
		t.Fatalf("unexpected rule ID: %s", result.RuleID)
	}

	if result.Severity != model.SeverityCritical {
		t.Fatalf("expected critical, got %s", result.Severity)
	}

	if result.ContainerName != "worker" {
		t.Fatalf("expected worker, got %s", result.ContainerName)
	}

	if len(result.Commands) != 4 {
		t.Fatalf("expected 4 commands, got %d", len(result.Commands))
	}
}

func TestPodOOMKilledRuleIgnoresNormalExit(t *testing.T) {
	pod := corev1.Pod{
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "worker",
					LastTerminationState: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							Reason:   "Completed",
							ExitCode: 0,
						},
					},
				},
			},
		},
	}

	rule := NewPodOOMKilledRule(true)

	results := rule.Evaluate(diagnosis.ClusterData{
		Pods: []corev1.Pod{pod},
	})

	if len(results) != 0 {
		t.Fatalf("expected no diagnosis, got %d", len(results))
	}
}
