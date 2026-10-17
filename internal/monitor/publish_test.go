package monitor

import "testing"

func TestPublishRejectsInvalidEvents(t *testing.T) {
	service := NewService(nil, 1)
	if service.Publish(Event{ID: "missing-fields"}) {
		t.Fatal("expected invalid event to be rejected")
	}
	if service.dropped.Load() != 1 {
		t.Fatal("expected rejected event to increment dropped counter")
	}
}
