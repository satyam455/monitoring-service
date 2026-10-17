package config

import "testing"

func TestDefaultsUseDevnet(t *testing.T) {
	c := Load()
	if c.SolanaNetwork != "devnet" || c.SolanaWSURL != "wss://api.devnet.solana.com" {
		t.Fatalf("unsafe defaults: %+v", c)
	}
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
}
