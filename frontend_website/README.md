# Travora Manager — Next.js

Versi Next.js dari aplikasi Flutter `admin_website`.

## Menjalankan

```bash
npm install
npm run dev
```

Buka `http://localhost:3000`. Backend secara default dibaca dari
`http://localhost:8080`. Untuk menggantinya, salin `.env.example` menjadi
`.env.local` lalu ubah `NEXT_PUBLIC_API_URL`.

Fitur forgot password memerlukan konfigurasi SMTP pada `backend/.env`.
Lihat contoh variabelnya di `backend/.env.example`. Jika memakai Gmail,
aktifkan verifikasi 2 langkah lalu gunakan App Password.

## Session admin

Autentikasi admin menggunakan cookie `HttpOnly` dan tabel `admin_sessions`.
Tambahkan URL frontend pada `backend/.env`:

```env
FRONTEND_URL=http://localhost:3000
GIN_MODE=debug
```

Untuk production, gunakan HTTPS dan `GIN_MODE=release` agar cookie diberi
atribut `Secure`. Jika **Remember this device** dipilih, sesi berlaku 7 hari;
tanpa pilihan tersebut cookie berakhir ketika browser ditutup.

## Build produksi

```bash
npm run build
npm start
```
