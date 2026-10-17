package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr        string
	Workers         int
	ShutdownTimeout time.Duration
	SolanaWSURL     string
	SolanaNetwork   string
	EVMWSURL        string
	EVMNetwork      string
	PostgresURL     string
	RedisAddress    string
	LargeWalletUSD  float64
	APIKey          string
}

func Load() Config {
	return Config{
		HTTPAddr:        env("HTTP_ADDR", ":8080"),
		Workers:         positiveInt("WORKERS", 4),
		ShutdownTimeout: positiveDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		SolanaWSURL:     env("SOLANA_WS_URL", "wss://api.devnet.solana.com"),
		SolanaNetwork:   env("SOLANA_NETWORK", "devnet"),
		EVMWSURL:        env("EVM_WS_URL", ""),
		EVMNetwork:      env("EVM_NETWORK", "ethereum"),
		PostgresURL:     env("POSTGRES_URL", ""),
		RedisAddress:    env("REDIS_ADDRESS", ""),
		LargeWalletUSD:  positiveFloat("LARGE_WALLET_USD", 100000),
		APIKey:          os.Getenv("API_KEY"),
	}
}

func (c Config) Validate() error {
	if c.HTTPAddr == "" {
		return fmt.Errorf("HTTP_ADDR cannot be empty")
	}
	if c.Workers < 1 {
		return fmt.Errorf("WORKERS must be positive")
	}
	if c.SolanaWSURL == "" || c.SolanaNetwork == "" {
		return fmt.Errorf("Solana websocket and network are required")
	}
	return nil
}

func positiveFloat(key string, fallback float64) float64 {
	value, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func positiveInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func positiveDuration(key string, fallback time.Duration) time.Duration {
	value, err := time.ParseDuration(os.Getenv(key))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
