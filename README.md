# BotIndodax

Bot trading Indodax berbasis Go yang dapat dijalankan di Termux. Bot memantau harga pair yang dikonfigurasi dan dapat mengirim **order market sungguhan** melalui API Indodax.

> **Peringatan:** Bot ini dapat menggunakan saldo asli dan menyebabkan kerugian. Periksa kode dan konfigurasi sebelum menjalankannya. Mulai dengan nominal kecil. Jangan pernah membagikan API key atau secret key.

## Persyaratan

- Android dengan Termux. Disarankan memasang Termux dari sumber resmi seperti F-Droid atau GitHub Termux.
- Koneksi internet.
- Akun Indodax dan API key dengan izin yang diperlukan untuk membaca saldo dan melakukan trading. **Jangan aktifkan izin penarikan dana.**
- Git dan Go di Termux.

## Instalasi

Buka Termux, lalu jalankan:

```sh
pkg update && pkg upgrade
pkg install git golang

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
```

## Konfigurasi

Buat salinan file konfigurasi contoh, lalu edit:

```sh
cp config.json config.local.json
nano config.json
```

Isi `apiKey` dan `secretKey` dengan kredensial API Indodax milik Anda. Jangan commit atau membagikan file yang berisi kredensial tersebut.

| Opsi | Keterangan |
| --- | --- |
| `apiKey` | API key Indodax |
| `secretKey` | Secret key untuk API key tersebut |
| `pair` | Pair pasar, misalnya `ethidr` |
| `buyTrigger` | Bot membeli jika status `KOSONG` dan harga sama dengan atau di bawah nilai ini |
| `sellTrigger` | Bot menjual jika status `PEGANG` dan harga sama dengan atau di atas nilai ini |
| `modalPerTrade` | Perkiraan modal IDR untuk setiap pembelian |
| `minQty` | Jumlah aset minimum agar pembelian diteruskan |
| `cooldownMs` | Jeda minimum antatransaksi dalam milidetik; `1800000` = 30 menit |
| `maxDailyTrades` | Batas transaksi harian |

Sesuaikan nilai contoh dengan strategi dan toleransi risiko Anda. Pastikan pair, harga pemicu, modal, serta batas transaksi sudah benar sebelum bot dijalankan.

## Menjalankan

> **Catatan kondisi kode saat ini:** beberapa file Go belum mengimpor package yang dipakai. Jika muncul `undefined: strconv` atau `undefined: strings`, tambahkan import yang disebutkan di bagian [Troubleshooting](#troubleshooting) sebelum menjalankan perintah berikut.

```sh
go run *.go
```

Bot akan mencetak harga dan status, lalu mengulang pemeriksaan setiap **1 jam**. Biarkan Termux tetap berjalan dan perangkat terhubung ke internet. Android dapat menghentikan aplikasi di latar belakang; pengaturan baterai perangkat mungkin perlu disesuaikan.

Hentikan bot dengan `Ctrl+C`.

## Status bot

Bot menyimpan status transaksi di `state.json` pada direktori kerja saat berjalan. Status awalnya `KOSONG`; setelah pembelian berhasil menjadi `PEGANG`, dan setelah penjualan berhasil kembali menjadi `KOSONG`. Jangan menghapus atau mengubah file ini tanpa memahami dampaknya: status lokal bot bisa tidak sesuai dengan aset yang sebenarnya ada di akun.

## Troubleshooting

### `undefined: strconv` saat kompilasi

Di `main.go`, tambahkan `"strconv"` ke daftar import.

### `undefined: strings` saat kompilasi

Di `order.go`, tambahkan `"strings"` ke daftar import.

Setelah memperbaiki import, coba lagi:

```sh
go run *.go
```

### Gagal membaca `config.json`

Pastikan perintah dijalankan dari direktori `BotIndodax`, dan file `config.json` berisi JSON valid.

### Gagal mengambil harga, saldo, atau mengirim order

Periksa koneksi internet, pair, API key, secret key, izin API, dan respons/error yang ditampilkan bot. Jangan mengirimkan kredensial API saat meminta bantuan.

## Disclaimer

Proyek ini disediakan apa adanya, bukan saran investasi dan tanpa jaminan keuntungan. Anda bertanggung jawab atas konfigurasi, keamanan kredensial, dan semua transaksi yang dilakukan akun Anda.
