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
	typeWriter("[•] ========================================", 40*time.Millisecond)
	typeWriter("[✓] BotIndodax starting | Bismillahirrahmanirrahim...", 40*time.Millisecond)

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

		log.Printf("[$] Config => Buy: Rp %s | Sell: Rp %s", formatRupiah(config.BuyTrigger), formatRupiah(config.SellTrigger))

		price, err := getTickerPrice(config.Pair)
		if err != nil {
			log.Printf("[×] Gagal mengambil harga: %v", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		log.Printf("[•] Pair: %s", strings.ToUpper(config.Pair))
		log.Printf("[$] Price: Rp %s", formatRupiah(price))
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

				switch {
				case cooldownActive:
					log.Println("[×] Cooldown aktif | Order skipped!")
				
				// EKSEKUSI BUY
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
						
						if result.OrderID > 0 {
							log.Printf("[✓] Order BUY sukses! ID: %d", result.OrderID)
							
							state.Status = "PEGANG"
							state.LastTradeAt = nowMillis
							
							if err := saveState(stateFile, state); err != nil {
								log.Printf("[×] State gagal disimpan: %v", err)
								return
							}
							
							_, assetAfter, _ := getBalances(config)
							log.Printf("[i] Cek saldo aset: %.8f", assetAfter)
							
						} else {
							log.Println("[×] Order BUY gagal!")
						}
					}

				// EKSEKUSI SELL & SATPAM DESYNC
				case sellSignal && assetBalance <= 0:
					log.Println(" [×] Saldo aset 0 | Reset status: KOSONG")
					state.Status = "KOSONG"
					state.LastTradeAt = 0
					
					if err := saveState(stateFile, state); err != nil {
						log.Printf("[×] Gagal menyimpan reset status: %v", err)
					}
				
				case sellSignal:
					result, orderErr := kirimOrder(config, "SELL", price, assetBalance)
					
					if orderErr != nil {
						log.Printf("[×] Order SELL gagal | Bot stopped! %v", orderErr)
						return
					}
					
					if result.OrderID > 0 {
						log.Printf("[✓] Order SELL sukses! ID: %d", result.OrderID)
						
						state.Status = "KOSONG"
						state.LastTradeAt = nowMillis
						
						if err := saveState(stateFile, state); err != nil {
							log.Printf("[×] State gagal disimpan: %v", err)
							return
						}
						
						idrAfter, _, _ := getBalances(config)
						log.Printf("[i] Cek Saldo: IDR saat ini = Rp %s", formatRupiah(idrAfter))
						
					} else {
						log.Println("[×] Order SELL gagal!")
					}
				}
			}
		} else {
			typeWriter("[÷] Sabar bos, kita ngopi dulu...", 40*time.Millisecond)
		}
		
		countdown(15 * time.Minute)
	}
}
