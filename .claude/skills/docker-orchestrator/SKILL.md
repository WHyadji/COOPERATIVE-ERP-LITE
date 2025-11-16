---
name: Docker Orchestrator
description: Create optimized Dockerfiles, multi-stage builds, docker-compose configurations, and implement best practices for containerization. Includes security scanning, size optimization, layer caching strategies, and production-ready configurations.
---

# Docker Orchestrator

## Overview

This skill generates optimized Docker configurations for any application stack. It creates multi-stage Dockerfiles, docker-compose setups, implements caching strategies, security best practices, and size optimization techniques for production-ready containers.

## Quick Start

### Generate Dockerfile
```bash
python scripts/docker_orchestrator.py --app node --output Dockerfile

python scripts/docker_orchestrator.py --app python --framework django --optimize
```

### Generate Docker Compose
```bash
python scripts/docker_orchestrator.py --stack fullstack --services web,api,db,redis --output docker-compose.yml

python scripts/docker_orchestrator.py --compose microservices --scale 3
```

### Optimize Existing Docker Setup
```bash
python scripts/docker_orchestrator.py --analyze . --optimize

python scripts/docker_orchestrator.py --scan Dockerfile --fix-security
```

## Dockerfile Templates

### Node.js Multi-Stage Build
```dockerfile
# Stage 1: Dependencies
FROM node:18-alpine AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

# Stage 2: Build
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 3: Production
FROM node:18-alpine AS runner
WORKDIR /app

ENV NODE_ENV production

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Copy built application
COPY --from=builder --chown=nodejs:nodejs /app/dist ./dist
COPY --from=deps --chown=nodejs:nodejs /app/node_modules ./node_modules
COPY --from=builder --chown=nodejs:nodejs /app/package*.json ./

USER nodejs

EXPOSE 3000

CMD ["node", "dist/index.js"]
```

### Python Multi-Stage Build
```dockerfile
# Stage 1: Build dependencies
FROM python:3.11-slim AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
COPY requirements.txt .
RUN pip install --user --no-cache-dir -r requirements.txt

# Stage 2: Production
FROM python:3.11-slim AS runner

WORKDIR /app

# Create non-root user
RUN groupadd -r appuser && useradd -r -g appuser appuser

# Copy dependencies from builder
COPY --from=builder --chown=appuser:appuser /root/.local /home/appuser/.local

# Copy application code
COPY --chown=appuser:appuser . .

# Set PATH to include user's local bin
ENV PATH=/home/appuser/.local/bin:$PATH

USER appuser

EXPOSE 8000

CMD ["gunicorn", "--bind", "0.0.0.0:8000", "--workers", "4", "app:application"]
```

### Go Multi-Stage Build
```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Production (distroless)
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/main .

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/main"]
```

### Java Spring Boot Multi-Stage
```dockerfile
# Stage 1: Build with Maven
FROM maven:3.9-eclipse-temurin-17 AS builder

WORKDIR /app

# Cache dependencies
COPY pom.xml .
RUN mvn dependency:go-offline

# Build application
COPY src ./src
RUN mvn clean package -DskipTests

# Stage 2: Extract layers
FROM eclipse-temurin:17-jre AS extractor

WORKDIR /app

COPY --from=builder /app/target/*.jar app.jar
RUN java -Djarmode=layertools -jar app.jar extract

# Stage 3: Production
FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

# Create non-root user
RUN addgroup -S spring && adduser -S spring -G spring

# Copy layers in order of change frequency
COPY --from=extractor --chown=spring:spring /app/dependencies/ ./
COPY --from=extractor --chown=spring:spring /app/spring-boot-loader/ ./
COPY --from=extractor --chown=spring:spring /app/snapshot-dependencies/ ./
COPY --from=extractor --chown=spring:spring /app/application/ ./

USER spring

EXPOSE 8080

ENTRYPOINT ["java", "org.springframework.boot.loader.JarLauncher"]
```

### Frontend (React/Vue/Angular)
```dockerfile
# Stage 1: Build
FROM node:18-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY package*.json ./
RUN npm ci

# Build application
COPY . .
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

## Docker Compose Configurations

### Full Stack Application
```yaml
version: '3.9'

services:
  # Frontend
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      cache_from:
        - ${REGISTRY}/frontend:latest
    image: ${REGISTRY}/frontend:${VERSION:-latest}
    ports:
      - "80:80"
    environment:
      - API_URL=http://backend:3000
    depends_on:
      - backend
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Backend API
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
      target: production
    image: ${REGISTRY}/backend:${VERSION:-latest}
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - DATABASE_URL=postgres://user:pass@postgres:5432/appdb
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    networks:
      - app-network
    volumes:
      - ./backend/uploads:/app/uploads
    restart: unless-stopped
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M

  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=appdb
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_INITDB_ARGS=--encoding=UTF-8
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./database/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d appdb"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis Cache
  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis-data:/data
    ports:
      - "6379:6379"
    networks:
      - app-network
    restart: unless-stopped

  # Nginx Reverse Proxy
  nginx:
    image: nginx:alpine
    ports:
      - "443:443"
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
    depends_on:
      - frontend
      - backend
    networks:
      - app-network
    restart: unless-stopped

networks:
  app-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.28.0.0/16

volumes:
  postgres-data:
    driver: local
  redis-data:
    driver: local
```

### Microservices Stack
```yaml
version: '3.9'

x-common-variables: &common-variables
  RABBITMQ_URL: amqp://guest:guest@rabbitmq:5672
  REDIS_URL: redis://redis:6379
  JAEGER_AGENT_HOST: jaeger
  JAEGER_AGENT_PORT: 6831

x-default-service: &default-service
  restart: unless-stopped
  networks:
    - microservices
  logging:
    driver: "json-file"
    options:
      max-size: "10m"
      max-file: "3"

services:
  # API Gateway
  gateway:
    <<: *default-service
    build: ./gateway
    ports:
      - "8080:8080"
    environment:
      <<: *common-variables
      SERVICES_AUTH: http://auth-service:3001
      SERVICES_USER: http://user-service:3002
      SERVICES_PRODUCT: http://product-service:3003
      SERVICES_ORDER: http://order-service:3004
    depends_on:
      - auth-service
      - user-service
      - product-service
      - order-service

  # Auth Service
  auth-service:
    <<: *default-service
    build: ./services/auth
    environment:
      <<: *common-variables
      PORT: 3001
      DB_HOST: auth-db
      SECRET_KEY: ${AUTH_SECRET}
    depends_on:
      - auth-db
      - rabbitmq
    deploy:
      replicas: 2

  # User Service
  user-service:
    <<: *default-service
    build: ./services/user
    environment:
      <<: *common-variables
      PORT: 3002
      DB_HOST: user-db
    depends_on:
      - user-db
      - rabbitmq
    deploy:
      replicas: 2

  # Product Service
  product-service:
    <<: *default-service
    build: ./services/product
    environment:
      <<: *common-variables
      PORT: 3003
      MONGO_URL: mongodb://product-db:27017/products
    depends_on:
      - product-db
      - rabbitmq

  # Order Service
  order-service:
    <<: *default-service
    build: ./services/order
    environment:
      <<: *common-variables
      PORT: 3004
      DB_HOST: order-db
    depends_on:
      - order-db
      - rabbitmq

  # Databases
  auth-db:
    <<: *default-service
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: auth
      POSTGRES_USER: auth
      POSTGRES_PASSWORD: ${AUTH_DB_PASSWORD}
    volumes:
      - auth-db-data:/var/lib/postgresql/data

  user-db:
    <<: *default-service
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: users
      POSTGRES_USER: user
      POSTGRES_PASSWORD: ${USER_DB_PASSWORD}
    volumes:
      - user-db-data:/var/lib/postgresql/data

  product-db:
    <<: *default-service
    image: mongo:6
    environment:
      MONGO_INITDB_DATABASE: products
    volumes:
      - product-db-data:/data/db

  order-db:
    <<: *default-service
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: orders
      POSTGRES_USER: order
      POSTGRES_PASSWORD: ${ORDER_DB_PASSWORD}
    volumes:
      - order-db-data:/var/lib/postgresql/data

  # Message Queue
  rabbitmq:
    <<: *default-service
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq-data:/var/lib/rabbitmq

  # Cache
  redis:
    <<: *default-service
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data

  # Monitoring
  jaeger:
    <<: *default-service
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"
      - "6831:6831/udp"

networks:
  microservices:
    driver: bridge

volumes:
  auth-db-data:
  user-db-data:
  product-db-data:
  order-db-data:
  rabbitmq-data:
  redis-data:
```

## Optimization Strategies

### Layer Caching Optimization
```dockerfile
# ❌ BAD: Breaks cache frequently
FROM node:18
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

# ✅ GOOD: Maximizes cache reuse
FROM node:18
WORKDIR /app

# Dependencies layer (changes rarely)
COPY package*.json ./
RUN npm ci --only=production

# Source code layer (changes frequently)
COPY . .
RUN npm run build
```

### Size Optimization Techniques
```dockerfile
# 1. Use Alpine or Distroless base images
FROM node:18-alpine  # 178MB vs 900MB for node:18

# 2. Multi-stage builds
FROM node:18 AS builder
# Build steps...
FROM node:18-alpine
COPY --from=builder /app/dist ./dist

# 3. Minimize layers
# ❌ BAD: Multiple RUN commands
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y git

# ✅ GOOD: Single RUN command
RUN apt-get update && apt-get install -y \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

# 4. Remove unnecessary files
RUN npm ci --only=production && \
    npm cache clean --force && \
    rm -rf /tmp/* /var/tmp/*

# 5. Use .dockerignore
# .dockerignore file:
node_modules
npm-debug.log
.git
.gitignore
README.md
.env
.vscode
.idea
coverage
.nyc_output
```

### Security Best Practices
```dockerfile
# 1. Don't run as root
RUN groupadd -r appuser && useradd -r -g appuser appuser
USER appuser

# 2. Use specific versions
# ❌ BAD
FROM node:latest

# ✅ GOOD
FROM node:18.17.1-alpine3.18

# 3. Scan for vulnerabilities
# Use tools like Trivy, Snyk, or Docker Scout
RUN apk add --no-cache \
    ca-certificates \
    && update-ca-certificates

# 4. Use secrets properly
# ❌ BAD
ENV API_KEY=secret123

# ✅ GOOD - Use BuildKit secrets
# docker build --secret id=api_key,src=api_key.txt .
RUN --mount=type=secret,id=api_key \
    API_KEY=$(cat /run/secrets/api_key) && \
    echo "Using secret for build"

# 5. Set security options
# In docker-compose.yml:
security_opt:
  - no-new-privileges:true
  - apparmor:docker-default
cap_drop:
  - ALL
cap_add:
  - NET_BIND_SERVICE
read_only: true
```

### Build Time Optimization
```dockerfile
# 1. Use BuildKit
# DOCKER_BUILDKIT=1 docker build .

# 2. Cache mount for package managers
# syntax=docker/dockerfile:1
FROM node:18
RUN --mount=type=cache,target=/root/.npm \
    npm ci --only=production

# 3. Parallel builds
# syntax=docker/dockerfile:1
FROM node:18 AS frontend-builder
WORKDIR /frontend
COPY frontend/ .
RUN npm ci && npm run build

FROM python:3.11 AS backend-builder
WORKDIR /backend
COPY backend/ .
RUN pip install -r requirements.txt

FROM nginx:alpine AS production
COPY --from=frontend-builder /frontend/dist /usr/share/nginx/html
COPY --from=backend-builder /backend /app
```

## Health Checks

### Dockerfile Health Check
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD node healthcheck.js || exit 1

# Or for web services
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1
```

### Docker Compose Health Check
```yaml
healthcheck:
  test: ["CMD-SHELL", "curl -f http://localhost:3000/health || exit 1"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

## Container Orchestration

### Docker Swarm Mode
```yaml
version: '3.9'

services:
  web:
    image: myapp:latest
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      placement:
        constraints:
          - node.role == worker
          - node.labels.type == web
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
    networks:
      - webnet
    secrets:
      - source: api_key
        target: /run/secrets/api_key
        mode: 0400

networks:
  webnet:
    driver: overlay
    attachable: true

secrets:
  api_key:
    external: true
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: app
        image: myapp:latest
        ports:
        - containerPort: 3000
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
        env:
        - name: NODE_ENV
          value: "production"
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
```

## Development Environment

### docker-compose.dev.yml
```yaml
version: '3.9'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
      - /app/node_modules
    environment:
      - NODE_ENV=development
    command: npm run dev
    ports:
      - "3000:3000"
      - "9229:9229"  # Node debugger

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=devdb
      - POSTGRES_USER=developer
      - POSTGRES_PASSWORD=devpass
    ports:
      - "5432:5432"
    volumes:
      - ./data/postgres:/var/lib/postgresql/data

  mailhog:
    image: mailhog/mailhog
    ports:
      - "1025:1025"
      - "8025:8025"
```

### Dockerfile.dev
```dockerfile
FROM node:18

WORKDIR /app

# Install development tools
RUN apt-get update && apt-get install -y \
    vim \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

# Install global packages
RUN npm install -g nodemon tsx

# Copy package files
COPY package*.json ./
RUN npm install

# Expose debugging port
EXPOSE 9229

# Run with hot reload
CMD ["nodemon", "--inspect=0.0.0.0:9229", "src/index.js"]
```

## CI/CD Integration

### GitHub Actions
```yaml
name: Docker Build and Push

on:
  push:
    branches: [main]
    tags: ['v*']

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Log in to Registry
        uses: docker/login-action@v2
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v4
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha
      
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

## Docker Commands Reference

### Build Commands
```bash
# Build with BuildKit
DOCKER_BUILDKIT=1 docker build -t app:latest .

# Multi-platform build
docker buildx build --platform linux/amd64,linux/arm64 -t app:latest .

# Build with cache
docker build --cache-from app:latest -t app:latest .

# Build specific stage
docker build --target production -t app:prod .
```

### Compose Commands
```bash
# Start services
docker-compose up -d

# Scale service
docker-compose up -d --scale web=3

# View logs
docker-compose logs -f web

# Execute command
docker-compose exec web bash

# Rebuild and restart
docker-compose up -d --build

# Stop and remove
docker-compose down -v
```

### Optimization Commands
```bash
# Analyze image layers
docker history app:latest

# Check image size
docker images app:latest

# Scan for vulnerabilities
docker scout cves app:latest

# Prune unused resources
docker system prune -a --volumes
```

## Best Practices Summary

1. **Use Multi-Stage Builds** - Reduce final image size
2. **Leverage Build Cache** - Order commands by change frequency
3. **Run as Non-Root** - Security best practice
4. **Use Specific Tags** - Avoid :latest in production
5. **Minimize Layers** - Combine RUN commands
6. **Use .dockerignore** - Exclude unnecessary files
7. **Health Checks** - Monitor container health
8. **Resource Limits** - Prevent resource exhaustion
9. **Secret Management** - Never hardcode secrets
10. **Regular Updates** - Keep base images updated