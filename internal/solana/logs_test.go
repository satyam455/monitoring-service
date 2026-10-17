package solana

import (
	"testing"
	"time"
)

func TestParseTokenLaunch(t *testing.T) {
	payload := []byte(`{"method":"logsNotification","params":{"result":{"context":{"slot":42},"value":{"signature":"sig-1","err":null,"logs":["Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss658c5DA invoke [1]","Program log: Instruction: InitializeMint"]}}}}`)
	event, ok, err := ParseTokenLaunch(payload, "devnet", time.Unix(0, 0))
	if err != nil || !ok {
		t.Fatalf("expected token launch, ok=%v err=%v", ok, err)
	}
	if event.ID != "sig-1" || event.Network != "devnet" || event.Metadata["slot"] != "42" {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestParseTokenLaunchIgnoresOtherLogs(t *testing.T) {
	payload := []byte(`{"method":"logsNotification","params":{"result":{"context":{"slot":42},"value":{"signature":"sig-1","err":null,"logs":["Program log: Instruction: Transfer"]}}}}`)
	_, ok, err := ParseTokenLaunch(payload, "devnet", time.Now())
	if err != nil || ok {
		t.Fatalf("expected non-launch, ok=%v err=%v", ok, err)
	}
}
