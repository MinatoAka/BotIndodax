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
		MinQty: 0.001, 
		CooldownMinutes: 10,
		PollingInterval: 15,
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

func TestSaveStateUsesPrivatePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := saveState(path, State{Status: "KOSONG"}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	
	if got := info.Mode().Perm(); got != 0600 {
		t.Fatalf("state permissions = %04o, want 0600", got)
	}
}

func TestBuildQueryStringDeterministic(t *testing.T) {
	got := buildQueryString(map[string]string{"z": "last", "a": "first", "empty": ""})
	if got != "a=first&z=last" {
		t.Fatalf("buildQueryString() = %q", got)
	}
}
