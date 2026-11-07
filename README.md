# Backend REST API Suara Rakyat: Kritik Pajak (Golang + PostgreSQL)

Backend REST API berbasis **Golang** untuk aplikasi web **"Suara Rakyat: Menolak Pajak Mencekik"**. Dibangun dengan menerapkan prinsip **Clean Architecture / Standard Go Project Layout** yang modular, terpisah antara layer handler (controller), service (business logic), repository (akses data), dan model.

---

## 📁 Struktur Folder Proyek

```text
pajak-prj/
├── cmd/
│   └── server/
│       └── main.go                 # Titik masuk utama aplikasi (entrypoint) & graceful shutdown
├── config/
│   └── config.go                   # Pemuatan environment variable & konfigurasi DSN database
├── database/
│   ├── postgres.go                 # Inisialisasi koneksi PostgreSQL & connection pooling
│   └── schema.sql                  # Skema DDL tabel dan data awal (seeding)
├── internal/
│   ├── handler/
│   │   ├── comment_handler.go      # HTTP Controller (GET /api/comments, POST /api/comments)
│   │   └── middleware.go           # Middleware CORS dan Logging request
│   ├── model/
│   │   └── comment.go              # Struct entitas data, DTO Request/Response, & validasi
│   ├── repository/
│   │   └── comment_repository.go   # Interface & implementasi query SQL database
│   └── service/
│       └── comment_service.go      # Interface & implementasi business logic
├── .env                            # Konfigurasi environment lokal (diabaikan dari VCS)
├── .env.example                    # Template environment variable
├── go.mod                          # Definisi Go module dan dependensi
├── go.sum                          # Checksum verifikasi dependensi Go
├── kritik_pajak.html               # File frontend HTML interaktif
└── README.md                       # Dokumentasi panduan lengkap
```

---

## 🗄️ Skema Database PostgreSQL

File skema berada di `database/schema.sql`. Anda dapat menjalankannya langsung di PostgreSQL:

```sql
-- Buat Tabel comments
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'Umum',
    comment TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk mempercepat query sorting berdasarkan waktu terbaru
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);

-- Data Dummy Awal
INSERT INTO comments (name, category, comment, created_at) VALUES
('Rina S.', 'Pajak Pembalut', 'Pembalut itu barang primer, menstruasi bukan pilihan! Sangat tidak etis kalau kebutuhan biologis bulanan wanita dibebani PPN.', NOW() - INTERVAL '2 hours'),
('Budi K.', 'Sektor Informal (Cukur/Pijat)', 'Tukang cukur dan pijat keliling itu kerja keras sambung hidup. Seharusnya dibantu permodalannya, bukan malah diincar pajaknya.', NOW() - INTERVAL '1 hour'),
('Dewi A.', 'PPN 12%', 'Daya beli sedang turun dan banyak PHK, kenaikan PPN 12% semakin mendongkrak harga kebutuhan harian keluarga.', NOW() - INTERVAL '30 minutes')
ON CONFLICT DO NOTHING;
```

---

## 🚀 Panduan Menjalankan Secara Lokal (Local Development)

### 1. Prasyarat
- **Go** versi 1.22 atau yang lebih baru (`go version`).
- **PostgreSQL** aktif dan berjalan (`localhost:5432`).

### 2. Persiapan Database
1. Buka terminal PostgreSQL (`psql`):
   ```bash
   psql -U postgres
   ```
2. Buat database baru:
   ```sql
   CREATE DATABASE pajak_db;
   \c pajak_db;
   ```
3. Eksekusi file `database/schema.sql`:
   ```bash
   psql -U postgres -d pajak_db -f database/schema.sql
   ```

### 3. Konfigurasi Environment Variable
Salin file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```
Sesuaikan konfigurasi pada file `.env` dengan kredensial PostgreSQL lokal Anda:
```env
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password_postgres_anda
DB_NAME=pajak_db
DB_SSLMODE=disable
```

### 4. Mengunduh Dependensi & Menjalankan Aplikasi
```bash
# Unduh dependensi yang diperlukan
go mod tidy

# Jalankan server
go run cmd/server/main.go
```
Jika berhasil, terminal akan menampilkan:
```
==================================================
 Memulai Suara Rakyat: Pajak Mencekik Backend API 
==================================================
[SUCCESS] Berhasil terhubung ke database PostgreSQL (localhost:5432/pajak_db)
[READY] Server berjalan di http://0.0.0.0:8080
[INFO] API Comments siap menerima request di /api/comments
```

### 5. Membuka Frontend
Buka file `kritik_pajak.html` di browser Anda (bisa langsung klik dua kali atau gunakan ekstensi Live Server). Frontend telah otomatis terhubung ke endpoint `http://localhost:8080/api/comments`.

---

## 📡 Dokumentasi Endpoint REST API

### 1. Healthcheck
- **Method:** `GET`
- **URL:** `/api/health` atau `/health`
- **Response (200 OK):**
  ```json
  {
    "status": "healthy",
    "service": "Pajak Mencekik Backend API"
  }
  ```

### 2. Mengambil Seluruh Kritik / Aspirasi
- **Method:** `GET`
- **URL:** `/api/comments` atau `/api/comments?limit=20`
- **Response (200 OK):**
  ```json
  {
    "success": true,
    "message": "Berhasil mengambil daftar kritik",
    "data": [
      {
        "id": 1,
        "name": "Rina S.",
        "category": "Pajak Pembalut",
        "comment": "Pembalut itu barang primer, menstruasi bukan pilihan!...",
        "created_at": "2026-11-05T09:24:00Z"
      }
    ],
    "total": 1
  }
  ```

### 3. Mengirimkan Kritik Baru
- **Method:** `POST`
- **URL:** `/api/comments`
- **Header:** `Content-Type: application/json`
- **Body Request:**
  ```json
  {
    "name": "Siti Rahma",
    "category": "Pajak Pembalut",
    "comment": "Jangan bebankan PPN pada barang higienitas esensial perempuan!"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "success": true,
    "message": "Aspirasi/kritik berhasil dipublikasikan",
    "data": {
      "id": 4,
      "name": "Siti Rahma",
      "category": "Pajak Pembalut",
      "comment": "Jangan bebankan PPN pada barang higienitas esensial perempuan!",
      "created_at": "2026-11-05T11:25:00Z"
    }
  }
  ```

---

## 🌐 Panduan Deployment ke Server Produksi (Linux / Nginx)

Panduan berikut mengasumsikan:
- Server OS: **Ubuntu / Debian Linux**.
- Domain dan DNS di Bind9 sudah mengarah ke IP Server Anda (misal domain: `suarapajak.id` atau subdomain `aspirasi.domainanda.com`).
- Nginx sudah terinstal di server.

---

### Langkah 1: Compile Binary Go untuk Linux
Lakukan build executable binary Golang target Linux:
```bash
# Dari komputer lokal (jika cross-compiling dari Windows):
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o pajak-server cmd/server/main.go

# Atau jika langsung di server Linux:
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pajak-server cmd/server/main.go
```

Kirimkan file binary `pajak-server`, folder `database/schema.sql`, file `.env`, dan `kritik_pajak.html` ke server (misal ke direktori `/var/www/pajak-prj`).

---

### Langkah 2: Setup Systemd Service di Server
Jalankan backend Go sebagai background daemon dengan Systemd agar aplikasi otomatis menyala saat server restart dan hidup kembali jika crash.

1. Buat file service:
   ```bash
   sudo nano /etc/systemd/system/pajak-backend.service
   ```
2. Masukkan konfigurasi berikut:
   ```ini
   [Unit]
   Description=Pajak Mencekik Backend API Golang
   After=network.target postgresql.service

   [Service]
   Type=simple
   User=www-data
   Group=www-data
   WorkingDirectory=/var/www/pajak-prj
   ExecStart=/var/www/pajak-prj/pajak-server
   Restart=always
   RestartSec=5
   EnvironmentFile=/var/www/pajak-prj/.env

   # Hardening Keamanan
   ProtectSystem=full
   NoNewPrivileges=true

   [Install]
   WantedBy=multi-user.target
   ```
3. Aktifkan dan jalankan service:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable pajak-backend
   sudo systemctl start pajak-backend
   sudo systemctl status pajak-backend
   ```

---

### Langkah 3: Konfigurasi Nginx (Web Server & Reverse Proxy)

1. Buat konfigurasi server block baru di Nginx:
   ```bash
   sudo nano /etc/nginx/sites-available/pajak.conf
   ```

2. Masukkan konfigurasi lengkap berikut:
   ```nginx
   # Upstream backend Golang yang berjalan di port 8080
   upstream golang_backend {
       server 127.0.0.1:8080;
       keepalive 32;
   }

   server {
       listen 80;
       server_name suarapajak.id www.suarapajak.id; # Ganti dengan domain yang di-setting pada DNS Bind9

       # Direktori tempat file frontend disimpan
       root /var/www/pajak-prj;
       index kritik_pajak.html index.html;

       # Log akses dan error
       access_log /var/log/nginx/pajak_access.log;
       error_log /var/log/nginx/pajak_error.log;

       # Batasan ukuran payload upload
       client_max_body_size 2M;

       # 1. Melayani Halaman Frontend HTML
       location / {
           try_files $uri $uri/ /kritik_pajak.html;
           add_header Cache-Control "no-cache, public, must-revalidate";
       }

       # 2. Reverse Proxy untuk REST API Golang (/api/ dan /health)
       location /api/ {
           proxy_pass http://golang_backend;
           
           # Header Standar Reverse Proxy
           proxy_http_version 1.1;
           proxy_set_header Connection "";
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;

           # Timeout proxy
           proxy_connect_timeout 10s;
           proxy_send_timeout 30s;
           proxy_read_timeout 30s;
       }

       location /health {
           proxy_pass http://golang_backend/health;
           proxy_set_header Host $host;
       }

       # Keamanan: Blokir akses ke file sensitif (.env, .git, .go)
       location ~ /\.(?!well-known).* {
           deny all;
       }

       location ~ \.(go|mod|sum|sql|env)$ {
           deny all;
       }
   }
   ```

3. Aktifkan konfigurasi server block dan reload Nginx:
   ```bash
   sudo ln -s /etc/nginx/sites-available/pajak.conf /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl reload nginx
   ```

---

### Langkah 4: Setup SSL / HTTPS Gratis (Certbot Let's Encrypt - Opsional tapi Direkomendasikan)
Setelah domain di Bind9 DNS mengarah ke IP server publik:
```bash
sudo apt update
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d suarapajak.id -d www.suarapajak.id
```
Certbot akan otomatis memperbarui file konfigurasi Nginx dengan sertifikat SSL dan mengalihkan traffic HTTP (Port 80) ke HTTPS (Port 443).

---

## 💡 Catatan Integrasi Frontend JavaScript

Pada `kritik_pajak.html`, variabel `API_BASE_URL` telah disetting secara otomatis:
- Saat dibuka langsung via `file://` atau port selain 8080 (misalnya localhost saat tes lokal), request akan menuju `http://localhost:8080/api/comments`.
- Saat dideploy di server produksi di balik Nginx, request akan menuju path relatif `/api/comments` yang langsung diteruskan oleh reverse proxy Nginx ke backend Golang port 8080 tanpa masalah CORS.

