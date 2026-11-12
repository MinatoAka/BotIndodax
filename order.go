package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MinatoAka/hmacgo"
)

type OrderResult struct {
	OrderID     int64
	ExecutedQty float64
}

type indodaxAPIError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func kirimOrder(config Config, side string, price float64, qty float64) (OrderResult, error) {
	side = strings.ToUpper(strings.TrimSpace(side))
	if side != "BUY" && side != "SELL" {
		return OrderResult{}, fmt.Errorf("invalid order side %q", side)
	}
	if !isFinitePositive(price) || !isFinitePositive(qty) {
		return OrderResult{}, fmt.Errorf("order price and quantity must be positive finite numbers")
	}

	params := map[string]string{
		"symbol":     strings.ToUpper(config.Pair),
		"side":       side,
		"type":       "MARKET",
		"timestamp":  strconv.FormatInt(time.Now().UnixMilli(), 10),
		"recvWindow": "5000",
	}
	if side == "BUY" {
		quoteAmount := math.Floor(qty * price)
		if quoteAmount < 1 || math.IsInf(quoteAmount, 0) || math.IsNaN(quoteAmount) {
			return OrderResult{}, fmt.Errorf("calculated quoteOrderQty is invalid: %.0f", quoteAmount)
		}
		params["quoteOrderQty"] = fmt.Sprintf("%.0f", quoteAmount)
	} else {
		params["quantity"] = formatQty(qty)
	}

	body := hmacgo.BuildQuery(params)
	signature := hmacgo.SignSHA256(body, config.SecretKey)
	req, err := http.NewRequest(http.MethodPost, "https://api.indodax.com/api/v2/order", bytes.NewBufferString(body))
	if err != nil {
		return OrderResult{}, fmt.Errorf("create order request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-APIKEY", config.ApiKey)
	req.Header.Set("Sign", signature)

	// Gunakan Smart Hub (HTTPClient dari config.go) yang sudah ada Retry, Rate Limit, & UA
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return OrderResult{}, fmt.Errorf("order request outcome is unknown; reconcile the order in Indodax before restarting: %w", err)
	}
	defer resp.Body.Close()
	
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return OrderResult{}, fmt.Errorf("order submitted but response could not be read; reconcile in Indodax: %w", err)
	}
	
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr indodaxAPIError
		if json.Unmarshal(responseBody, &apiErr) == nil && apiErr.Msg != "" {
			return OrderResult{}, fmt.Errorf("order rejected (HTTP %d, code %d): %s", resp.StatusCode, apiErr.Code, apiErr.Msg)
		}
		return OrderResult{}, fmt.Errorf("order rejected (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var response struct {
		OrderID     int64  `json:"orderId"`
		Symbol      string `json:"symbol"`
		Side        string `json:"side"`
		ExecutedQty string `json:"executedQty"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return OrderResult{}, fmt.Errorf("order may have been accepted but response is invalid; reconcile in Indodax: %w", err)
	}
	if response.OrderID <= 0 || response.Symbol == "" || response.Side == "" {
		return OrderResult{}, fmt.Errorf("order may have been accepted but confirmation is incomplete; reconcile in Indodax: %s", strings.TrimSpace(string(responseBody)))
	}
	executedQty, err := strconv.ParseFloat(response.ExecutedQty, 64)
	if err != nil || executedQty < 0 || math.IsNaN(executedQty) || math.IsInf(executedQty, 0) {
		return OrderResult{}, fmt.Errorf("order %d may have been accepted but executed quantity is invalid; reconcile in Indodax", response.OrderID)
	}
	return OrderResult{OrderID: response.OrderID, ExecutedQty: executedQty}, nil
}
