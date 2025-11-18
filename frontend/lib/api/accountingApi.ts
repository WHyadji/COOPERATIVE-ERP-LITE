// ============================================================================
// Accounting API - CRUD Operations for Accounts & Transactions
// Type-safe API calls for accounting management
// ============================================================================

import apiClient from './client';
import type {
  Akun,
  CreateAkunRequest,
  UpdateAkunRequest,
  Transaksi,
  CreateTransaksiRequest,
  TransaksiListFilters,
  APIResponse,
  PaginatedResponse,
  TipeAkun,
  LedgerEntry,
} from '@/types';

// ============================================================================
// Chart of Accounts (Akun) API Functions
// ============================================================================

/**
 * Get all accounts with optional filters
 */
export const getAccounts = async (
  tipeAkun?: TipeAkun | 'all',
  statusAktif?: boolean
): Promise<Akun[]> => {
  const params: Record<string, string> = {};

  if (tipeAkun && tipeAkun !== 'all') {
    params.tipeAkun = tipeAkun;
  }

  if (statusAktif !== undefined) {
    params.statusAktif = statusAktif.toString();
  }

  const response = await apiClient.get<APIResponse<Akun[]>>('/akun', {
    params,
  });

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch accounts');
  }

  return response.data.data;
};

/**
 * Get account by ID
 */
export const getAccountById = async (id: string): Promise<Akun> => {
  const response = await apiClient.get<APIResponse<Akun>>(`/akun/${id}`);

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch account');
  }

  return response.data.data;
};

/**
 * Get account balance
 */
export const getAccountBalance = async (
  id: string,
  tanggalPer?: string
): Promise<number> => {
  const params: Record<string, string> = {};

  if (tanggalPer) {
    params.tanggalPer = tanggalPer;
  }

  const response = await apiClient.get<APIResponse<{ saldo: number }>>(
    `/akun/${id}/saldo`,
    { params }
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch account balance');
  }

  return response.data.data.saldo;
};

/**
 * Create new account
 */
export const createAccount = async (data: CreateAkunRequest): Promise<Akun> => {
  const response = await apiClient.post<APIResponse<Akun>>('/akun', data);

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to create account');
  }

  return response.data.data;
};

/**
 * Update existing account
 */
export const updateAccount = async (
  id: string,
  data: UpdateAkunRequest
): Promise<Akun> => {
  const response = await apiClient.put<APIResponse<Akun>>(`/akun/${id}`, data);

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to update account');
  }

  return response.data.data;
};

/**
 * Delete account
 */
export const deleteAccount = async (id: string): Promise<void> => {
  const response = await apiClient.delete<APIResponse>(`/akun/${id}`);

  if (!response.data.success) {
    throw new Error('Failed to delete account');
  }
};

/**
 * Seed default Chart of Accounts
 */
export const seedDefaultCOA = async (): Promise<void> => {
  const response = await apiClient.post<APIResponse>('/akun/seed-coa');

  if (!response.data.success) {
    throw new Error('Failed to seed Chart of Accounts');
  }
};

// ============================================================================
// Transaction (Transaksi) API Functions
// ============================================================================

/**
 * Get paginated list of transactions
 */
export const getTransactions = async (
  filters?: TransaksiListFilters
): Promise<PaginatedResponse<Transaksi>> => {
  const params: Record<string, string | number> = {
    page: filters?.page || 1,
    pageSize: filters?.pageSize || 20,
  };

  if (filters?.tanggalMulai) {
    params.tanggalMulai = filters.tanggalMulai;
  }

  if (filters?.tanggalAkhir) {
    params.tanggalAkhir = filters.tanggalAkhir;
  }

  if (filters?.tipeTransaksi) {
    params.tipeTransaksi = filters.tipeTransaksi;
  }

  const response = await apiClient.get<PaginatedResponse<Transaksi>>(
    '/transaksi',
    {
      params,
    }
  );

  return response.data;
};

/**
 * Get transaction by ID
 */
export const getTransactionById = async (id: string): Promise<Transaksi> => {
  const response = await apiClient.get<APIResponse<Transaksi>>(
    `/transaksi/${id}`
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch transaction');
  }

  return response.data.data;
};

/**
 * Create new transaction (journal entry)
 */
export const createTransaction = async (
  data: CreateTransaksiRequest
): Promise<Transaksi> => {
  const response = await apiClient.post<APIResponse<Transaksi>>(
    '/transaksi',
    data
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to create transaction');
  }

  return response.data.data;
};

/**
 * Update existing transaction (journal entry)
 */
export const updateTransaction = async (
  id: string,
  data: CreateTransaksiRequest
): Promise<Transaksi> => {
  const response = await apiClient.put<APIResponse<Transaksi>>(
    `/transaksi/${id}`,
    data
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to update transaction');
  }

  return response.data.data;
};

/**
 * Delete transaction
 */
export const deleteTransaction = async (id: string): Promise<void> => {
  const response = await apiClient.delete<APIResponse>(`/transaksi/${id}`);

  if (!response.data.success) {
    throw new Error('Failed to delete transaction');
  }
};

/**
 * Get account ledger (buku besar)
 */
export const getAccountLedger = async (
  idAkun: string,
  tanggalMulai?: string,
  tanggalAkhir?: string
): Promise<LedgerEntry[]> => {
  const params: Record<string, string> = {
    idAkun,
  };

  if (tanggalMulai) {
    params.tanggalMulai = tanggalMulai;
  }

  if (tanggalAkhir) {
    params.tanggalAkhir = tanggalAkhir;
  }

  const response = await apiClient.get<APIResponse<LedgerEntry[]>>('/laporan/buku-besar', {
    params,
  });

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch account ledger');
  }

  return response.data.data;
};

// ============================================================================
// Export all accounting API functions
// ============================================================================

const accountingApi = {
  // Account functions
  getAccounts,
  getAccountById,
  getAccountBalance,
  createAccount,
  updateAccount,
  deleteAccount,
  seedDefaultCOA,

  // Transaction functions
  getTransactions,
  getTransactionById,
  createTransaction,
  updateTransaction,
  deleteTransaction,

  // Ledger functions
  getAccountLedger,
};

export default accountingApi;
