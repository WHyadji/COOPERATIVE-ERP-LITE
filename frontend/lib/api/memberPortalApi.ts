// ============================================================================
// Member Portal API - API calls for member self-service portal
// Handles member balance, transaction history, and profile
// ============================================================================

import apiClient from './client';
import type { APIResponse, SaldoSimpananAnggota, Member, PaginatedResponse } from '@/types';

// ============================================================================
// Member Balance API
// ============================================================================

/**
 * Get member's share capital balance (Pokok, Wajib, Sukarela)
 */
export const getMemberBalance = async (): Promise<SaldoSimpananAnggota> => {
  const response = await apiClient.get<APIResponse<SaldoSimpananAnggota>>(
    '/portal/saldo'
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch member balance');
  }

  return response.data.data;
};

// ============================================================================
// Transaction History API
// ============================================================================

export interface RiwayatTransaksiAnggota {
  id: string;
  tanggalTransaksi: string;
  tipeSimpanan: 'pokok' | 'wajib' | 'sukarela';
  jumlah: number;
  keterangan: string;
  nomorReferensi: string;
}

export interface MemberTransactionFilters {
  tipeSimpanan?: 'pokok' | 'wajib' | 'sukarela' | 'all';
  tanggalMulai?: string;
  tanggalAkhir?: string;
  page?: number;
  pageSize?: number;
}

/**
 * Get member's transaction history
 */
export const getMemberTransactions = async (
  filters?: MemberTransactionFilters
): Promise<RiwayatTransaksiAnggota[]> => {
  const params: Record<string, string | number> = {
    page: filters?.page || 1,
    pageSize: filters?.pageSize || 20,
  };

  const response = await apiClient.get<PaginatedResponse<RiwayatTransaksiAnggota>>(
    '/portal/riwayat',
    { params }
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch transaction history');
  }

  return response.data.data;
};

// ============================================================================
// Member Profile API
// ============================================================================

/**
 * Get current member's profile
 */
export const getMemberProfile = async (): Promise<Member> => {
  const response = await apiClient.get<APIResponse<Member>>(
    '/portal/profile'
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch member profile');
  }

  return response.data.data;
};

/**
 * Update member profile (limited fields)
 * Note: Member portal uses PIN authentication, profile updates may not be available
 */
export const updateMemberProfile = async (
  data: Partial<Member>
): Promise<Member> => {
  // This endpoint may not exist in backend - members typically cannot update their own profile
  // For now, we'll keep this but it should be reviewed
  const response = await apiClient.put<APIResponse<Member>>(
    '/portal/profile',
    data
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to update member profile');
  }

  return response.data.data;
};

// ============================================================================
// Member Dashboard Summary
// ============================================================================

export interface MemberDashboardSummary {
  saldoSimpanan: SaldoSimpananAnggota;
  transaksiTerbaru: RiwayatTransaksiAnggota[];
  totalTransaksi: number;
}

/**
 * Get member dashboard summary
 */
export const getMemberDashboard = async (): Promise<MemberDashboardSummary> => {
  // Backend doesn't have a dedicated dashboard endpoint
  // We'll fetch balance and recent transactions separately
  const [saldoSimpanan, transaksiResponse] = await Promise.all([
    getMemberBalance(),
    getMemberTransactions({ page: 1, pageSize: 5 }),
  ]);

  return {
    saldoSimpanan,
    transaksiTerbaru: transaksiResponse,
    totalTransaksi: transaksiResponse.length,
  };
};
