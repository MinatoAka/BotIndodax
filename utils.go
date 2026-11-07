package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type State struct {
	Status       string `json:"status"`
	LastTradeAt  int64  `json:"lastTradeAt"`
	DailyTrades  int    `json:"dailyTrades"`
	LastTradeDay string `json:"lastTradeDay"`
}

func loadState(filename string) (State, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return State{Status: "KOSONG"}, nil
		}
		return State{}, fmt.Errorf("read state: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}, fmt.Errorf("parse state: %w", err)
	}
	state.Status = strings.ToUpper(strings.TrimSpace(state.Status))
	if state.Status == "" {
		state.Status = "KOSONG"
	}
	if state.Status != "KOSONG" && state.Status != "PEGANG" {
		return State{}, fmt.Errorf("status state tidak dikenal: %q", state.Status)
	}
	if state.DailyTrades < 0 || state.LastTradeAt < 0 {
		return State{}, fmt.Errorf("state berisi nilai transaksi negatif")
	}
	return state, nil
}

func saveState(filename string, state State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}
	
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("write state: %w", err)
	}
	return nil
}

func getIndodaxSignature(payloadString, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	_, _ = h.Write([]byte(payloadString))
	return hex.EncodeToString(h.Sum(nil))
}

func buildQueryString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	values := url.Values{}
	for _, key := range keys {
		if params[key] != "" {
			values.Set(key, params[key])
		}
	}
	return values.Encode()
}

func formatQty(qty float64) string {
	if qty <= 0 {
		return "0"
	}
	formatted := fmt.Sprintf("%.8f", qty)
	formatted = strings.TrimRight(formatted, "0")
	return strings.TrimRight(formatted, ".")
}

func countdown(d time.Duration) {
	endTime := time.Now().Add(d)

	var durationStr string
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 && minutes == 0 {
		durationStr = fmt.Sprintf("%d Jam", hours)
	} else if hours > 0 {
		durationStr = fmt.Sprintf("%d Jam %d Menit", hours, minutes)
	} else {
		durationStr = fmt.Sprintf("%d Menit", minutes)
	}

	log.Printf("⏳ Menunggu %s untuk iterasi berikutnya...", durationStr)

	for {
		remaining := time.Until(endTime)
		if remaining <= 0 {
			fmt.Println("\n✅ Waktu tunggu selesai. Lanjut cek pasar...")
			break
		}

		h := int(remaining.Hours())
		m := int(remaining.Minutes()) % 60
		s := int(remaining.Seconds()) % 60

		fmt.Printf("\r⏳ Sisa waktu: %02d:%02d:%02d \033[K", h, m, s)

		time.Sleep(1 * time.Second)
	}
}
