# E2E Test Data Seeding

This directory contains scripts to seed test data for E2E (end-to-end) testing of the Member Portal.

## Overview

The seeding scripts create:
- **1 Test Cooperative**: "Koperasi Test E2E" (TEST-E2E-001)
- **1 Test Member**: A001 with PIN 123456
- **1 Balance Record**: Initial balance with Simpanan Pokok, Wajib, and Sukarela
- **8 Sample Transactions**: Mix of Pokok, Wajib, and Sukarela deposits

## Test Credentials

After running the seeding script, you can use these credentials for testing:

```
Nomor Anggota: A001
PIN: 123456
```

## Method 1: Go Seeding Script (Recommended)

The Go script properly handles PIN hashing and is the recommended method.

### Prerequisites

- PostgreSQL running on `localhost:5432`
- Database `koperasi_erp` created
- Database credentials: `postgres:postgres`

### Steps

```bash
# From the backend directory
cd backend

# Run the seeding script
go run cmd/seed-test-data/main.go
```

### Expected Output

```
===========================================
E2E Test Data Seeding Script
===========================================
✓ Connected to database

1. Creating test cooperative...
   ✓ Koperasi ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

2. Creating test member...
   ✓ Member: Test Member Portal (A001)
   ✓ PIN: 123456 (hashed)

3. Creating initial balance...
   ✓ Simpanan Pokok: Rp 1,000,000
   ✓ Simpanan Wajib: Rp 2,500,000
   ✓ Simpanan Sukarela: Rp 500,000

4. Creating sample transactions...
   ✓ Created 8 transactions

===========================================
✓ Test data seeding completed successfully!
===========================================

Test Credentials:
  Nomor Anggota: A001
  PIN: 123456

You can now run E2E tests with:
  cd frontend && npx playwright test
===========================================
```

### Idempotency

The Go script is idempotent - you can run it multiple times safely. It will:
- Reuse existing test cooperative if found
- Reuse existing test member if found
- Reuse existing balance and transactions if found

## Method 2: SQL Script (Alternative)

⚠️ **Warning**: The SQL script uses a placeholder PIN hash. Use the Go script for proper PIN hashing.

### Steps

```bash
# From the backend directory
cd backend/cmd/seed-test-data

# Run the SQL script
psql -U postgres -d koperasi_erp -f seed_e2e_data.sql
```

### After Running SQL Script

You'll need to manually update the PIN hash:

```sql
-- Connect to database
psql -U postgres -d koperasi_erp

-- Update PIN with properly hashed value (generated from Go bcrypt)
UPDATE anggota
SET pin = '$2a$10$<actual_bcrypt_hash_here>'
WHERE nomor_anggota = 'A001';
```

Or better yet, just use the Go seeding script instead.

## Verifying Test Data

After seeding, verify the data was created correctly:

```bash
# Connect to database
psql -U postgres -d koperasi_erp

# Check cooperative
SELECT * FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001';

# Check member
SELECT nomor_anggota, nama_lengkap, status FROM anggota WHERE nomor_anggota = 'A001';

# Check balance
SELECT * FROM saldo_simpanan_anggota
WHERE id_anggota = (SELECT id FROM anggota WHERE nomor_anggota = 'A001');

# Check transactions
SELECT nomor_referensi, tipe_simpanan, jumlah, tanggal_transaksi
FROM simpanan
WHERE id_anggota = (SELECT id FROM anggota WHERE nomor_anggota = 'A001')
ORDER BY tanggal_transaksi;
```

## Running E2E Tests

After seeding the test data:

```bash
# From the frontend directory
cd ../../frontend

# Install Playwright browsers (first time only)
npx playwright install

# Run all E2E tests
npx playwright test

# Run tests in UI mode (interactive)
npx playwright test --ui

# Run specific test file
npx playwright test e2e/member-portal.spec.ts

# Run tests in specific browser
npx playwright test --project=chromium

# Show test report
npx playwright show-report
```

## Cleaning Up Test Data

To remove all test data:

```sql
-- Connect to database
psql -U postgres -d koperasi_erp

-- Delete in correct order (respecting foreign keys)
DELETE FROM simpanan
WHERE id_koperasi IN (SELECT id FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001');

DELETE FROM saldo_simpanan_anggota
WHERE id_koperasi IN (SELECT id FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001');

DELETE FROM anggota
WHERE id_koperasi IN (SELECT id FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001');

DELETE FROM koperasi
WHERE nomor_badan_hukum = 'TEST-E2E-001';
```

## Troubleshooting

### "Failed to connect to database"

Check that PostgreSQL is running:
```bash
# Start PostgreSQL with Docker
docker start koperasi-postgres

# Or check status
docker ps | grep postgres
```

### "Database does not exist"

Create the database:
```bash
docker exec -it koperasi-postgres psql -U postgres -c "CREATE DATABASE koperasi_erp;"
```

### "PIN authentication failed" in E2E tests

The PIN hash might not be correct. Always use the Go seeding script for proper PIN hashing.

### E2E tests timeout

Make sure the backend API server is running:
```bash
cd backend
go run cmd/api/main.go
```

And the frontend dev server (Playwright will auto-start, but you can manually start it):
```bash
cd frontend
npm run dev
```

## Test Data Details

### Cooperative
- **Name**: Koperasi Test E2E
- **Number**: TEST-E2E-001
- **Founded**: 2024-01-01

### Member
- **Number**: A001
- **Name**: Test Member Portal
- **NIK**: 1234567890123456
- **Gender**: L (Laki-laki)
- **PIN**: 123456
- **Status**: Aktif
- **Joined**: 2024-01-01

### Balance
- **Simpanan Pokok**: Rp 1,000,000
- **Simpanan Wajib**: Rp 2,500,000
- **Simpanan Sukarela**: Rp 500,000
- **Total**: Rp 4,000,000

### Transactions (8 total)
1. **SP-2024-001**: Simpanan Pokok - Rp 1,000,000 (2024-01-01)
2. **SW-2024-001**: Simpanan Wajib - Rp 500,000 (2024-01-15)
3. **SW-2024-002**: Simpanan Wajib - Rp 500,000 (2024-02-15)
4. **SW-2024-003**: Simpanan Wajib - Rp 500,000 (2024-03-15)
5. **SW-2024-004**: Simpanan Wajib - Rp 500,000 (2024-04-15)
6. **SW-2024-005**: Simpanan Wajib - Rp 500,000 (2024-05-15)
7. **SS-2024-001**: Simpanan Sukarela - Rp 200,000 (2024-02-01)
8. **SS-2024-002**: Simpanan Sukarela - Rp 300,000 (2024-03-20)

## Integration with CI/CD

For automated testing in CI/CD pipelines:

```bash
# Setup script example
#!/bin/bash

# 1. Start PostgreSQL
docker run -d --name test-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=koperasi_erp \
  -p 5432:5432 \
  postgres:15

# 2. Wait for PostgreSQL to be ready
sleep 5

# 3. Run migrations (if you have them)
# go run cmd/migrate/main.go up

# 4. Seed test data
go run cmd/seed-test-data/main.go

# 5. Run E2E tests
cd frontend && npx playwright test

# 6. Cleanup
docker stop test-postgres
docker rm test-postgres
```

## Notes

- The seeding script is safe to run multiple times (idempotent)
- Test data uses realistic Indonesian cooperative member data
- Transaction dates span multiple months for filtering tests
- All monetary amounts use Indonesian Rupiah (IDR)
- The test cooperative ID is automatically generated (UUID)
