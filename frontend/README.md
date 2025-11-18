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
│       ├── members/        # Member management
│       └── page.tsx        # Dashboard home
├── components/layout/      # Sidebar & Header
├── lib/
│   ├── api/               # API integration
│   └── context/           # Auth context
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

### ✅ Dashboard Layout
- Responsive sidebar navigation
- Mobile drawer
- User profile menu
- Role-based menu items

## Default Credentials (Development)

Contact backend team for test credentials.

## API Endpoints

All endpoints use `/api/v1` prefix:

- `POST /auth/login` - Authentication
- `GET /auth/profile` - User profile
- `GET /anggota` - List members
- `POST /anggota` - Create member
- `GET /anggota/:id` - Get member
- `PUT /anggota/:id` - Update member
- `DELETE /anggota/:id` - Delete member

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

Placeholder pages to implement:
- Savings Management (`/dashboard/simpanan`)
- POS/Sales (`/dashboard/pos`)
- Products (`/dashboard/products`)
- Accounting (`/dashboard/accounting`)
- Reports (`/dashboard/reports`)
- Settings (`/dashboard/settings`)

## Documentation

See main project docs:
- `/docs/mvp-action-plan.md` - Feature roadmap
- `/docs/quick-start-guide.md` - Setup guide
- `/CLAUDE.md` - Development guidelines
