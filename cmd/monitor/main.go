package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smartx-monitor/internal/config"
	"smartx-monitor/internal/evm"
	"smartx-monitor/internal/monitor"
	"smartx-monitor/internal/signals"
	"smartx-monitor/internal/solana"
	"smartx-monitor/internal/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	service := monitor.NewService(logger, cfg.Workers)
	service.SetAnalyzer(signals.NewAnalyzer(cfg.LargeWalletUSD))
	service.SetAPIKey(cfg.APIKey)
	var postgres *storage.PostgresStore
	if cfg.PostgresURL != "" {
		var err error
		postgres, err = storage.NewPostgres(ctx, cfg.PostgresURL)
		if err != nil {
			logger.Error("postgres unavailable", "error", err)
		} else if err = postgres.Migrate(ctx); err != nil {
			logger.Error("postgres migration failed", "error", err)
		} else {
			service.SetStore(postgres)
		}
	}
	if cfg.RedisAddress != "" {
		service.SetCache(storage.NewRedis(cfg.RedisAddress, 15*time.Minute))
	}
	go service.Run(ctx)
	go func() {
		adapter := &solana.Adapter{URL: cfg.SolanaWSURL, Network: cfg.SolanaNetwork, Publisher: service, Logger: logger}
		adapter.RunForever(ctx)
	}()
	if cfg.EVMWSURL != "" {
		go func() {
			adapter := &evm.Adapter{URL: cfg.EVMWSURL, Network: cfg.EVMNetwork, Publisher: service, Logger: logger}
			adapter.RunForever(ctx)
		}()
	}

	server := &http.Server{Addr: cfg.HTTPAddr, Handler: service.Handler()}
	go func() {
		logger.Info("monitor service started", "addr", cfg.HTTPAddr, "workers", cfg.Workers)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server stopped", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
	service.Stop()
	if postgres != nil {
		postgres.Close()
	}
	logger.Info("monitor service stopped")
}
