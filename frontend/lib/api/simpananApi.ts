// ============================================================================
// Simpanan API - Operations for Share Capital (Simpanan)
// Type-safe API calls for share capital management
// ============================================================================

import apiClient from './client';
import type {
  Simpanan,
  CreateSimpananRequest,
  SaldoSimpananAnggota,
  RingkasanSimpanan,
  APIResponse,
  PaginatedResponse,
  SimpananListFilters,
} from '@/types';

// ============================================================================
// Simpanan API Functions
// ============================================================================

/**
 * Get paginated list of simpanan transactions
 */
export const getSimpananList = async (
  filters?: SimpananListFilters
): Promise<PaginatedResponse<Simpanan>> => {
  const params: Record<string, string | number> = {
    page: filters?.page || 1,
    pageSize: filters?.pageSize || 20,
  };

  if (filters?.tipeSimpanan && filters.tipeSimpanan !== 'all') {
    params.tipeSimpanan = filters.tipeSimpanan;
  }

  if (filters?.idAnggota) {
    params.idAnggota = filters.idAnggota;
  }

  if (filters?.tanggalMulai) {
    params.tanggalMulai = filters.tanggalMulai;
  }

  if (filters?.tanggalAkhir) {
    params.tanggalAkhir = filters.tanggalAkhir;
  }

  const response = await apiClient.get<PaginatedResponse<Simpanan>>('/simpanan', {
    params,
  });

  return response.data;
};

/**
 * Create new simpanan deposit
 */
export const createSimpanan = async (
  data: CreateSimpananRequest
): Promise<Simpanan> => {
  const response = await apiClient.post<APIResponse<Simpanan>>('/simpanan', data);

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to create simpanan');
  }

  return response.data.data;
};

/**
 * Get simpanan balance for a specific member
 */
export const getSaldoAnggota = async (
  idAnggota: string
): Promise<SaldoSimpananAnggota> => {
  const response = await apiClient.get<APIResponse<SaldoSimpananAnggota>>(
    `/simpanan/anggota/${idAnggota}/saldo`
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch member balance');
  }

  return response.data.data;
};

/**
 * Get simpanan summary (total by type)
 */
export const getRingkasan = async (): Promise<RingkasanSimpanan> => {
  const response = await apiClient.get<APIResponse<RingkasanSimpanan>>(
    '/simpanan/ringkasan'
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch simpanan summary');
  }

  return response.data.data;
};

/**
 * Get balance report for all members
 */
export const getLaporanSaldo = async (): Promise<SaldoSimpananAnggota[]> => {
  const response = await apiClient.get<APIResponse<SaldoSimpananAnggota[]>>(
    '/simpanan/laporan-saldo'
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch balance report');
  }

  return response.data.data;
};

// ============================================================================
// Export all simpanan API functions
// ============================================================================

const simpananApi = {
  getSimpananList,
  createSimpanan,
  getSaldoAnggota,
  getRingkasan,
  getLaporanSaldo,
};

export default simpananApi;
