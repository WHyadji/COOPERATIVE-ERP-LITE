# Quick Start Guide - Development Setup

## Get Started in 30 Minutes

This guide will get you from zero to coding in 30 minutes.

---

## Prerequisites

### Install Required Software

**1. Go 1.21+**
```bash
# macOS
brew install go

# Verify
go version  # Should show 1.21 or higher
```

**2. Node.js 20+**
```bash
# macOS
brew install node

# Verify
node --version  # Should show 20.x or higher
npm --version
```

**3. PostgreSQL 15**
```bash
# macOS - Using Docker (Recommended)
docker pull postgres:15

# OR install directly
brew install postgresql@15

# Verify
psql --version
```

**4. Git**
```bash
# macOS
brew install git

# Configure
git config --global user.name "Your Name"
git config --global user.email "your@email.com"
```

**5. VS Code (Recommended IDE)**
```bash
# Download from https://code.visualstudio.com

# Install Extensions:
- Go (by Go Team at Google)
- ES7+ React/Redux/React-Native snippets
- Prettier - Code formatter
- GitLens
```

---

## Project Setup (15 Minutes)

### Step 1: Create Project Structure

```bash
# Navigate to your workspace
cd ~/Documents/VISI-DIGITAL-TERPADU

# Backend (Go)
mkdir -p COOPERATIVE-ERP-LITE/backend
cd COOPERATIVE-ERP-LITE/backend

# Initialize Go module
go mod init github.com/yourusername/koperasi-erp

# Install dependencies
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/joho/godotenv
go get github.com/go-playground/validator/v10

# Create folder structure
mkdir -p cmd/api
mkdir -p internal/{models,handlers,middleware,services,database,config}
mkdir -p pkg/{utils,logger}
```

### Step 2: Create Backend Files

**`cmd/api/main.go`**
```go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/yourusername/koperasi-erp/internal/config"
    "github.com/yourusername/koperasi-erp/internal/database"
    "github.com/yourusername/koperasi-erp/internal/handlers"
)

func main() {
    // Load config
    cfg := config.Load()

    // Connect to database
    db := database.Connect(cfg.DatabaseURL)
    database.AutoMigrate(db)

    // Setup router
    r := gin.Default()

    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // API routes
    api := r.Group("/api")
    {
        // Auth routes
        auth := api.Group("/auth")
        {
            authHandler := handlers.NewAuthHandler(db)
            auth.POST("/login", authHandler.Login)
            auth.POST("/logout", authHandler.Logout)
        }

        // Protected routes will go here
    }

    // Start server
    port := cfg.Port
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    r.Run(":" + port)
}
```

**`internal/config/config.go`**
```go
package config

import (
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL string
    JWTSecret   string
    Port        string
}

func Load() *Config {
    godotenv.Load()

    return &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/koperasi_erp?sslmode=disable"),
        JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
        Port:        getEnv("PORT", "8080"),
    }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
```

**`internal/database/database.go`**
```go
package database

import (
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/yourusername/koperasi-erp/internal/models"
)

func Connect(databaseURL string) *gorm.DB {
    db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    log.Println("Database connected successfully")
    return db
}

func AutoMigrate(db *gorm.DB) {
    err := db.AutoMigrate(
        &models.Cooperative{},
        &models.User{},
        &models.Member{},
        &models.ShareCapital{},
        &models.Account{},
        &models.Transaction{},
        &models.TransactionLine{},
        &models.Product{},
        &models.Sale{},
        &models.SaleItem{},
    )

    if err != nil {
        log.Fatal("Migration failed:", err)
    }

    log.Println("Database migrated successfully")
}
```

**`internal/models/user.go`**
```go
package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    CooperativeID uuid.UUID `gorm:"type:uuid;not null"`
    Username      string    `gorm:"unique;not null"`
    PasswordHash  string    `gorm:"not null"`
    FullName      string    `gorm:"not null"`
    Role          string    `gorm:"not null"` // admin, treasurer, cashier, member
    IsActive      bool      `gorm:"default:true"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

**`.env`**
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/koperasi_erp?sslmode=disable
JWT_SECRET=change-this-in-production
PORT=8080
GIN_MODE=debug
```

**`.gitignore`**
```
# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
vendor/

# Environment
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log
```

### Step 3: Create Frontend

```bash
# Go back to project root
cd ..

# Create Next.js app
npx create-next-app@latest frontend --typescript --tailwind --app --no-src-dir

cd frontend

# Install dependencies
npm install axios react-hook-form @hookform/resolvers zod
npm install @mui/material @emotion/react @emotion/styled @mui/icons-material
npm install js-cookie
npm install -D @types/js-cookie
```

**`frontend/.env.local`**
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

**`frontend/lib/api.ts`**
```typescript
import axios from 'axios';
import Cookies from 'js-cookie';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests
api.interceptors.request.use((config) => {
  const token = Cookies.get('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      Cookies.remove('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
```

**`frontend/app/login/page.tsx`**
```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import api from '@/lib/api';
import Cookies from 'js-cookie';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    try {
      const response = await api.post('/auth/login', {
        username,
        password,
      });

      Cookies.set('token', response.data.token);
      router.push('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Login failed');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-96">
        <h1 className="text-2xl font-bold mb-6 text-center">
          Koperasi ERP Lite
        </h1>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-gray-700 text-sm font-bold mb-2">
              Username
            </label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:border-blue-500"
              required
            />
          </div>

          <div className="mb-6">
            <label className="block text-gray-700 text-sm font-bold mb-2">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:border-blue-500"
              required
            />
          </div>

          <button
            type="submit"
            className="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition"
          >
            Login
          </button>
        </form>
      </div>
    </div>
  );
}
```

### Step 4: Setup Database

```bash
# Start PostgreSQL with Docker
docker run --name koperasi-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=koperasi_erp \
  -p 5432:5432 \
  -d postgres:15

# Verify it's running
docker ps

# Connect to database (optional, to check)
docker exec -it koperasi-postgres psql -U postgres -d koperasi_erp
# Type \q to exit
```

### Step 5: Run Everything

**Terminal 1: Backend**
```bash
cd backend
go run cmd/api/main.go

# Should see:
# Database connected successfully
# Database migrated successfully
# Server starting on port 8080
```

**Terminal 2: Frontend**
```bash
cd frontend
npm run dev

# Should see:
# - ready started server on 0.0.0.0:3000
# - Local: http://localhost:3000
```

**Terminal 3: Test API**
```bash
# Health check
curl http://localhost:8080/health

# Should return: {"status":"ok"}
```

### Step 6: Verify Setup

1. Open browser: http://localhost:3000
2. You should see Next.js default page
3. Navigate to http://localhost:3000/login
4. You should see login page (won't work yet, no users)

---

## Create First User (Seed Data)

**`backend/cmd/seed/main.go`**
```go
package main

import (
    "log"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "github.com/yourusername/koperasi-erp/internal/config"
    "github.com/yourusername/koperasi-erp/internal/database"
    "github.com/yourusername/koperasi-erp/internal/models"
)

func main() {
    cfg := config.Load()
    db := database.Connect(cfg.DatabaseURL)

    // Create test cooperative
    coop := models.Cooperative{
        ID:   uuid.New(),
        Name: "Koperasi Test",
        Address: "Jl. Test No. 123",
    }
    db.Create(&coop)

    // Create admin user
    password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
    user := models.User{
        ID:            uuid.New(),
        CooperativeID: coop.ID,
        Username:      "admin",
        PasswordHash:  string(password),
        FullName:      "Admin User",
        Role:          "admin",
        IsActive:      true,
    }
    db.Create(&user)

    log.Println("Seed data created successfully!")
    log.Println("Username: admin")
    log.Println("Password: admin123")
}
```

**Run seed:**
```bash
cd backend
go run cmd/seed/main.go
```

**Now you can login!**
- Username: admin
- Password: admin123

---

## Next Steps (First Week)

### Day 1: Authentication
- [ ] Implement JWT token generation
- [ ] Create login handler
- [ ] Create logout handler
- [ ] Test login flow

### Day 2: User Management
- [ ] Create user list endpoint
- [ ] Create user creation endpoint
- [ ] Build user management UI
- [ ] Test CRUD operations

### Day 3: Dashboard
- [ ] Create dashboard layout
- [ ] Add navigation
- [ ] Show basic stats
- [ ] Protected routes

### Day 4: Member Module Start
- [ ] Design member form
- [ ] Create member endpoints
- [ ] Member list UI
- [ ] Member creation flow

### Day 5: Testing & Cleanup
- [ ] Test all features
- [ ] Fix bugs
- [ ] Clean up code
- [ ] Document progress

---

## Useful Commands

### Backend (Go)
```bash
# Run server
go run cmd/api/main.go

# Run with auto-reload (install air)
go install github.com/cosmtrek/air@latest
air

# Run tests
go test ./...

# Format code
go fmt ./...

# Build binary
go build -o bin/api cmd/api/main.go
```

### Frontend (Next.js)
```bash
# Development
npm run dev

# Build
npm run build

# Start production
npm start

# Lint
npm run lint
```

### Database
```bash
# Start PostgreSQL
docker start koperasi-postgres

# Stop PostgreSQL
docker stop koperasi-postgres

# Connect to DB
docker exec -it koperasi-postgres psql -U postgres -d koperasi_erp

# Backup database
docker exec koperasi-postgres pg_dump -U postgres koperasi_erp > backup.sql

# Restore database
docker exec -i koperasi-postgres psql -U postgres koperasi_erp < backup.sql
```

---

## Troubleshooting

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Database Connection Failed
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check logs
docker logs koperasi-postgres

# Restart container
docker restart koperasi-postgres
```

### CORS Errors
- Make sure CORS middleware is configured in backend
- Check API URL in frontend .env.local
- Clear browser cache

---

## Project Structure After Setup

```
COOPERATIVE-ERP-LITE/
├── backend/
│   ├── cmd/
│   │   ├── api/main.go
│   │   └── seed/main.go
│   ├── internal/
│   │   ├── config/config.go
│   │   ├── database/database.go
│   │   ├── models/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── services/
│   ├── pkg/
│   ├── go.mod
│   ├── go.sum
│   └── .env
├── frontend/
│   ├── app/
│   │   ├── login/page.tsx
│   │   └── dashboard/page.tsx
│   ├── lib/
│   │   └── api.ts
│   ├── package.json
│   └── .env.local
├── docs/
│   ├── mvp-action-plan.md
│   ├── quick-start-guide.md (this file)
│   └── ...
└── README.md
```

---

## Resources

### Documentation
- Go: https://go.dev/doc/
- Gin: https://gin-gonic.com/docs/
- GORM: https://gorm.io/docs/
- Next.js: https://nextjs.org/docs
- PostgreSQL: https://www.postgresql.org/docs/

### Learning
- Go by Example: https://gobyexample.com/
- Next.js Tutorial: https://nextjs.org/learn
- GORM Guide: https://gorm.io/docs/index.html

---

## Getting Help

1. Check documentation first
2. Search Stack Overflow
3. Ask in team chat
4. Create GitHub issue

---

**You're all set! Start coding! 🚀**

**First commit:**
```bash
git init
git add .
git commit -m "Initial project setup"
git remote add origin <your-repo-url>
git push -u origin main
```
