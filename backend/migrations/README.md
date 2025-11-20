# Database Migrations

This directory contains SQL migration scripts for the Cooperative ERP Lite database schema.

## Migration Naming Convention

Migrations should be named using the following format:

```
<sequence>_<descriptive_name>.sql
```

Examples:
- `001_fix_anggota_unique_constraint.sql`
- `002_add_product_categories.sql`
- `003_create_audit_log_table.sql`

## How to Apply Migrations

### Development (Local PostgreSQL via Docker)

```bash
# Apply a single migration
docker exec -i cooperative-erp-postgres psql -U postgres -d koperasi_erp < migrations/001_fix_anggota_unique_constraint.sql

# Apply all migrations in order
for migration in migrations/*.sql; do
  echo "Applying $migration..."
  docker exec -i cooperative-erp-postgres psql -U postgres -d koperasi_erp < "$migration"
done
```

### Production (Cloud SQL)

For production deployments, migrations should be applied carefully:

1. **Backup the database first**
2. **Test the migration in staging**
3. **Apply during maintenance window if possible**
4. **Monitor for errors**

```bash
# Connect to Cloud SQL proxy, then:
psql -h localhost -U postgres -d koperasi_erp < migrations/001_fix_anggota_unique_constraint.sql
```

## Migration Best Practices

1. **Always use transactions** - Wrap changes in `BEGIN` and `COMMIT`
2. **Make migrations idempotent** - Use `IF EXISTS`, `IF NOT EXISTS`, or `ON CONFLICT`
3. **Add rollback instructions** - Document how to undo the migration
4. **Test thoroughly** - Test on a copy of production data
5. **Document changes** - Add comments explaining why the migration is needed
6. **Keep migrations small** - One logical change per migration

## Migration Template

```sql
-- ============================================================================
-- Migration: <Descriptive Title>
-- Date: YYYY-MM-DD
-- Description: Brief description of what this migration does and why
-- ============================================================================

-- ISSUE/CONTEXT:
-- Explain the problem this migration solves

-- CHANGES:
-- List the specific changes being made

BEGIN;

-- Your migration SQL here
-- ...

-- Verification query (optional but recommended)
SELECT 'Migration completed successfully' as status;

COMMIT;

-- ============================================================================
-- ROLLBACK INSTRUCTIONS (if applicable)
-- ============================================================================
-- BEGIN;
-- -- Rollback steps here
-- COMMIT;
```

## Applied Migrations

| Migration | Date Applied | Description |
|-----------|--------------|-------------|
| 001_fix_anggota_unique_constraint.sql | 2025-11-20 | Fixed unique constraint on anggota table to support multi-tenant (id_koperasi, nomor_anggota) |
| 002_fix_remaining_multitenant_constraints.sql | 2025-11-20 | Fixed remaining 5 multi-tenant constraints on akun, transaksi, produk, penjualan, and pengguna tables |
| 003_implement_row_level_security.sql | 2025-11-20 | Implemented PostgreSQL Row-Level Security (RLS) for database-level multi-tenant isolation |
| 004_add_financial_constraints.sql | 2025-11-20 | Added CHECK constraints to prevent negative amounts in financial transactions (7 constraints) |
| 005_add_product_constraints.sql | 2025-11-20 | Added CHECK constraints for product prices, stock, and sale item validation (7 constraints) |
| 006_add_updated_at_columns.sql | 2025-11-20 | Added tanggal_diubah (updated_at) column to all tables with automatic trigger updates (complete audit trail) |
| 007_add_enum_constraints.sql | 2025-11-20 | Added CHECK constraints for enum validation on status, role, type fields (8 constraints) |
| 008_add_performance_indexes.sql | 2025-11-20 | Added composite and partial indexes for query optimization (6 performance indexes) |

## Future Migration Tool

For Phase 2+, consider using a migration tool like:
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)
- [Atlas](https://atlasgo.io/)

These tools provide:
- Automatic tracking of applied migrations
- Rollback support
- Version control
- Better error handling
- Integration with CI/CD pipelines
