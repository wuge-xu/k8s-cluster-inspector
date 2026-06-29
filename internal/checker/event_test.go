package checker

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func TestCheckEventsAggregatesWarningEventCounts(t *testing.T) {
	events := []corev1.Event{
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "event-1",
			},
			Type:   corev1.EventTypeWarning,
			Reason: "FailedScheduling",
			Count:  3,
			InvolvedObject: corev1.ObjectReference{
				Kind: "Pod",
				Name: "pod-a",
			},
			Message: "failed scheduling",
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "event-2",
			},
			Type:   corev1.EventTypeWarning,
			Reason: "FailedScheduling",
			Count:  1,
			InvolvedObject: corev1.ObjectReference{
				Kind: "Pod",
				Name: "pod-b",
			},
			Message: "failed scheduling again",
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "monitoring",
				Name:      "event-3",
			},
			Type:   corev1.EventTypeWarning,
			Reason: "FailedMount",
			Count:  2,
			InvolvedObject: corev1.ObjectReference{
				Kind: "Pod",
				Name: "pod-c",
			},
			Message: "failed mount",
		},
		{
			Type:   corev1.EventTypeNormal,
			Reason: "Pulled",
			Count:  1,
		},
	}

	var r model.Report
	CheckEvents(events, &r)

	if r.Events.Total != 4 {
		t.Fatalf("expected total events 4, got %d", r.Events.Total)
	}

	if r.Events.Warning != 3 {
		t.Fatalf("expected warning events 3, got %d", r.Events.Warning)
	}

	if r.Events.Normal != 1 {
		t.Fatalf("expected normal events 1, got %d", r.Events.Normal)
	}

	if len(r.EventReasons) != 2 {
		t.Fatalf("expected 2 event reasons, got %d", len(r.EventReasons))
	}

	if r.EventReasons[0].Reason != "FailedScheduling" || r.EventReasons[0].Count != 4 {
		t.Fatalf("expected FailedScheduling count 4, got %+v", r.EventReasons[0])
	}

	if len(r.AbnormalEvents) != 3 {
		t.Fatalf("expected 3 abnormal events, got %d", len(r.AbnormalEvents))
	}
}
