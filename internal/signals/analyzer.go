package signals

import (
	"sort"
	"sync"
	"time"
)

type Event struct {
	Asset, Wallet, Type         string
	Amount, PriceUSD, Liquidity float64
	ObservedAt                  time.Time
}

type Insight struct {
	Asset       string    `json:"asset"`
	Signal      string    `json:"signal"`
	Score       int       `json:"score"`
	LargeWallet bool      `json:"large_wallet"`
	Events      uint64    `json:"events"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Analyzer struct {
	mu        sync.RWMutex
	threshold float64
	insights  map[string]Insight
}

func NewAnalyzer(threshold float64) *Analyzer {
	if threshold <= 0 {
		threshold = 100000
	}
	return &Analyzer{threshold: threshold, insights: make(map[string]Insight)}
}
func (a *Analyzer) Observe(event Event) {
	asset := event.Asset
	if asset == "" {
		asset = event.Wallet
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	i := a.insights[asset]
	i.Asset, i.Events, i.UpdatedAt = asset, i.Events+1, event.ObservedAt
	if event.Type == "token_launch" {
		i.Score += 20
	}
	if event.Liquidity > 0 {
		i.Score += 10
	}
	i.LargeWallet = event.Amount*event.PriceUSD >= a.threshold
	if i.LargeWallet {
		i.Score += 15
	}
	if i.Score > 100 {
		i.Score = 100
	}
	if i.Score >= 35 {
		i.Signal = "strong"
	} else if i.Score >= 15 {
		i.Signal = "watch"
	} else {
		i.Signal = "neutral"
	}
	a.insights[asset] = i
}
func (a *Analyzer) Snapshot() []Insight {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make([]Insight, 0, len(a.insights))
	for _, i := range a.insights {
		result = append(result, i)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].Asset < result[j].Asset
		}
		return result[i].Score > result[j].Score
	})
	return result
}
