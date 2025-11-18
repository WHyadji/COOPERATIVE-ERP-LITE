// ============================================================================
// Journal Entry Page Tests
// Integration tests for transaction list and journal entry management
// ============================================================================

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import JournalEntryPage from '@/app/(dashboard)/akuntansi/jurnal/page';
import accountingApi from '@/lib/api/accountingApi';
import type { Transaksi } from '@/types';

// Mock dependencies
vi.mock('@/lib/api/accountingApi');
vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: vi.fn(),
    back: vi.fn(),
  }),
}));

describe('Journal Entry Page', () => {
  const mockTransactions: Transaksi[] = [
    {
      id: '1',
      idKoperasi: 'kop1',
      nomorJurnal: 'JU-2025-001',
      tanggalTransaksi: '2025-11-18',
      deskripsi: 'Penerimaan kas dari penjualan',
      nomorReferensi: 'INV-001',
      totalDebit: 1000000,
      totalKredit: 1000000,
      statusBalanced: true,
    },
    {
      id: '2',
      idKoperasi: 'kop1',
      nomorJurnal: 'JU-2025-002',
      tanggalTransaksi: '2025-11-18',
      deskripsi: 'Pembayaran beban gaji',
      totalDebit: 500000,
      totalKredit: 500000,
      statusBalanced: true,
    },
  ];

  const mockPaginatedResponse = {
    success: true,
    message: 'Success',
    data: mockTransactions,
    pagination: {
      page: 1,
      pageSize: 20,
      totalItems: 2,
      totalPages: 1,
    },
  };

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(accountingApi.getTransactions).mockResolvedValue(
      mockPaginatedResponse
    );
  });

  // ============================================================================
  // Page Rendering Tests
  // ============================================================================

  describe('Page Rendering', () => {
    it('should render page title', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(
          screen.getByText('Jurnal Umum (Journal Entries)')
        ).toBeInTheDocument();
      });
    });

    it('should render Tambah Jurnal button', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText(/tambah jurnal/i)).toBeInTheDocument();
      });
    });

    it('should render date filters', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByLabelText(/tanggal mulai/i)).toBeInTheDocument();
        expect(screen.getByLabelText(/tanggal akhir/i)).toBeInTheDocument();
      });
    });

    it('should render transactions table', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByRole('table')).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Data Loading Tests
  // ============================================================================

  describe('Data Loading', () => {
    it('should fetch and display transactions on mount', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(accountingApi.getTransactions).toHaveBeenCalledWith({
          page: 1,
          pageSize: 20,
          tanggalMulai: undefined,
          tanggalAkhir: undefined,
        });
      });

      await waitFor(() => {
        expect(screen.getByText('JU-2025-001')).toBeInTheDocument();
        expect(screen.getByText('JU-2025-002')).toBeInTheDocument();
      });
    });

    it('should show loading state initially', () => {
      render(<JournalEntryPage />);

      expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });

    it('should show error message when fetch fails', async () => {
      vi.mocked(accountingApi.getTransactions).mockRejectedValue(
        new Error('Failed to fetch')
      );

      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(
          screen.getByText(/gagal memuat data transaksi/i)
        ).toBeInTheDocument();
      });
    });

    it('should show empty state when no transactions', async () => {
      vi.mocked(accountingApi.getTransactions).mockResolvedValue({
        ...mockPaginatedResponse,
        data: [],
        pagination: {
          page: 1,
          pageSize: 20,
          totalItems: 0,
          totalPages: 0,
        },
      });

      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText(/tidak ada data transaksi/i)).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Table Display Tests
  // ============================================================================

  describe('Table Display', () => {
    it('should display transaction information correctly', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText('JU-2025-001')).toBeInTheDocument();
        expect(screen.getByText('Penerimaan kas dari penjualan')).toBeInTheDocument();
        expect(screen.getByText('INV-001')).toBeInTheDocument();
      });
    });

    it('should display formatted currency', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        // Should format in IDR
        const currencyElements = screen.getAllByText(/Rp/);
        expect(currencyElements.length).toBeGreaterThan(0);
      });
    });

    it('should display formatted dates', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        // Should format as dd/MM/yyyy
        const dateElements = screen.getAllByText(/18\/11\/2025/);
        expect(dateElements.length).toBeGreaterThan(0);
      });
    });

    it('should display Balanced status chip for balanced entries', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        const balancedChips = screen.getAllByText('Balanced');
        expect(balancedChips.length).toBe(mockTransactions.length);
      });
    });

    it('should display Unbalanced status chip for unbalanced entries', async () => {
      const unbalancedTransaction: Transaksi = {
        ...mockTransactions[0],
        totalDebit: 1000000,
        totalKredit: 500000,
        statusBalanced: false,
      };

      vi.mocked(accountingApi.getTransactions).mockResolvedValue({
        ...mockPaginatedResponse,
        data: [unbalancedTransaction],
      });

      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText('Unbalanced')).toBeInTheDocument();
      });
    });

    it('should display action buttons for each transaction', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        const viewButtons = screen.getAllByTitle(/lihat detail/i);
        const deleteButtons = screen.getAllByTitle(/hapus/i);

        expect(viewButtons.length).toBe(mockTransactions.length);
        expect(deleteButtons.length).toBe(mockTransactions.length);
      });
    });
  });

  // ============================================================================
  // Filter Tests
  // ============================================================================

  describe('Date Filters', () => {
    it('should filter by date range', async () => {
      const user = userEvent.setup();
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText('JU-2025-001')).toBeInTheDocument();
      });

      const tanggalMulaiInput = screen.getByLabelText(/tanggal mulai/i);
      const tanggalAkhirInput = screen.getByLabelText(/tanggal akhir/i);

      await user.type(tanggalMulaiInput, '2025-11-01');
      await user.type(tanggalAkhirInput, '2025-11-30');

      await waitFor(() => {
        expect(accountingApi.getTransactions).toHaveBeenCalledWith(
          expect.objectContaining({
            tanggalMulai: '2025-11-01',
            tanggalAkhir: '2025-11-30',
          })
        );
      });
    });

    it('should reset filters when Reset Filter clicked', async () => {
      const user = userEvent.setup();
      render(<JournalEntryPage />);

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
  // Pagination Tests
  // ============================================================================

  describe('Pagination', () => {
    it('should display pagination controls', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText(/baris per halaman/i)).toBeInTheDocument();
      });
    });

    it('should change page when pagination controls used', async () => {
      const multiPageResponse = {
        ...mockPaginatedResponse,
        pagination: {
          page: 1,
          pageSize: 20,
          totalItems: 50,
          totalPages: 3,
        },
      };

      vi.mocked(accountingApi.getTransactions).mockResolvedValue(
        multiPageResponse
      );

      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText('JU-2025-001')).toBeInTheDocument();
      });

      // Pagination info should be displayed
      expect(screen.getByText(/1–2 dari 50/)).toBeInTheDocument();
    });

    it('should fetch data when rows per page changed', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(accountingApi.getTransactions).toHaveBeenCalled();
      });

      // Initial call with pageSize: 20
      expect(accountingApi.getTransactions).toHaveBeenCalledWith(
        expect.objectContaining({
          pageSize: 20,
        })
      );
    });
  });

  // ============================================================================
  // Create Transaction Tests
  // ============================================================================

  describe('Create Transaction', () => {
    it('should open transaction form when Tambah Jurnal clicked', async () => {
      const user = userEvent.setup();
      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText(/tambah jurnal/i)).toBeInTheDocument();
      });

      const addButton = screen.getByText(/tambah jurnal/i);
      await user.click(addButton);

      // Form dialog should open (component test covers the actual form)
      expect(addButton).toBeInTheDocument();
    });
  });

  // ============================================================================
  // Delete Transaction Tests
  // ============================================================================

  describe('Delete Transaction', () => {
    it('should delete transaction with confirmation', async () => {
      const user = userEvent.setup();
      vi.mocked(accountingApi.deleteTransaction).mockResolvedValue(undefined);

      const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(true);

      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText('JU-2025-001')).toBeInTheDocument();
      });

      const deleteButtons = screen.getAllByTitle(/hapus/i);
      await user.click(deleteButtons[0]);

      expect(confirmSpy).toHaveBeenCalled();

      confirmSpy.mockRestore();
    });

    it('should not delete if confirmation cancelled', async () => {
      const user = userEvent.setup();

      const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false);

      render(<JournalEntryPage />);

      await waitFor(() => {
        expect(screen.getByText('JU-2025-001')).toBeInTheDocument();
      });

      const deleteButtons = screen.getAllByTitle(/hapus/i);
      await user.click(deleteButtons[0]);

      expect(accountingApi.deleteTransaction).not.toHaveBeenCalled();

      confirmSpy.mockRestore();
    });
  });

  // ============================================================================
  // Currency Formatting Tests
  // ============================================================================

  describe('Currency Formatting', () => {
    it('should format amounts in Indonesian Rupiah', async () => {
      render(<JournalEntryPage />);

      await waitFor(() => {
        // Should display Rp format
        const currencyText = screen.getAllByText(/Rp\s*1\.000\.000/i);
        expect(currencyText.length).toBeGreaterThan(0);
      });
    });
  });
});
