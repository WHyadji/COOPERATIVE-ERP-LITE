// ============================================================================
// Accounting API Client Tests
// Tests for Chart of Accounts and Transaction API functions
// ============================================================================

import { describe, it, expect, vi, beforeEach } from "vitest";
import accountingApi from "@/lib/api/accountingApi";
import apiClient from "@/lib/api/client";
import type { Akun, Transaksi, TipeAkun, CreateAkunRequest } from "@/types";

// Mock axios client
vi.mock("@/lib/api/client");

describe("Accounting API Client", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  // ============================================================================
  // Chart of Accounts (Akun) Tests
  // ============================================================================

  describe("getAccounts", () => {
    it("should fetch all accounts successfully", async () => {
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
      ];

      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: mockAccounts,
        },
      });

      const result = await accountingApi.getAccounts();

      expect(apiClient.get).toHaveBeenCalledWith("/akun", { params: {} });
      expect(result).toEqual(mockAccounts);
      expect(result).toHaveLength(2);
    });

    it("should fetch accounts with type filter", async () => {
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
      ];

      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: mockAccounts,
        },
      });

      const result = await accountingApi.getAccounts("aset", undefined);

      expect(apiClient.get).toHaveBeenCalledWith("/akun", {
        params: { tipeAkun: "aset" },
      });
      expect(result).toEqual(mockAccounts);
    });

    it("should fetch accounts with status filter", async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: [],
        },
      });

      await accountingApi.getAccounts("all", true);

      expect(apiClient.get).toHaveBeenCalledWith("/akun", {
        params: { statusAktif: "true" },
      });
    });

    it("should throw error when API fails", async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: false,
          message: "Error",
          data: null,
        },
      });

      await expect(accountingApi.getAccounts()).rejects.toThrow(
        "Failed to fetch accounts"
      );
    });
  });

  describe("getAccountById", () => {
    it("should fetch account by ID successfully", async () => {
      const mockAccount: Akun = {
        id: "1",
        idKoperasi: "kop1",
        kodeAkun: "1101",
        namaAkun: "Kas",
        tipeAkun: "aset",
        normalSaldo: "debit",
        statusAktif: true,
      };

      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: mockAccount,
        },
      });

      const result = await accountingApi.getAccountById("1");

      expect(apiClient.get).toHaveBeenCalledWith("/akun/1");
      expect(result).toEqual(mockAccount);
      expect(result.kodeAkun).toBe("1101");
    });
  });

  describe("createAccount", () => {
    it("should create account successfully", async () => {
      const newAccount: CreateAkunRequest = {
        kodeAkun: "1101",
        namaAkun: "Kas",
        tipeAkun: "aset",
        normalSaldo: "debit",
      };

      const createdAccount: Akun = {
        id: "1",
        idKoperasi: "kop1",
        ...newAccount,
        statusAktif: true,
      };

      vi.mocked(apiClient.post).mockResolvedValue({
        data: {
          success: true,
          message: "Account created",
          data: createdAccount,
        },
      });

      const result = await accountingApi.createAccount(newAccount);

      expect(apiClient.post).toHaveBeenCalledWith("/akun", newAccount);
      expect(result).toEqual(createdAccount);
      expect(result.id).toBe("1");
    });
  });

  describe("updateAccount", () => {
    it("should update account successfully", async () => {
      const updateData = {
        kodeAkun: "1101",
        namaAkun: "Kas Updated",
        tipeAkun: "aset" as TipeAkun,
        normalSaldo: "debit" as const,
        statusAktif: false,
      };

      const updatedAccount: Akun = {
        id: "1",
        idKoperasi: "kop1",
        ...updateData,
      };

      vi.mocked(apiClient.put).mockResolvedValue({
        data: {
          success: true,
          message: "Account updated",
          data: updatedAccount,
        },
      });

      const result = await accountingApi.updateAccount("1", updateData);

      expect(apiClient.put).toHaveBeenCalledWith("/akun/1", updateData);
      expect(result.namaAkun).toBe("Kas Updated");
    });
  });

  describe("deleteAccount", () => {
    it("should delete account successfully", async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({
        data: {
          success: true,
          message: "Account deleted",
        },
      });

      await accountingApi.deleteAccount("1");

      expect(apiClient.delete).toHaveBeenCalledWith("/akun/1");
    });
  });

  describe("getAccountBalance", () => {
    it("should get account balance successfully", async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: { saldo: 1000000 },
        },
      });

      const balance = await accountingApi.getAccountBalance("1");

      expect(apiClient.get).toHaveBeenCalledWith("/akun/1/saldo", {
        params: {},
      });
      expect(balance).toBe(1000000);
    });

    it("should get account balance with date filter", async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: { saldo: 500000 },
        },
      });

      const balance = await accountingApi.getAccountBalance("1", "2025-11-18");

      expect(apiClient.get).toHaveBeenCalledWith("/akun/1/saldo", {
        params: { tanggalPer: "2025-11-18" },
      });
      expect(balance).toBe(500000);
    });
  });

  describe("seedDefaultCOA", () => {
    it("should seed default Chart of Accounts successfully", async () => {
      vi.mocked(apiClient.post).mockResolvedValue({
        data: {
          success: true,
          message: "COA seeded",
        },
      });

      await accountingApi.seedDefaultCOA();

      expect(apiClient.post).toHaveBeenCalledWith("/akun/seed-coa");
    });
  });

  // ============================================================================
  // Transaction (Transaksi) Tests
  // ============================================================================

  describe("getTransactions", () => {
    it("should fetch transactions with pagination", async () => {
      const mockTransactions: Transaksi[] = [
        {
          id: "1",
          idKoperasi: "kop1",
          nomorJurnal: "JU-2025-001",
          tanggalTransaksi: "2025-11-18",
          deskripsi: "Penerimaan kas",
          totalDebit: 1000000,
          totalKredit: 1000000,
          statusBalanced: true,
        },
      ];

      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: mockTransactions,
          pagination: {
            page: 1,
            pageSize: 20,
            totalItems: 1,
            totalPages: 1,
          },
        },
      });

      const result = await accountingApi.getTransactions({
        page: 1,
        pageSize: 20,
      });

      expect(apiClient.get).toHaveBeenCalledWith("/transaksi", {
        params: { page: 1, pageSize: 20 },
      });
      expect(result.data).toEqual(mockTransactions);
      expect(result.pagination.totalItems).toBe(1);
    });

    it("should fetch transactions with date filters", async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: [],
          pagination: {
            page: 1,
            pageSize: 20,
            totalItems: 0,
            totalPages: 0,
          },
        },
      });

      await accountingApi.getTransactions({
        tanggalMulai: "2025-11-01",
        tanggalAkhir: "2025-11-30",
      });

      expect(apiClient.get).toHaveBeenCalledWith("/transaksi", {
        params: {
          page: 1,
          pageSize: 20,
          tanggalMulai: "2025-11-01",
          tanggalAkhir: "2025-11-30",
        },
      });
    });
  });

  describe("createTransaction", () => {
    it("should create transaction with balanced entries", async () => {
      const newTransaction = {
        nomorJurnal: "JU-2025-001",
        tanggalTransaksi: "2025-11-18",
        deskripsi: "Penerimaan kas",
        barisTransaksi: [
          {
            idAkun: "1",
            jumlahDebit: 1000000,
            jumlahKredit: 0,
          },
          {
            idAkun: "2",
            jumlahDebit: 0,
            jumlahKredit: 1000000,
          },
        ],
      };

      const createdTransaction: Transaksi = {
        id: "1",
        idKoperasi: "kop1",
        ...newTransaction,
        totalDebit: 1000000,
        totalKredit: 1000000,
        statusBalanced: true,
      };

      vi.mocked(apiClient.post).mockResolvedValue({
        data: {
          success: true,
          message: "Transaction created",
          data: createdTransaction,
        },
      });

      const result = await accountingApi.createTransaction(newTransaction);

      expect(apiClient.post).toHaveBeenCalledWith("/transaksi", newTransaction);
      expect(result.statusBalanced).toBe(true);
      expect(result.totalDebit).toBe(result.totalKredit);
    });

    it("should throw error for unbalanced transaction", async () => {
      vi.mocked(apiClient.post).mockResolvedValue({
        data: {
          success: false,
          message: "Total debit harus sama dengan total kredit",
        },
      });

      await expect(
        accountingApi.createTransaction({
          nomorJurnal: "JU-2025-001",
          tanggalTransaksi: "2025-11-18",
          deskripsi: "Unbalanced entry",
          barisTransaksi: [
            {
              idAkun: "1",
              jumlahDebit: 1000000,
              jumlahKredit: 0,
            },
            {
              idAkun: "2",
              jumlahDebit: 0,
              jumlahKredit: 500000,
            },
          ],
        })
      ).rejects.toThrow();
    });
  });

  describe("getTransactionById", () => {
    it("should fetch transaction by ID with line items", async () => {
      const mockTransaction: Transaksi = {
        id: "1",
        idKoperasi: "kop1",
        nomorJurnal: "JU-2025-001",
        tanggalTransaksi: "2025-11-18",
        deskripsi: "Penerimaan kas",
        totalDebit: 1000000,
        totalKredit: 1000000,
        statusBalanced: true,
        barisTransaksi: [
          {
            id: "b1",
            idAkun: "1",
            kodeAkun: "1101",
            namaAkun: "Kas",
            jumlahDebit: 1000000,
            jumlahKredit: 0,
          },
          {
            id: "b2",
            idAkun: "2",
            kodeAkun: "4101",
            namaAkun: "Pendapatan",
            jumlahDebit: 0,
            jumlahKredit: 1000000,
          },
        ],
      };

      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: mockTransaction,
        },
      });

      const result = await accountingApi.getTransactionById("1");

      expect(apiClient.get).toHaveBeenCalledWith("/transaksi/1");
      expect(result.barisTransaksi).toHaveLength(2);
      expect(result.statusBalanced).toBe(true);
    });
  });

  describe("deleteTransaction", () => {
    it("should delete transaction successfully", async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({
        data: {
          success: true,
          message: "Transaction deleted",
        },
      });

      await accountingApi.deleteTransaction("1");

      expect(apiClient.delete).toHaveBeenCalledWith("/transaksi/1");
    });
  });

  // ============================================================================
  // Ledger Tests
  // ============================================================================

  describe("getAccountLedger", () => {
    it("should fetch account ledger successfully", async () => {
      const mockLedger = [
        {
          tanggal: "2025-11-18",
          nomorJurnal: "JU-2025-001",
          deskripsi: "Penerimaan kas",
          debit: 1000000,
          kredit: 0,
          saldo: 1000000,
        },
      ];

      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: mockLedger,
        },
      });

      const result = await accountingApi.getAccountLedger("1");

      expect(apiClient.get).toHaveBeenCalledWith("/laporan/buku-besar", {
        params: { idAkun: "1" },
      });
      expect(result).toEqual(mockLedger);
    });

    it("should fetch ledger with date range", async () => {
      vi.mocked(apiClient.get).mockResolvedValue({
        data: {
          success: true,
          message: "Success",
          data: [],
        },
      });

      await accountingApi.getAccountLedger("1", "2025-11-01", "2025-11-30");

      expect(apiClient.get).toHaveBeenCalledWith("/laporan/buku-besar", {
        params: {
          idAkun: "1",
          tanggalMulai: "2025-11-01",
          tanggalAkhir: "2025-11-30",
        },
      });
    });
  });
});
