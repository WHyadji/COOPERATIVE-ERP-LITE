# Nginx Configuration - Cooperative ERP Lite

## Overview

Nginx 1.28.0 (Stable) configured as reverse proxy with production-ready best practices including:
- Rate limiting
- Security headers
- SSL/TLS ready
- Performance optimizations
- Error handling
- Health checks

## Architecture

```
Client → Nginx (Port 80/443) → Backend API (Port 8080) → PostgreSQL
```

## Features

### Security
- ✅ Security headers (X-Frame-Options, X-XSS-Protection, etc.)
- ✅ Rate limiting (100 req/min for API, 5 req/min for login)
- ✅ CORS configuration
- ✅ Hidden Nginx version
- ✅ SSL/TLS ready (TLS 1.2+)

### Performance
- ✅ Gzip compression
- ✅ Connection keep-alive
- ✅ Upstream load balancing ready
- ✅ File descriptor caching
- ✅ Buffer optimizations

### Monitoring
- ✅ Custom access logs with timing
- ✅ Separate error logs
- ✅ Health check endpoint

## Directory Structure

```
nginx/
├── nginx.conf          # Main Nginx configuration
├── conf.d/
│   └── api.conf       # API reverse proxy configuration
├── ssl/               # SSL certificates (add your own)
└── logs/              # Nginx logs (auto-created)
```

## Quick Start

### 1. Start with Docker Compose

```bash
# Build and start all services
docker compose up -d

# Check Nginx status
docker compose ps nginx

# View Nginx logs
docker compose logs -f nginx
```

### 2. Test Endpoints

```bash
# Health check (no rate limit)
curl http://localhost/health

# API endpoint (with rate limiting)
curl http://localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"test","kataSandi":"test"}'

# Test rate limiting (run 10+ times quickly)
for i in {1..15}; do curl -s http://localhost/api/test; done
# Should get 429 (Too Many Requests) after burst limit
```

### 3. Check Configuration

```bash
# Test Nginx configuration
docker compose exec nginx nginx -t

# Reload configuration (without downtime)
docker compose exec nginx nginx -s reload
```

## Configuration Files

### nginx.conf
Main configuration with:
- Worker processes auto-tuning
- Connection limits
- Gzip compression
- Rate limiting zones
- Upstream backend pool

### conf.d/api.conf
Site-specific configuration with:
- HTTP/HTTPS server blocks
- Reverse proxy to backend
- Security headers
- Rate limiting rules
- Custom error pages

## Rate Limiting

| Endpoint | Rate | Burst | Description |
|----------|------|-------|-------------|
| `/health` | No limit | - | Health checks |
| `/api/*` | 100/min | 20 | General API |
| `/api/v1/auth/login` | 5/min | 3 | Login endpoint |

## SSL/TLS Setup

### Development (Self-signed)

```bash
# Generate self-signed certificate
mkdir -p nginx/ssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/privkey.pem \
  -out nginx/ssl/fullchain.pem \
  -subj "/CN=localhost"

# Uncomment SSL lines in nginx/conf.d/api.conf
# Restart Nginx
docker compose restart nginx
```

### Production (Let's Encrypt)

```bash
# On VPS with domain
sudo certbot --nginx -d api.yourdomain.com

# Certificates will be in /etc/letsencrypt/live/
# Update nginx/conf.d/api.conf with certificate paths
```

## Performance Tuning

### Worker Processes
```nginx
worker_processes auto;  # Auto-detect CPU cores
worker_rlimit_nofile 65535;  # Max open files
worker_connections 4096;  # Connections per worker
```

### Buffering
```nginx
client_body_buffer_size 128k;
client_max_body_size 10m;  # Max upload size
```

### Timeouts
```nginx
client_body_timeout 12;
client_header_timeout 12;
send_timeout 10;
```

## Monitoring

### Access Logs
```bash
# Real-time access logs
docker compose logs -f nginx | grep access

# Request timing
docker compose exec nginx tail -f /var/log/nginx/api-access.log
```

### Error Logs
```bash
# Real-time error logs
docker compose logs -f nginx | grep error

# Error details
docker compose exec nginx tail -f /var/log/nginx/api-error.log
```

### Metrics
Access log format includes:
- Response time (`rt`)
- Upstream connect time (`uct`)
- Upstream header time (`uht`)
- Upstream response time (`urt`)

## Troubleshooting

### Configuration Test Failed
```bash
# Check syntax
docker compose exec nginx nginx -t

# Check logs
docker compose logs nginx
```

### 502 Bad Gateway
```bash
# Check backend is running
docker compose ps backend

# Check backend logs
docker compose logs backend

# Check network connectivity
docker compose exec nginx ping backend
```

### Rate Limit Issues
```bash
# Check rate limit zones
docker compose exec nginx cat /etc/nginx/nginx.conf | grep limit_req_zone

# Adjust limits in nginx/nginx.conf
# Reload configuration
docker compose exec nginx nginx -s reload
```

### SSL Certificate Issues
```bash
# Verify certificate
openssl x509 -in nginx/ssl/fullchain.pem -text -noout

# Check certificate expiry
openssl x509 -in nginx/ssl/fullchain.pem -noout -dates
```

## Best Practices Applied

1. **Security**
   - Security headers on all responses
   - Rate limiting to prevent abuse
   - Hidden server version
   - SSL/TLS 1.2+ only

2. **Performance**
   - Gzip compression for text content
   - Connection keep-alive
   - File descriptor caching
   - Optimized buffer sizes

3. **Reliability**
   - Health checks
   - Auto-restart on failure
   - Proper error handling
   - Upstream failover ready

4. **Observability**
   - Detailed access logs with timing
   - Separate error logs
   - Request ID tracking
   - Health check monitoring

## Production Checklist

- [ ] Update `server_name` in api.conf
- [ ] Install valid SSL certificates
- [ ] Update rate limits for your needs
- [ ] Configure CORS allowed origins
- [ ] Set up log rotation
- [ ] Configure monitoring/alerting
- [ ] Test failover scenarios
- [ ] Document custom configurations

## References

- [Nginx 1.28.0 Release Notes](https://nginx.org/en/CHANGES-1.28)
- [Nginx Best Practices](https://www.nginx.com/blog/nginx-best-practices/)
- [Security Headers](https://securityheaders.com/)
- [SSL Configuration](https://ssl-config.mozilla.org/)

## Support

For issues or questions:
1. Check configuration with `nginx -t`
2. Review logs in `nginx/logs/`
3. Consult Nginx documentation
4. Contact team

---

**Version**: 1.0.0
**Last Updated**: 2025-11-17
**Nginx Version**: 1.28.0 (stable)
