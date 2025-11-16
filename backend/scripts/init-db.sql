-- ============================================
-- COOPERATIVE ERP LITE - Database Initialization
-- ============================================
-- This script runs automatically when PostgreSQL container starts
-- It creates the database and sets up extensions

-- Create database if not exists (handled by POSTGRES_DB env var)
-- This file is mainly for extensions and initial setup

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone to Jakarta
SET timezone = 'Asia/Jakarta';

-- Create custom types if needed in the future
-- (Currently handled by GORM)

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE koperasi_erp TO postgres;

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'Database initialization completed successfully!';
    RAISE NOTICE 'UUID extension enabled';
    RAISE NOTICE 'Timezone set to Asia/Jakarta';
END $$;
