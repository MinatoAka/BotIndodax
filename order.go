package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

func kirimOrder(config Config, side string, price float64, qty float64) bool {
	params := map[string]string{
		"symbol":     config.Pair,
		"side":       side,
		"type":       "MARKET",
		"timestamp":  strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10),
		"recvWindow": "5000",
	}

	if side == "BUY" {
		quoteAmount := math.Floor(qty * price)
		params["quoteOrderQty"] = strconv.FormatFloat(quoteAmount, 'f', -1, 64)
	} else {
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
