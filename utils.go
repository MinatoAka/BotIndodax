package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"os"
	"sort"
	"math"
	"strings"
	"strconv"
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
	
	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}
	
	tempFile, err := os.CreateTemp(filepath.Dir(filename), ".state-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary state file: %w", err)
	}
	tempName := tempFile.Name()
	defer os.Remove(tempName) 
	if err := tempFile.Chmod(0600); err != nil {
		tempFile.Close()
		return fmt.Errorf("secure temporary state file: %w", err)
	}
	
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return fmt.Errorf("write temporary state file: %w", err)
	}
	
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		return fmt.Errorf("sync temporary state file: %w", err)
	}
	
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temporary state file: %w", err)
	}
	
	if err := os.Rename(tempName, filename); err != nil {
		return fmt.Errorf("replace state file: %w", err)
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

func formatRupiah(n float64) string {
	val := int64(math.Round(n))
	sign := ""
	if val < 0 {
		sign = "-"
		val = -val
	}
	s := strconv.FormatInt(val, 10)
	var b strings.Builder
	b.WriteString(sign)
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(c)
	}
	return b.String()
}

func countdown(d time.Duration) {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")
	
	if d <= 0 {
		fmt.Printf("\r⏳ Menunggu 0 Detik | [%s] 100%% | Sisa: 00:00:00 \033[K", strings.Repeat("\033[32m█\033[0m", 20))
		fmt.Println("\n✅ Waktu tunggu selesai | Lanjut cek market...")
		return
	}

	startTime := time.Now()
	endTime := startTime.Add(d)

	totalSeconds := int(math.Ceil(d.Seconds()))
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	var durationStr string
	switch {
	case hours > 0 && minutes > 0:
		durationStr = fmt.Sprintf("%d Jam %d Menit", hours, minutes)
	case hours > 0:
		durationStr = fmt.Sprintf("%d Jam", hours)
	case minutes > 0 && seconds > 0:
		durationStr = fmt.Sprintf("%d Menit %d Detik", minutes, seconds)
	case minutes > 0:
		durationStr = fmt.Sprintf("%d Menit", minutes)
	default:
		durationStr = fmt.Sprintf("%d Detik", seconds)
	}

	barLength := 20

	for {
		remaining := time.Until(endTime)
		if remaining < 0 {
			remaining = 0
		}

		percent := (1 - remaining.Seconds()/d.Seconds()) * 100
		if percent < 0 {
			percent = 0
		}
		if percent > 100 {
			percent = 100
		}

		remainingSeconds := int(math.Ceil(remaining.Seconds()))
		h := remainingSeconds / 3600
		m := (remainingSeconds % 3600) / 60
		s := remainingSeconds % 60

		if remaining == 0 {
			percent = 100
		}

		filled := int((percent / 100) * float64(barLength))
		if filled > barLength {
			filled = barLength
		}
		if filled < 0 {
			filled = 0
		}

		filledStr := strings.Repeat("\033[32m█\033[0m", filled)
		emptyStr := strings.Repeat("\033[90m░\033[0m", barLength-filled)
		bar := filledStr + emptyStr

		fmt.Printf("\r⏳ Menunggu %s | [%s] %.0f%% | Sisa: %02d:%02d:%02d \033[K", durationStr, bar, percent, h, m, s)

		if remaining == 0 {
			fmt.Println("\n✅ Waktu tunggu selesai | Lanjut cek market...")
			return
		}

		wait := time.Second
		if remaining < wait {
			wait = remaining
		}
		time.Sleep(wait)
	}
}

func resetDailyCounterIfNeeded(state *State, stateFile string, maxTrades int) {
	now := time.Now()
	today := now.Add(-7 * time.Hour).Format("2006-01-02")

	if state.LastTradeDay != today {
		state.DailyTrades = 0
		state.LastTradeDay = today
		log.Printf("🔄 Hari trading baru (UTC) terdeteksi (%s) | Limit harian di-reset ke 0.", today)
		
		if err := saveState(stateFile, *state); err != nil {
			log.Printf("❌ Gagal menyimpan reset limit harian: %v", err)
		}
	} else {
		sisa := maxTrades - state.DailyTrades
		if sisa < 0 {
			sisa = 0
		}
		log.Printf("📊 Limit harian: %d/%d terpakai | Sisa kuota: %d trade.", state.DailyTrades, maxTrades, sisa)
	}
}

func clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func printBanner() {
	cyan := "\033[36m"
	white := "\033[37m"
	bold := "\033[1m"
	reset := "\033[0m"

	fmt.Printf("%s", cyan)
	fmt.Println(`
     $$$$$$\  $$\   $$\  PROJECT™
    $$  __$$\ $$ | $$  | 
    $$ /  $$ |$$ |$$  / $$$$$$\  
    $$$$$$$$ |$$$$$  /  \____$$\ 
    $$  __$$ |$$  $$<   $$$$$$$ |
    $$ |  $$ |$$ |\$$\ $$  __$$ |
    $$ |  $$ |$$ | \$$\\$$$$$$$ |
    \__|  \__|\__|  \__|\_______|`)
    fmt.Printf("%s", reset)
	fmt.Printf("%s%s  BotIndodax — Trading Engine%s\n\n", white, bold, reset)
}
