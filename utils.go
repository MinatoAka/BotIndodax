package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"
)

type State struct {
	Status      string `json:"status"`
	LastTradeAt int64  `json:"lastTradeAt"`
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c 
		
		fmt.Print("\033[?25h") 
		os.Exit(0)
	}()
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
	if state.LastTradeAt < 0 {
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

func parseFormattedQty(qty float64) (float64, error) {
	parsed, err := strconv.ParseFloat(formatQty(qty), 64)
	if err != nil || parsed <= 0 {
		return 0, err
	}
	return parsed, nil
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

func typeWriter(text string, speed time.Duration) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("%s ", timestamp)
	
	chars := []rune(text)
	for i, char := range chars {
		fmt.Printf("%c", char)
		
		if i == len(chars)-1 {
			continue
		}
		if unicode.IsSpace(char) {
			time.Sleep(speed / 2)
		} else {
			time.Sleep(speed)
		}
	}
	fmt.Println()
}

func countdown(d time.Duration) {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")
	
	if d <= 0 {
		fmt.Printf("\r[÷] Time remaining:  100%% [%s] 00:00:00 \033[K", strings.Repeat("\033[32m█\033[0m", 20))
		time.Sleep(1 * time.Second)
		fmt.Print("\r\033[2K")
		typeWriter("[•] =============================================", 40*time.Millisecond)
		typeWriter("[✓] Waktu tunggu selesai | Lanjut cek market...", 40*time.Millisecond)
		typeWriter("[•] =============================================", 40*time.Millisecond)
		return
	}

	startTime := time.Now()
	endTime := startTime.Add(d)

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

		fmt.Printf("\r[÷] Time remaining: %.0f%% [%s] %02d:%02d:%02d \033[K", percent, bar, h, m, s)

		if remaining == 0 {
			time.Sleep(1 * time.Second)
			fmt.Print("\r\033[2K")
			typeWriter("[•] =============================================", 40*time.Millisecond)
			typeWriter("[✓] Waktu tunggu selesai | Lanjut cek market...", 40*time.Millisecond)
			typeWriter("[•] =============================================", 40*time.Millisecond)
			return
		}

		wait := time.Second
		if remaining < wait {
			wait = remaining
		}
		time.Sleep(wait)
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
	fmt.Printf("%s%s  BotIndodax — Trading Engine V1.1%s\n\n", white, bold, reset)
}
