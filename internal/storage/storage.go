package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"smartx-monitor/internal/monitor"
)

type PostgresStore struct{ pool *pgxpool.Pool }

func NewPostgres(ctx context.Context, url string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &PostgresStore{pool: pool}, nil
}
func (s *PostgresStore) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}
func (s *PostgresStore) Migrate(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS monitor_events (id TEXT PRIMARY KEY, chain TEXT NOT NULL, network TEXT NOT NULL, type TEXT NOT NULL, asset TEXT, wallet TEXT, amount DOUBLE PRECISION, price_usd DOUBLE PRECISION, liquidity_usd DOUBLE PRECISION, metadata JSONB NOT NULL DEFAULT '{}'::jsonb, observed_at TIMESTAMPTZ NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now())`)
	return err
}
func (s *PostgresStore) Save(ctx context.Context, event monitor.Event) error {
	metadata, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `INSERT INTO monitor_events (id,chain,network,type,asset,wallet,amount,price_usd,liquidity_usd,metadata,observed_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT (id) DO NOTHING`, event.ID, event.Chain, event.Network, event.Type, event.Asset, event.Wallet, event.Amount, event.PriceUSD, event.Liquidity, metadata, event.ObservedAt)
	return err
}

type RedisCache struct {
	address string
	ttl     time.Duration
}

func NewRedis(address string, ttl time.Duration) *RedisCache {
	return &RedisCache{address: address, ttl: ttl}
}
func (r *RedisCache) Set(ctx context.Context, event monitor.Event) error {
	if r.address == "" {
		return nil
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", r.address)
	if err != nil {
		return err
	}
	defer conn.Close()
	seconds := int(r.ttl / time.Second)
	if seconds < 1 {
		seconds = 60
	}
	if _, err = conn.Write(resp("SET", "smartx:event:"+event.ID, string(payload), "EX", strconv.Itoa(seconds))); err != nil {
		return err
	}
	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	return err
}
func resp(parts ...string) []byte {
	result := fmt.Sprintf("*%d\r\n", len(parts))
	for _, part := range parts {
		result += fmt.Sprintf("$%d\r\n%s\r\n", len(part), part)
	}
	return []byte(result)
}
