# Scripts Directory

This directory contains utility scripts for database seeding, testing, and debugging.

## Directory Structure

```
scripts/
├── seed/           # Database seeding scripts
├── test/           # Testing scripts
├── debug/          # Debugging and diagnostic scripts
└── README.md       # This file
```

## Seed Scripts (`seed/`)

Scripts for populating the database with initial and test data.

| Script | Description | Usage |
|--------|-------------|-------|
| `seed-initial.sh` | Seed initial cooperative, users, and basic data | `./scripts/seed/seed-initial.sh` |
| `seed-coa.sh` | Seed Chart of Accounts (SAK ETAP) | `./scripts/seed/seed-coa.sh` |
| `test-data.sh` | Generate test data for development | `./scripts/seed/test-data.sh` |
| `create-test-data-sequential.sh` | Create test data sequentially | `./scripts/seed/create-test-data-sequential.sh` |
| `cleanup-test-data.sh` | Remove all test data from database | `./scripts/seed/cleanup-test-data.sh` |

### Common Seed Workflow

```bash
# 1. Seed initial data
./scripts/seed/seed-initial.sh

# 2. Seed Chart of Accounts
./scripts/seed/seed-coa.sh

# 3. (Optional) Add test data for development
./scripts/seed/test-data.sh
```

## Test Scripts (`test/`)

Scripts for testing API endpoints and system functionality.

| Script | Description | Usage |
|--------|-------------|-------|
| `test-member.sh` | Test member management endpoints | `./scripts/test/test-member.sh` |
| `test-deposit.sh` | Test simpanan (deposit) endpoints | `./scripts/test/test-deposit.sh` |
| `test-concurrent.sh` | Test concurrent transactions (race conditions) | `./scripts/test/test-concurrent.sh` |
| `test-api-summary.sh` | Run comprehensive API test suite | `./scripts/test/test-api-summary.sh` |
| `test-edge-cases.sh` | Test edge cases and error handling | `./scripts/test/test-edge-cases.sh` |

### Running Tests

```bash
# Run all API tests
./scripts/test/test-api-summary.sh

# Test specific functionality
./scripts/test/test-member.sh
./scripts/test/test-deposit.sh

# Test concurrency and race conditions
./scripts/test/test-concurrent.sh

# Test edge cases
./scripts/test/test-edge-cases.sh
```

## Debug Scripts (`debug/`)

Scripts for debugging and diagnostics.

| Script | Description | Usage |
|--------|-------------|-------|
| `debug-member.sh` | Debug member-related issues | `./scripts/debug/debug-member.sh` |
| `check-balance.sh` | Check account balances and trial balance | `./scripts/debug/check-balance.sh` |

### Debugging Workflow

```bash
# Check member data
./scripts/debug/debug-member.sh

# Verify accounting balances
./scripts/debug/check-balance.sh
```

## Prerequisites

All scripts require:
- Backend server running on `http://localhost:8080`
- PostgreSQL database accessible
- Valid authentication token (scripts handle login automatically)

## Environment Variables

Some scripts may use these environment variables:

```bash
API_URL=http://localhost:8080/api/v1
DB_HOST=localhost
DB_PORT=5432
DB_NAME=koperasi_erp
DB_USER=postgres
DB_PASSWORD=postgres
```

## Script Conventions

All scripts follow these conventions:
- Executable permissions: `chmod +x script-name.sh`
- Use bash shebang: `#!/bin/bash`
- Return non-zero exit codes on failure
- Print progress and results to stdout
- Print errors to stderr

## Development Notes

- Scripts are version-controlled in git
- Always test scripts in development before using in production
- Update this README when adding new scripts
- Use meaningful script names that describe their purpose

## Quick Reference

### Full Development Setup
```bash
# Start services
make up

# Seed database
./scripts/seed/seed-initial.sh
./scripts/seed/seed-coa.sh

# Run tests
./scripts/test/test-api-summary.sh

# Check system health
./scripts/debug/check-balance.sh
```

### Clean Slate
```bash
# Remove test data
./scripts/seed/cleanup-test-data.sh

# Re-seed fresh data
./scripts/seed/seed-initial.sh
./scripts/seed/seed-coa.sh
```

## Troubleshooting

**Problem**: Script returns "Permission denied"
```bash
# Solution: Make script executable
chmod +x scripts/seed/script-name.sh
```

**Problem**: "Connection refused" errors
```bash
# Solution: Ensure backend server is running
make dev
# or
docker-compose up backend
```

**Problem**: Database errors
```bash
# Solution: Check database is running and accessible
make db-shell
# or
docker ps | grep postgres
```

## Contributing

When adding new scripts:
1. Place in appropriate directory (`seed/`, `test/`, or `debug/`)
2. Make executable: `chmod +x script-name.sh`
3. Add documentation to this README
4. Test thoroughly before committing
5. Follow existing script conventions

---

For more information, see the main project [README.md](../README.md) or [CLAUDE.md](../CLAUDE.md).
