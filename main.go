package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ==========================================
// 1. STRUCTS
// ==========================================
type Config struct {
	ApiKey         string  `json:"apiKey"`
	SecretKey      string  `json:"secretKey"`
	Pair           string  `json:"pair"`
	BuyTrigger     float64 `json:"buyTrigger"`
	SellTrigger    float64 `json:"sellTrigger"`
	ModalPerTrade  float64 `json:"modalPerTrade"`
	MinQty         float64 `json:"minQty"`
	CooldownMs     int64   `json:"cooldownMs"`
	MaxDailyTrades int     `json:"maxDailyTrades"`
}

type State struct {
	Status       string `json:"status"` // "KOSONG" atau "PEGANG"
	LastTradeAt  int64  `json:"lastTradeAt"`
	DailyTrades  int    `json:"dailyTrades"`
	LastTradeDay string `json:"lastTradeDay"`
}

// ==========================================
// 2. UTILS
// ==========================================
func readJSONFile(filename string, target interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func writeJSONFile(filename string, data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0644)
}

// Membuat Signature HMAC-SHA256 ala Indodax
func getIndodaxSignature(payloadString, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(payloadString))
	return hex.EncodeToString(h.Sum(nil))
}

// Build Query String berurut abjad
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

// ==========================================
// 3. MARKET
// ==========================================
func getTickerPrice(pair string) (float64, error) {
	url := fmt.Sprintf("https://indodax.com/api/ticker/%s", pair)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Pastikan status 200 OK
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP Error: %d", resp.StatusCode)
	}

	// Parse JSON
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	ticker := data["ticker"].(map[string]interface{})
	lastPrice, _ := strconv.ParseFloat(ticker["last"].(string), 64)
	return lastPrice, nil
}

func getBalances(config Config) (float64, float64, error) {
	params := map[string]string{
		"timestamp":  strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
		"recvWindow": "5000",
	}
	query := buildQueryString(params)
	sign := getIndodaxSignature(query, config.SecretKey)

	req, _ := http.NewRequest("GET", "https://api.indodax.com/api/v2/account?"+query, nil)
	req.Header.Set("X-APIKEY", config.ApiKey)
	req.Header.Set("Sign", sign)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var resData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&resData)

	var idrSaldo, asetSaldo float64
	if balances, ok := resData["balances"].([]interface{}); ok {
		for _, item := range balances {
			b := item.(map[string]interface{})
			asset := strings.ToLower(b["asset"].(string))
			freeStr := b["free"].(string)
			free, _ := strconv.ParseFloat(freeStr, 64)

			if asset == "idr" {
				idrSaldo = free
			}
			baseAsset := strings.Replace(config.Pair, "idr", "", 1)
			if asset == baseAsset {
				asetSaldo = free
			}
		}
	}
	return idrSaldo, asetSaldo, nil
}

// ==========================================
// 4. BOT LOGIC
// ==========================================
func kirimOrder(config Config, side string, price float64, qty float64) bool {
	params := map[string]string{
		"symbol":     config.Pair,
		"side":       side,
		"type":       "MARKET",
		"timestamp":  strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
		"recvWindow": "5000",
	}

	if side == "BUY" {
		// BUY MARKET butuh nominal IDR
		quoteAmount := math.Floor(qty * price)
		params["quoteOrderQty"] = strconv.FormatFloat(quoteAmount, 'f', -1, 64)
	} else {
		// SELL MARKET butuh jumlah koin
		params["quantity"] = strconv.FormatFloat(qty, 'f', 8, 64)
	}

	query := buildQueryString(params)
	sign := getIndodaxSignature(query, config.SecretKey)

	req, _ := http.NewRequest("POST", "https://api.indodax.com/api/v2/order", bytes.NewBufferString(query))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-APIKEY", config.ApiKey)
	req.Header.Set("Sign", sign)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Error sistem saat kirim order:", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("✅ Order %s terkirim. Response: %s\n", side, string(body))
	return true
}

// ==========================================
// 5. MAIN LOOP (PENGGANTI Trigger.gs)
// ==========================================
func main() {
	fmt.Println("🤖 Bot Aji Mumpung (Golang Version) Menyala!")
	
	for {
		// 1. Baca Config & State (Selalu baca ulang biar bisa diubah tanpa restart)
		var config Config
		if err := readJSONFile("config.json", &config); err != nil {
			fmt.Println("❌ Gagal baca config.json:", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		var state State
		readJSONFile("state.json", &state)
		if state.Status == "" {
			state.Status = "KOSONG"
		}

		// 2. Ambil Harga
		hargaSekarang, err := getTickerPrice(config.Pair)
		if err != nil {
			fmt.Println("❌ Gagal ambil harga:", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		fmt.Printf("📝 Log: Pair=%s | Harga=Rp %.0f | Status=%s\n", strings.ToUpper(config.Pair), hargaSekarang, state.Status)

		// 3. Cek Sinyal
		isBuySignal := state.Status == "KOSONG" && hargaSekarang <= config.BuyTrigger
		isSellSignal := state.Status == "PEGANG" && hargaSekarang >= config.SellTrigger

		if !isBuySignal && !isSellSignal {
			fmt.Println("🟡 Sabar bos kita tunggu momen yang pas...")
		} else {
			fmt.Println("🔔 Sinyal terdeteksi! Mengecek saldo...")
			idrSaldo, asetSaldo, err := getBalances(config)
			if err != nil {
				fmt.Println("❌ Gagal ambil saldo:", err)
			} else {
				// Cek Cooldown & Daily Limit
				now := time.Now().UnixMilli()
				today := time.Now().Format("2006-01-02")
				if state.LastTradeDay != today {
					state.DailyTrades = 0 // Reset harian
				}

				isCooldown := (now - state.LastTradeAt) < config.CooldownMs
				isLimitReached := state.DailyTrades >= config.MaxDailyTrades

				if isCooldown {
					fmt.Println("🛡️ Cooldown aktif, belum bisa trade.")
				} else if isLimitReached {
					fmt.Println("🛡️ Limit harian tercapai.")
				} else {
					// EKSEKUSI
					if isBuySignal && idrSaldo >= config.ModalPerTrade {
						qty := config.ModalPerTrade / hargaSekarang
						if qty >= config.MinQty {
							if kirimOrder(config, "BUY", hargaSekarang, qty) {
								state.Status = "PEGANG"
								state.LastTradeAt = now
								state.DailyTrades++
								state.LastTradeDay = today
								writeJSONFile("state.json", state)
							}
						} else {
							fmt.Println("⚠️ Qty di bawah minimum Indodax.")
						}
					} else if isSellSignal && asetSaldo > 0 {
						if kirimOrder(config, "SELL", hargaSekarang, asetSaldo) {
							state.Status = "KOSONG"
							state.LastTradeAt = now
							state.DailyTrades++
							state.LastTradeDay = today
							writeJSONFile("state.json", state)
						}
					}
				}
			}
		}

		// 4. Jeda sebelum loop berikutnya (Bisa disesuaikan)
		fmt.Println("⏳ Menunggu 1 jam untuk iterasi berikutnya...\n")
		time.Sleep(1 * time.Hour)
	}
}
