package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"smartx-monitor/internal/monitor"
)

// PostgresStore is the persistence boundary. A PostgreSQL driver can be
// registered here without changing the monitoring pipeline contract.
type PostgresStore struct{ url string }

func NewPostgres(_ context.Context, url string) (*PostgresStore, error) {
	return nil, fmt.Errorf("postgres driver is not bundled in the Go 1.22 build: %s", url)
}
func (s *PostgresStore) Close() {}
func (s *PostgresStore) Migrate(context.Context) error {
	return fmt.Errorf("postgres driver is unavailable")
}
func (s *PostgresStore) Save(_ context.Context, event monitor.Event) error {
	return fmt.Errorf("postgres driver is unavailable for event %s", event.ID)
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
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}
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
