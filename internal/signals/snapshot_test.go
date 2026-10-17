package signals

import (
	"testing"
	"time"
)

func TestSnapshotIsSortedAndScoreCapped(t *testing.T) {
	a := NewAnalyzer(1)
	for i := 0; i < 20; i++ {
		a.Observe(Event{Asset: "B", Type: "token_launch", Amount: 2, PriceUSD: 2, ObservedAt: time.Now()})
	}
	a.Observe(Event{Asset: "A", Type: "token_launch", ObservedAt: time.Now()})
	s := a.Snapshot()
	if s[0].Asset != "B" || s[0].Score != 100 {
		t.Fatalf("unexpected snapshot: %+v", s)
	}
}
