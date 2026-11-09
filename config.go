package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
)

type Config struct {
	ApiKey          string  `json:"apiKey"`
	SecretKey       string  `json:"secretKey"`
	Pair            string  `json:"pair"`
	BuyTrigger      float64 `json:"buyTrigger"`
	SellTrigger     float64 `json:"sellTrigger"`
	ModalPerTrade   float64 `json:"modalPerTrade"`
	MinQty          float64 `json:"minQty"`
	CooldownMinutes int     `json:"cooldownMinutes"`
	PollingInterval int     `json:"pollingInterval"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("read config: %w", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("parse config: %w", err)
	}

	config.ApiKey = strings.TrimSpace(config.ApiKey)
	config.SecretKey = strings.TrimSpace(config.SecretKey)
	config.Pair = strings.ToLower(strings.TrimSpace(config.Pair))
	if err := config.validate(); err != nil {
		return config, err
	}
	return config, nil
}

func (c Config) validate() error {
	if c.ApiKey == "" || c.SecretKey == "" {
		return fmt.Errorf("apiKey dan secretKey wajib diisi")
	}
	if !strings.HasSuffix(c.Pair, "idr") || len(c.Pair) <= len("idr") {
		return fmt.Errorf("pair harus berupa pair IDR yang valid, misalnya ethidr")
	}
	if !isFinitePositive(c.BuyTrigger) || !isFinitePositive(c.SellTrigger) {
		return fmt.Errorf("buyTrigger dan sellTrigger harus berupa angka positif")
	}
	if c.SellTrigger <= c.BuyTrigger {
		return fmt.Errorf("sellTrigger harus lebih tinggi dari buyTrigger")
	}
	if !isFinitePositive(c.ModalPerTrade) || !isFinitePositive(c.MinQty) {
		return fmt.Errorf("modalPerTrade dan minQty harus berupa angka positif")
	}
	if c.CooldownMinutes < 0 {
		return fmt.Errorf("cooldownMinutes tidak boleh negatif")
	}
	if c.PollingInterval <= 0 {
		return fmt.Errorf("pollingInterval harus lebih besar dari 0")
	}
	return nil
}

func isFinitePositive(value float64) bool {
	return value > 0 && !math.IsNaN(value) && !math.IsInf(value, 0)
}
