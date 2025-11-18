// ============================================================================
// Reports API - Financial Reports Operations
// Type-safe API calls for generating financial reports
// ============================================================================

import apiClient from './client';
import type {
  LaporanPosisiKeuangan,
  LaporanLabaRugi,
  LaporanArusKas,
  NeracaSaldo,
  LaporanTransaksiHarian,
  BukuBesarResponse,
  APIResponse,
} from '@/types';

// ============================================================================
// Reports (Laporan) API Functions
// ============================================================================

/**
 * Get Balance Sheet (Neraca / Laporan Posisi Keuangan)
 * Shows assets, liabilities, and equity at a specific date
 *
 * @param tanggalPer - Optional date (YYYY-MM-DD), defaults to today
 */
export const getBalanceSheet = async (
  tanggalPer?: string
): Promise<LaporanPosisiKeuangan> => {
  const params: Record<string, string> = {};

  if (tanggalPer) {
    params.tanggalPer = tanggalPer;
  }

  const response = await apiClient.get<APIResponse<LaporanPosisiKeuangan>>(
    '/laporan/neraca',
    { params }
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch balance sheet');
  }

  return response.data.data;
};

/**
 * Get Income Statement (Laporan Laba Rugi)
 * Shows revenue, expenses, and net profit/loss for a period
 *
 * @param tanggalMulai - Start date (YYYY-MM-DD) - required
 * @param tanggalAkhir - End date (YYYY-MM-DD) - required
 */
export const getIncomeStatement = async (
  tanggalMulai: string,
  tanggalAkhir: string
): Promise<LaporanLabaRugi> => {
  if (!tanggalMulai || !tanggalAkhir) {
    throw new Error('Start date and end date are required');
  }

  const params = {
    tanggalMulai,
    tanggalAkhir,
  };

  const response = await apiClient.get<APIResponse<LaporanLabaRugi>>(
    '/laporan/laba-rugi',
    { params }
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch income statement');
  }

  return response.data.data;
};

/**
 * Get Cash Flow Statement (Laporan Arus Kas)
 * Shows cash flows from operating, investing, and financing activities
 *
 * @param tanggalMulai - Start date (YYYY-MM-DD) - required
 * @param tanggalAkhir - End date (YYYY-MM-DD) - required
 */
export const getCashFlow = async (
  tanggalMulai: string,
  tanggalAkhir: string
): Promise<LaporanArusKas> => {
  if (!tanggalMulai || !tanggalAkhir) {
    throw new Error('Start date and end date are required');
  }

  const params = {
    tanggalMulai,
    tanggalAkhir,
  };

  const response = await apiClient.get<APIResponse<LaporanArusKas>>(
    '/laporan/arus-kas',
    { params }
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch cash flow statement');
  }

  return response.data.data;
};

/**
 * Get Trial Balance (Neraca Saldo)
 * Shows all accounts with debit and credit balances
 * Serves as member balances report (includes all accounts)
 *
 * @param tanggalPer - Optional date (YYYY-MM-DD), defaults to today
 */
export const getTrialBalance = async (
  tanggalPer?: string
): Promise<NeracaSaldo> => {
  const params: Record<string, string> = {};

  if (tanggalPer) {
    params.tanggalPer = tanggalPer;
  }

  const response = await apiClient.get<APIResponse<NeracaSaldo>>(
    '/laporan/neraca-saldo',
    { params }
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch trial balance');
  }

  return response.data.data;
};

/**
 * Get Daily Transaction Report (Laporan Transaksi Harian)
 * Summary of cash flows and transaction counts for a specific day
 *
 * @param tanggal - Optional date (YYYY-MM-DD), defaults to today
 */
export const getDailyTransactionReport = async (
  tanggal?: string
): Promise<LaporanTransaksiHarian> => {
  const params: Record<string, string> = {};

  if (tanggal) {
    params.tanggal = tanggal;
  }

  const response = await apiClient.get<APIResponse<LaporanTransaksiHarian>>(
    '/laporan/transaksi-harian',
    { params }
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch daily transaction report');
  }

  return response.data.data;
};

/**
 * Get General Ledger (Buku Besar)
 * Detailed transaction history for a specific account
 *
 * @param idAkun - Account ID (required)
 * @param tanggalMulai - Optional start date (YYYY-MM-DD)
 * @param tanggalAkhir - Optional end date (YYYY-MM-DD)
 */
export const getGeneralLedger = async (
  idAkun: string,
  tanggalMulai?: string,
  tanggalAkhir?: string
): Promise<BukuBesarResponse> => {
  if (!idAkun) {
    throw new Error('Account ID is required');
  }

  const params: Record<string, string> = {
    idAkun,
  };

  if (tanggalMulai) {
    params.tanggalMulai = tanggalMulai;
  }

  if (tanggalAkhir) {
    params.tanggalAkhir = tanggalAkhir;
  }

  const response = await apiClient.get<APIResponse<BukuBesarResponse>>(
    '/laporan/buku-besar',
    { params }
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch general ledger');
  }

  return response.data.data;
};

// Export all functions as default object for convenience
export default {
  getBalanceSheet,
  getIncomeStatement,
  getCashFlow,
  getTrialBalance,
  getDailyTransactionReport,
  getGeneralLedger,
};
