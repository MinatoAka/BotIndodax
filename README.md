<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.27.3%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.27.3+"></a>
  <a href="https://termux.dev/"><img src="https://img.shields.io/badge/Platform-Termux-2E8B57?style=for-the-badge&logo=android&logoColor=white" alt="Platform: Termux"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License: MIT"></a>
  <a href="https://saweria.co/AKaAturDompet"><img src="https://img.shields.io/badge/Support-Saweria-FF5722?style=for-the-badge&logo=buy-me-a-coffee&logoColor=white" alt="Support via Saweria"></a>
</p>

# BotIndodax

Bot trading otomatis Indodax berbasis Go yang dirancang untuk dijalankan di Termux. Bot memantau satu pair IDR dan mengirim order **market** berdasarkan trigger harga yang ditentukan di konfigurasi.

> [!CAUTION]
> Bot dapat mengirim transaksi nyata menggunakan saldo akun Anda. Trading memiliki risiko kerugian dan tidak ada jaminan keuntungan. Tinjau kode dan konfigurasi sebelum digunakan. Buat API key dengan izin baca dan spot trading saja—**jangan aktifkan izin penarikan**. Gunakan dengan risiko Anda sendiri.

## Fitur

- Memantau harga ticker publik dan membaca saldo akun melalui Indodax API v2.
- Mengirim order market BUY ketika status tersimpan `KOSONG` dan harga ticker mencapai atau berada di bawah `buyTrigger`.
- Mengirim order market SELL ketika status tersimpan `PEGANG` dan harga ticker mencapai atau berada di atas `sellTrigger`; jumlah jual memakai seluruh saldo aset tersedia.
- Menggunakan `modalPerTrade` untuk menghitung jumlah pembelian dan memeriksa saldo IDR serta `minQty`.
- Menerapkan jeda transaksi (`cooldownMs`) dan menyimpan status serta waktu transaksi terakhir di `state.json`.
- Membaca ulang `config.json` setiap iterasi. Pemeriksaan market dilakukan setiap 15 menit.
- Berhenti setelah hasil order tidak dapat dipastikan agar order yang mungkin sudah diterima dapat diperiksa terlebih dahulu.

> **Batasan:** Bot mengelola satu pair dan satu status posisi lokal (`KOSONG` atau `PEGANG`). Tidak ada pengaturan batas transaksi harian. Status lokal tidak menggantikan pemeriksaan order dan saldo di akun Indodax.

## Persyaratan

- Android dengan Termux dan koneksi internet.
- Go versi **1.27.3 atau lebih baru** (sesuai `go.mod`).
- Akun Indodax dengan API key TAPI v2 berizin baca dan spot trading. Jangan aktifkan `WITHDRAWAL`.

## Instalasi

```sh
pkg update
pkg upgrade
pkg install git golang nano

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
```

## Konfigurasi

Edit `config.json` dan ganti nilai API key serta secret dengan kredensial Anda. File konfigurasi contoh berada di repository; jangan masukkan kredensial asli ke commit atau membagikannya.

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
| `secretKey` | Secret key pasangan API key. Wajib diisi dan harus dirahasiakan. |
| `pair` | Pair IDR huruf kecil, misalnya `ethidr`. |
| `buyTrigger` | Harga maksimum pemicu BUY. |
| `sellTrigger` | Harga minimum pemicu SELL; harus lebih besar dari `buyTrigger`. |
| `modalPerTrade` | Nominal IDR yang digunakan untuk menghitung nilai pembelian. Saldo IDR harus mencukupi. |
| `minQty` | Kuantitas aset minimum untuk order BUY. Sesuaikan dengan aturan pair. |
| `cooldownMs` | Jeda minimum antara transaksi sukses, dalam milidetik. Nilai `0` menonaktifkan jeda. |

Konfigurasi divalidasi saat bot berjalan. Jika konfigurasi tidak valid, bot mencatat kesalahan dan mencoba kembali setelah satu menit. Jika kredensial pernah terungkap, cabut API key tersebut dan buat yang baru.

## Build dan jalankan

Jalankan bot dari direktori repository agar `config.json` dan `state.json` ditemukan.

```sh
go build -o bot .
./bot
```

Atau jalankan langsung:

```sh
go run .
```

Di Termux, `termux-wake-lock` dapat membantu mencegah perangkat tidur selama bot berjalan. Hentikan dengan `Ctrl+C`.

Setelah order mengalami timeout atau respons tidak lengkap, jangan langsung menjalankan ulang bot. Periksa riwayat order dan saldo Indodax terlebih dahulu karena order mungkin sudah diterima.

## State dan pengujian

Bot menyimpan `status` dan `lastTradeAt` di `state.json`. Jika file belum ada, status awal adalah `KOSONG`. Jangan menghapus atau mengubah state tanpa memastikan posisi dan saldo akun agar keputusan bot tidak berbeda dari kondisi akun.

Jalankan pengujian dengan:

```sh
go test ./...
```

Tes yang tersedia bersifat lokal dan tidak mengirim order ke Indodax. **Catatan:** test pada repository saat ini masih merujuk field `MaxDailyTrades`, sementara field tersebut tidak ada pada konfigurasi aplikasi; karena itu pengujian perlu diselaraskan sebelum dapat dianggap lulus.

## Struktur repository

| File | Kegunaan |
| --- | --- |
| `main.go` | Loop utama, keputusan transaksi, dan interval pemeriksaan. |
| `config.go` | Membaca dan memvalidasi konfigurasi. |
| `market.go` | Mengambil ticker dan saldo akun. |
| `order.go` | Membuat dan mengirim order market. |
| `utils.go` | State, signature API, format kuantitas, dan tampilan terminal. |
| `config.json` | Contoh konfigurasi lokal. |
| `utils_test.go` | Pengujian unit. |
| `LICENSE` | Lisensi MIT. |

## Dukungan

[♡ Dukung pengembangan via Saweria](https://saweria.co/AKaAturDompet)

## Disclaimer

Perangkat lunak ini disediakan apa adanya, tanpa jaminan, dan bukan nasihat investasi. Anda bertanggung jawab atas keamanan API key, konfigurasi, pemantauan bot, serta seluruh keputusan dan risiko trading.
