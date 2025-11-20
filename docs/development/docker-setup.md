# Docker Setup Guide - Cooperative ERP Lite

Production-ready Docker configuration for running the full stack (Frontend + Backend + Database) locally.

## 📋 Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Configuration](#configuration)
- [Development Workflow](#development-workflow)
- [Production Deployment](#production-deployment)
- [Troubleshooting](#troubleshooting)

---

## 🎯 Overview

The Docker setup includes:

1. **Frontend**: Next.js 15 with standalone build (~229MB image)
2. **Backend**: Go 1.25 API server
3. **Database**: PostgreSQL 15
4. **Nginx**: Reverse proxy routing frontend and backend

### Key Features

✅ **Production-Ready**
- Multi-stage builds for optimal image sizes
- Non-root users for security
- Health checks for all services
- Proper logging and monitoring

✅ **Zero Technical Debt**
- Follows Next.js 15 best practices
- Standalone output mode (75% size reduction)
- Optimized layer caching
- Reproducible builds

✅ **Developer-Friendly**
- Single command setup: `make quick-start`
- Hot reload in development mode
- Separate dev and prod workflows

---

## 📦 Prerequisites

- **Docker**: 20.10+ (with Docker Compose V2)
- **Make**: GNU Make (pre-installed on macOS/Linux)
- **Go**: 1.25.4 (for local development)
- **Node.js**: 20+ (for local development)

---

## 🚀 Quick Start

### Option 1: Full Docker Stack (Recommended)

```bash
# Complete setup and start all services
make quick-start

# Access the application
# Frontend: http://localhost
# Backend API: http://localhost/api/v1
# Swagger Docs: http://localhost/swagger/index.html
```

### Option 2: Development Mode (without Docker)

```bash
# Install dependencies
make setup

# Terminal 1: Start backend
make dev

# Terminal 2: Start frontend
make dev-frontend

# Or run both in parallel
make dev-all
```

---

## 🏗️ Architecture

### Service Architecture

```
┌─────────────────────────────────────────────────┐
│                   Nginx (:80)                   │
│  Routes traffic to frontend and backend         │
└───────────┬─────────────────────┬───────────────┘
            │                     │
    ┌───────▼────────┐    ┌──────▼────────┐
    │   Frontend     │    │    Backend    │
    │   Next.js      │    │      Go       │
    │   Port: 3000   │    │   Port: 8080  │
    └────────────────┘    └───────┬───────┘
                                  │
                          ┌───────▼────────┐
                          │   PostgreSQL   │
                          │   Port: 5432   │
                          └────────────────┘
```

### Routing Rules

| Path | Service | Description |
|------|---------|-------------|
| `/` | Frontend | Main application (Next.js) |
| `/api/*` | Backend | REST API endpoints |
| `/swagger/*` | Backend | API documentation |
| `/health` | Backend | Health check endpoint |
| `/_next/*` | Frontend | Next.js assets |

---

## ⚙️ Configuration

### Frontend Configuration

**File**: `frontend/.env.example`

```bash
# Copy and customize for local development
cp frontend/.env.example frontend/.env.local

# Edit .env.local
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
```

**Build-time vs Runtime**:
- `NEXT_PUBLIC_*` variables are baked into JS at build time
- Set via `build.args` in docker-compose.yml for Docker builds
- For local dev, use `.env.local`

### Backend Configuration

**File**: `backend/.env` (created during setup)

Default values are suitable for Docker Compose. No changes needed.

### Nginx Configuration

**File**: `nginx/conf.d/app.conf`

- Routes `/api` and `/swagger` to backend
- Routes everything else to frontend
- Includes rate limiting and security headers
- Commented HTTPS section for production

---

## 💻 Development Workflow

### Available Make Commands

```bash
# Setup and Installation
make help              # Show all available commands
make setup             # Initial project setup
make frontend-setup    # Install frontend dependencies only

# Docker Operations
make build             # Build all Docker images
make up                # Start all services
make down              # Stop all services
make restart           # Restart all services
make logs              # View logs from all services
make status            # Show service status

# Development Mode (no Docker)
make dev               # Run backend only
make dev-frontend      # Run frontend only
make dev-all           # Run both in parallel

# Frontend-Specific
make frontend-build    # Build frontend for production
make frontend-logs     # View frontend container logs
make frontend-rebuild  # Rebuild frontend container only

# Database
make db-shell          # Connect to PostgreSQL shell
make db-migrate        # Run migrations (handled by GORM)

# Cleanup
make clean             # Remove all build artifacts and volumes
```

### Typical Development Flow

```bash
# 1. Initial setup (first time only)
make quick-start

# 2. Make changes to code
# Frontend: frontend/app/**/*.tsx
# Backend: backend/**/*.go

# 3. Rebuild specific service
make frontend-rebuild  # If frontend changed
make restart           # If backend changed

# 4. View logs
make frontend-logs     # Frontend logs
make logs              # All logs

# 5. Cleanup when done
make down
```

### Hot Reload in Development

For faster development with hot reload:

```bash
# Run services locally (outside Docker)
make dev-all

# Frontend: http://localhost:3000 (with HMR)
# Backend: http://localhost:8080 (with Air - auto-reload)
```

---

## 🌐 Production Deployment

### Docker Image Sizes

| Service | Without Optimization | With Optimization | Reduction |
|---------|---------------------|-------------------|-----------|
| Frontend | ~892 MB | ~229 MB | **75%** |
| Backend | ~892 MB | ~20 MB | **98%** |

### Optimization Techniques

**Frontend (Next.js)**:
- ✅ Standalone output mode (`output: 'standalone'`)
- ✅ Multi-stage build (4 stages)
- ✅ Alpine Linux base image
- ✅ Non-root user (nextjs:1001)
- ✅ Layer caching via .dockerignore

**Backend (Go)**:
- ✅ Multi-stage build (2 stages)
- ✅ Static binary compilation
- ✅ Scratch base image (minimal)
- ✅ Non-root user (appuser:1000)

### Cloud Run Deployment

The Docker images are Cloud Run ready:

```bash
# Build for Cloud Run
docker build -t gcr.io/PROJECT_ID/frontend:latest ./frontend
docker build -t gcr.io/PROJECT_ID/backend:latest .

# Push to Google Container Registry
docker push gcr.io/PROJECT_ID/frontend:latest
docker push gcr.io/PROJECT_ID/backend:latest

# Deploy to Cloud Run
gcloud run deploy frontend \
  --image gcr.io/PROJECT_ID/frontend:latest \
  --set-env-vars NEXT_PUBLIC_API_BASE_URL=https://api.your-domain.com/api/v1

gcloud run deploy backend \
  --image gcr.io/PROJECT_ID/backend:latest \
  --set-env-vars DB_HOST=CLOUD_SQL_PROXY
```

### Environment Variables for Production

**Frontend**:
```bash
NODE_ENV=production
NEXT_PUBLIC_API_BASE_URL=https://your-domain.com/api/v1
```

**Backend**:
```bash
GIN_MODE=release
DB_HOST=your-cloud-sql-instance
JWT_SECRET=your-secure-secret-key
```

---

## 🔧 Troubleshooting

### Common Issues

#### 1. Port Already in Use

**Error**: `Bind for 0.0.0.0:80 failed: port is already allocated`

**Solution**:
```bash
# Check what's using port 80
sudo lsof -i :80

# Stop conflicting service or change port in docker-compose.yml
ports:
  - "8080:80"  # Changed from 80:80
```

#### 2. Frontend Not Building

**Error**: `Module not found` or `npm install fails`

**Solution**:
```bash
# Clean and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install

# Or rebuild Docker image
make frontend-rebuild
```

#### 3. Database Connection Failed

**Error**: `connection refused` to PostgreSQL

**Solution**:
```bash
# Check if PostgreSQL is healthy
docker compose ps postgres

# View PostgreSQL logs
docker compose logs postgres

# Restart database
docker compose restart postgres
```

#### 4. API Not Accessible

**Error**: 404 when accessing `/api/v1/*`

**Solution**:
```bash
# Check nginx configuration
docker compose exec nginx nginx -t

# View nginx error logs
docker compose logs nginx

# Ensure backend is running
docker compose ps backend
```

#### 5. Frontend Shows "API Error"

**Error**: Frontend can't connect to backend

**Solution**:
```bash
# Check environment variable
docker compose exec frontend env | grep NEXT_PUBLIC_API_BASE_URL

# Should be: http://localhost/api/v1 (for browser)
# If wrong, update docker-compose.yml and rebuild
make frontend-rebuild
```

### Health Checks

All services have health checks:

```bash
# View health status
docker compose ps

# Healthy services show: Up (healthy)
# Unhealthy services show: Up (unhealthy)
```

### Viewing Logs

```bash
# All services
make logs

# Specific service
docker compose logs -f frontend
docker compose logs -f backend
docker compose logs -f postgres
docker compose logs -f nginx
```

### Container Shell Access

```bash
# Frontend container
docker compose exec frontend sh

# Backend container
docker compose exec backend sh

# PostgreSQL shell
make db-shell
```

---

## 📚 Additional Resources

- [Next.js Docker Documentation](https://nextjs.org/docs/deployment#docker-image)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Nginx Configuration Guide](https://nginx.org/en/docs/)

---

## 🔐 Security Considerations

1. **Non-root Users**: All containers run as non-root users
2. **Secrets Management**: Never commit `.env` files
3. **Rate Limiting**: Nginx implements rate limiting on login endpoints
4. **HTTPS**: Enable SSL in production (nginx/conf.d/app.conf)
5. **CORS**: Restrict origins in production

---

## 📝 File Structure

```
COOPERATIVE-ERP-LITE/
├── frontend/
│   ├── Dockerfile              # Multi-stage Next.js build
│   ├── .dockerignore          # Optimized for layer caching
│   ├── .env.example           # Environment variable template
│   └── next.config.ts         # Standalone output enabled
├── backend/
│   ├── Dockerfile             # Exists (Go multi-stage build)
│   └── .env                   # Created by make setup
├── nginx/
│   ├── conf.d/
│   │   └── app.conf           # Main routing configuration
│   └── nginx.conf             # Main nginx config
├── docker-compose.yml         # All services defined here
└── Makefile                   # Development commands
```

---

## 🎉 Success Indicators

After running `make quick-start`, you should see:

```
✅ Cooperative ERP Lite is running!
📍 Frontend: http://localhost
📍 Backend API: http://localhost/api/v1
📍 Swagger Docs: http://localhost/swagger/index.html
```

Visit http://localhost and you should see the Next.js frontend with the login page.

---

**Last Updated**: 2025-11-18
**Docker Version**: 20.10+
**Docker Compose**: V2
