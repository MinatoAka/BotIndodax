<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.22+"></a>
  <a href="https://termux.dev/"><img src="https://img.shields.io/badge/Platform-Termux-2E8B57?style=for-the-badge&logo=android&logoColor=white" alt="Platform: Termux"></a>
  <a href="https://saweria.co/AKaAturDompet"><img src="https://img.shields.io/badge/Support-Saweria-FF5722?style=for-the-badge&logo=buy-me-a-coffee&logoColor=white" alt="Support via Saweria"></a>
</p>

# BotIndodax

Bot trading Indodax berbasis Go untuk Termux. Bot memantau satu pair IDR dan mengirim order **market** otomatis berdasarkan batas harga yang dikonfigurasi.

> [!CAUTION]
> Bot dapat melakukan transaksi dengan dana asli dan menyebabkan kerugian. Tidak ada jaminan keuntungan. Periksa kode dan konfigurasi sebelum menjalankan. Gunakan API key dengan izin baca dan spot trading saja—**jangan aktifkan izin penarikan**. Anda bertanggung jawab atas seluruh transaksi.

## Fitur

- Memantau ticker dan saldo akun Indodax.
- Membeli saat status bot `KOSONG` dan harga mencapai atau berada di bawah `buyTrigger`.
- Menjual saat status bot `PEGANG` dan harga mencapai atau berada di atas `sellTrigger`.
- Menyimpan status transaksi dan batas transaksi harian di `state.json`.
- Memuat ulang `config.json` setiap iterasi.
- Menghentikan bot jika hasil pengiriman order tidak dapat dipastikan.

## Persyaratan

- Android dengan Termux dan koneksi internet.
- Go 1.22 atau lebih baru.
- Akun Indodax dan API key TAPI v2 dengan izin baca serta spot trading. Jangan beri izin `WITHDRAWAL`.

## Instalasi

```sh
pkg update
pkg upgrade
pkg install git golang nano

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
go mod tidy
```

## Konfigurasi

Edit `config.json` dan masukkan kredensial API Anda:

```sh
nano config.json
```

Contoh konfigurasi:

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
| `secretKey` | Secret key pasangan API key; simpan secara rahasia. |
| `pair` | Pair IDR dalam huruf kecil, misalnya `ethidr`. |
| `buyTrigger` | Batas harga untuk sinyal beli. |
| `sellTrigger` | Batas harga untuk sinyal jual; harus lebih besar dari `buyTrigger`. |
| `modalPerTrade` | Nominal IDR yang direncanakan untuk pembelian. |
| `minQty` | Kuantitas aset minimum sesuai aturan pair. |
| `cooldownMs` | Jeda minimum antartransaksi dalam milidetik. |
| `maxDailyTrades` | Batas jumlah transaksi harian. |

Pastikan pair serta batas minimum jumlah dan nominal sesuai aturan pasar Indodax. Jangan unggah atau commit `config.json` yang berisi kredensial. Jika kredensial terungkap, cabut API key tersebut dan buat yang baru.

## Pengujian

Pengujian tidak mengirim order ke Indodax.

```sh
go test ./...
```

## Build dan jalankan

Jalankan dari direktori repository agar `config.json` dan `state.json` dapat ditemukan.

```sh
go build -o bot .
./bot
```

Untuk menjalankan tanpa membuat binary:

```sh
go run .
```

Di Termux, gunakan `termux-wake-lock` agar perangkat tidak masuk mode tidur selama bot berjalan. Hentikan bot dengan `Ctrl+C`.

Bot mengirim order market secara nyata saat kondisi konfigurasi terpenuhi. Jika request order mengalami timeout atau responsnya tidak jelas, periksa order dan saldo akun Indodax sebelum menjalankan bot kembali.

## Dukungan

[♡ Dukung pengembangan via Saweria](https://saweria.co/AKaAturDompet)

## Disclaimer

Proyek ini disediakan apa adanya, tanpa jaminan, dan bukan nasihat investasi. Anda bertanggung jawab atas keamanan API key, konfigurasi, pemantauan bot, serta seluruh keputusan dan risiko trading.
