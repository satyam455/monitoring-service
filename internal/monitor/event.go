package monitor

import (
	"fmt"
	"time"
)

type Chain string

const (
	ChainSolana Chain = "solana"
	ChainEVM    Chain = "evm"
)

type EventType string

const (
	EventTokenLaunch     EventType = "token_launch"
	EventPriceMovement   EventType = "price_movement"
	EventLiquidityChange EventType = "liquidity_change"
	EventWalletTransfer  EventType = "wallet_transfer"
)

// Event is the chain-neutral shape all adapters must emit.
// Raw chain payloads should stay inside their adapter and never reach consumers.
type Event struct {
	ID         string            `json:"id"`
	Chain      Chain             `json:"chain"`
	Network    string            `json:"network"`
	Type       EventType         `json:"type"`
	Asset      string            `json:"asset,omitempty"`
	Wallet     string            `json:"wallet,omitempty"`
	Amount     float64           `json:"amount,omitempty"`
	PriceUSD   float64           `json:"price_usd,omitempty"`
	Liquidity  float64           `json:"liquidity_usd,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	ObservedAt time.Time         `json:"observed_at"`
}

func (e Event) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("event id is required")
	}
	if e.Chain != ChainSolana && e.Chain != ChainEVM {
		return fmt.Errorf("unsupported chain: %q", e.Chain)
	}
	if e.Type != EventTokenLaunch && e.Type != EventPriceMovement && e.Type != EventLiquidityChange && e.Type != EventWalletTransfer {
		return fmt.Errorf("unsupported event type: %q", e.Type)
	}
	if e.ObservedAt.IsZero() {
		return fmt.Errorf("observed_at is required")
	}
	return nil
}
