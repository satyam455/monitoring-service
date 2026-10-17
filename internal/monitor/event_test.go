package monitor

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
)

func TestServicePublishesEvents(t *testing.T) {
	service := NewService(slog.New(slog.NewTextHandler(io.Discard, nil)), 2)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { service.Run(ctx); close(done) }()

	if !service.Publish(Event{ID: "evt-1", Chain: ChainSolana, Type: EventTokenLaunch, ObservedAt: time.Now()}) {
		t.Fatal("expected event to be accepted")
	}
	service.Stop()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("service did not stop")
	}
}
