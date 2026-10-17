package solana

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"smartx-monitor/internal/monitor"
)

const (
	TokenProgram     = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss658c5DA"
	Token2022Program = "TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"
)

type logNotification struct {
	Method string `json:"method"`
	Params struct {
		Result struct {
			Context struct {
				Slot uint64 `json:"slot"`
			} `json:"context"`
			Value struct {
				Signature string   `json:"signature"`
				Err       any      `json:"err"`
				Logs      []string `json:"logs"`
			} `json:"value"`
		} `json:"result"`
	} `json:"params"`
}

func ParseTokenLaunch(data []byte, network string, now time.Time) (monitor.Event, bool, error) {
	var notification logNotification
	if err := json.Unmarshal(data, &notification); err != nil {
		return monitor.Event{}, false, err
	}
	if notification.Method != "logsNotification" || notification.Params.Result.Value.Err != nil {
		return monitor.Event{}, false, nil
	}
	logs := notification.Params.Result.Value.Logs
	if !containsMintInstruction(logs) {
		return monitor.Event{}, false, nil
	}
	return monitor.Event{
		ID:         notification.Params.Result.Value.Signature,
		Chain:      monitor.ChainSolana,
		Network:    network,
		Type:       monitor.EventTokenLaunch,
		ObservedAt: now,
		Metadata: map[string]string{
			"slot":    stringValue(notification.Params.Result.Context.Slot),
			"source":  "logsSubscribe",
			"program": tokenProgramFromLogs(logs),
		},
	}, true, nil
}

func containsMintInstruction(logs []string) bool {
	for _, log := range logs {
		if strings.Contains(strings.ToLower(log), "instruction: initializemint") {
			return true
		}
	}
	return false
}

func tokenProgramFromLogs(logs []string) string {
	for _, log := range logs {
		if strings.Contains(log, Token2022Program) {
			return Token2022Program
		}
		if strings.Contains(log, TokenProgram) {
			return TokenProgram
		}
	}
	return "unknown"
}

func stringValue(value uint64) string {
	return fmt.Sprintf("%d", value)
}
