# Cooperative ERP Lite - Frontend

Next.js 15.5 frontend application for the Cooperative ERP Lite system - Indonesia's first cooperative operating system.

## Tech Stack

- **Framework**: Next.js 15.5 with App Router & Turbopack
- **Language**: TypeScript 5.7
- **UI Library**: Material-UI (MUI) 6.3
- **Forms**: React Hook Form 7.54 + Zod 3.24
- **HTTP Client**: Axios 1.8
- **State Management**: React Context API
- **Date Handling**: date-fns 4.1

## Quick Start

### 1. Install Dependencies

```bash
npm install
```

### 2. Configure Environment

Create `.env.local` file:

```bash
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
```

### 3. Run Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000)

### 4. Build for Production

```bash
npm run build
npm start
```

## Project Structure

```
frontend/
├── app/                    # Next.js App Router
│   ├── (auth)/login        # Login page
│   └── (dashboard)/        # Protected dashboard
│       ├── anggota/        # Member management
│       ├── simpanan/       # Savings/Share capital
│       ├── akuntansi/      # Accounting module
│       │   ├── akun/       # Chart of accounts
│       │   └── jurnal/     # Journal entries
│       └── page.tsx        # Dashboard home
├── components/
│   ├── accounting/         # Accounting components
│   └── layout/            # Sidebar & Header
├── lib/
│   ├── api/               # API integration
│   └── context/           # Auth & Toast contexts
└── types/                 # TypeScript definitions
```

## Features Implemented

### ✅ Authentication

- Login with JWT
- Protected routes
- Auto-redirect on token expiration
- Role-based access control

### ✅ Member Management

- **List**: Paginated table with search & filters
- **Create**: Multi-section form with validation
- **View**: Member detail page
- **Edit**: Inline editing with cancel
- **Delete**: With confirmation

### ✅ Savings Management (Simpanan)

- **List**: Paginated table with filters by type & date range
- **Create**: Transaction form with member selection
- **View**: Balance summaries by member and type
- **Types**: Simpanan Pokok, Wajib, Sukarela
- **Reports**: Member balance summaries

### ✅ Accounting Module (Akuntansi)

- **Chart of Accounts (Akun)**:
  - List with hierarchical display
  - Create/Edit account forms
  - Filter by account type (Aset, Kewajiban, Modal, Pendapatan, Beban)
  - Account balance tracking
  - Seed default Indonesian cooperative COA

- **Journal Entries (Jurnal)**:
  - Create/Edit journal entries with line items
  - Double-entry validation (Debit = Kredit)
  - Transaction detail view with audit trail
  - Edit functionality with version tracking
  - Balance status indicators
  - Date range filtering
  - Print-friendly transaction views
  - Audit trail: Creator and updater tracking

### ✅ User Experience

- **Toast Notifications**: Professional feedback for all CRUD operations
- **Audit Trail**: Complete tracking of who created/updated records
- **Responsive Design**: Mobile-first Material-UI components
- **Loading States**: Skeleton screens and spinners
- **Error Handling**: User-friendly error messages

### ✅ Dashboard Layout

- Responsive sidebar navigation
- Mobile drawer
- User profile menu
- Role-based menu items

## Default Credentials (Development)

Contact backend team for test credentials.

## API Endpoints

All endpoints use `/api/v1` prefix:

### Authentication

- `POST /auth/login` - Authentication
- `GET /auth/profile` - User profile

### Members (Anggota)

- `GET /anggota` - List members with pagination
- `POST /anggota` - Create member
- `GET /anggota/:id` - Get member details
- `PUT /anggota/:id` - Update member
- `DELETE /anggota/:id` - Delete member
- `GET /anggota/statistik` - Member statistics

### Savings (Simpanan)

- `GET /simpanan` - List savings transactions
- `POST /simpanan` - Create savings transaction
- `GET /simpanan/:id` - Get transaction details
- `DELETE /simpanan/:id` - Delete transaction
- `GET /simpanan/saldo/:idAnggota` - Get member balance
- `GET /simpanan/ringkasan` - Get savings summary

### Accounting (Akuntansi)

- **Chart of Accounts**:
  - `GET /akun` - List accounts with filters
  - `POST /akun` - Create account
  - `GET /akun/:id` - Get account details
  - `PUT /akun/:id` - Update account
  - `DELETE /akun/:id` - Delete account
  - `GET /akun/:id/saldo` - Get account balance
  - `POST /akun/seed-coa` - Seed default COA

- **Transactions**:
  - `GET /transaksi` - List transactions with pagination
  - `POST /transaksi` - Create journal entry
  - `GET /transaksi/:id` - Get transaction details
  - `PUT /transaksi/:id` - Update journal entry
  - `DELETE /transaksi/:id` - Delete transaction

- **Reports**:
  - `GET /laporan/buku-besar` - Get account ledger

## User Roles

- **Admin**: Full access
- **Bendahara**: Finance & members
- **Kasir**: POS only
- **Anggota**: Read-only (future)

## Multi-tenant

All requests automatically include cooperative context from JWT token.

## Development

```bash
# Development with auto-reload
npm run dev

# Type checking
npx tsc --noEmit

# Linting
npm run lint

# Build
npm run build
```

## Next Steps

### Completed Features ✅

- ✅ Member Management (`/anggota`)
- ✅ Savings Management (`/simpanan`)
- ✅ Chart of Accounts (`/akuntansi/akun`)
- ✅ Journal Entries (`/akuntansi/jurnal`)
- ✅ Toast Notification System
- ✅ Audit Trail Tracking

### In Progress 🚧

- Account Ledger View (`/akuntansi/buku-besar`)
- Financial Reports (`/laporan`)

### Upcoming Features 📋

- POS/Sales (`/pos`)
- Products (`/produk`)
- Inventory Management
- Advanced Reports (Balance Sheet, P&L)
- Settings (`/pengaturan`)
- User Management

## Documentation

See main project docs:

- `/docs/mvp-action-plan.md` - Feature roadmap
- `/docs/quick-start-guide.md` - Setup guide
- `/CLAUDE.md` - Development guidelines
