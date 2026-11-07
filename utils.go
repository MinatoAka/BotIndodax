package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
)

// ==========================================
// STATE MANAGEMENT
// ==========================================
type State struct {
	Status       string `json:"status"` // "KOSONG" atau "PEGANG"
	LastTradeAt  int64  `json:"lastTradeAt"`
	DailyTrades  int    `json:"dailyTrades"`
	LastTradeDay string `json:"lastTradeDay"`
}

func loadState(filename string) (State, error) {
	var state State
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		// Kalau file belum ada, return default state
		return State{Status: "KOSONG"}, nil
	}
	err = json.Unmarshal(data, &state)
	if state.Status == "" {
		state.Status = "KOSONG"
	}
	return state, err
}

func saveState(filename string, state State) error {
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0644)
}

// ==========================================
// CRYPTO & UTILS
// ==========================================
func getIndodaxSignature(payloadString, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payloadString))
	return hex.EncodeToString(h.Sum(nil))
}

func buildQueryString(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryParts []string
	for _, k := range keys {
		if params[k] != "" {
			queryParts = append(queryParts, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
		}
	}
	return strings.Join(queryParts, "&")
}
