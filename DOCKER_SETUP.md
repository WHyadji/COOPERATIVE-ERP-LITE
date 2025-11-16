# Docker Setup Guide - Cooperative ERP Lite

Panduan lengkap untuk menjalankan Cooperative ERP Lite menggunakan Docker.

## 📋 Prerequisites

Pastikan sudah terinstall di sistem Anda:
- **Docker** (v20.10+): [Download Docker](https://docs.docker.com/get-docker/)
- **Docker Compose** (v2.0+): Biasanya sudah include dengan Docker Desktop
- **Make** (opsional, untuk command shortcuts): Pre-installed di macOS/Linux
- **Git**: Untuk clone repository

Cek instalasi:
```bash
docker --version
docker-compose --version
make --version  # Opsional
```

## 🚀 Quick Start (Untuk Pertama Kali)

Jika Anda baru pertama kali setup, jalankan satu command ini:

```bash
make quick-start
```

Command di atas akan:
1. ✅ Copy `.env.example` ke `.env`
2. ✅ Install Go dependencies
3. ✅ Install swag CLI untuk Swagger
4. ✅ Generate Swagger documentation
5. ✅ Build Docker images
6. ✅ Start semua services (PostgreSQL, Backend, Adminer)
7. ✅ Seed database dengan data sample

**Output:**
```
🎉 Everything is ready!

📡 Services:
   - API: http://localhost:8080/api/v1
   - Swagger UI: http://localhost:8080/swagger/index.html
   - Adminer: http://localhost:8081
   - Health Check: http://localhost:8080/health

🔑 Login Credentials:
   - Admin: username=admin, password=admin123
   - Bendahara: username=bendahara, password=bendahara123
   - Kasir: username=kasir, password=kasir123
```

## 📖 Manual Setup (Step-by-Step)

Jika ingin setup manual step by step:

### 1. Clone Repository

```bash
git clone https://github.com/your-org/cooperative-erp-lite.git
cd cooperative-erp-lite
```

### 2. Setup Environment Variables

```bash
# Copy environment template
cp backend/.env.example backend/.env

# Edit .env file dengan editor favorit
nano backend/.env  # atau vim, code, dll
```

**Konfigurasi penting di `.env`:**
```env
# Database (gunakan 'postgres' untuk Docker, 'localhost' untuk local)
DB_HOST=postgres

# JWT Secret (GANTI di production!)
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# CORS (tambahkan frontend URL)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

### 3. Build Docker Images

```bash
make build
# Atau manual:
docker-compose build
```

### 4. Start Services

```bash
make up
# Atau manual:
docker-compose up -d
```

### 5. Seed Database (Data Sample)

```bash
make seed
# Atau manual:
docker-compose exec backend go run cmd/seed/main.go
```

### 6. Verify Installation

```bash
# Check health
curl http://localhost:8080/health

# Should return:
{
  "success": true,
  "message": "Server berjalan dengan baik",
  "data": {
    "status": "healthy",
    "env": "debug"
  }
}
```

## 🎯 Available Services

Setelah `make up`, services berikut akan running:

| Service | URL | Description |
|---------|-----|-------------|
| **Backend API** | http://localhost:8080 | Go Gin REST API |
| **Swagger UI** | http://localhost:8080/swagger/index.html | Interactive API Documentation |
| **Health Check** | http://localhost:8080/health | Server health status |
| **Adminer** | http://localhost:8081 | Database management UI |
| **PostgreSQL** | localhost:5432 | Database (internal) |

### Adminer Login
- **System**: PostgreSQL
- **Server**: postgres
- **Username**: postgres
- **Password**: postgres
- **Database**: koperasi_erp

## 🛠️ Common Commands

### Docker Operations

```bash
# Start services (detached)
make up

# Stop services
make down

# Restart services
make restart

# View all logs (follow mode)
make logs

# View backend logs only
make logs-backend

# View PostgreSQL logs only
make logs-postgres

# Show running containers
make ps
```

### Development

```bash
# Run backend with auto-reload (requires Air)
make dev

# Run backend directly (without Docker)
make run

# Generate Swagger documentation
make swagger

# Seed database
make seed
```

### Database Operations

```bash
# Connect to PostgreSQL shell
make db-connect

# Backup database
make db-backup

# Restore from latest backup
make db-restore

# Drop database (WARNING: destructive!)
make db-drop
```

### Testing & Quality

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Format code
make fmt

# Tidy go modules
make tidy
```

### Cleanup

```bash
# Clean build artifacts and volumes
make clean

# Deep clean (including Docker images)
make clean-all
```

## 🧪 Testing API dengan curl

### 1. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "namaPengguna": "admin",
    "kataSandi": "admin123"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "pengguna": {
      "id": "...",
      "namaPengguna": "admin",
      "namaLengkap": "Administrator",
      "peran": "admin"
    }
  }
}
```

**Simpan token untuk request selanjutnya!**

### 2. Get Member List

```bash
export TOKEN="your-jwt-token-here"

curl -X GET "http://localhost:8080/api/v1/anggota?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Create Member

```bash
curl -X POST http://localhost:8080/api/v1/anggota \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "namaLengkap": "John Doe",
    "tempatLahir": "Jakarta",
    "tanggalLahir": "1990-01-01",
    "jenisKelamin": "L",
    "alamat": "Jl. Sudirman No. 1",
    "telepon": "08123456789",
    "email": "john@example.com",
    "pekerjaan": "Wiraswasta"
  }'
```

### 4. Get Financial Reports

```bash
# Neraca (Balance Sheet)
curl -X GET "http://localhost:8080/api/v1/laporan/neraca?tanggalPer=2025-01-16" \
  -H "Authorization: Bearer $TOKEN"

# Laba Rugi (Income Statement)
curl -X GET "http://localhost:8080/api/v1/laporan/laba-rugi?tanggalMulai=2025-01-01&tanggalAkhir=2025-01-16" \
  -H "Authorization: Bearer $TOKEN"
```

## 📚 Using Swagger UI

Swagger UI menyediakan dokumentasi interaktif untuk semua endpoints:

1. **Buka browser**: http://localhost:8080/swagger/index.html
2. **Authorize**: Click tombol "Authorize" di kanan atas
3. **Input Token**: Masukkan `Bearer your-jwt-token` (hasil login)
4. **Try Endpoints**: Click "Try it out" untuk test endpoint

## 🔧 Troubleshooting

### Port Already in Use

Jika port 8080 atau 5432 sudah digunakan:

```bash
# Check what's using the port
lsof -i :8080
lsof -i :5432

# Stop the process or change port di docker-compose.yml
```

### Database Connection Error

```bash
# Check if PostgreSQL is running
make ps

# View PostgreSQL logs
make logs-postgres

# Restart services
make restart
```

### Cannot Connect to Database

Jika backend tidak bisa connect ke database:

1. Check `backend/.env`:
   ```env
   DB_HOST=postgres  # Harus 'postgres' untuk Docker
   ```

2. Restart services:
   ```bash
   make restart
   ```

### Permission Denied

Jika ada error permission denied:

```bash
# Fix file permissions
sudo chown -R $USER:$USER .

# Or run Docker with proper permissions
sudo usermod -aG docker $USER
```

### Swagger Generation Failed

Jika Swagger docs tidak ter-generate:

```bash
# Install swag manually
go install github.com/swaggo/swag/cmd/swag@latest

# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Generate docs
cd backend
swag init -g cmd/api/main.go -o docs
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────┐
│         Docker Compose Network          │
│                                         │
│  ┌──────────┐      ┌──────────────┐   │
│  │          │      │              │   │
│  │  Adminer │◄─────┤  PostgreSQL  │   │
│  │  :8081   │      │    :5432     │   │
│  │          │      │              │   │
│  └──────────┘      └──────▲───────┘   │
│                           │           │
│                    ┌──────┴───────┐   │
│                    │              │   │
│                    │   Backend    │   │
│                    │   Go API     │   │
│                    │   :8080      │   │
│                    │              │   │
│                    └──────────────┘   │
│                                       │
└───────────────────────────────────────┘
           │
           │ HTTP
           ▼
    ┌──────────────┐
    │   Frontend   │
    │   Next.js    │
    │   :3000      │
    └──────────────┘
```

## 📊 Database Schema

Auto-migration akan create tables berikut:
- `koperasi` - Cooperative entities
- `pengguna` - Users (Admin, Bendahara, Kasir)
- `anggota` - Members
- `akun` - Chart of Accounts (31 accounts)
- `transaksi` - Journal entries
- `baris_transaksi` - Transaction line items
- `simpanan` - Share capital deposits
- `produk` - Products for POS
- `penjualan` - Sales transactions
- `item_penjualan` - Sale line items

## 🔐 Security Notes

**IMPORTANT untuk Production:**

1. **Change JWT Secret**:
   ```env
   # Generate strong secret
   openssl rand -base64 32

   # Update .env
   JWT_SECRET=<generated-secret>
   ```

2. **Use SSL for Database**:
   ```env
   DB_SSLMODE=require
   ```

3. **Set Gin to Release Mode**:
   ```env
   GIN_MODE=release
   ```

4. **Disable Swagger in Production**:
   Swagger will auto-disable when `GIN_MODE=release`

5. **Use Environment Variables**:
   Jangan commit `.env` ke Git!

## 📝 Next Steps

Setelah setup berhasil:

1. ✅ Test all API endpoints via Swagger UI
2. ✅ Customize seed data di `backend/cmd/seed/main.go`
3. ✅ Setup frontend Next.js (coming soon)
4. ✅ Configure production environment
5. ✅ Setup CI/CD pipeline

## 🆘 Getting Help

- **Documentation**: `/docs` directory
- **Issues**: [GitHub Issues](https://github.com/your-org/cooperative-erp-lite/issues)
- **Email**: support@cooperativeerp.com

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details
