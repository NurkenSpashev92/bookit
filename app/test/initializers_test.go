package test

import (
	"testing"

	"github.com/nurkenspashev92/bookit/internal/initializers"
)

func TestNewFiberConfig(t *testing.T) {
	cfg := initializers.NewFiberConfig()

	if cfg.ServerHeader != "Bookit" {
		t.Errorf("ServerHeader = %q, want Bookit", cfg.ServerHeader)
	}
	if cfg.AppName != "Bookit App v0.1-beta" {
		t.Errorf("AppName = %q, want Bookit App v0.1-beta", cfg.AppName)
	}
	if !cfg.CaseSensitive {
		t.Error("CaseSensitive should be true")
	}
}

func TestNewLogger(t *testing.T) {
	handler := initializers.NewLogger()
	if handler == nil {
		t.Fatal("NewLogger returned nil")
	}
}
