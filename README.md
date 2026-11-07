---

```markdown
# BotIndodax

Bot trading Indodax berbasis Go untuk dijalankan melalui terminal Android (Termux).
Bot memantau satu pair IDR dan dapat mengirim order market menggunakan API key akun Anda secara otomatis.

> [!CAUTION]
> Bot dapat melakukan transaksi dengan saldo asli dan menyebabkan kerugian.
> Tidak ada jaminan keuntungan.
> Periksa kode, konfigurasi, dan izin API key sebelum menjalankan bot.
> Jangan aktifkan izin penarikan (withdrawal) pada API key.
> DWYOR (Do With Your Own Risk)

## Fitur

- Membaca harga ticker dan saldo akun Indodax secara *real-time*.
- Membeli ketika status lokal `KOSONG` dan harga sama dengan atau di bawah `buyTrigger`.
- Menjual ketika status lokal `PEGANG` dan harga sama dengan atau di atas `sellTrigger`.
- Menyimpan status transaksi, waktu transaksi terakhir, dan batas transaksi harian di `state.json`.
- Memuat ulang `config.json` di setiap iterasi.
- Memvalidasi konfigurasi serta menangani kegagalan API dengan lebih aman (HTTP Timeout).
- Menghentikan bot bila hasil pengiriman order tidak pasti agar tidak otomatis mengirim order duplikat.
- Menampilkan *countdown* (hitung mundur) yang informatif di terminal saat menunggu iterasi berikutnya.

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

# Download dependencies (jika ada) dan siapkan go.mod
go mod tidy
```

## Konfigurasi

Bot membaca config.json dari direktori kerja. Buka file konfigurasi dan masukkan API key/secret Anda:

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

Field Keterangan
apiKey API key Indodax TAPI v2 dengan izin baca dan spot trading.
secretKey Secret key pasangan API key; rahasiakan.
pair Pair IDR, huruf kecil, misalnya ethidr.
buyTrigger Batas harga maksimum untuk sinyal beli.
sellTrigger Batas harga minimum untuk sinyal jual; harus lebih besar dari buyTrigger.
modalPerTrade Nominal IDR yang direncanakan untuk pembelian.
minQty Kuantitas aset minimum sesuai aturan pair.
cooldownMs Jeda minimum antar transaksi dalam milidetik.
maxDailyTrades Batas transaksi harian yang dikonfigurasi.

Pastikan pair dan aturan jumlah/nominal minimum sesuai ketentuan pasar Indodax. Jangan unggah atau commit config.json berisi kredensial. Jika kredensial pernah terpublikasi, segera cabut dan buat ulang API key.

Menjalankan pengujian

Tes ini memeriksa validasi konfigurasi, perilaku pembacaan dan penyimpanan state, serta konsistensi query string. Tes tidak mengirim order ke Indodax.

```sh
go test -v ./...
```

Jika ingin menjalankan satu tes tertentu:

```sh
go test -run TestBuildQueryStringDeterministic -v
```

Build dan menjalankan bot

Sangat disarankan untuk melakukan build (compile) menjadi satu binary agar bot berjalan lebih ringan dan hemat RAM di HP.

```sh
# Compile semua file .go di folder ini menjadi file bernama 'bot'
go build -o bot .

# Jalankan bot
./bot
```

Atau jalankan langsung tanpa membuat binary (cocok untuk tahap testing kode):

```sh
go run .
```

## ⚠️ PENTING UNTUK ANDROID:
Jalankan bot dari direktori repository agar config.json dan state.json ditemukan. Sebelum menjalankan bot, sangat disarankan untuk mengetik perintah berikut di Termux agar Android tidak mematikan bot saat layar mati:

```sh
termux-wake-lock
```

Hentikan bot dengan Ctrl + C.

Bot mengirim order market secara nyata jika sinyal dan batas konfigurasi terpenuhi. Jika request order mengalami timeout atau hasil respons tidak dapat dipastikan, periksa order dan saldo di akun Indodax sebelum menjalankan ulang.

State transaksi

state.json dibuat secara otomatis saat pertama kali transaksi berhasil disimpan. Jika file belum ada, bot memulai dengan status KOSONG. Jika file ada tetapi tidak valid atau tidak dapat dibaca, bot berhenti demi keamanan—jangan hapus atau ubah state tanpa terlebih dahulu membandingkannya dengan saldo dan riwayat akun.

Status lokal bukan sumber kebenaran saldo atau posisi akun. Periksa akun Indodax secara manual setelah gangguan, order yang hasilnya tidak pasti, atau kegagalan menyimpan state.

Struktur project

```text
BotIndodax/
├── README.md
├── config.go       # pemuatan dan validasi konfigurasi
├── config.json     # konfigurasi lokal dan kredensial (Tidak di-push ke Git)
├── go.mod          # Definisi Go module
├── main.go         # loop utama dan keputusan transaksi
├── market.go       # ticker dan saldo melalui API
├── order.go        # pengiriman dan pemeriksaan respons order
├── utils.go        # state, signature, query string, dan helper kuantitas
└── utils_test.go   # tes helper (tidak masuk ke binary aplikasi)
```

## Troubleshooting
Konfigurasi ditolak

Periksa JSON, API key/secret, pair berakhiran idr, nilai trigger/modal/kuantitas positif, sellTrigger lebih tinggi dari buyTrigger, serta maxDailyTrades lebih besar dari nol.

Gagal membaca harga/saldo atau order

Periksa koneksi internet, waktu perangkat, pair, izin API key, saldo, dan pesan kesalahan API. Jika hasil order tidak pasti, jangan langsung menjalankan ulang bot—periksa riwayat/order dan saldo di akun Indodax terlebih dahulu.

Bot mati sendiri saat layar HP mati

Ini adalah perilaku manajemen baterai Android. Pastikan Anda sudah menjalankan perintah termux-wake-lock sebelum menjalankan bot.

Jangan pernah membagikan API key atau secret key saat meminta bantuan.

## Disclaimer
Proyek ini disediakan apa adanya, tanpa jaminan apa pun dan bukan nasihat investasi.
Anda bertanggung jawab atas keamanan API key, konfigurasi, pemantauan bot, dan seluruh keputusan trading.
