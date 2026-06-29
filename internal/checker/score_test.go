package checker

import (
	"testing"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func TestCalculateScoreHealthy(t *testing.T) {
	score := CalculateScore(model.Report{})

	if score != 100 {
		t.Fatalf("expected score 100, got %d", score)
	}
}

func TestCalculateScoreWithIssues(t *testing.T) {
	r := model.Report{
		Nodes: model.Nodes{
			NotReady: 1,
		},
		Pods: model.Pods{
			Pending:         1,
			HighRestartPods: 2,
		},
		Deployments: model.Deployments{
			Unavailable: 1,
		},
		PVCs: model.PVCs{
			Pending: 1,
		},
		Events: model.Events{
			Warning: 3,
		},
	}

	score := CalculateScore(r)
	expected := 50

	if score != expected {
		t.Fatalf("expected score %d, got %d", expected, score)
	}
}

func TestCalculateScoreFloor(t *testing.T) {
	r := model.Report{
		Nodes: model.Nodes{
			NotReady: 10,
		},
	}

	score := CalculateScore(r)

	if score != 0 {
		t.Fatalf("expected score 0, got %d", score)
	}
}
