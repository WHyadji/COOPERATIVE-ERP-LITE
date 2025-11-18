// ============================================================================
// POS API - Sales Transaction Operations
// Type-safe API calls for point of sale and sales management
// ============================================================================

import apiClient from './client';
import type {
  CreatePenjualanRequest,
  PenjualanResponse,
  Penjualan,
  PenjualanListFilters,
  StrukPenjualan,
  RingkasanPenjualanHariIni,
  TopProduk,
  APIResponse,
  PaginatedResponse,
} from '@/types';

// ============================================================================
// POS / Sales (Penjualan) API Functions
// ============================================================================

/**
 * Create a new sale transaction
 * Processes POS sale, reduces stock, and posts to accounting automatically
 */
export const createSale = async (
  data: CreatePenjualanRequest
): Promise<PenjualanResponse> => {
  const response = await apiClient.post<APIResponse<PenjualanResponse>>(
    '/penjualan',
    data
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to create sale');
  }

  return response.data.data;
};

/**
 * Get all sales with pagination and filters
 */
export const getSales = async (
  filters?: PenjualanListFilters
): Promise<PaginatedResponse<Penjualan>> => {
  const params: Record<string, string> = {};

  if (filters?.search) {
    params.search = filters.search;
  }

  if (filters?.tanggalMulai) {
    params.tanggalMulai = filters.tanggalMulai;
  }

  if (filters?.tanggalAkhir) {
    params.tanggalAkhir = filters.tanggalAkhir;
  }

  if (filters?.idAnggota) {
    params.idAnggota = filters.idAnggota;
  }

  if (filters?.idKasir) {
    params.idKasir = filters.idKasir;
  }

  if (filters?.page) {
    params.page = filters.page.toString();
  }

  if (filters?.pageSize) {
    params.pageSize = filters.pageSize.toString();
  }

  const response = await apiClient.get<PaginatedResponse<Penjualan>>(
    '/penjualan',
    { params }
  );

  if (!response.data.success) {
    throw new Error('Failed to fetch sales');
  }

  return response.data;
};

/**
 * Get sale by ID
 */
export const getSaleById = async (id: string): Promise<Penjualan> => {
  const response = await apiClient.get<APIResponse<Penjualan>>(
    `/penjualan/${id}`
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch sale');
  }

  return response.data.data;
};

/**
 * Get digital receipt for a sale
 */
export const getReceipt = async (id: string): Promise<StrukPenjualan> => {
  const response = await apiClient.get<APIResponse<StrukPenjualan>>(
    `/penjualan/${id}/struk`
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch receipt');
  }

  return response.data.data;
};

/**
 * Get today's sales summary
 */
export const getTodaySummary = async (): Promise<RingkasanPenjualanHariIni> => {
  const response = await apiClient.get<APIResponse<RingkasanPenjualanHariIni>>(
    '/penjualan/hari-ini'
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch today\'s summary');
  }

  return response.data.data;
};

/**
 * Get top selling products
 */
export const getTopProducts = async (
  limit?: number
): Promise<TopProduk[]> => {
  const params: Record<string, string> = {};

  if (limit) {
    params.limit = limit.toString();
  }

  const response = await apiClient.get<APIResponse<TopProduk[]>>(
    '/penjualan/top-produk',
    { params }
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch top products');
  }

  return response.data.data;
};

// Export all functions as a single object (optional, for convenience)
const posApi = {
  createSale,
  getSales,
  getSaleById,
  getReceipt,
  getTodaySummary,
  getTopProducts,
};

export default posApi;
