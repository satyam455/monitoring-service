package monitor

import "testing"

func TestNewServiceProvidesDefaultLogger(t *testing.T) {
	service := NewService(nil, 1)
	if service.logger == nil {
		t.Fatal("expected default logger")
	}
}
