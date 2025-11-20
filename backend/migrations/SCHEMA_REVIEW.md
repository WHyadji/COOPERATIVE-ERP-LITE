# Cooperative ERP Lite - Database Schema Review

**Review Date:** 2025-11-20
**Database:** PostgreSQL 15
**Environment:** Development (Pre-Production)
**Reviewer:** Schema Analysis Tool

---

## Executive Summary

### 🎯 Overall Status: **NEEDS ATTENTION**

- ✅ **Good:** Well-structured relational design with proper foreign keys
- ✅ **Good:** Soft delete pattern implemented across all tables
- ✅ **Good:** Proper indexing on foreign keys and query patterns
- ⚠️ **Critical:** **5 multi-tenant unique constraint violations found**
- ⚠️ **Warning:** Some indexes could be optimized for query performance
- ⚠️ **Info:** Missing some composite indexes for common query patterns

---

## 📊 Schema Overview

### Tables Summary

| Table | Rows | Size | Purpose | Multi-Tenant |
|-------|------|------|---------|--------------|
| koperasi | - | 48 KB | Cooperative organizations (tenant root) | N/A |
| anggota | 1 | 128 KB | Cooperative members | ✅ Fixed |
| simpanan | 8 | 184 KB | Member savings/deposits | ✅ Good |
| akun | - | 96 KB | Chart of Accounts | ❌ **ISSUE** |
| transaksi | - | 144 KB | Accounting journal entries | ❌ **ISSUE** |
| baris_transaksi | - | 96 KB | Journal entry line items | ✅ Good |
| produk | - | 40 KB | Product catalog (POS) | ❌ **ISSUE** |
| penjualan | - | 64 KB | Sales transactions | ❌ **ISSUE** |
| item_penjualan | - | 32 KB | Sales line items | ✅ Good |
| pengguna | - | 80 KB | System users | ❌ **ISSUE** |

**Total Database Size:** ~912 KB (development, minimal data)

---

## 🚨 Critical Issues Found

### Multi-Tenant Unique Constraint Violations

The following tables have unique constraints that are **NOT scoped by `id_koperasi`**, which breaks multi-tenant functionality:

| # | Table | Current Constraint | Column(s) | Impact |
|---|-------|-------------------|-----------|--------|
| 1 | ✅ anggota | `idx_koperasi_nomor` | `(id_koperasi, nomor_anggota)` | **FIXED** |
| 2 | ❌ akun | `idx_koperasi_kode_akun` | `(kode_akun)` | Different cooperatives cannot use same account codes |
| 3 | ❌ transaksi | `idx_koperasi_nomor_jurnal` | `(nomor_jurnal)` | Different cooperatives cannot use same journal numbers |
| 4 | ❌ produk | `idx_koperasi_kode_produk` | `(kode_produk)` | Different cooperatives cannot use same product codes |
| 5 | ❌ penjualan | `idx_koperasi_nomor_penjualan` | `(nomor_penjualan)` | Different cooperatives cannot use same sales numbers |
| 6 | ❌ pengguna | `idx_koperasi_username` | `(nama_pengguna)` | Different cooperatives cannot use same usernames |

**Risk Level:** 🔴 **CRITICAL**

**Business Impact:**
- Prevents multiple cooperatives from using the same account codes (e.g., "1010" for Cash)
- Prevents cooperatives from having independent sequential numbering
- Breaks fundamental multi-tenant isolation
- Will cause errors when onboarding second cooperative

---

## 📐 Entity Relationship Diagram

```
┌─────────────┐
│  KOPERASI   │ (Tenant Root)
│ (id)        │
└──────┬──────┘
       │
       │ 1:N relationships to all tenant tables
       │
       ├─────────────────┬─────────────────┬─────────────────┬─────────────────┐
       │                 │                 │                 │                 │
       ▼                 ▼                 ▼                 ▼                 ▼
┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│   ANGGOTA   │   │    AKUN     │   │  TRANSAKSI  │   │   PRODUK    │   │  PENGGUNA   │
│ (Members)   │   │ (Accounts)  │   │ (Journals)  │   │ (Products)  │   │   (Users)   │
└──────┬──────┘   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘   └──────┬──────┘
       │                 │                 │                 │                 │
       │ 1:N             │ 1:N             │ 1:N             │ 1:N             │ 1:N
       ▼                 ▼                 │                 │                 ▼
┌─────────────┐   ┌─────────────┐        │                 │           ┌─────────────┐
│  SIMPANAN   │   │   BARIS     │◄───────┘                 │           │  PENJUALAN  │
│ (Deposits)  │   │  TRANSAKSI  │                          │           │   (Sales)   │
└─────────────┘   │(Lines)      │                          │           └──────┬──────┘
                  └─────────────┘                          │                  │
                                                           │                  │ 1:N
                  ┌──────────────────────────┐            │                  ▼
                  │  Self-referential (tree) │            │           ┌─────────────┐
                  │  for sub-accounts        │            │           │    ITEM     │
                  └──────────────────────────┘            └───────────┤  PENJUALAN  │
                                                                      │ (Line Items)│
                                                                      └─────────────┘

Legend:
┌─────┐
│Table│  Primary entity
└─────┘
   │
   ▼    Relationship direction
1:N     One-to-Many relationship
```

### Key Relationships

**Koperasi (Parent Entity - Multi-Tenant Root):**
- Has many: Anggota, Akun, Transaksi, Produk, Penjualan, Pengguna, Simpanan
- All queries MUST filter by `id_koperasi` for data isolation

**Anggota (Members):**
- Belongs to: Koperasi
- Has many: Simpanan, Penjualan (optional, for member purchases)

**Akun (Chart of Accounts):**
- Belongs to: Koperasi
- Self-referential: Sub-accounts can have parent accounts
- Has many: Baris Transaksi (journal entry lines)

**Transaksi (Journal Entries):**
- Belongs to: Koperasi
- Has many: Baris Transaksi (lines)
- Referenced by: Simpanan, Penjualan (for accounting integration)

**Produk (Products):**
- Belongs to: Koperasi
- Referenced by: Item Penjualan

**Penjualan (Sales):**
- Belongs to: Koperasi
- References: Anggota (optional), Pengguna (kasir), Transaksi
- Has many: Item Penjualan

---

## 🔍 Index Analysis

### Well-Indexed Patterns ✅

1. **Foreign Key Indexes:** All foreign keys have corresponding indexes
2. **Soft Delete Indexes:** All `tanggal_dihapus` columns are indexed
3. **Tenant Scoping:** Most queries can use `(id_koperasi, *)` composite indexes
4. **Date Range Queries:** Transaction dates are indexed

### Missing/Suboptimal Indexes ⚠️

| Table | Suggested Index | Reason |
|-------|----------------|--------|
| anggota | `(id_koperasi, status, tanggal_dihapus)` | Common query: active members by cooperative |
| simpanan | `(id_koperasi, tipe_simpanan, tanggal_transaksi)` | Reports filter by cooperative and type |
| transaksi | `(id_koperasi, status_balanced, tanggal_transaksi)` | Finding unbalanced entries |
| penjualan | `(id_koperasi, tanggal_penjualan, metode_pembayaran)` | Daily sales reports by payment method |

### Over-Indexed? 🤔

- `simpanan` has 8 indexes - some may be redundant
- Consider if all are necessary based on actual query patterns

---

## 🏗️ Schema Design Patterns

### Strengths

✅ **Soft Delete Pattern**
- All tables have `tanggal_dihapus` (deleted_at)
- Allows audit trail and data recovery
- Properly indexed

✅ **Audit Timestamps**
- `tanggal_dibuat` and `tanggal_diperbarui` on all tables
- Good for compliance and debugging

✅ **UUID Primary Keys**
- Distributed-friendly
- No ID collision risk across cooperatives
- Good for eventual microservices migration

✅ **Proper Foreign Keys**
- All relationships enforced with FK constraints
- Some use `ON DELETE CASCADE` appropriately (line items)

✅ **JSONB for Settings**
- `koperasi.pengaturan` uses JSONB for flexible configuration
- PostgreSQL native JSON support

### Areas for Improvement

⚠️ **Denormalization for Performance**
- `item_penjualan.nama_produk` is denormalized (good for historical accuracy)
- Consider similar pattern for prices (capture at transaction time)

⚠️ **Missing Check Constraints**
- No validation for `status` enum values
- No checks for positive amounts
- No validation for email format

⚠️ **Missing Triggers**
- Could use triggers to auto-update `total_debit`/`total_kredit` in `transaksi`
- Auto-update `stok` in `produk` after sales

---

## 📝 Recommended Migrations

### Priority 1: Fix Multi-Tenant Constraints (CRITICAL)

See: `002_fix_remaining_multitenant_constraints.sql`

### Priority 2: Add Check Constraints (HIGH)

```sql
-- Add check constraints for data integrity
ALTER TABLE simpanan
  ADD CONSTRAINT chk_simpanan_jumlah_positive
  CHECK (jumlah_setoran > 0);

ALTER TABLE transaksi
  ADD CONSTRAINT chk_transaksi_balanced
  CHECK (total_debit = total_kredit OR status_balanced = false);

ALTER TABLE produk
  ADD CONSTRAINT chk_produk_harga_positive
  CHECK (harga >= 0 AND (harga_beli IS NULL OR harga_beli >= 0));
```

### Priority 3: Add Missing Indexes (MEDIUM)

```sql
-- Composite indexes for common queries
CREATE INDEX idx_anggota_koperasi_status_active
  ON anggota (id_koperasi, status)
  WHERE status = 'aktif' AND tanggal_dihapus IS NULL;

CREATE INDEX idx_simpanan_koperasi_type_date
  ON simpanan (id_koperasi, tipe_simpanan, tanggal_transaksi DESC);

CREATE INDEX idx_penjualan_daily_report
  ON penjualan (id_koperasi, tanggal_penjualan, metode_pembayaran)
  WHERE tanggal_dihapus IS NULL;
```

### Priority 4: Add Database Functions (LOW)

```sql
-- Helper function to calculate member balance
CREATE OR REPLACE FUNCTION get_member_balance(
  p_id_koperasi UUID,
  p_id_anggota UUID
)
RETURNS TABLE (
  simpanan_pokok NUMERIC,
  simpanan_wajib NUMERIC,
  simpanan_sukarela NUMERIC,
  total_simpanan NUMERIC
) AS $$
BEGIN
  RETURN QUERY
  SELECT
    COALESCE(SUM(CASE WHEN tipe_simpanan = 'pokok' THEN jumlah_setoran ELSE 0 END), 0),
    COALESCE(SUM(CASE WHEN tipe_simpanan = 'wajib' THEN jumlah_setoran ELSE 0 END), 0),
    COALESCE(SUM(CASE WHEN tipe_simpanan = 'sukarela' THEN jumlah_setoran ELSE 0 END), 0),
    COALESCE(SUM(jumlah_setoran), 0)
  FROM simpanan
  WHERE id_koperasi = p_id_koperasi
    AND id_anggota = p_id_anggota
    AND tanggal_dihapus IS NULL;
END;
$$ LANGUAGE plpgsql;
```

---

## 🔐 Security Considerations

### Current State

✅ **Multi-Tenancy:** Cooperative-level isolation via `id_koperasi`
✅ **Soft Deletes:** Data not permanently removed
✅ **Password Hashing:** Using bcrypt for both users and member PINs
✅ **Foreign Key Constraints:** Referential integrity enforced

### Recommendations

1. **Row-Level Security (RLS):**
   ```sql
   -- Enable RLS on all tenant tables
   ALTER TABLE anggota ENABLE ROW LEVEL SECURITY;

   -- Policy: Users can only see their cooperative's data
   CREATE POLICY anggota_tenant_isolation ON anggota
     USING (id_koperasi = current_setting('app.current_koperasi_id')::UUID);
   ```

2. **Audit Logging:**
   - Add `dibuat_oleh` to all tables (already on some)
   - Consider audit trigger for sensitive changes

3. **Data Encryption:**
   - Consider encrypting PII fields (NIK, email, phone)
   - Use PostgreSQL's pgcrypto extension

---

## 📈 Performance Optimization

### Current Performance Profile

- **Small Dataset:** Development stage, minimal performance concerns
- **Indexed Lookups:** Fast on foreign keys and dates
- **Query Patterns:** Mostly OLTP (transaction processing)

### Scalability Recommendations

1. **Partitioning (Future):**
   - Partition large tables by `id_koperasi` when >10 cooperatives
   - Partition transaction tables by date range (monthly/yearly)

2. **Materialized Views:**
   ```sql
   -- For dashboard/reporting queries
   CREATE MATERIALIZED VIEW mv_daily_sales_summary AS
   SELECT
     id_koperasi,
     DATE(tanggal_penjualan) as tanggal,
     COUNT(*) as jumlah_transaksi,
     SUM(total_belanja) as total_penjualan,
     AVG(total_belanja) as rata_rata_per_transaksi
   FROM penjualan
   WHERE tanggal_dihapus IS NULL
   GROUP BY id_koperasi, DATE(tanggal_penjualan);

   -- Refresh daily
   CREATE INDEX ON mv_daily_sales_summary (id_koperasi, tanggal);
   ```

3. **Connection Pooling:**
   - Use PgBouncer or similar for connection management
   - Essential when scaling to 10+ concurrent cooperatives

---

## 🎯 Next Steps

### Immediate (This Week)

1. ✅ **DONE:** Fix `anggota` multi-tenant constraint
2. 🔴 **TODO:** Apply migration `002_fix_remaining_multitenant_constraints.sql`
3. 🔴 **TODO:** Test multi-tenant isolation with 2+ cooperatives
4. 🟡 **TODO:** Add check constraints for data validation

### Short Term (Before MVP Launch)

1. Add missing composite indexes
2. Implement database functions for common calculations
3. Set up automated backups
4. Create schema documentation (ER diagrams)
5. Performance testing with realistic data volume

### Long Term (Post-MVP)

1. Implement Row-Level Security (RLS)
2. Set up read replicas for reporting
3. Consider partitioning strategy for scale
4. Implement materialized views for dashboards
5. Full audit logging implementation

---

## 📚 Schema Best Practices Compliance

| Practice | Status | Notes |
|----------|--------|-------|
| Normalized structure | ✅ Good | 3NF in most tables |
| Proper data types | ✅ Good | Using appropriate types |
| Primary keys | ✅ Good | UUID on all tables |
| Foreign keys | ✅ Good | All relationships enforced |
| Indexes on FKs | ✅ Good | All FKs indexed |
| NOT NULL constraints | ⚠️ Partial | Some columns should be NOT NULL |
| CHECK constraints | ❌ Missing | Need data validation |
| Default values | ✅ Good | Sensible defaults |
| Timestamps | ✅ Good | Created/updated on all |
| Soft deletes | ✅ Good | Implemented consistently |
| Multi-tenant isolation | ⚠️ **In Progress** | **5 constraints need fixing** |

---

## 🔧 Migration Scripts Generated

1. ✅ `001_fix_anggota_unique_constraint.sql` - Applied
2. 🔴 `002_fix_remaining_multitenant_constraints.sql` - **Needs Review & Application**
3. 🟡 `003_add_check_constraints.sql` - Recommended
4. 🟡 `004_add_performance_indexes.sql` - Recommended

---

## 📞 Contact & Questions

For schema-related questions:
- Review team: Backend architects
- DBA: Database administrator
- Reference: CLAUDE.md for project context

---

**Generated by:** Database Schema Analysis Tool
**Review Status:** Complete
**Action Required:** Apply Priority 1 migration immediately
