package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigValidateRejectsInvalidTradingValues(t *testing.T) {
	config := Config{
		ApiKey: "key", SecretKey: "secret", Pair: "ethidr",
		BuyTrigger: 100, SellTrigger: 110, ModalPerTrade: 1000,
		MinQty: 0.001, CooldownMs: 1000, MaxDailyTrades: 2,
	}
	if err := config.validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	config.SellTrigger = config.BuyTrigger
	if err := config.validate(); err == nil {
		t.Fatal("expected invalid trigger relationship to be rejected")
	}
}

func TestLoadStateOnlyDefaultsWhenFileIsMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "state.json")
	state, err := loadState(missing)
	if err != nil || state.Status != "KOSONG" {
		t.Fatalf("missing state should default safely; got state=%+v err=%v", state, err)
	}

	corrupt := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(corrupt, []byte("{"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := loadState(corrupt); err == nil {
		t.Fatal("expected malformed state to return an error")
	}
}

func TestBuildQueryStringDeterministic(t *testing.T) {
	got := buildQueryString(map[string]string{"z": "last", "a": "first", "empty": ""})
	if got != "a=first&z=last" {
		t.Fatalf("buildQueryString() = %q", got)
	}
}
