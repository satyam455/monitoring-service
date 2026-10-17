package solana

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"smartx-monitor/internal/monitor"
)

type Publisher interface{ Publish(monitor.Event) bool }

type Adapter struct {
	URL       string
	Network   string
	Publisher Publisher
	Logger    *slog.Logger
}

// RunForever reconnects after transient provider/network failures.
func (a *Adapter) RunForever(ctx context.Context) {
	backoff := time.Second
	for {
		err := a.Run(ctx)
		if ctx.Err() != nil {
			return
		}
		if a.Logger != nil {
			a.Logger.Warn("solana stream disconnected; reconnecting", "error", err, "backoff", backoff)
		}
		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

func (a *Adapter) Run(ctx context.Context) error {
	if a.URL == "" || a.Publisher == nil {
		return fmt.Errorf("solana adapter requires URL and publisher")
	}
	connection, _, err := websocket.DefaultDialer.DialContext(ctx, a.URL, nil)
	if err != nil {
		return fmt.Errorf("connect solana websocket: %w", err)
	}
	defer connection.Close()

	subscription := map[string]any{"jsonrpc": "2.0", "id": 1, "method": "logsSubscribe", "params": []any{
		map[string]any{"mentions": []string{TokenProgram, Token2022Program}}, map[string]any{"commitment": "confirmed"},
	}}
	if err := connection.WriteJSON(subscription); err != nil {
		return fmt.Errorf("subscribe to solana logs: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		_, message, err := connection.ReadMessage()
		if err != nil {
			return fmt.Errorf("read solana websocket: %w", err)
		}
		event, ok, err := ParseTokenLaunch(message, a.Network, time.Now().UTC())
		if err != nil {
			a.Logger.Warn("invalid solana notification", "error", err)
			continue
		}
		if ok {
			if err := event.Validate(); err != nil { a.Logger.Warn("invalid normalized solana event", "error", err); continue }
		}
		if ok && !a.Publisher.Publish(event) {
			a.Logger.Warn("dropping solana event: pipeline is full", "event_id", event.ID)
		}
	}
}
