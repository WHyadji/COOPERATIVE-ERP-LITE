// ============================================================================
// Account Ledger Page Tests
// Integration tests for account ledger (buku besar) view
// ============================================================================

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AccountLedgerPage from '@/app/(dashboard)/akuntansi/ledger/[id]/page';
import accountingApi from '@/lib/api/accountingApi';
import type { Akun } from '@/types';

// Mock dependencies
vi.mock('@/lib/api/accountingApi');
vi.mock('next/navigation', () => ({
  useParams: () => ({
    id: '1',
  }),
  useRouter: () => ({
    push: vi.fn(),
    back: vi.fn(),
  }),
}));

describe('Account Ledger Page', () => {
  const mockAccount: Akun = {
    id: '1',
    idKoperasi: 'kop1',
    kodeAkun: '1101',
    namaAkun: 'Kas',
    tipeAkun: 'aset',
    normalSaldo: 'debit',
    statusAktif: true,
  };

  const mockLedgerEntries = [
    {
      tanggal: '2025-11-15',
      nomorJurnal: 'JU-2025-001',
      deskripsi: 'Penerimaan kas awal',
      debit: 5000000,
      kredit: 0,
      saldo: 5000000,
    },
    {
      tanggal: '2025-11-16',
      nomorJurnal: 'JU-2025-002',
      deskripsi: 'Pembayaran beban gaji',
      debit: 0,
      kredit: 2000000,
      saldo: 3000000,
    },
    {
      tanggal: '2025-11-18',
      nomorJurnal: 'JU-2025-003',
      deskripsi: 'Penerimaan dari penjualan',
      debit: 1000000,
      kredit: 0,
      saldo: 4000000,
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(accountingApi.getAccountById).mockResolvedValue(mockAccount);
    vi.mocked(accountingApi.getAccountLedger).mockResolvedValue(mockLedgerEntries);
  });

  // ============================================================================
  // Page Rendering Tests
  // ============================================================================

  describe('Page Rendering', () => {
    it('should render page title', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Buku Besar (General Ledger)')).toBeInTheDocument();
      });
    });

    it('should render breadcrumbs', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Bagan Akun')).toBeInTheDocument();
        expect(screen.getByText('Buku Besar')).toBeInTheDocument();
      });
    });

    it('should display account information', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText(/1101 - Kas/)).toBeInTheDocument();
        expect(screen.getByText(/Normal: DEBIT/)).toBeInTheDocument();
      });
    });

    it('should render date filters', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByLabelText(/tanggal mulai/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/tanggal akhir/i)).toBeInTheDocument();
      });
    });

    it('should render back button', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Kembali')).toBeInTheDocument();
      });
    });

    it('should render print button', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Cetak')).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Data Loading Tests
  // ============================================================================

  describe('Data Loading', () => {
    it('should fetch account details on mount', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(accountingApi.getAccountById).toHaveBeenCalledWith('1');
      });
    });

    it('should fetch ledger entries on mount', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(accountingApi.getAccountLedger).toHaveBeenCalledWith(
          '1',
          undefined,
          undefined
        );
      });
    });

    it('should show loading state initially', () => {
      render(<AccountLedgerPage />);

      expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });

    it('should show error message when fetch fails', async () => {
      vi.mocked(accountingApi.getAccountLedger).mockRejectedValue(
        new Error('Failed to fetch')
      );

      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(
          screen.getByText(/gagal memuat data buku besar/i)
        ).toBeInTheDocument();
      });
    });

    it('should show empty state when no entries', async () => {
      vi.mocked(accountingApi.getAccountLedger).mockResolvedValue([]);

      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(
          screen.getByText(/tidak ada transaksi untuk akun ini/i)
        ).toBeInTheDocument();
      });
    });

    it('should show error when account not found', async () => {
      vi.mocked(accountingApi.getAccountById).mockRejectedValue(
        new Error('Account not found')
      );

      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText(/akun tidak ditemukan/i)).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Ledger Table Tests
  // ============================================================================

  describe('Ledger Table', () => {
    it('should display all ledger entries', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('JU-2025-001')).toBeInTheDocument();
        expect(screen.getByText('JU-2025-002')).toBeInTheDocument();
        expect(screen.getByText('JU-2025-003')).toBeInTheDocument();
      });
    });

    it('should display transaction descriptions', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Penerimaan kas awal')).toBeInTheDocument();
        expect(screen.getByText('Pembayaran beban gaji')).toBeInTheDocument();
        expect(screen.getByText('Penerimaan dari penjualan')).toBeInTheDocument();
      });
    });

    it('should display formatted dates', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('15/11/2025')).toBeInTheDocument();
        expect(screen.getByText('16/11/2025')).toBeInTheDocument();
        expect(screen.getByText('18/11/2025')).toBeInTheDocument();
      });
    });

    it('should display debit amounts correctly', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        const debitCells = screen.getAllByText(/Rp\s*5\.000\.000/i);
        expect(debitCells.length).toBeGreaterThan(0);
      });
    });

    it('should display kredit amounts correctly', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        const kreditCells = screen.getAllByText(/Rp\s*2\.000\.000/i);
        expect(kreditCells.length).toBeGreaterThan(0);
      });
    });

    it('should display running balance correctly', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText(/Rp\s*4\.000\.000/i)).toBeInTheDocument();
      });
    });

    it('should show "-" for zero amounts', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        const dashes = screen.getAllByText('-');
        expect(dashes.length).toBeGreaterThan(0);
      });
    });
  });

  // ============================================================================
  // Totals Display Tests
  // ============================================================================

  describe('Totals Display', () => {
    it('should display TOTAL row', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('TOTAL')).toBeInTheDocument();
      });
    });

    it('should calculate and display total debit', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        // Total debit = 5,000,000 + 1,000,000 = 6,000,000
        expect(screen.getByText(/Rp\s*6\.000\.000/i)).toBeInTheDocument();
      });
    });

    it('should calculate and display total kredit', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        // Total kredit = 2,000,000
        expect(screen.getByText(/Rp\s*2\.000\.000/i)).toBeInTheDocument();
      });
    });

    it('should display final balance', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        // Final balance = 4,000,000
        expect(screen.getByText(/Rp\s*4\.000\.000/i)).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Summary Section Tests
  // ============================================================================

  describe('Summary Section', () => {
    it('should display summary section', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Ringkasan')).toBeInTheDocument();
      });
    });

    it('should show total debit in summary', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Total Debit')).toBeInTheDocument();
      });
    });

    it('should show total kredit in summary', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Total Kredit')).toBeInTheDocument();
      });
    });

    it('should show saldo akhir in summary', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Saldo Akhir')).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Filter Tests
  // ============================================================================

  describe('Date Filters', () => {
    it('should filter ledger by date range', async () => {
      const user = userEvent.setup();
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('JU-2025-001')).toBeInTheDocument();
      });

      const tanggalMulaiInput = screen.getByLabelText(/tanggal mulai/i);
      const tanggalAkhirInput = screen.getByLabelText(/tanggal akhir/i);

      await user.type(tanggalMulaiInput, '2025-11-01');
      await user.type(tanggalAkhirInput, '2025-11-30');

      await waitFor(() => {
        expect(accountingApi.getAccountLedger).toHaveBeenCalledWith(
          '1',
          '2025-11-01',
          '2025-11-30'
        );
      });
    });

    it('should reset filters when Reset Filter clicked', async () => {
      const user = userEvent.setup();
      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText(/reset filter/i)).toBeInTheDocument();
      });

      const resetButton = screen.getByText(/reset filter/i);
      await user.click(resetButton);

      const tanggalMulaiInput = screen.getByLabelText(/tanggal mulai/i);
      const tanggalAkhirInput = screen.getByLabelText(/tanggal akhir/i);

      expect(tanggalMulaiInput).toHaveValue('');
      expect(tanggalAkhirInput).toHaveValue('');
    });
  });

  // ============================================================================
  // Navigation Tests
  // ============================================================================

  describe('Navigation', () => {
    it('should have back button', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        const backButton = screen.getByText('Kembali');
        expect(backButton).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Print Functionality Tests
  // ============================================================================

  describe('Print Functionality', () => {
    it('should call window.print when Cetak clicked', async () => {
      const user = userEvent.setup();
      const printSpy = vi.spyOn(window, 'print').mockImplementation(() => {});

      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText('Cetak')).toBeInTheDocument();
      });

      const printButton = screen.getByText('Cetak');
      await user.click(printButton);

      expect(printSpy).toHaveBeenCalled();

      printSpy.mockRestore();
    });
  });

  // ============================================================================
  // Currency Formatting Tests
  // ============================================================================

  describe('Currency Formatting', () => {
    it('should format amounts in Indonesian Rupiah', async () => {
      render(<AccountLedgerPage />);

      await waitFor(() => {
        const currencyText = screen.getAllByText(/Rp/i);
        expect(currencyText.length).toBeGreaterThan(0);
      });
    });

    it('should handle negative balances with (CR) indicator', async () => {
      const negativeBalanceLedger = [
        {
          tanggal: '2025-11-15',
          nomorJurnal: 'JU-2025-001',
          deskripsi: 'Test',
          debit: 0,
          kredit: 1000000,
          saldo: -1000000,
        },
      ];

      vi.mocked(accountingApi.getAccountLedger).mockResolvedValue(
        negativeBalanceLedger
      );

      render(<AccountLedgerPage />);

      await waitFor(() => {
        expect(screen.getByText(/(CR)/)).toBeInTheDocument();
      });
    });
  });
});
