# Critical Fix for E2E Tests

## Problem
E2E tests failing because:
1. Database is empty (no koperasi, no anggota)
2. Frontend uses invalid UUID for idKoperasi ("default-koperasi-id")
3. Backend expects valid UUID query parameter

## Solution

### Step 1: Seed Database
Run this SQL to create test data:

```sql
-- Create test koperasi
INSERT INTO koperasi (id, nama_koperasi, created_at, updated_at)
VALUES ('11111111-1111-1111-1111-111111111111', 'Koperasi Test E2E', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Create test member (A001 with PIN 123456)
-- bcrypt hash of "123456"
INSERT INTO anggota (
    id,
    id_koperasi,
    nomor_anggota,
    nama_lengkap,
    tanggal_bergabung,
    status,
    pin_portal,
    created_at,
    updated_at
)
VALUES (
    '22222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    'A001',
    'Test Member Portal',
    '2024-01-01',
    'aktif',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye92cC12dkXQOaX8y2lUjJwXTkLKJ.CSC',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO NOTHING;
```

### Step 2: Update docker-compose.yml
Change koperasi ID in frontend environment:

```yaml
frontend:
  environment:
    - NEXT_PUBLIC_DEFAULT_KOPERASI_ID=11111111-1111-1111-1111-111111111111
```

### Step 3: Rebuild Frontend
```bash
docker compose up -d --build frontend
```

### Step 4: Run E2E Tests
```bash
cd frontend && npx playwright test --project=chromium --reporter=list
```

## Test Credentials
- Nomor Anggota: A001
- PIN: 123456
- Koperasi ID: 11111111-1111-1111-1111-111111111111
