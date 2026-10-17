package monitor

import (
	"testing"
	"time"
)

func TestEventValidation(t *testing.T) {
	event := Event{ID: "x", Chain: ChainSolana, Type: EventTokenLaunch, ObservedAt: time.Now()}
	if err := event.Validate(); err != nil {
		t.Fatal(err)
	}
	event.ID = ""
	if err := event.Validate(); err == nil {
		t.Fatal("expected missing id error")
	}
}
