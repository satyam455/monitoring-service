package evm

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"smartx-monitor/internal/monitor"
)

const (
	ERC20TransferTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	PairCreatedTopic   = "0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9"
)

type logMessage struct {
	Method string `json:"method"`
	Params struct {
		Result struct {
			Address string   `json:"address"`
			Topics  []string `json:"topics"`
			Data    string   `json:"data"`
			Block   string   `json:"blockNumber"`
			Tx      string   `json:"transactionHash"`
		} `json:"result"`
	} `json:"params"`
}

func ParseLog(data []byte, network string, now time.Time) (monitor.Event, bool, error) {
	var message logMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return monitor.Event{}, false, err
	}
	if message.Method != "eth_subscription" || len(message.Params.Result.Topics) == 0 {
		return monitor.Event{}, false, nil
	}

	topic := strings.ToLower(message.Params.Result.Topics[0])
	eventType := monitor.EventType("")
	source := ""
	switch topic {
	case ERC20TransferTopic:
		eventType, source = monitor.EventWalletTransfer, "erc20_transfer"
	case PairCreatedTopic:
		eventType, source = monitor.EventTokenLaunch, "pair_created"
	default:
		return monitor.Event{}, false, nil
	}

	return monitor.Event{
		ID:         message.Params.Result.Tx,
		Chain:      monitor.ChainEVM,
		Network:    network,
		Type:       eventType,
		Asset:      message.Params.Result.Address,
		ObservedAt: now,
		Metadata: map[string]string{
			"block":  message.Params.Result.Block,
			"source": source,
			"topics": fmt.Sprintf("%d", len(message.Params.Result.Topics)),
		},
	}, true, nil
}
