package monitor

import (
	"testing"
	"time"
)

func TestPublishDeduplicatesEventIDs(t *testing.T) {
	s := NewService(nil, 1)
	e := Event{ID: "same", Chain: ChainSolana, Type: EventTokenLaunch, ObservedAt: time.Now()}
	if !s.Publish(e) || !s.Publish(e) {
		t.Fatal("expected duplicate publish calls to be accepted")
	}
	if got := len(s.events); got != 1 {
		t.Fatalf("expected one queued event, got %d", got)
	}
}
