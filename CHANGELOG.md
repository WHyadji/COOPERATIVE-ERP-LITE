# Changelog

All notable changes to the Cooperative ERP Lite project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added - 2025-01-18

#### Accounting Module
- **Chart of Accounts (Akun)**:
  - Complete CRUD operations for account management
  - Hierarchical account display with parent-child relationships
  - Filter by account type (Aset, Kewajiban, Modal, Pendapatan, Beban)
  - Account balance tracking and calculation
  - Seed default Indonesian cooperative Chart of Accounts
  - Account form with validation and parent account selection

- **Journal Entries (Jurnal)**:
  - Create journal entries with multiple line items
  - Edit existing journal entries with full transaction update
  - View transaction details with comprehensive information display
  - Double-entry validation (Total Debit = Total Kredit)
  - Balance status indicators (Balanced/Unbalanced chips)
  - Date range filtering for transaction list
  - Print-friendly transaction detail view
  - Transaction line items table with account selection
  - Automatic total calculation and validation

#### User Experience Improvements
- **Toast Notification System**:
  - Global toast notification context using MUI Snackbar
  - Success, error, info, and warning message types
  - Auto-dismiss after 6 seconds
  - Replaced all native `alert()` calls with professional toast notifications
  - Consistent feedback across all CRUD operations

- **Audit Trail Tracking**:
  - Track creator information (user ID and name)
  - Track last updater information (user ID and name)
  - Display creation and update timestamps
  - Show audit trail in transaction detail view
  - Backend support for audit fields in all transaction operations

#### Backend Enhancements
- **Transaction Update Endpoint**:
  - `PUT /api/v1/transaksi/:id` for updating journal entries
  - Transaction-based update with atomic operations
  - Delete old line items and create new ones for consistency
  - Preserve audit trail information (creator, updater)
  - Validation for double-entry accounting rules

- **Audit Trail Fields**:
  - Added `diperbaruiOleh` (updater UUID) field to Transaksi model
  - Enhanced TransaksiResponse with audit fields:
    - `dibuatOleh` (creator UUID)
    - `namaDibuatOleh` (creator name)
    - `diperbaruiOleh` (updater UUID)
    - `namaDiperbaruiOleh` (updater name)
    - `tanggalDibuat` (creation timestamp)
    - `tanggalDiperbarui` (update timestamp)
  - Automatic population of audit fields in service layer

### Changed

#### Frontend
- Updated TransactionForm component to support both create and edit modes
- Enhanced journal entry list page with Edit button in actions column
- Improved transaction detail page with Edit functionality
- Refactored metadata section to display comprehensive audit trail
- Updated TypeScript types to include audit trail fields

#### Backend
- Modified `PerbaruiTransaksi` service method for proper transaction updates
- Enhanced error handling for transaction update operations
- Improved validation messages for double-entry accounting

### Technical Details

#### Dependencies
- No new dependencies added for toast notifications (using existing MUI Snackbar)
- Leveraging existing Material-UI components for consistent design

#### Database Schema
- Added `diperbarui_oleh` column to `transaksi` table (UUID, nullable)
- Maintains backward compatibility with existing records

#### API Changes
- Enhanced `GET /api/v1/transaksi/:id` response with audit fields
- Implemented `PUT /api/v1/transaksi/:id` endpoint
- All responses include creator and updater information when available

### Files Modified

#### Frontend
- `frontend/lib/context/ToastContext.tsx` (NEW) - Toast notification system
- `frontend/app/(dashboard)/layout.tsx` - Added ToastProvider
- `frontend/app/(dashboard)/akuntansi/jurnal/page.tsx` - Edit functionality
- `frontend/app/(dashboard)/akuntansi/jurnal/[id]/page.tsx` - Enhanced detail view
- `frontend/components/accounting/TransactionForm.tsx` - Edit mode support
- `frontend/types/index.ts` - Added audit trail fields
- `frontend/lib/api/accountingApi.ts` - Added updateTransaction function
- `frontend/README.md` - Updated documentation

#### Backend
- `backend/internal/models/transaksi.go` - Added audit trail fields
- `backend/internal/services/transaksi_service.go` - PerbaruiTransaksi method
- `backend/internal/handlers/transaksi_handler.go` - Update handler implementation

## [0.1.0] - Previous Release

### Added
- Member Management (Anggota) - Complete CRUD operations
- Savings Management (Simpanan) - Simpanan Pokok, Wajib, Sukarela
- Authentication system with JWT
- Dashboard layout with responsive sidebar
- Role-based access control (Admin, Bendahara, Kasir, Anggota)
- Multi-tenant architecture with cooperative isolation

### Features by Module

#### Authentication
- Login with username/password
- JWT token-based authentication
- Auto-redirect on session expiration
- Protected routes with middleware

#### Member Management
- Paginated member list with search and filters
- Multi-section member form with validation
- Member detail view
- Edit and delete functionality
- Member statistics dashboard

#### Savings Management
- Create savings transactions (Pokok, Wajib, Sukarela)
- View member balance summaries
- Transaction history with pagination
- Date range filtering
- Balance calculations by savings type

#### Chart of Accounts
- Pre-configured Indonesian cooperative COA
- Account hierarchy with parent-child relationships
- Account type classification
- Create and edit accounts
- Account status management (active/inactive)

---

## Future Releases

### Planned for v0.2.0
- Account Ledger (Buku Besar) view
- Financial Reports (Balance Sheet, Income Statement)
- Trial Balance report
- POS/Sales module
- Product management
- Inventory tracking

### Planned for v0.3.0
- Advanced reporting with export to PDF/Excel
- User management for administrators
- System settings and configuration
- Backup and restore functionality
- Activity logs and audit reports
