<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.27.3%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.27.3+"></a>
  <a href="https://termux.dev/"><img src="https://img.shields.io/badge/Platform-Termux-2E8B57?style=for-the-badge&logo=android&logoColor=white" alt="Platform: Termux"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="MIT License"></a>
  <a href="https://saweria.co/AKaAturDompet"><img src="https://img.shields.io/badge/Support-Saweria-FF5722?style=for-the-badge&logo=buy-me-a-coffee&logoColor=white" alt="Support via Saweria"></a>
</p>

# BotIndodax

Bot trading otomatis Indodax berbasis Go yang dirancang untuk dijalankan di Termux. Bot memantau satu pair IDR, membaca harga pasar, dan mengirim order market berdasarkan trigger harga yang dikonfigurasi.

> [!CAUTION]
> Bot dapat mengirim transaksi nyata dengan saldo akun Anda. Trading memiliki risiko kerugian dan tidak ada jaminan keuntungan. Tinjau kode dan konfigurasi sebelum digunakan. Gunakan API key dengan izin baca dan spot trading saja—jangan aktifkan izin penarikan. Gunakan dengan risiko Anda sendiri.

## Ringkasan

BotIndodax adalah bot trading berbasis terminal yang berjalan di Android dengan Termux. Bot membaca harga ticker, mengambil saldo akun dari API Indodax, lalu menentukan aksi beli/jual berdasarkan state lokal (`KOSONG` atau `PEGANG`).

Bot ini sengaja dibuat sederhana dan fokus pada satu pair. Data state disimpan di `state.json` agar bot dapat melanjutkan siklus trading secara konsisten.

## Fitur

- Memantau harga ticker dari API Indodax.
- Membaca saldo IDR dan saldo aset dari akun Indodax.
- Mengecek state lokal: `KOSONG` atau `PEGANG`.
- Mengirim order BUY saat status `KOSONG` dan harga `<= buyTrigger`.
- Mengirim order SELL saat status `PEGANG` dan harga `>= sellTrigger`.
- Menghitung jumlah pembelian menggunakan `modalPerTrade`.
- Memeriksa `minQty` sebelum order BUY.
- Menyimpan `lastTradeAt` dan `status` di `state.json`.
- Menerapkan `cooldownMs` agar transaksi tidak berulang terlalu cepat.
- Membaca ulang `config.json` pada setiap iterasi.
- Menjaga keamanan dengan menghentikan bot bila hasil order tidak dapat dipastikan.

## Batasan

- Bot hanya memantau satu pair.
- State lokal adalah keputusan lokal, bukan pengganti pemeriksaan saldo dan order di akun Indodax.
- Saat ticker gagal diambil, bot menunggu beberapa menit sebelum mencoba lagi.

## Persyaratan

- Android dengan Termux dan koneksi internet.
- Go versi 1.27.3 atau lebih baru.
- Akun Indodax dengan API key TAPI v2.
- API key harus memiliki izin baca dan spot trading. Jangan aktifkan izin `WITHDRAWAL`.

## Instalasi

```sh
pkg update
pkg upgrade
pkg install git golang nano

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
```

## Konfigurasi

Edit file `config.json` di root repository, lalu isi API key dan secret Anda.

```sh
nano config.json
```

Contoh konfigurasi:

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
| `pair` | Pair IDR dalam huruf kecil, contohnya `ethidr`. |
| `buyTrigger` | Harga maksimum yang memicu BUY. |
| `sellTrigger` | Harga minimum yang memicu SELL; harus lebih besar dari `buyTrigger`. |
| `modalPerTrade` | Nominal IDR yang dipakai untuk menghitung kuantitas BUY. |
| `minQty` | Kuantitas minimum aset yang diperbolehkan untuk order BUY. |
| `cooldownMs` | Jeda minimum antar transaksi, dalam milidetik. |

Catatan penting:
- Jangan commit atau bagikan kredensial asli ke repository publik.
- Jika API key pernah bocor, cabut dan buat yang baru.
- Konfigurasi akan divalidasi saat bot dijalankan.

## Build dan jalankan

Jalankan dari directory repository agar `config.json` dan `state.json` dapat ditemukan.

```sh
go build -o bot .
./bot
```

Atau jalankan langsung tanpa build binary:

```sh
go run .
```

Di Termux, Anda juga dapat menggunakan perintah berikut agar perangkat tidak masuk sleep saat bot berjalan:

```sh
termux-wake-lock
```

Hentikan bot dengan:

```sh
Ctrl + C
```

## State dan pengujian

Bot menyimpan status posisi dan waktu transaksi terakhir di `state.json`.

- Jika file tidak ada, bot akan mulai dengan status `KOSONG`.
- Jangan mengubah state secara sembarangan tanpa memastikan kondisi akun Anda masih sesuai.

Jalankan pengujian lokal:

```sh
go test ./...
```

Catatan: pada saat ini, tes di repo masih perlu diselaraskan dengan konfigurasi yang benar-benar dipakai aplikasi. Namun fungsi utama bot tetap sesuai dengan implementasi yang ada.

## Struktur repository

```text
BotIndodax/
├── LICENSE
├── README.md
├── config.go
├── config.json
├── go.mod
├── main.go
├── market.go
├── order.go
├── utils.go
├── utils_test.go
└── state.json   # dibuat saat runtime
```

## Dukungan

[♡ Dukung pengembangan via Saweria](https://saweria.co/AKaAturDompet)

## Lisensi

Proyek ini dilisensikan di bawah MIT License. Lihat [LICENSE](LICENSE) untuk detailnya.

## Disclaimer

Proyek ini disediakan apa adanya, tanpa jaminan, dan bukan nasihat investasi. Anda bertanggung jawab atas keamanan API key, konfigurasi, pemantauan bot, serta seluruh keputusan dan risiko trading yang Anda ambil.
