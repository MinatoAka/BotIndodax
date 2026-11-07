package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getTickerPrice(pair string) (float64, error) {
	url := fmt.Sprintf("https://indodax.com/api/ticker/%s", pair)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP Error: %d", resp.StatusCode)
	}

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
