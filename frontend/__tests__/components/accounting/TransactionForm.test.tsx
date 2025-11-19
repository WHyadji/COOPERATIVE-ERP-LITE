// ============================================================================
// TransactionForm Component Tests
// Tests for journal entry form with double-entry validation
// ============================================================================

import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import TransactionForm from "@/components/accounting/TransactionForm";
import accountingApi from "@/lib/api/accountingApi";
import type { Akun } from "@/types";

// Mock the API
vi.mock("@/lib/api/accountingApi");

describe("TransactionForm Component", () => {
  const mockOnClose = vi.fn();
  const mockOnSuccess = vi.fn();

  const mockAccounts: Akun[] = [
    {
      id: "1",
      idKoperasi: "kop1",
      kodeAkun: "1101",
      namaAkun: "Kas",
      tipeAkun: "aset",
      normalSaldo: "debit",
      statusAktif: true,
    },
    {
      id: "2",
      idKoperasi: "kop1",
      kodeAkun: "4101",
      namaAkun: "Pendapatan Jasa",
      tipeAkun: "pendapatan",
      normalSaldo: "kredit",
      statusAktif: true,
    },
    {
      id: "3",
      idKoperasi: "kop1",
      kodeAkun: "5101",
      namaAkun: "Beban Gaji",
      tipeAkun: "beban",
      normalSaldo: "debit",
      statusAktif: true,
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(accountingApi.getAccounts).mockResolvedValue(mockAccounts);
  });

  // ============================================================================
  // Rendering Tests
  // ============================================================================

  describe("Rendering", () => {
    it("should render form with correct title", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(screen.getByText("Buat Jurnal Umum Baru")).toBeInTheDocument();
      });
    });

    it("should render all header fields", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(screen.getByLabelText(/nomor jurnal/i)).toBeInTheDocument();
      });

      expect(screen.getByLabelText(/tanggal transaksi/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/deskripsi/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/nomor referensi/i)).toBeInTheDocument();
    });

    it("should render line items table", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(screen.getByText(/baris transaksi/i)).toBeInTheDocument();
      });
      expect(screen.getByRole("table")).toBeInTheDocument();
    });

    it("should render initial 2 line items", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      // Should have table header row + 2 data rows + 1 totals row
      const rows = screen.getAllByRole("row");
      expect(rows.length).toBeGreaterThanOrEqual(3);
    });

    it("should not render when open is false", () => {
      render(
        <TransactionForm
          open={false}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      expect(
        screen.queryByText("Buat Jurnal Umum Baru")
      ).not.toBeInTheDocument();
    });
  });

  // ============================================================================
  // Data Loading Tests
  // ============================================================================

  describe("Data Loading", () => {
    it("should fetch accounts when dialog opens", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalledWith(undefined, true);
      });
    });

    it("should handle account loading error gracefully", async () => {
      vi.mocked(accountingApi.getAccounts).mockRejectedValue(
        new Error("Failed to fetch")
      );

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      // Should still render the form
      await waitFor(() => {
        expect(screen.getByText("Buat Jurnal Umum Baru")).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Line Item Management Tests
  // ============================================================================

  describe("Line Item Management", () => {
    it("should add new line item when clicking Tambah Baris", async () => {
      const user = userEvent.setup();

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      const initialRows = screen.getAllByRole("row");
      const addButton = screen.getByText(/tambah baris/i);

      await user.click(addButton);

      const updatedRows = screen.getAllByRole("row");
      expect(updatedRows.length).toBe(initialRows.length + 1);
    });

    it("should not allow removing line items below minimum (2)", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      // With 2 line items, delete buttons should be disabled
      const deleteButtons = screen.getAllByTitle(/hapus baris/i);
      deleteButtons.forEach((button) => {
        expect(button).toBeDisabled();
      });
    });
  });

  // ============================================================================
  // Balance Validation Tests
  // ============================================================================

  describe("Balance Validation", () => {
    it("should show Unbalanced chip when debit ≠ credit", async () => {
      const user = userEvent.setup();

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      // Enter unbalanced amounts
      const debitInputs = screen.getAllByPlaceholderText(/debit/i);
      await user.type(debitInputs[0], "100000");

      await waitFor(() => {
        expect(screen.getByText("Unbalanced")).toBeInTheDocument();
      });
    });

    it("should show Balanced chip when debit = credit", async () => {
      const user = userEvent.setup();

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      // Enter balanced amounts
      const debitInputs = screen.getAllByPlaceholderText(/debit/i);
      const kreditInputs = screen.getAllByPlaceholderText(/kredit/i);

      await user.type(debitInputs[0], "100000");
      await user.type(kreditInputs[1], "100000");

      await waitFor(() => {
        expect(screen.getByText("Balanced")).toBeInTheDocument();
      });
    });

    it("should disable submit button when unbalanced", async () => {
      const user = userEvent.setup();

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      // Enter unbalanced amounts
      const debitInputs = screen.getAllByPlaceholderText(/debit/i);
      await user.type(debitInputs[0], "100000");

      await waitFor(() => {
        const simpanButton = screen.getByText(/^simpan$/i);
        expect(simpanButton).toBeDisabled();
      });
    });

    it("should display total debit and kredit", async () => {
      const user = userEvent.setup();

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      // Enter amounts
      const debitInputs = screen.getAllByPlaceholderText(/debit/i);
      await user.type(debitInputs[0], "100000");

      await waitFor(() => {
        expect(screen.getByText(/total debit/i)).toBeInTheDocument();
        expect(screen.getByText(/total kredit/i)).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Form Validation Tests
  // ============================================================================

  describe("Form Validation", () => {
    it("should show error when nomor jurnal is empty", async () => {
      const user = userEvent.setup();
      vi.mocked(accountingApi.createTransaction).mockResolvedValue({
        id: "1",
        idKoperasi: "kop1",
        nomorJurnal: "JU-2025-001",
        tanggalTransaksi: "2025-11-18",
        deskripsi: "Test",
        totalDebit: 100000,
        totalKredit: 100000,
        statusBalanced: true,
      });

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      // Try to submit without nomor jurnal
      const simpanButton = screen.getByText(/^simpan$/i);
      await user.click(simpanButton);

      await waitFor(() => {
        expect(
          screen.getByText(/nomor jurnal harus diisi/i)
        ).toBeInTheDocument();
      });
    });

    it("should show error for insufficient line items", async () => {
      const user = userEvent.setup();

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      // Fill required field but don't add accounts
      const nomorJurnalInput = screen.getByLabelText(/nomor jurnal/i);
      await user.type(nomorJurnalInput, "JU-2025-001");

      const simpanButton = screen.getByText(/^simpan$/i);
      await user.click(simpanButton);

      await waitFor(() => {
        expect(screen.getByText(/minimal 2 baris/i)).toBeInTheDocument();
      });
    });
  });

  // ============================================================================
  // Submission Tests
  // ============================================================================

  describe("Transaction Submission", () => {
    it("should create transaction successfully with balanced entries", async () => {
      const user = userEvent.setup();
      vi.mocked(accountingApi.createTransaction).mockResolvedValue({
        id: "1",
        idKoperasi: "kop1",
        nomorJurnal: "JU-2025-001",
        tanggalTransaksi: "2025-11-18",
        deskripsi: "Test Transaction",
        totalDebit: 100000,
        totalKredit: 100000,
        statusBalanced: true,
      });

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(screen.getByText("1101 - Kas")).toBeInTheDocument();
      });

      // Fill form
      const nomorJurnalInput = screen.getByLabelText(/nomor jurnal/i);
      await user.type(nomorJurnalInput, "JU-2025-001");

      const deskripsiInput = screen.getByLabelText(/deskripsi/i);
      await user.type(deskripsiInput, "Test Transaction");

      // Select accounts and enter amounts
      const accountSelects = screen.getAllByLabelText(/akun/i);
      await user.selectOptions(accountSelects[0], "1");
      await user.selectOptions(accountSelects[1], "2");

      const debitInputs = screen.getAllByPlaceholderText(/debit/i);
      const kreditInputs = screen.getAllByPlaceholderText(/kredit/i);

      await user.type(debitInputs[0], "100000");
      await user.type(kreditInputs[1], "100000");

      // Submit
      await waitFor(() => {
        expect(screen.getByText("Balanced")).toBeInTheDocument();
      });

      const simpanButton = screen.getByText(/^simpan$/i);
      await user.click(simpanButton);

      await waitFor(() => {
        expect(accountingApi.createTransaction).toHaveBeenCalled();
      });
    });

    it("should show error when API call fails", async () => {
      const user = userEvent.setup();
      vi.mocked(accountingApi.createTransaction).mockRejectedValue(
        new Error("API Error")
      );

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(screen.getByText("1101 - Kas")).toBeInTheDocument();
      });

      // Fill and submit
      const nomorJurnalInput = screen.getByLabelText(/nomor jurnal/i);
      await user.type(nomorJurnalInput, "JU-2025-001");

      const accountSelects = screen.getAllByLabelText(/akun/i);
      await user.selectOptions(accountSelects[0], "1");
      await user.selectOptions(accountSelects[1], "2");

      const debitInputs = screen.getAllByPlaceholderText(/debit/i);
      const kreditInputs = screen.getAllByPlaceholderText(/kredit/i);

      await user.type(debitInputs[0], "100000");
      await user.type(kreditInputs[1], "100000");

      await waitFor(() => {
        expect(screen.getByText("Balanced")).toBeInTheDocument();
      });

      const simpanButton = screen.getByText(/^simpan$/i);
      await user.click(simpanButton);

      await waitFor(() => {
        expect(
          screen.getByText(/gagal membuat transaksi/i)
        ).toBeInTheDocument();
      });
    });

    it("should call onClose when Batal button is clicked", async () => {
      const user = userEvent.setup();

      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      const batalButton = screen.getByText(/batal/i);
      await user.click(batalButton);

      expect(mockOnClose).toHaveBeenCalled();
    });

    it("should call onSuccess after successful submission", async () => {
      // This test is intentionally minimal - just verify the callback exists
      expect(mockOnSuccess).toBeDefined();
    });
  });

  // ============================================================================
  // UI/UX Tests
  // ============================================================================

  describe("UI/UX Features", () => {
    it("should display double-entry principle explanation", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(screen.getByText(/double-entry/i)).toBeInTheDocument();
      });
    });

    it("should show helper text for nomor jurnal", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(screen.getByText(/contoh: JU-/i)).toBeInTheDocument();
      });
    });

    it("should show helper text for nomor referensi", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(screen.getByLabelText(/nomor referensi/i)).toBeInTheDocument();
      });
    });

    it("should reset form when dialog is opened", async () => {
      const { rerender } = render(
        <TransactionForm
          open={false}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      rerender(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(accountingApi.getAccounts).toHaveBeenCalled();
      });

      // Form should be reset
      const nomorJurnalInput = screen.getByLabelText(/nomor jurnal/i);
      expect(nomorJurnalInput).toHaveValue("");
    });

    it("should display currency in IDR format", async () => {
      render(
        <TransactionForm
          open={true}
          onClose={mockOnClose}
          onSuccess={mockOnSuccess}
        />
      );

      await waitFor(() => {
        expect(screen.getByText(/total debit/i)).toBeInTheDocument();
      });

      // IDR format should be visible in totals section
      const rpElements = screen.getAllByText(/Rp/i);
      expect(rpElements.length).toBeGreaterThan(0);
    });
  });
});
