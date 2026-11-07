package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 15 * time.Second}

func getTickerPrice(pair string) (float64, error) {
	endpoint := "https://indodax.com/api/ticker/" + url.PathEscape(strings.ToLower(pair))
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("create ticker request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request ticker: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return 0, fmt.Errorf("ticker API HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var response struct {
		Ticker struct {
			Last string `json:"last"`
		} `json:"ticker"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&response); err != nil {
		return 0, fmt.Errorf("decode ticker response: %w", err)
	}
	price, err := strconv.ParseFloat(response.Ticker.Last, 64)
	if err != nil || !isFinitePositive(price) {
		return 0, fmt.Errorf("ticker API returned invalid last price %q", response.Ticker.Last)
	}
	return price, nil
}

func getBalances(config Config) (float64, float64, error) {
	params := map[string]string{
		"timestamp":  strconv.FormatInt(time.Now().UnixMilli(), 10),
		"recvWindow": "5000",
	}
	query := buildQueryString(params)
	signature := getIndodaxSignature(query, config.SecretKey)
	endpoint := "https://api.indodax.com/api/v2/account?" + query
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("create account request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-APIKEY", config.ApiKey)
	req.Header.Set("Sign", signature)

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("request account balances: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return 0, 0, fmt.Errorf("read account response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, 0, fmt.Errorf("account API HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var response struct {
		Balances []struct {
			Asset string `json:"asset"`
			Free  string `json:"free"`
		} `json:"balances"`
		Code *int   `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, 0, fmt.Errorf("decode account response: %w", err)
	}
	if response.Code != nil || response.Balances == nil {
		return 0, 0, fmt.Errorf("account API error (code=%v): %s", response.Code, response.Msg)
	}

	baseAsset := strings.ToUpper(strings.TrimSuffix(config.Pair, "idr"))
	var idrBalance, assetBalance float64
	for _, balance := range response.Balances {
		asset := strings.ToUpper(balance.Asset)
		free, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil || free < 0 || math.IsNaN(free) || math.IsInf(free, 0) {
			return 0, 0, fmt.Errorf("account API returned invalid free balance for %q: %q", asset, balance.Free)
		}
		switch asset {
		case "IDR":
			idrBalance = free
		case baseAsset:
			assetBalance = free
		}
	}
	return idrBalance, assetBalance, nil
}
