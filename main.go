package main

import (
	"log"
	"strings"
	"time"
)

const stateFile = "state.json"

func main() {
	clearTerminal()
	setupGracefulShutdown()
	printBanner()
	
	loc, errLoc := time.LoadLocation("Asia/Jakarta")
	if errLoc == nil {
		time.Local = loc
	}
	typeWriter("[✓] BotIndodax starting...", 40*time.Millisecond)

	for {
		config, err := loadConfig("config.json")
		if err != nil {
			log.Printf("[×] Konfigurasi tidak valid (tunggu 1 menit): %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		state, err := loadState(stateFile)
		if err != nil {
			log.Printf("[×] State tidak dapat dibaca | Bot stopped! %v", err)
			return
		}

		resetDailyCounterIfNeeded(&state, stateFile, config.MaxDailyTrades)
		log.Printf("[$] Config => Buy: Rp %s | Sell: Rp %s", formatRupiah(config.BuyTrigger), formatRupiah(config.SellTrigger))

		price, err := getTickerPrice(config.Pair)
		if err != nil {
			log.Printf("[×] Gagal mengambil harga: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		log.Printf("[•] Pair: %s", strings.ToUpper(config.Pair))
		log.Printf("[$] Harga: Rp %s", formatRupiah(price))
		log.Printf("[•] Status: %s", state.Status)

		now := time.Now()
		nowMillis := now.UnixMilli()

		buySignal := state.Status == "KOSONG" && price <= config.BuyTrigger
		sellSignal := state.Status == "PEGANG" && price >= config.SellTrigger

		if buySignal || sellSignal {
			idrBalance, assetBalance, err := getBalances(config)
			if err != nil {
				log.Printf("[×] Gagal mengambil saldo | Order stopped! %v", err)
			} else {
				cooldownActive := state.LastTradeAt > 0 && nowMillis-state.LastTradeAt < config.CooldownMs
				limitReached := state.DailyTrades >= config.MaxDailyTrades

				switch {
				case cooldownActive:
					log.Println("[×] Cooldown aktif | Order skipped!")
				case limitReached:
					log.Println("[×] Batas harian tercapai | Order skipped!")
				case buySignal && idrBalance < config.ModalPerTrade:
					log.Println("[×] Saldo IDR tidak cukup | Order skipped!")
				case buySignal:
					qty := config.ModalPerTrade / price
					qty, err = parseFormattedQty(qty)
					if err != nil {
						log.Printf("[×] Kuantitas pembelian tidak valid: %v", err)
					} else if qty < config.MinQty {
						log.Println("[×] Kuantitas di bawah minimum!")
					} else {
						result, orderErr := kirimOrder(config, "BUY", price, qty)
						if orderErr != nil {
							log.Printf("[×] Order BUY gagal | Bot stopped! %v", orderErr)
							return
						}
						state.LastTradeAt = nowMillis
						state.DailyTrades++
						if result.ExecutedQty > 0 {
							state.Status = "PEGANG"
						}
						if err := saveState(stateFile, state); err != nil {
							log.Printf("[×] Order BUY diterima tapi state tidak tersimpan | Bot stopped! %v", err)
							return
						}
						if result.ExecutedQty <= 0 {
							log.Printf("[×] Order %d diterima tanpa kuantitas | Bot stopped!", result.OrderID)
							return
						}
						log.Printf("[✓] Order BUY %d diterima | Kuantitas: %.8f", result.OrderID, result.ExecutedQty)
					}
				case sellSignal && assetBalance <= 0:
					log.Println("[×] Saldo tidak cukup | Order SELL skipped!")
				case sellSignal:
					result, orderErr := kirimOrder(config, "SELL", price, assetBalance)
					if orderErr != nil {
						log.Printf("[×] Order SELL gagal | Bot stopped! %v", orderErr)
						return
					}
					state.LastTradeAt = nowMillis
					state.DailyTrades++
					if result.ExecutedQty >= assetBalance {
						state.Status = "KOSONG"
					}
					if err := saveState(stateFile, state); err != nil {
						log.Printf("[×] Order SELL diterima tapi state tidak tersimpan | Bot stopped! %v", err)
						return
					}
					if result.ExecutedQty <= 0 {
						log.Printf("[×] Order SELL %d diterima kuantitas: 0 | Bot stopped", result.OrderID)
						return
					}
					log.Printf("[✓] Order SELL %d diterima | Kuantitas: %.8f", result.OrderID, result.ExecutedQty)
				}
			}
		} else {
			typeWriter("[÷] Sabar bos, kita ngopi dulu...", 40*time.Millisecond)
		}
		countdown(30 * time.Minute)
	}
}
