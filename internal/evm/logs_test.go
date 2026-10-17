package evm

import (
	"testing"
	"time"
)

func TestParseTransferLog(t *testing.T) {
	payload := []byte(`{"method":"eth_subscription","params":{"result":{"address":"0xtoken","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0xfrom","0xto"],"data":"0x01","blockNumber":"0x2a","transactionHash":"0xtx"}}}`)
	event, ok, err := ParseLog(payload, "ethereum", time.Unix(0, 0))
	if err != nil || !ok || event.Type != "wallet_transfer" || event.ID != "0xtx" {
		t.Fatalf("unexpected result: event=%+v ok=%v err=%v", event, ok, err)
	}
}

func TestParsePairCreatedLog(t *testing.T) {
	payload := []byte(`{"method":"eth_subscription","params":{"result":{"address":"0xfactory","topics":["0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9"],"transactionHash":"0xtx"}}}`)
	event, ok, err := ParseLog(payload, "base", time.Now())
	if err != nil || !ok || event.Type != "token_launch" {
		t.Fatalf("unexpected result: event=%+v ok=%v err=%v", event, ok, err)
	}
}
