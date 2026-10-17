package signals

import (
	"testing"
	"time"
)

func TestAnalyzerFlagsLargeWallet(t *testing.T) {
	a := NewAnalyzer(1000)
	a.Observe(Event{Asset: "TOKEN", Type: "wallet_transfer", Amount: 2, PriceUSD: 600, ObservedAt: time.Now()})
	s := a.Snapshot()
	if len(s) != 1 || !s[0].LargeWallet || s[0].Signal != "watch" {
		t.Fatalf("unexpected snapshot: %+v", s)
	}
}
