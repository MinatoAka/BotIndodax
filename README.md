# BotIndodax

Bot trading Indodax berbasis Go untuk dijalankan melalui terminal Android (Termux). Bot memantau satu pair IDR dan dapat mengirim order market menggunakan API key akun Anda.

> [!CAUTION]
> Bot dapat melakukan transaksi dengan saldo asli dan menyebabkan kerugian. Tidak ada jaminan keuntungan. Periksa kode, konfigurasi, dan izin API key sebelum menjalankan bot. Jangan aktifkan izin penarikan pada API key. Gunakan dengan risiko Anda sendiri.

## Fitur

- Membaca harga ticker dan saldo akun Indodax.
- Membeli ketika status lokal `KOSONG` dan harga sama dengan atau di bawah `buyTrigger`.
- Menjual ketika status lokal `PEGANG` dan harga sama dengan atau di atas `sellTrigger`.
- Menyimpan status transaksi, waktu transaksi terakhir, dan batas transaksi harian di `state.json`.
- Memvalidasi konfigurasi serta menangani kegagalan API dengan lebih aman.
- Menghentikan bot bila hasil pengiriman order tidak pasti agar tidak otomatis mengirim order duplikat.
- Menunggu satu jam di antara pemeriksaan normal. Jika pengambilan harga gagal, bot mencoba lagi setelah lima menit.

## Persyaratan

- Android dengan Termux dan koneksi internet.
- Go terpasang (`pkg install golang`).
- Akun Indodax dan API key TAPI v2 dengan izin baca dan spot trading. Jangan beri izin penarikan.

## Instalasi

```sh
pkg update
pkg upgrade
pkg install git golang nano

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
go version
```

## Konfigurasi

Bot membaca `config.json` dari direktori kerja. Buka file konfigurasi yang sudah ada dan masukkan API key/secret Anda:

```sh
nano config.json
```

Contoh format:

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
| `apiKey` | API key Indodax TAPI v2 dengan izin baca dan spot trading. |
| `secretKey` | Secret key pasangan API key; rahasiakan. |
| `pair` | Pair IDR, huruf kecil, misalnya `ethidr`. |
| `buyTrigger` | Batas harga maksimum untuk sinyal beli. |
| `sellTrigger` | Batas harga minimum untuk sinyal jual; harus lebih besar dari `buyTrigger`. |
| `modalPerTrade` | Nominal IDR yang direncanakan untuk pembelian. |
| `minQty` | Kuantitas aset minimum sesuai aturan pair. |
| `cooldownMs` | Jeda minimum antar transaksi dalam milidetik. |
| `maxDailyTrades` | Batas transaksi harian yang dikonfigurasi. |

Pastikan pair dan aturan jumlah/nominal minimum sesuai ketentuan pasar Indodax. Jangan unggah atau commit `config.json` berisi kredensial. Jika kredensial pernah terpublikasi, segera cabut dan buat ulang API key.

## Menjalankan pengujian

Tes ini memeriksa validasi konfigurasi, perilaku pembacaan dan penyimpanan state, permission file state, serta konsistensi query string. Tes tidak mengirim order ke Indodax.

Repositori belum menggunakan `go.mod`, jadi jalankan tes dalam mode GOPATH:

```sh
GO111MODULE=off go test -v
```

Jika ingin menjalankan satu tes tertentu:

```sh
GO111MODULE=off go test -run TestBuildQueryStringDeterministic -v
```

## Build dan menjalankan bot

Build hanya file sumber aplikasi. Jangan gunakan `go build *.go`, karena wildcard juga akan memasukkan `utils_test.go`.

```sh
go build -o bot main.go config.go market.go order.go utils.go
./bot
```

Atau jalankan tanpa membuat binary:

```sh
go run main.go config.go market.go order.go utils.go
```

Jalankan bot dari direktori repository agar `config.json` dan `state.json` ditemukan. Hentikan dengan `Ctrl + C`.

Bot mengirim order market secara nyata jika sinyal dan batas konfigurasi terpenuhi. Jika request order mengalami timeout atau hasil respons tidak dapat dipastikan, periksa order dan saldo di akun Indodax sebelum menjalankan ulang.

## State transaksi

`state.json` dibuat saat pertama kali disimpan. Jika file belum ada, bot memulai dengan status `KOSONG`. Jika file ada tetapi tidak valid atau tidak dapat dibaca, bot berhenti demi keamanan—jangan hapus atau ubah state tanpa terlebih dahulu membandingkannya dengan saldo dan riwayat akun.

Status lokal bukan sumber kebenaran saldo atau posisi akun. Periksa akun Indodax secara manual setelah gangguan, order yang hasilnya tidak pasti, atau kegagalan menyimpan state.

## Struktur project

```text
BotIndodax/
├── README.md
├── config.go       # pemuatan dan validasi konfigurasi
├── config.json     # konfigurasi lokal dan kredensial
├── main.go         # loop utama dan keputusan transaksi
├── market.go       # ticker dan saldo melalui API
├── order.go        # pengiriman dan pemeriksaan respons order
├── utils.go        # state, signature, query string, dan helper kuantitas
└── utils_test.go   # tes helper; tidak masuk ke binary aplikasi
```

## Troubleshooting

### `go test` menyatakan `go.mod` tidak ditemukan

Project ini belum memakai Go modules. Jalankan tes dengan mode GOPATH:

```sh
GO111MODULE=off go test -v
```

### Build gagal atau `utils_test.go` ikut disebut

Pastikan build menggunakan file aplikasi secara eksplisit, bukan `*.go`:

```sh
go build -o bot main.go config.go market.go order.go utils.go
```

### Konfigurasi ditolak

Periksa JSON, API key/secret, pair berakhiran `idr`, nilai trigger/modal/kuantitas positif, `sellTrigger` lebih tinggi dari `buyTrigger`, serta `maxDailyTrades` lebih besar dari nol.

### Gagal membaca harga/saldo atau order

Periksa koneksi internet, waktu perangkat, pair, izin API key, saldo, dan pesan kesalahan API. Jika hasil order tidak pasti, jangan langsung menjalankan ulang bot—periksa riwayat/order dan saldo di akun Indodax terlebih dahulu.

Jangan pernah membagikan API key atau secret key saat meminta bantuan.

## Disclaimer

Proyek ini disediakan apa adanya, tanpa jaminan apa pun dan bukan nasihat investasi. Anda bertanggung jawab atas keamanan API key, konfigurasi, pemantauan bot, dan seluruh keputusan trading.
