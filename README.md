# BotIndodax

[![Go](https://img.shields.io/badge/Go-100%25-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Termux](https://img.shields.io/badge/Platform-Termux-7F5AF0)](https://termux.dev/)
[![Status](https://img.shields.io/badge/Status-Experimental-orange)](https://github.com/MinatoAka/BotIndodax)

Bot trading Indodax berbasis Go yang dirancang untuk dijalankan dari terminal Android melalui Termux. Bot memantau satu pair tertentu, membaca saldo, dan mengeksekusi order market saat kondisi pasar memenuhi trigger yang telah ditentukan.

> [!CAUTION]
> Bot dapat melakukan transaksi dengan saldo asli. Order market dapat langsung dieksekusi dan menimbulkan kerugian; tidak ada jaminan keuntungan. Tinjau kode, konfigurasi, dan API key dengan hati-hati. Gunakan akun dan modal yang siap ditanggung risikonya. Jangan pernah memberikan API key atau secret key kepada pihak yang tidak dipercaya.

## Fitur utama

- Memantau harga ticker untuk pair yang diatur di `config.json`.
- Membeli saat status `KOSONG` dan harga berada pada atau di bawah `buyTrigger`.
- Menjual saat status `PEGANG` dan harga berada pada atau di atas `sellTrigger`.
- Menyimpan status lokal dan batas transaksi harian di `state.json`.
- Beroperasi dalam loop dengan jeda satu jam antar pengecekan.

## Persyaratan

- Perangkat Android dengan Termux.
- Koneksi internet aktif.
- Akun Indodax dengan API key yang memiliki izin baca saldo dan trading.
- Paket `git` dan `golang` di Termux.

> Penting: Jangan beri izin penarikan pada API key.

## Instalasi di Termux

Buka Termux lalu jalankan:

```sh
pkg update
pkg upgrade
pkg install git golang nano

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
```

Pastikan Go siap dipakai:

```sh
go version
```

## Menyiapkan konfigurasi

Program membaca konfigurasi dari `config.json` di direktori kerja. Buat salinan contoh lalu edit file yang dibaca bot:

```sh
cp config.json config.json.example
nano config.json
```

Contoh struktur file konfigurasi:

```json
{
  "apiKey": "API_KEY_ANDA",
  "secretKey": "SECRET_KEY_ANDA",
  "pair": "ethidr",
  "buyTrigger": 41333000,
  "sellTrigger": 42572990,
  "modalPerTrade": 15000,
  "minQty": 0.0001,
  "cooldownMs": 1800000,
  "maxDailyTrades": 6
}
```

| Field | Keterangan |
| --- | --- |
| `apiKey` | API key Indodax dengan izin yang diperlukan. |
| `secretKey` | Secret key pasangan API key. Simpan dengan aman. |
| `pair` | Pair pasar, misalnya `ethidr`. |
| `buyTrigger` | Harga pembelian maksimum yang dipertimbangkan saat `KOSONG`. |
| `sellTrigger` | Harga penjualan minimum yang dipertimbangkan saat `PEGANG`. |
| `modalPerTrade` | Nominal IDR per pembelian yang direncanakan. |
| `minQty` | Kuantitas minimum untuk transaksi. |
| `cooldownMs` | Jeda minimum antar transaksi dalam milidetik. |
| `maxDailyTrades` | Batas jumlah transaksi harian. |

Pastikan pair, trigger, modal, dan batas transaksi sudah benar sebelum aktif. File konfigurasi berisi kredensial dalam bentuk teks biasa, jadi jangan mengunggah atau membagikan file tersebut ke repositori publik.

## Kompilasi di Termux

Setelah file konfigurasi siap, kompilasi program menjadi executable bernama `bot`:

```sh
go build -o bot *.go
```

Validasi hasilnya:

```sh
ls -l bot
```

Jika build berhasil, program siap dijalankan di lingkungan Termux saat ini. Kompilasi ulang diperlukan jika ada perubahan pada kode sumber.

## Menjalankan bot

Pastikan berada di direktori repo yang benar, lalu jalankan:

```sh
./bot
```

Atau tanpa membuat file executable dahulu:

```sh
go run *.go
```

Bot akan membaca `config.json` dan `state.json`, mengambil harga pasar, lalu menunggu satu jam sebelum iterasi berikutnya. Hentikan proses dengan `Ctrl + C` di terminal.

Jika ingin meminimalkan kemungkinan proses terhenti ketika Termux berjalan di latar belakang, Anda dapat menjalankan:

```sh
termux-wake-lock
```

Setelah selesai, lepaskan dengan:

```sh
termux-wake-unlock
```

Catatan: pengaturan baterai Android bisa tetap membatasi proses latar belakang; `wake-lock` hanya membantu, bukan jaminan proses selalu berjalan terus-menerus.

## Status dan perilaku transaksi

`state.json` dibuat dan dipelihara di direktori kerja saat bot berjalan. Jika file belum ada atau tidak dapat dibaca, bot akan memulai dengan status `KOSONG`.

Status lokal ini digunakan untuk mengontrol:

- apakah bot siap membeli atau menjual,
- cooldown antar transaksi,
- counter transaksi harian,
- tanggal terakhir transaksi.

Jangan menghapus atau mengubah `state.json` sembarangan. Jika status lokal tidak sesuai dengan saldo atau riwayat akun yang sebenarnya, keputusan bot bisa berbeda dari yang diharapkan.

## Struktur file

```text
BotIndodax/
├── README.md
├── config.go
├── config.json
├── main.go
├── market.go
├── order.go
├── utils.go
├── state.json        # dibuat saat runtime
└── bot               # hasil build Go (jika sudah dikompilasi)
```

Penjelasan tiap file:

- `main.go` — loop utama program. Membaca konfigurasi, mengambil harga, mengecek sinyal, dan menjalankan logika trading.
- `config.go` — membaca dan mengurai file `config.json` ke struktur `Config`.
- `market.go` — berisi fungsi untuk mengambil harga ticker dan saldo akun dari API Indodax.
- `order.go` — menyiapkan dan mengirim order BUY/SELL ke API Indodax dengan tanda tangan HMAC.
- `utils.go` — fungsi bantu seperti state management, pembuatan query string, normalisasi teks, dan format kuantitas.
- `config.json` — file konfigurasi utama berisi API key, pair, trigger, serta parameter trading.
- `state.json` — status lokal bot seperti apakah posisi sedang kosong atau menahan aset, cooldown, dan jumlah transaksi harian.
- `README.md` — dokumentasi project.

## Troubleshooting

### `undefined: strings` saat kompilasi

Pastikan blok import di `order.go` sudah memasukkan `"strings"`:

```go
import (
    "bytes"
    "fmt"
    "io/ioutil"
    "math"
    "net/http"
    "strconv"
    "strings"
    "time"
)
```

Lalu jalankan kembali:

```sh
go build -o bot *.go
```

### Gagal membaca `config.json`

Pastikan perintah dijalankan dari direktori `BotIndodax`, nama file tepat `config.json`, dan JSON valid. Periksa tanda kutip, koma, serta tipe data.

### Gagal mengambil harga/saldo atau mengirim order

Periksa:

- koneksi internet,
- pair yang digunakan,
- waktu perangkat,
- saldo akun,
- izin API key,
- respons error dari Indodax.

Jangan membagikan API key atau secret key saat meminta bantuan.

## Disclaimer

Proyek ini disediakan apa adanya, tanpa jaminan apa pun, dan bukan nasihat investasi. Anda bertanggung jawab atas keamanan API key, konfigurasi, pemantauan proses, serta seluruh keputusan trading yang dilakukan.

Semua operasi ini dilakukan sesuai dengan keputusan pengguna dan risiko yang ditanggung sendiri.
