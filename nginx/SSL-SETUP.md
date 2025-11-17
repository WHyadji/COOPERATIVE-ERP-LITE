# SSL/TLS Certificate Setup Guide

## Overview

Complete guide for setting up SSL/TLS certificates for Cooperative ERP Lite with best practices for both development and production environments.

---

## Development Setup (Self-Signed Certificate)

### ✅ Already Configured!

Your development environment is already configured with:
- Self-signed SSL certificate (valid for 1 year)
- TLS 1.2 and TLS 1.3 protocols
- Modern cipher suites
- HTTP/2 enabled
- Perfect Forward Secrecy (DH parameters)

### Certificate Details

```bash
# View certificate information
openssl x509 -in nginx/ssl/fullchain.pem -noout -text

# Check expiration date
openssl x509 -in nginx/ssl/fullchain.pem -noout -dates

# Verify private key matches certificate
openssl x509 -noout -modulus -in nginx/ssl/fullchain.pem | openssl md5
openssl rsa -noout -modulus -in nginx/ssl/privkey.pem | openssl md5
```

### Testing HTTPS

```bash
# Test HTTPS endpoint
curl -k https://localhost/health

# Test HTTP to HTTPS redirect
curl -I http://localhost/api/v1/test

# Check SSL/TLS configuration
openssl s_client -connect localhost:443 -tls1_2 -servername localhost

# Test HTTP/2 support
curl -k -I https://localhost/health --http2
```

---

## Production Setup (Let's Encrypt)

### Prerequisites

1. **Domain Name**: Registered and pointing to your VPS
   - Example: `api.yourdomain.com`

2. **DNS A Record**: Set up to point to your server IP
   ```
   api.yourdomain.com  →  123.456.789.0 (your VPS IP)
   ```

3. **Ports Open**: Ensure firewall allows ports 80 and 443
   ```bash
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   ```

4. **Server Access**: SSH access to your VPS

### Method 1: Certbot with Nginx (Recommended)

#### Step 1: Install Certbot

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install certbot python3-certbot-nginx

# CentOS/RHEL
sudo yum install certbot python3-certbot-nginx
```

#### Step 2: Update Nginx Configuration

Before running Certbot, ensure your `nginx/conf.d/api.conf` has the correct `server_name`:

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name api.yourdomain.com;  # Change this!

    location / {
        return 301 https://$server_name$request_uri;
    }
}
```

#### Step 3: Stop Docker Nginx (Temporary)

```bash
docker compose stop nginx
```

#### Step 4: Install System Nginx (Temporary)

```bash
sudo apt install nginx
sudo systemctl start nginx
```

#### Step 5: Get Certificate

```bash
# Obtain certificate
sudo certbot --nginx -d api.yourdomain.com

# Follow prompts:
# - Enter email address
# - Agree to Terms of Service
# - Choose redirect HTTP to HTTPS (recommended)
```

#### Step 6: Copy Certificates to Project

```bash
# Copy certificates
sudo cp /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/api.yourdomain.com/privkey.pem nginx/ssl/
sudo chown $USER:$USER nginx/ssl/*.pem
sudo chmod 644 nginx/ssl/fullchain.pem
sudo chmod 600 nginx/ssl/privkey.pem
```

#### Step 7: Stop System Nginx, Start Docker

```bash
# Stop system nginx
sudo systemctl stop nginx
sudo systemctl disable nginx

# Start Docker nginx
docker compose start nginx
```

#### Step 8: Set Up Auto-Renewal

```bash
# Test renewal
sudo certbot renew --dry-run

# Set up cron job for auto-renewal
sudo crontab -e

# Add this line (runs twice daily):
0 */12 * * * certbot renew --quiet --deploy-hook "cp /etc/letsencrypt/live/api.yourdomain.com/*.pem /path/to/nginx/ssl/ && docker compose restart nginx"
```

### Method 2: Certbot Standalone Mode

#### Step 1: Stop Docker Services

```bash
docker compose down
```

#### Step 2: Get Certificate

```bash
sudo certbot certonly --standalone -d api.yourdomain.com
```

#### Step 3: Copy Certificates

```bash
sudo cp /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/api.yourdomain.com/privkey.pem nginx/ssl/
sudo chown $USER:$USER nginx/ssl/*.pem
```

#### Step 4: Update docker-compose.yml

Add certificate renewal volume:

```yaml
nginx:
  volumes:
    - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    - ./nginx/conf.d:/etc/nginx/conf.d:ro
    - ./nginx/ssl:/etc/nginx/ssl:ro
    - /etc/letsencrypt:/etc/letsencrypt:ro  # Add this
```

#### Step 5: Start Services

```bash
docker compose up -d
```

#### Step 6: Set Up Auto-Renewal

```bash
# Create renewal script
sudo nano /usr/local/bin/renew-cert.sh
```

Add:
```bash
#!/bin/bash
docker compose -f /path/to/docker-compose.yml down
certbot renew
docker compose -f /path/to/docker-compose.yml up -d
```

Make executable and add to cron:
```bash
sudo chmod +x /usr/local/bin/renew-cert.sh
sudo crontab -e
# Add:
0 0 * * * /usr/local/bin/renew-cert.sh >> /var/log/cert-renewal.log 2>&1
```

### Method 3: Manual DNS Challenge (For Wildcard)

```bash
# Get wildcard certificate
sudo certbot certonly --manual --preferred-challenges dns \
  -d "*.yourdomain.com" -d "yourdomain.com"

# Follow prompts to add TXT records to DNS
# Then copy certificates as in Method 2
```

---

## SSL/TLS Best Practices Configuration

### Current Configuration Includes:

1. **Protocols**: TLS 1.2 and 1.3 only (secure)
2. **Ciphers**: Modern cipher suite (Mozilla Modern)
3. **Perfect Forward Secrecy**: DH parameters (2048-bit)
4. **HTTP/2**: Enabled for performance
5. **HSTS**: Strict-Transport-Security header
6. **Session Management**: Optimized caching and no tickets

### SSL Test Tools

```bash
# Test SSL configuration
openssl s_client -connect yourdomain.com:443 -tls1_2

# Check cipher strength
nmap --script ssl-enum-ciphers -p 443 yourdomain.com

# Online SSL test
# Visit: https://www.ssllabs.com/ssltest/
```

### Security Headers Verification

```bash
# Check security headers
curl -I https://yourdomain.com/health

# Should include:
# - Strict-Transport-Security
# - X-Frame-Options: DENY
# - X-Content-Type-Options: nosniff
# - Content-Security-Policy
# - Referrer-Policy
```

---

## Certificate Renewal

### Let's Encrypt Certificates

- **Validity**: 90 days
- **Renewal Window**: 30 days before expiration
- **Auto-Renewal**: Set up via cron (recommended)

### Manual Renewal

```bash
# Check expiration
sudo certbot certificates

# Renew manually
sudo certbot renew

# Force renewal (testing)
sudo certbot renew --force-renewal
```

### Automatic Renewal Script

Create `/usr/local/bin/ssl-renew.sh`:

```bash
#!/bin/bash

# Stop services
docker compose -f /path/to/docker-compose.yml stop nginx

# Renew certificate
certbot renew --quiet

# Copy new certificates
cp /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem /path/to/nginx/ssl/
cp /etc/letsencrypt/live/api.yourdomain.com/privkey.pem /path/to/nginx/ssl/

# Restart services
docker compose -f /path/to/docker-compose.yml start nginx

# Log renewal
echo "Certificate renewed on $(date)" >> /var/log/ssl-renewal.log
```

Make executable:
```bash
chmod +x /usr/local/bin/ssl-renew.sh
```

Add to crontab:
```bash
0 0 1 * * /usr/local/bin/ssl-renew.sh
```

---

## Troubleshooting

### Certificate Not Found

```bash
# Check certificate files
ls -la nginx/ssl/

# Verify permissions
ls -l nginx/ssl/*.pem

# Should be:
# -rw-r--r-- fullchain.pem
# -rw------- privkey.pem
```

### Nginx Won't Start

```bash
# Test configuration
docker compose exec nginx nginx -t

# Check logs
docker compose logs nginx

# Common issues:
# - Certificate path wrong
# - Permission denied
# - Certificate expired
```

### Certificate Expired

```bash
# Check expiration
openssl x509 -in nginx/ssl/fullchain.pem -noout -dates

# Renew
sudo certbot renew --force-renewal

# Update files
sudo cp /etc/letsencrypt/live/api.yourdomain.com/*.pem nginx/ssl/
docker compose restart nginx
```

### Browser Shows "Not Secure"

1. **Self-signed certificate**: Browser warning is normal
   - Click "Advanced" → "Proceed to site"
   - Or add exception in browser

2. **Production certificate**: Check these:
   ```bash
   # Verify certificate chain
   openssl s_client -connect yourdomain.com:443 -showcerts

   # Check intermediate certificates
   curl -I https://yourdomain.com
   ```

### SSL Labs Grade

Target: A+ rating

```bash
# Test at: https://www.ssllabs.com/ssltest/

# Current configuration should achieve:
# - A+ rating
# - 100% for Key Exchange
# - 100% for Cipher Strength
# - HSTS enabled
```

---

## Security Checklist

Production SSL/TLS checklist:

- [ ] Valid SSL certificate installed
- [ ] Certificate auto-renewal configured
- [ ] Only TLS 1.2 and 1.3 enabled
- [ ] Strong ciphers configured
- [ ] HSTS header enabled
- [ ] HTTP redirects to HTTPS
- [ ] Certificate chain complete
- [ ] Private key permissions secure (600)
- [ ] OCSP stapling enabled (production)
- [ ] Certificate monitoring set up

---

## Monitoring

### Certificate Expiry Monitoring

```bash
# Check days until expiration
openssl x509 -enddate -noout -in nginx/ssl/fullchain.pem

# Create monitoring script
cat > /usr/local/bin/check-cert-expiry.sh << 'EOF'
#!/bin/bash
CERT_FILE="/path/to/nginx/ssl/fullchain.pem"
DAYS_WARN=30

EXPIRY=$(openssl x509 -enddate -noout -in $CERT_FILE | cut -d= -f2)
EXPIRY_EPOCH=$(date -d "$EXPIRY" +%s)
NOW_EPOCH=$(date +%s)
DAYS_LEFT=$(( ($EXPIRY_EPOCH - $NOW_EPOCH) / 86400 ))

if [ $DAYS_LEFT -lt $DAYS_WARN ]; then
    echo "WARNING: Certificate expires in $DAYS_LEFT days!"
    # Send alert (email, slack, etc.)
fi
EOF

chmod +x /usr/local/bin/check-cert-expiry.sh

# Add to daily cron
echo "0 8 * * * /usr/local/bin/check-cert-expiry.sh" | crontab -
```

---

## References

- [Mozilla SSL Configuration Generator](https://ssl-config.mozilla.org/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [SSL Labs Server Test](https://www.ssllabs.com/ssltest/)
- [Certbot Documentation](https://certbot.eff.org/docs/)
- [Nginx SSL Module](https://nginx.org/en/docs/http/ngx_http_ssl_module.html)

---

**Last Updated**: 2025-11-17
**Nginx Version**: 1.28.0
**TLS Version**: 1.2, 1.3
