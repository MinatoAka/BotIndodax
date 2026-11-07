package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const stateFile = "state.json"

func main() {
	log.Println("BotIndodax starting")
	config, err := loadConfig("config.json")
	if err != nil {
		log.Printf("Konfigurasi tidak valid: %v", err)
		return
	}

	for {
		state, err := loadState(stateFile)
		if err != nil {
			log.Printf("State tidak dapat dibaca; bot dihentikan demi keamanan: %v", err)
			return
		}

		price, err := getTickerPrice(config.Pair)
		if err != nil {
			log.Printf("Gagal mengambil harga: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		log.Printf("Pair=%s | Harga=Rp %.0f | Status=%s", strings.ToUpper(config.Pair), price, state.Status)

		buySignal := state.Status == "KOSONG" && price <= config.BuyTrigger
		sellSignal := state.Status == "PEGANG" && price >= config.SellTrigger
		if buySignal || sellSignal {
			idrBalance, assetBalance, err := getBalances(config)
			if err != nil {
				log.Printf("Gagal mengambil saldo; tidak ada order yang dikirim: %v", err)
			} else {
				now := time.Now()
				today := now.Format("2006-01-02")
				if state.LastTradeDay != today {
					state.DailyTrades = 0
				}
				nowMillis := now.UnixMilli()
				cooldownActive := state.LastTradeAt > 0 && nowMillis-state.LastTradeAt < config.CooldownMs
				limitReached := state.DailyTrades >= config.MaxDailyTrades

				switch {
				case cooldownActive:
					log.Println("Cooldown aktif; order dilewati")
				case limitReached:
					log.Println("Batas transaksi harian tercapai; order dilewati")
				case buySignal && idrBalance < config.ModalPerTrade:
					log.Println("Saldo IDR tidak cukup; order dilewati")
				case buySignal:
					qty := config.ModalPerTrade / price
					qty, err = parseFormattedQty(qty)
					if err != nil {
						log.Printf("Kuantitas pembelian tidak valid: %v", err)
					} else if qty < config.MinQty {
						log.Println("Kuantitas di bawah minimum yang dikonfigurasi")
					} else {
						result, orderErr := kirimOrder(config, "BUY", price, qty)
						if orderErr != nil {
							log.Printf("Order BUY gagal atau hasilnya tidak pasti; bot dihentikan untuk mencegah duplikasi: %v", orderErr)
							return
						}
						state.LastTradeAt = nowMillis
						state.DailyTrades++
						state.LastTradeDay = today
						if result.ExecutedQty > 0 {
							state.Status = "PEGANG"
						}
						if err := saveState(stateFile, state); err != nil {
							log.Printf("Order diterima tetapi state gagal disimpan; bot dihentikan, periksa akun dan state manual: %v", err)
							return
						}
						if result.ExecutedQty <= 0 {
							log.Printf("Order %d diterima tanpa kuantitas eksekusi; bot dihentikan untuk rekonsiliasi manual", result.OrderID)
							return
						}
						log.Printf("Order BUY %d diterima; kuantitas terisi %.8f", result.OrderID, result.ExecutedQty)
					}
				case sellSignal && assetBalance <= 0:
					log.Println("Saldo aset tidak cukup; order dilewati")
				case sellSignal:
					result, orderErr := kirimOrder(config, "SELL", price, assetBalance)
					if orderErr != nil {
						log.Printf("Order SELL gagal atau hasilnya tidak pasti; bot dihentikan untuk mencegah duplikasi: %v", orderErr)
						return
					}
					state.LastTradeAt = nowMillis
					state.DailyTrades++
					state.LastTradeDay = today
					if result.ExecutedQty >= assetBalance {
						state.Status = "KOSONG"
					}
					if err := saveState(stateFile, state); err != nil {
						log.Printf("Order diterima tetapi state gagal disimpan; bot dihentikan, periksa akun dan state manual: %v", err)
						return
					}
					if result.ExecutedQty <= 0 {
						log.Printf("Order %d diterima tanpa kuantitas eksekusi; bot dihentikan untuk rekonsiliasi manual", result.OrderID)
						return
					}
					log.Printf("Order SELL %d diterima; kuantitas terisi %.8f", result.OrderID, result.ExecutedQty)
				}
			}
		}

		log.Println("Menunggu 1 jam untuk iterasi berikutnya")
		time.Sleep(time.Hour)
	}
}

func parseFormattedQty(qty float64) (float64, error) {
	parsed, err := strconv.ParseFloat(formatQty(qty), 64)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("kuantitas %.8f tidak dapat diformat", qty)
	}
	return parsed, nil
}
