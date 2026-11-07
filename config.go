package main

import (
	"encoding/json"
	"io/ioutil"
)

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

func loadConfig(filename string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	
	config.Pair = normalizeText(config.Pair)
	
	return config, err
}
