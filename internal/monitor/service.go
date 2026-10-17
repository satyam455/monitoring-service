package monitor

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"

	"log/slog"
	"smartx-monitor/internal/signals"
)

type Service struct {
	logger    *slog.Logger
	events    chan Event
	workers   int
	processed atomic.Uint64
	dropped   atomic.Uint64
	stopOnce  sync.Once
	stop      chan struct{}
	store     EventStore
	cache     interface {
		Set(context.Context, Event) error
	}
	analyzer *signals.Analyzer
	apiKey   string
}

type EventStore interface {
	Save(context.Context, Event) error
}

func (s *Service) SetStore(store EventStore) { s.store = store }
func (s *Service) SetCache(cache interface {
	Set(context.Context, Event) error
}) {
	s.cache = cache
}
func (s *Service) SetAnalyzer(analyzer *signals.Analyzer) { s.analyzer = analyzer }
func (s *Service) SetAPIKey(apiKey string)                { s.apiKey = apiKey }

func NewService(logger *slog.Logger, workers int) *Service {
	if workers < 1 {
		workers = 1
	}
	return &Service{logger: logger, events: make(chan Event, workers*16), workers: workers, stop: make(chan struct{})}
}

func (s *Service) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for {
				select {
				case event := <-s.events:
					s.process(worker, event)
				case <-ctx.Done():
					return
				case <-s.stop:
					return
				}
			}
		}(i)
	}
	wg.Wait()
}

func (s *Service) Publish(event Event) bool {
	select {
	case s.events <- event:
		return true
	case <-s.stop:
		return false
	default:
		s.dropped.Add(1)
		return false
	}
}

func (s *Service) Stop() { s.stopOnce.Do(func() { close(s.stop) }) }

func (s *Service) process(worker int, event Event) {
	s.processed.Add(1)
	if s.analyzer != nil {
		s.analyzer.Observe(signals.Event{Asset: event.Asset, Wallet: event.Wallet, Type: string(event.Type), Amount: event.Amount, PriceUSD: event.PriceUSD, Liquidity: event.Liquidity, ObservedAt: event.ObservedAt})
	}
	if s.store != nil {
		if err := s.store.Save(context.Background(), event); err != nil {
			s.logger.Error("persist event", "event_id", event.ID, "error", err)
		}
	}
	if s.cache != nil {
		if err := s.cache.Set(context.Background(), event); err != nil {
			s.logger.Error("cache event", "event_id", event.ID, "error", err)
		}
	}
	s.logger.Debug("event processed", "worker", worker, "event_id", event.ID, "type", event.Type, "chain", event.Chain)
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	})
	mux.HandleFunc("GET /v1/signals", func(w http.ResponseWriter, r *http.Request) {
		if s.apiKey != "" && r.Header.Get("X-API-Key") != s.apiKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if s.analyzer == nil {
			http.Error(w, "signals unavailable", http.StatusServiceUnavailable)
			return
		}
		_ = json.NewEncoder(w).Encode(s.analyzer.Snapshot())
	})
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = w.Write([]byte("# HELP smartx_events_processed_total Events successfully processed.\n# TYPE smartx_events_processed_total counter\nsmartx_events_processed_total " + strconv.FormatUint(s.processed.Load(), 10) + "\n# HELP smartx_events_dropped_total Events rejected because the pipeline was full.\n# TYPE smartx_events_dropped_total counter\nsmartx_events_dropped_total " + strconv.FormatUint(s.dropped.Load(), 10) + "\n"))
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ready", "processed": s.processed.Load()})
	})
	return mux
}
