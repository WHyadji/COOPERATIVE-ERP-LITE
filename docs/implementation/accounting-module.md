# Accounting Module Documentation

## Overview

The Accounting Module (Akuntansi) provides comprehensive double-entry bookkeeping functionality for Indonesian cooperatives, following SAK ETAP (Indonesian GAAP for cooperatives) standards.

## Features

### 1. Chart of Accounts (Bagan Akun)

#### Account Types
The system supports five main account types following Indonesian accounting standards:

- **Aset** (Assets) - Resources owned by the cooperative
- **Kewajiban** (Liabilities) - Obligations and debts
- **Modal** (Equity) - Owner's equity and capital
- **Pendapatan** (Revenue) - Income and sales
- **Beban** (Expenses) - Operating costs and expenses

#### Account Structure
- Hierarchical account structure with parent-child relationships
- Account code system (e.g., 1-1000, 1-1100, 1-1110)
- Normal balance designation (Debit/Kredit)
- Active/Inactive status management
- Description field for detailed account information

#### Default Chart of Accounts
The system comes pre-configured with Indonesian cooperative standard accounts:

```
1-xxxx: Aset (Assets)
  1-1xxx: Aset Lancar (Current Assets)
    1-1000: Kas (Cash)
    1-1100: Bank (Bank Accounts)
    1-1200: Piutang Anggota (Member Receivables)
    1-1300: Persediaan (Inventory)

  1-2xxx: Aset Tetap (Fixed Assets)
    1-2000: Tanah (Land)
    1-2100: Bangunan (Buildings)
    1-2200: Peralatan (Equipment)

2-xxxx: Kewajiban (Liabilities)
  2-1xxx: Kewajiban Lancar (Current Liabilities)
    2-1000: Hutang Usaha (Accounts Payable)
    2-1100: Simpanan Anggota (Member Savings)

3-xxxx: Modal (Equity)
  3-1000: Modal Dasar (Basic Capital)
  3-2000: Cadangan (Reserves)
  3-3000: SHU Tahun Berjalan (Current Year SHU)

4-xxxx: Pendapatan (Revenue)
  4-1000: Pendapatan Penjualan (Sales Revenue)
  4-2000: Pendapatan Jasa (Service Revenue)

5-xxxx: Beban (Expenses)
  5-1000: Beban Gaji (Salary Expenses)
  5-2000: Beban Operasional (Operating Expenses)
```

### 2. Journal Entries (Jurnal Umum)

#### Creating Journal Entries

Journal entries follow the double-entry bookkeeping principle:
- **Total Debit = Total Kredit** (mandatory validation)
- Each line item can have EITHER debit OR kredit (not both)
- Minimum 2 line items required per transaction
- Each line item must have:
  - Account selection (from Chart of Accounts)
  - Amount (Debit OR Kredit)
  - Optional description/note

#### Transaction Header Fields
- **Nomor Jurnal** (Journal Number): Unique identifier (e.g., JU-2025-001)
- **Tanggal Transaksi** (Transaction Date): Date of the transaction
- **Deskripsi** (Description): Main transaction description
- **Nomor Referensi** (Reference Number): External reference/voucher number (optional)
- **Tipe Transaksi** (Transaction Type): Type classification (default: "manual")

#### Transaction Line Items
Each journal entry contains multiple line items (baris transaksi):
- Account selection from active accounts
- Debit amount (for debit entries)
- Kredit amount (for kredit entries)
- Line item description (optional, inherits from header if blank)

#### Example Journal Entry

**Transaction**: Cash sale of Rp 1,000,000

| Akun | Kode | Deskripsi | Debit | Kredit |
|------|------|-----------|-------|--------|
| Kas | 1-1000 | Penerimaan dari penjualan | 1,000,000 | - |
| Pendapatan Penjualan | 4-1000 | Penjualan tunai | - | 1,000,000 |
| **TOTAL** | | | **1,000,000** | **1,000,000** |

Status: ✓ **Balanced** (Debit = Kredit)

### 3. Transaction Management

#### Create Transaction
1. Click "Tambah Jurnal" button
2. Fill transaction header (Nomor Jurnal, Tanggal, Deskripsi)
3. Add line items (minimum 2):
   - Select account
   - Enter debit OR kredit amount
   - Add description (optional)
4. System validates:
   - All required fields filled
   - Total Debit = Total Kredit
   - Each line has account and amount
5. Click "Simpan Jurnal" to save

#### Edit Transaction
1. From list page: Click Edit icon in actions column
2. From detail page: Click "Edit" button
3. Form opens with pre-populated data
4. Modify header or line items as needed
5. System re-validates on save
6. Audit trail tracks updater information

#### View Transaction Details
Transaction detail page shows:
- Transaction header information
- Complete line items table with account codes and names
- Total debit and kredit calculations
- Balance status indicator (Balanced/Unbalanced)
- Audit trail information:
  - Created by [User Name] on [Date]
  - Last updated by [User Name] on [Date]
- Action buttons: Back, Edit, Delete, Print

#### Delete Transaction
1. Click Delete icon/button
2. Confirmation dialog appears
3. Confirm to permanently delete
4. Toast notification confirms deletion

### 4. Validation Rules

#### Double-Entry Validation
```typescript
// Rule 1: Total Debit must equal Total Kredit
totalDebit === totalKredit

// Rule 2: Each line must have EITHER debit OR kredit (not both)
(debit > 0 && kredit === 0) || (debit === 0 && kredit > 0)

// Rule 3: Each line must have an amount
debit > 0 || kredit > 0

// Rule 4: Each line must have an account selected
idAkun !== null && idAkun !== ''

// Rule 5: Minimum 2 line items required
filledLineItems.length >= 2
```

#### Error Messages
- "Total debit harus sama dengan total kredit" - Unbalanced entry
- "Satu baris tidak boleh memiliki debit dan kredit sekaligus" - Both debit and kredit filled
- "Setiap baris harus memiliki nilai debit atau kredit" - No amount entered
- "Setiap baris harus memilih akun" - Account not selected
- "Minimal 2 baris transaksi harus diisi" - Insufficient line items

### 5. User Interface

#### List Page Features
- Paginated table (10/20/50/100 rows per page)
- Date range filter (Tanggal Mulai - Tanggal Akhir)
- Columns:
  - Nomor Jurnal (Journal Number)
  - Tanggal (Date)
  - Deskripsi (Description)
  - Nomor Referensi (Reference)
  - Total Debit
  - Total Kredit
  - Status (Balanced/Unbalanced chip)
  - Actions (View, Edit, Delete)

#### Detail Page Features
- Breadcrumb navigation
- Transaction header information card
- Line items table
- Summary cards (Total Debit, Total Kredit, Difference)
- Balance status alert
- Audit trail information box
- Print functionality (hides action buttons)

#### Form Features
- Dynamic title: "Buat Jurnal Umum Baru" vs "Edit Jurnal Umum"
- Add/Remove line item buttons
- Real-time total calculation
- Balance status indicator
- Account dropdown with code and name
- Date picker for transaction date
- Validation on submit
- Toast notifications for success/error

### 6. API Integration

#### Endpoints

**Chart of Accounts:**
```
GET    /api/v1/akun                 - List all accounts
POST   /api/v1/akun                 - Create account
GET    /api/v1/akun/:id             - Get account details
PUT    /api/v1/akun/:id             - Update account
DELETE /api/v1/akun/:id             - Delete account
GET    /api/v1/akun/:id/saldo       - Get account balance
POST   /api/v1/akun/seed-coa        - Seed default COA
```

**Transactions:**
```
GET    /api/v1/transaksi            - List transactions (paginated)
POST   /api/v1/transaksi            - Create journal entry
GET    /api/v1/transaksi/:id        - Get transaction details
PUT    /api/v1/transaksi/:id        - Update journal entry
DELETE /api/v1/transaksi/:id        - Delete transaction
```

**Reports:**
```
GET    /api/v1/laporan/buku-besar   - Get account ledger
```

#### Request/Response Examples

**Create Transaction:**
```json
POST /api/v1/transaksi

Request:
{
  "nomorJurnal": "JU-2025-001",
  "tanggalTransaksi": "2025-01-18",
  "deskripsi": "Penjualan tunai barang",
  "nomorReferensi": "INV-001",
  "tipeTransaksi": "manual",
  "barisTransaksi": [
    {
      "idAkun": "uuid-kas",
      "jumlahDebit": 1000000,
      "jumlahKredit": 0,
      "keterangan": "Penerimaan kas"
    },
    {
      "idAkun": "uuid-pendapatan",
      "jumlahDebit": 0,
      "jumlahKredit": 1000000,
      "keterangan": "Pendapatan penjualan"
    }
  ]
}

Response:
{
  "success": true,
  "message": "Transaksi berhasil dibuat",
  "data": {
    "id": "uuid",
    "nomorJurnal": "JU-2025-001",
    "tanggalTransaksi": "2025-01-18T00:00:00Z",
    "deskripsi": "Penjualan tunai barang",
    "totalDebit": 1000000,
    "totalKredit": 1000000,
    "statusBalanced": true,
    "dibuatOleh": "uuid",
    "namaDibuatOleh": "Admin Koperasi",
    "tanggalDibuat": "2025-01-18T10:30:00Z",
    "barisTransaksi": [...]
  }
}
```

### 7. Audit Trail System

#### Tracked Information
Every transaction records:
- **Creator**: User who created the transaction
  - `dibuatOleh` (UUID)
  - `namaDibuatOleh` (Full name)
  - `tanggalDibuat` (Timestamp)

- **Last Updater**: User who last modified the transaction
  - `diperbaruiOleh` (UUID)
  - `namaDiperbaruiOleh` (Full name)
  - `tanggalDiperbarui` (Timestamp)

#### Display Format
```
Informasi Audit Trail
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ID Transaksi: 123e4567-e89b-12d3-a456-426614174000

Dibuat oleh Admin Koperasi pada 18/01/2025

Terakhir diperbarui oleh Bendahara Koperasi pada 18/01/2025
```

### 8. Multi-tenant Isolation

All accounting data is automatically filtered by cooperative:
- Each cooperative has isolated Chart of Accounts
- Transactions are scoped to the cooperative
- Account balances calculated per cooperative
- No cross-cooperative data leakage

The `idKoperasi` (cooperative ID) is automatically extracted from the JWT token and applied to all queries.

### 9. Indonesian Accounting Standards (SAK ETAP)

The module follows SAK ETAP guidelines:
- Double-entry bookkeeping
- Five main account classifications
- Proper revenue and expense recognition
- Balance sheet and income statement preparation
- Cooperative-specific accounts (Simpanan, Modal, SHU)

### 10. Future Enhancements

#### Planned Features
- Account Ledger (Buku Besar) view with running balance
- Trial Balance report
- Balance Sheet (Neraca)
- Income Statement (Laba Rugi)
- Cash Flow Statement (Arus Kas)
- Journal Entry reversal functionality
- Recurring journal entries
- Journal entry templates
- Multi-currency support
- Fiscal year management
- Bank reconciliation
- Account closing process (year-end)

#### Advanced Features (Phase 2+)
- Budget tracking and variance analysis
- Cost center allocation
- Project/department segmentation
- Approval workflow for transactions
- Batch journal entry import
- Integration with bank feeds
- Automated bank reconciliation
- Financial dashboard and KPIs
- Custom report builder
- Export to accounting software (Excel, PDF)

## Best Practices

### 1. Account Naming
- Use clear, descriptive account names in Indonesian
- Follow standard account code patterns
- Group related accounts under parent accounts
- Use consistent naming conventions

### 2. Journal Entry Description
- Write clear, concise descriptions
- Include external reference numbers when available
- Specify transaction source (e.g., "Penjualan Tunai", "Pembelian Kredit")
- Add line item notes for additional detail

### 3. Account Organization
- Keep Chart of Accounts organized and clean
- Deactivate unused accounts instead of deleting
- Use parent accounts for grouping
- Regular review and maintenance of account list

### 4. Data Entry
- Double-check amounts before saving
- Ensure balanced entries (Debit = Kredit)
- Select correct accounts for each transaction type
- Verify transaction dates
- Use reference numbers for traceability

### 5. Security
- Implement role-based access control
- Limit edit/delete permissions to authorized users
- Regular backup of accounting data
- Audit trail review for compliance

## Troubleshooting

### Common Issues

**"Total debit harus sama dengan total kredit"**
- Check all line items have correct amounts
- Verify no typing errors in amounts
- Ensure decimal points are correct
- Use calculator to verify totals

**"Minimal 2 baris transaksi harus diisi"**
- Add at least one more line item
- Ensure both lines have account and amount selected
- Check that amounts are greater than 0

**Transaction not saving**
- Check all required fields are filled
- Verify balance validation passes
- Check browser console for errors
- Ensure backend is running and accessible

**Account not appearing in dropdown**
- Check account status is "Active"
- Verify account belongs to your cooperative
- Refresh the page to reload account list
- Check filter settings in account list

## Technical Implementation

### Frontend Stack
- **Framework**: Next.js 15 with App Router
- **Language**: TypeScript 5.7
- **UI Library**: Material-UI 6.3
- **Forms**: React Hook Form + controlled components
- **State**: React Context (ToastContext, AuthContext)
- **API Client**: Axios with interceptors
- **Date Handling**: date-fns 4.1

### Backend Stack
- **Language**: Go 1.25
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Authentication**: JWT tokens

### Database Schema

**transaksi table:**
```sql
CREATE TABLE transaksi (
    id UUID PRIMARY KEY,
    id_koperasi UUID NOT NULL,
    nomor_jurnal VARCHAR(50) NOT NULL,
    tanggal_transaksi DATE NOT NULL,
    deskripsi TEXT NOT NULL,
    nomor_referensi VARCHAR(50),
    tipe_transaksi VARCHAR(20),
    total_debit DECIMAL(15,2) NOT NULL,
    total_kredit DECIMAL(15,2) NOT NULL,
    status_balanced BOOLEAN NOT NULL,
    dibuat_oleh UUID NOT NULL,
    diperbarui_oleh UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(id_koperasi, nomor_jurnal)
);

CREATE TABLE baris_transaksi (
    id UUID PRIMARY KEY,
    id_transaksi UUID NOT NULL REFERENCES transaksi(id),
    id_akun UUID NOT NULL REFERENCES akun(id),
    jumlah_debit DECIMAL(15,2) DEFAULT 0,
    jumlah_kredit DECIMAL(15,2) DEFAULT 0,
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    CHECK (jumlah_debit >= 0 AND jumlah_kredit >= 0),
    CHECK (NOT (jumlah_debit > 0 AND jumlah_kredit > 0))
);
```

### Component Architecture

```
AccountingModule/
├── ChartOfAccounts/
│   ├── AccountList (page)
│   ├── AccountForm (component)
│   └── AccountDetail (component)
│
├── JournalEntries/
│   ├── TransactionList (page)
│   ├── TransactionForm (component)
│   ├── TransactionDetail (page)
│   └── LineItemRow (component)
│
├── Reports/
│   ├── AccountLedger (page)
│   ├── TrialBalance (component)
│   └── FinancialStatements (component)
│
└── Shared/
    ├── AccountSelector (component)
    ├── BalanceIndicator (component)
    └── AuditTrail (component)
```

## Conclusion

The Accounting Module provides a robust, user-friendly double-entry bookkeeping system tailored for Indonesian cooperatives. It follows SAK ETAP standards, provides comprehensive audit trails, and offers a modern web interface for efficient accounting operations.

For support or feature requests, please contact the development team.
