// ============================================================================
// Member Portal API - API calls for member self-service portal
// Handles member balance, transaction history, and profile
// ============================================================================

import apiClient from './client';
import type { APIResponse, SaldoSimpananAnggota, Simpanan, Member } from '@/types';

// ============================================================================
// Member Balance API
// ============================================================================

/**
 * Get member's share capital balance (Pokok, Wajib, Sukarela)
 */
export const getMemberBalance = async (): Promise<SaldoSimpananAnggota> => {
  const response = await apiClient.get<APIResponse<SaldoSimpananAnggota>>(
    '/member-portal/balance'
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch member balance');
  }

  return response.data.data;
};

// ============================================================================
// Transaction History API
// ============================================================================

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
): Promise<Simpanan[]> => {
  const response = await apiClient.get<APIResponse<Simpanan[]>>(
    '/member-portal/transactions',
    { params: filters }
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
    '/member-portal/profile'
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch member profile');
  }

  return response.data.data;
};

/**
 * Update member profile (limited fields)
 */
export const updateMemberProfile = async (
  data: Partial<Member>
): Promise<Member> => {
  const response = await apiClient.put<APIResponse<Member>>(
    '/member-portal/profile',
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
  transaksiTerbaru: Simpanan[];
  totalTransaksi: number;
}

/**
 * Get member dashboard summary
 */
export const getMemberDashboard = async (): Promise<MemberDashboardSummary> => {
  const response = await apiClient.get<APIResponse<MemberDashboardSummary>>(
    '/member-portal/dashboard'
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || 'Failed to fetch dashboard data');
  }

  return response.data.data;
};
