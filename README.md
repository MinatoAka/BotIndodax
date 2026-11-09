<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.27.3%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.27.3+"></a>
  <a href="https://termux.dev/"><img src="https://img.shields.io/badge/Platform-Termux-2E8B57?style=for-the-badge&logo=android&logoColor=white" alt="Platform: Termux"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License: MIT"></a>
  <a href="https://saweria.co/AKaAturDompet"><img src="https://img.shields.io/badge/Support-Saweria-FF5722?style=for-the-badge&logo=buy-me-a-coffee&logoColor=white" alt="Support via Saweria"></a>
</p>

<p align="center"><a href="#bahasa-indonesia">Bahasa Indonesia</a> · <a href="#english">English</a></p>

# BotIndodax

## Bahasa Indonesia

Bot trading otomatis Indodax berbasis Go yang dirancang untuk Termux. Bot memantau satu pair IDR dan mengirim order **market** berdasarkan trigger harga yang ditentukan dalam konfigurasi.

> [!CAUTION]
> Bot dapat mengirim transaksi nyata menggunakan saldo akun Anda. Trading berisiko menyebabkan kerugian dan tidak ada jaminan keuntungan. Tinjau kode dan konfigurasi sebelum digunakan. Gunakan API key dengan izin baca dan spot trading saja—**jangan aktifkan izin penarikan**. Gunakan dengan risiko Anda sendiri.

### Fitur

- Memantau ticker publik dan membaca saldo akun melalui Indodax API v2.
- Mengirim order market BUY saat status lokal `KOSONG` dan harga mencapai atau berada di bawah `buyTrigger`.
- Mengirim order market SELL saat status lokal `PEGANG` dan harga mencapai atau berada di atas `sellTrigger`; order menggunakan seluruh saldo aset tersedia.
- Menghitung jumlah BUY dari `modalPerTrade`, memeriksa saldo IDR, dan membandingkan kuantitas dengan `minQty`.
- Menyimpan status posisi (`KOSONG` atau `PEGANG`) dan waktu transaksi terakhir di `state.json`.
- Menerapkan jeda transaksi berdasarkan `cooldownMs` dan membaca ulang `config.json` setiap iterasi.

**Batasan:** Bot hanya memantau satu pair dan menggunakan status posisi lokal. Tidak ada batas transaksi harian. Status lokal tidak menggantikan pemeriksaan order dan saldo di akun Indodax. Pemeriksaan market normal dilakukan setiap 15 menit; jika ticker gagal diambil, bot menunggu 5 menit sebelum mencoba kembali.

### Persyaratan

- Android dengan Termux dan koneksi internet.
- Go versi **1.27.3 atau lebih baru**, sesuai `go.mod`.
- Akun Indodax dengan API key TAPI v2 berizin baca dan spot trading. Jangan aktifkan izin `WITHDRAWAL`.

### Instalasi

```sh
pkg update
pkg upgrade
pkg install git golang nano

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
```

### Konfigurasi

Edit `config.json`, lalu isi API key dan secret Anda:

```sh
nano config.json
```

Contoh format:

```json
{
  "apiKey": "ISI_API_KEY_LO",
  "secretKey": "ISI_SECRET_KEY_LO",
  "pair": "ethidr",
  "buyTrigger": 33450000,
  "sellTrigger": 34453500,
  "modalPerTrade": 15000,
  "minQty": 0.0001,
  "cooldownMs": 1800000
}
```

| Field | Keterangan |
| --- | --- |
| `apiKey` | API key Indodax TAPI v2. Wajib diisi. |
| `secretKey` | Secret key pasangan API key. Wajib diisi dan dirahasiakan. |
| `pair` | Pair IDR dalam huruf kecil, misalnya `ethidr`. |
| `buyTrigger` | Harga maksimum pemicu BUY. |
| `sellTrigger` | Harga minimum pemicu SELL; harus lebih besar dari `buyTrigger`. |
| `modalPerTrade` | Nominal IDR yang digunakan untuk menghitung jumlah pembelian. Saldo harus mencukupi. |
| `minQty` | Kuantitas aset minimum untuk order BUY; sesuaikan dengan aturan pair. |
| `cooldownMs` | Jeda minimum setelah transaksi sukses, dalam milidetik. Nilai `0` menonaktifkan jeda. |

Jangan unggah atau commit kredensial asli. Jika API key atau secret terungkap, cabut key tersebut dan buat yang baru.

### Build dan jalankan

Jalankan bot dari direktori repository agar `config.json` dan `state.json` ditemukan.

```sh
go build -o bot .
./bot
```

Atau jalankan langsung:

```sh
go run .
```

Di Termux, `termux-wake-lock` dapat membantu mencegah perangkat tidur selama bot berjalan. Hentikan bot dengan `Ctrl+C`.

Jika order mengalami timeout atau respons tidak lengkap, jangan langsung menjalankan ulang bot. Periksa riwayat order dan saldo Indodax terlebih dahulu karena order mungkin sudah diterima.

### State dan pengujian

Jika `state.json` belum ada, bot memulai dengan status `KOSONG`. Jangan menghapus atau mengubah state tanpa memastikan posisi dan saldo akun agar keputusan bot tidak berbeda dari kondisi akun.

Jalankan pengujian dengan:

```sh
go test ./...
```

**Catatan:** Test yang tersedia saat ini merujuk field `MaxDailyTrades`, tetapi field tersebut tidak ada pada konfigurasi aplikasi. Karena itu, test suite perlu diselaraskan sebelum dapat dikompilasi dan dijalankan.

### Struktur repository

| File | Kegunaan |
| --- | --- |
| `main.go` | Loop utama, keputusan transaksi, dan interval pemeriksaan. |
| `config.go` | Membaca dan memvalidasi konfigurasi. |
| `market.go` | Mengambil ticker dan saldo akun. |
| `order.go` | Membuat dan mengirim order market. |
| `utils.go` | State, signature API, format kuantitas, dan tampilan terminal. |
| `config.json` | Contoh konfigurasi. |
| `utils_test.go` | Pengujian unit. |
| `LICENSE` | Lisensi MIT. |

### Dukungan

[♡ Dukung pengembangan via Saweria](https://saweria.co/AKaAturDompet)

### Disclaimer

Perangkat lunak ini disediakan apa adanya, tanpa jaminan, dan bukan nasihat investasi. Anda bertanggung jawab atas keamanan API key, konfigurasi, pemantauan bot, serta seluruh keputusan dan risiko trading.

---

## English

An automated Indodax trading bot written in Go and designed for Termux. It monitors one IDR trading pair and submits **market** orders based on the configured price triggers.

> [!CAUTION]
> The bot can place real trades using your account balance. Trading can result in losses; profits are not guaranteed. Review the code and configuration before use. Use an API key with read and spot-trading permissions only—**do not enable withdrawals**. Use at your own risk.

### Features

- Monitors the public ticker and retrieves account balances through the Indodax API v2.
- Submits a market BUY when the local state is `KOSONG` and the price is at or below `buyTrigger`.
- Submits a market SELL when the local state is `PEGANG` and the price is at or above `sellTrigger`; the order uses the available asset balance.
- Calculates the BUY quantity from `modalPerTrade`, checks the IDR balance, and compares the quantity against `minQty`.
- Stores the position state (`KOSONG` or `PEGANG`) and the timestamp of the last trade in `state.json`.
- Applies a trade cooldown using `cooldownMs` and reloads `config.json` on each iteration.

**Limitations:** The bot monitors one pair and relies on a local position state. It does not impose a daily trade limit. Local state does not replace checking orders and balances in your Indodax account. Under normal operation, market checks run every 15 minutes; if the ticker request fails, the bot waits 5 minutes before retrying.

### Requirements

- Android with Termux and an internet connection.
- Go **1.27.3 or later**, as specified in `go.mod`.
- An Indodax account and a TAPI v2 API key with read and spot-trading permissions. Do not enable `WITHDRAWAL` permission.

### Installation

```sh
pkg update
pkg upgrade
pkg install git golang nano

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
```

### Configuration

Edit `config.json` and enter your API key and secret:

```sh
nano config.json
```

Example:

```json
{
  "apiKey": "YOUR_API_KEY",
  "secretKey": "YOUR_SECRET_KEY",
  "pair": "ethidr",
  "buyTrigger": 33450000,
  "sellTrigger": 34453500,
  "modalPerTrade": 15000,
  "minQty": 0.0001,
  "cooldownMs": 1800000
}
```

| Field | Description |
| --- | --- |
| `apiKey` | Indodax TAPI v2 API key. Required. |
| `secretKey` | Secret associated with the API key. Required; keep it private. |
| `pair` | IDR pair in lowercase, for example `ethidr`. |
| `buyTrigger` | Maximum price that triggers a BUY. |
| `sellTrigger` | Minimum price that triggers a SELL; must be greater than `buyTrigger`. |
| `modalPerTrade` | IDR amount used to calculate the BUY quantity. The account balance must be sufficient. |
| `minQty` | Minimum asset quantity for a BUY order; follow the pair's trading rules. |
| `cooldownMs` | Minimum delay after a successful trade, in milliseconds. Set to `0` to disable. |

Never commit or share real credentials. If an API key or secret is exposed, revoke it and create a replacement.

### Build and run

Run the bot from the repository directory so it can find `config.json` and `state.json`.

```sh
go build -o bot .
./bot
```

Alternatively, run it without building a binary:

```sh
go run .
```

On Termux, `termux-wake-lock` can help prevent the device from sleeping while the bot is running. Stop the bot with `Ctrl+C`.

If an order request times out or returns an incomplete response, do not immediately restart the bot. Check your Indodax order history and balances first, since the order may have been accepted.

### State and tests

If `state.json` does not exist, the bot starts with the `KOSONG` state. Do not delete or edit the state without first confirming your account position and balance, or the bot's local state may no longer match the account.

Run the tests with:

```sh
go test ./...
```

**Note:** The current test code refers to a `MaxDailyTrades` field that is not present in the application's configuration. The test suite must be brought in sync before it can compile and run.

### Repository structure

| File | Purpose |
| --- | --- |
| `main.go` | Main loop, trading decisions, and check intervals. |
| `config.go` | Loads and validates configuration. |
| `market.go` | Retrieves ticker prices and account balances. |
| `order.go` | Builds and submits market orders. |
| `utils.go` | State management, API signatures, quantity formatting, and terminal output. |
| `config.json` | Example configuration. |
| `utils_test.go` | Unit tests. |
| `LICENSE` | MIT License. |

### Support

[♡ Support development via Saweria](https://saweria.co/AKaAturDompet)

### Disclaimer

This software is provided “as is,” without warranty, and is not investment advice. You are responsible for API key security, configuration, monitoring the bot, and all trading decisions and risks.
