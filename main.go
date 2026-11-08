package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

const stateFile = "state.json"

func main() {
	loc, errLoc := time.LoadLocation("Asia/Jakarta")
	if errLoc == nil {
		time.Local = loc
	}
	log.Println("BotIndodax starting")

	for {
		config, err := loadConfig("config.json")
		if err != nil {
			log.Printf("Konfigurasi tidak valid (tunggu 1 menit): %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		state, err := loadState(stateFile)
		if err != nil {
			log.Printf("State tidak dapat dibaca | Bot stopped! %v", err)
			return
		}

		resetDailyCounterIfNeeded(&state, stateFile, config.MaxDailyTrades)
		log.Printf("📊 Config: Buy = Rp %s | Sell = Rp %s", formatRupiah(config.BuyTrigger), formatRupiah(config.SellTrigger))

		price, err := getTickerPrice(config.Pair)
		if err != nil {
			log.Printf("Gagal mengambil harga: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		log.Printf("📊 Pair = %s | Harga = Rp %s | Status = %s", strings.ToUpper(config.Pair), formatRupiah(price))
		log.Printf("📊 Status = %s", state.Status)

		now := time.Now()
		nowMillis := now.UnixMilli()

		buySignal := state.Status == "KOSONG" && price <= config.BuyTrigger
		sellSignal := state.Status == "PEGANG" && price >= config.SellTrigger

		if buySignal || sellSignal {
			idrBalance, assetBalance, err := getBalances(config)
			if err != nil {
				log.Printf("Gagal mengambil saldo | Order stopped! %v", err)
			} else {
				cooldownActive := state.LastTradeAt > 0 && nowMillis-state.LastTradeAt < config.CooldownMs
				limitReached := state.DailyTrades >= config.MaxDailyTrades

				switch {
				case cooldownActive:
					log.Println("Cooldown aktif | Order skipped!")
				case limitReached:
					log.Println("Batas transaksi harian tercapai | Order skipped!")
				case buySignal && idrBalance < config.ModalPerTrade:
					log.Println("Saldo IDR tidak cukup | Order skipped!")
				case buySignal:
					qty := config.ModalPerTrade / price
					qty, err = parseFormattedQty(qty)
					if err != nil {
						log.Printf("Kuantitas pembelian tidak valid: %v", err)
					} else if qty < config.MinQty {
						log.Println("Kuantitas di bawah minimum!")
					} else {
						result, orderErr := kirimOrder(config, "BUY", price, qty)
						if orderErr != nil {
							log.Printf("Order BUY gagal | Bot stopped! %v", orderErr)
							return
						}
						state.LastTradeAt = nowMillis
						state.DailyTrades++
						if result.ExecutedQty > 0 {
							state.Status = "PEGANG"
						}
						if err := saveState(stateFile, state); err != nil {
							log.Printf("Order BUY diterima tapi state tidak tersimpan | Bot stopped! %v", err)
							return
						}
						if result.ExecutedQty <= 0 {
							log.Printf("Order %d diterima tanpa kuantitas | Bot stopped!", result.OrderID)
							return
						}
						log.Printf("Order BUY %d diterima | Kuantitas terisi %.8f", result.OrderID, result.ExecutedQty)
					}
				case sellSignal && assetBalance <= 0:
					log.Println("Saldo aset tidak cukup | Order SELL skipped!")
				case sellSignal:
					result, orderErr := kirimOrder(config, "SELL", price, assetBalance)
					if orderErr != nil {
						log.Printf("Order SELL gagal | Bot stopped! %v", orderErr)
						return
					}
					state.LastTradeAt = nowMillis
					state.DailyTrades++
					if result.ExecutedQty >= assetBalance {
						state.Status = "KOSONG"
					}
					if err := saveState(stateFile, state); err != nil {
						log.Printf("Order SELL diterima tapi state tidak tersimpan | Bot stopped! %v", err)
						return
					}
					if result.ExecutedQty <= 0 {
						log.Printf("Order SELL %d diterima tanpa kuantitas | Bot stopped", result.OrderID)
						return
					}
					log.Printf("Order SELL %d diterima | Kuantitas terisi %.8f", result.OrderID, result.ExecutedQty)
				}
			}
		} else {
			log.Println("⏳ Sabar bos, kita ngopi dulu...")
		}
		countdown(1 * time.Hour)
	}
}

func parseFormattedQty(qty float64) (float64, error) {
	parsed, err := strconv.ParseFloat(formatQty(qty), 64)
	if err != nil || parsed <= 0 {
		return 0, err
	}
	return parsed, nil
}
