# BotIndodax

Bot trading Indodax berbasis Go yang dirancang untuk dijalankan dari terminal Android melalui Termux. Bot memantau satu pair, memeriksa saldo saat sinyal terpenuhi, lalu dapat mengirim order **market sungguhan** ke akun Indodax.

> [!CAUTION]
> Bot dapat melakukan transaksi dengan saldo asli. Order market dapat langsung dieksekusi dan menimbulkan kerugian; tidak ada jaminan keuntungan. Tinjau kode dan konfigurasi, uji dengan hati-hati, dan gunakan dana yang siap Anda tanggung risikonya. Jangan pernah mengaktifkan izin penarikan pada API key. Jangan bagikan atau commit API key, secret key, maupun file konfigurasi berisi kredensial.

## Fitur

- Memantau harga ticker untuk pair yang diatur di `config.json`.
- Membeli ketika status `KOSONG` dan harga berada pada atau di bawah `buyTrigger`.
- Menjual ketika status `PEGANG` dan harga berada pada atau di atas `sellTrigger`.
- Menyimpan status dan batas transaksi harian di `state.json`.

Bot memeriksa kondisi setiap satu jam. Ini bukan sistem trading frekuensi tinggi dan bukan layanan yang berjalan di server.

## Persyaratan

- Perangkat Android dengan Termux. Pasang Termux dari sumber resmi yang masih dipelihara, seperti F-Droid atau repositori GitHub resmi Termux.
- Koneksi internet.
- Akun Indodax dan API key dengan izin baca saldo serta trading. **Jangan beri izin penarikan.**
- Paket `git` dan `golang` di Termux.

## Instalasi di Termux

Buka Termux dan jalankan:

```sh
pkg update
pkg upgrade
pkg install git golang nano

git clone https://github.com/MinatoAka/BotIndodax.git
cd BotIndodax
```

Pastikan Go sudah tersedia:

```sh
go version
```

## Siapkan konfigurasi

Program membaca konfigurasi langsung dari `config.json` pada direktori kerja. Buat cadangan file contoh, kemudian edit file yang dibaca bot:

```sh
cp config.json config.json.example
nano config.json
```

Isi `apiKey` dan `secretKey` dengan API key Indodax Anda. Jangan memasukkan spasi tambahan atau mengubah nama field. Contoh struktur (gunakan nilai Anda sendiri; angka di bawah hanya contoh dan **bukan rekomendasi trading**):

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
| `secretKey` | Secret key yang berpasangan dengan API key. Simpan rahasia. |
| `pair` | Pair pasar, misalnya `ethidr`; gunakan format pair yang didukung API. |
| `buyTrigger` | Bot mempertimbangkan pembelian jika status `KOSONG` dan harga ticker sama dengan atau lebih rendah dari angka ini. |
| `sellTrigger` | Bot mempertimbangkan penjualan jika status `PEGANG` dan harga ticker sama dengan atau lebih tinggi dari angka ini. |
| `modalPerTrade` | Perkiraan jumlah IDR untuk pembelian. Pastikan saldo cukup dan pahami biaya transaksi. |
| `minQty` | Kuantitas minimum agar bot meneruskan pembelian. |
| `cooldownMs` | Jeda minimum antara transaksi, dalam milidetik. `1800000` berarti 30 menit. |
| `maxDailyTrades` | Batas transaksi harian yang dicatat bot. |

Sebelum melanjutkan, periksa kembali pair, trigger harga, modal, dan limit. File konfigurasi berisi kredensial dalam teks biasa; jangan unggah, commit, atau kirimkan file tersebut. File `config.json` contoh berada di repository, jadi berhati-hatilah agar perubahan lokal berisi rahasia tidak ikut di-commit.

## Kompilasi di Termux

Pada versi kode saat ini, `order.go` menggunakan `strings` tetapi belum mengimpornya. Tambahkan import tersebut sebelum kompilasi:

```sh
nano order.go
```

Di blok `import`, tambahkan baris `"strings"`, lalu simpan (`Ctrl+O`, Enter) dan keluar (`Ctrl+X`). Blok import harus mencakup package berikut:

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

Kompilasi seluruh file Go menjadi executable bernama `bot`:

```sh
go build -o bot *.go
```

Jika berhasil, periksa file hasil kompilasi:

```sh
ls -l bot
```

Executable dibuat untuk lingkungan perangkat/arsitektur Termux yang sedang digunakan. Kompilasi ulang setelah memperbarui kode sumber.

## Menjalankan

Pastikan masih berada di direktori repository, konfigurasi sudah benar, lalu jalankan:

```sh
./bot
```

Atau jalankan tanpa membuat executable terlebih dahulu:

```sh
go run *.go
```

Saat mulai, bot membaca `config.json` dan `state.json`, mengambil harga, serta menampilkan status. Setelah satu putaran, bot menunggu satu jam sebelum pemeriksaan berikutnya. Hentikan dengan `Ctrl+C`.

Agar proses lebih kecil kemungkinannya dihentikan ketika Termux berada di latar belakang, Anda dapat menjalankan:

```sh
termux-wake-lock
```

Lepaskan wakelock setelah selesai dengan `termux-wake-unlock`. Pengaturan baterai Android tetap dapat membatasi proses latar belakang; wakelock tidak menjamin bot selalu berjalan. Pastikan koneksi internet dan Termux tetap aktif.

## Status dan perilaku transaksi

`state.json` berada pada direktori kerja saat bot dijalankan. Jika file belum ada atau tidak dapat dibaca, status awal yang digunakan adalah `KOSONG`. Setelah bot menganggap order berhasil dikirim, status berubah menjadi `PEGANG` untuk pembelian atau `KOSONG` untuk penjualan. File ini juga menyimpan waktu transaksi terakhir dan penghitung transaksi harian.

Jangan menghapus atau mengubah `state.json` sembarangan: status lokal yang tidak cocok dengan saldo akun dapat membuat keputusan bot berbeda dari yang diharapkan. Periksa saldo dan riwayat order di Indodax secara mandiri. Batas harian dihitung berdasarkan tanggal lokal perangkat.

## Troubleshooting

### `undefined: strings` saat kompilasi

Tambahkan import `"strings"` ke blok import di `order.go` seperti dijelaskan pada bagian kompilasi, lalu ulangi:

```sh
go build -o bot *.go
```

### Gagal membaca `config.json`

Pastikan perintah dijalankan dari direktori `BotIndodax`, file bernama tepat `config.json`, dan isinya JSON valid. Periksa tanda kutip, koma, serta tipe nilai.

### Gagal mengambil harga, saldo, atau mengirim order

Periksa internet, pair, waktu perangkat, saldo, izin API, serta pesan dan respons yang ditampilkan. Jangan mengirimkan API key atau secret key saat meminta bantuan. Periksa juga saldo dan riwayat order langsung di akun Indodax; jangan hanya mengandalkan status lokal bot.

## Disclaimer

Proyek ini disediakan apa adanya, tanpa jaminan apa pun dan bukan nasihat investasi. Anda bertanggung jawab atas keamanan API key, konfigurasi, pemantauan proses, serta seluruh keputusan dan transaksi pada akun Anda.
