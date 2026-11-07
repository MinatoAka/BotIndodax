package main

import (
	"fmt"
	"strings"
	"strconv"
	"time"
)

func main() {
	fmt.Println("🤖 Bot Aji Mumpung (Modular Version) Menyala!")

	for {
		// 1. Baca Config & State
		config, err := loadConfig("config.json")
		if err != nil {
			fmt.Println("❌ Gagal baca config.json:", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		state, _ := loadState("state.json")

		// 2. Ambil Harga
		hargaSekarang, err := getTickerPrice(config.Pair)
		if err != nil {
			fmt.Println("❌ Gagal ambil harga:", err)
			time.Sleep(5 * time.Minute)
			continue
		}

		fmt.Printf("📝 Log: Pair=%s | Harga=Rp %.0f | Status=%s\n", strings.ToUpper(config.Pair), hargaSekarang, state.Status)

		// 3. Cek Sinyal
		isBuySignal := state.Status == "KOSONG" && hargaSekarang <= config.BuyTrigger
		isSellSignal := state.Status == "PEGANG" && hargaSekarang >= config.SellTrigger

		if !isBuySignal && !isSellSignal {
			fmt.Println("🟡 Sabar bos kita tunggu momen yang pas...")
		} else {
			fmt.Println("🔔 Sinyal terdeteksi! Mengecek saldo...")
			idrSaldo, asetSaldo, err := getBalances(config)
			if err != nil {
				fmt.Println("❌ Gagal ambil saldo:", err)
			} else {
				// Cek Cooldown & Daily Limit
				now := time.Now().UnixMilli()
				today := time.Now().Format("2006-01-02")
				if state.LastTradeDay != today {
					state.DailyTrades = 0
				}

				isCooldown := (now - state.LastTradeAt) < config.CooldownMs
				isLimitReached := state.DailyTrades >= config.MaxDailyTrades

				if isCooldown {
					fmt.Println("🛡️ Cooldown aktif, belum bisa trade.")
				} else if isLimitReached {
					fmt.Println("🛡️ Limit harian tercapai.")
				} else {
					// EKSEKUSI
					if isBuySignal && idrSaldo >= config.ModalPerTrade {
						rawQty := config.ModalPerTrade / hargaSekarang
						qtyStr := formatQty(rawQty) 
						qty, _ := strconv.ParseFloat(qtyStr, 64)
						if qty >= config.MinQty {
							if kirimOrder(config, "BUY", hargaSekarang, qty) {
								state.Status = "PEGANG"
								state.LastTradeAt = now
								state.DailyTrades++
								state.LastTradeDay = today
								saveState("state.json", state)
							}
						} else {
							fmt.Println("⚠️ Qty di bawah minimum Indodax.")
						}
					} else if isSellSignal && asetSaldo > 0 {
						if kirimOrder(config, "SELL", hargaSekarang, asetSaldo) {
							state.Status = "KOSONG"
							state.LastTradeAt = now
							state.DailyTrades++
							state.LastTradeDay = today
							saveState("state.json", state)
						}
					}
				}
			}
		}

		// 4. Jeda
		fmt.Println("⏳ Menunggu 1 jam untuk iterasi berikutnya...\n")
		time.Sleep(1 * time.Hour)
	}
}
