// ============================================================================
// Member API - CRUD Operations for Members (Anggota)
// Type-safe API calls for member management
// ============================================================================

import apiClient from './client';
import type {
  Member,
  CreateMemberRequest,
  UpdateMemberRequest,
  MemberStatistics,
  APIResponse,
  PaginatedResponse,
  MemberListFilters,
} from '@/types';

// ============================================================================
// Member API Functions
// ============================================================================

/**
 * Get paginated list of members
 */
export const getMembers = async (
  filters?: MemberListFilters
): Promise<PaginatedResponse<Member>> => {
  const params: Record<string, string | number> = {
    page: filters?.page || 1,
    pageSize: filters?.pageSize || 20,
  };

  if (filters?.search) {
    params.search = filters.search;
  }

  if (filters?.status && filters.status !== 'all') {
    params.status = filters.status;
  }

  const response = await apiClient.get<PaginatedResponse<Member>>('/anggota', {
    params,
  });

  return response.data;
};

/**
 * Get member by ID
 */
export const getMemberById = async (id: string): Promise<Member> => {
  const response = await apiClient.get<APIResponse<Member>>(`/anggota/${id}`);

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch member');
  }

  return response.data.data;
};

/**
 * Get member by member number
 */
export const getMemberByNumber = async (nomorAnggota: string): Promise<Member> => {
  const response = await apiClient.get<APIResponse<Member>>(
    `/anggota/nomor/${nomorAnggota}`
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch member');
  }

  return response.data.data;
};

/**
 * Create new member
 */
export const createMember = async (
  data: CreateMemberRequest
): Promise<Member> => {
  const response = await apiClient.post<APIResponse<Member>>('/anggota', data);

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to create member');
  }

  return response.data.data;
};

/**
 * Update existing member
 */
export const updateMember = async (
  id: string,
  data: UpdateMemberRequest
): Promise<Member> => {
  const response = await apiClient.put<APIResponse<Member>>(
    `/anggota/${id}`,
    data
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to update member');
  }

  return response.data.data;
};

/**
 * Delete member
 */
export const deleteMember = async (id: string): Promise<void> => {
  const response = await apiClient.delete<APIResponse>(`/anggota/${id}`);

  if (!response.data.success) {
    throw new Error('Failed to delete member');
  }
};

/**
 * Get member statistics
 */
export const getMemberStatistics = async (): Promise<MemberStatistics> => {
  const response = await apiClient.get<APIResponse<MemberStatistics>>(
    '/anggota/statistik'
  );

  if (!response.data.success || !response.data.data) {
    throw new Error('Failed to fetch member statistics');
  }

  return response.data.data;
};

/**
 * Set member portal PIN
 */
export const setMemberPin = async (
  id: string,
  pin: string
): Promise<void> => {
  const response = await apiClient.post<APIResponse>(
    `/anggota/${id}/set-pin`,
    { pin }
  );

  if (!response.data.success) {
    throw new Error('Failed to set member PIN');
  }
};

/**
 * Validate member PIN
 */
export const validateMemberPin = async (
  nomorAnggota: string,
  pin: string
): Promise<boolean> => {
  try {
    const response = await apiClient.post<APIResponse<{ valid: boolean }>>(
      '/anggota/validate-pin',
      { nomorAnggota, pin }
    );

    return response.data.data?.valid || false;
  } catch (error) {
    console.error('PIN validation error:', error);
    return false;
  }
};

// ============================================================================
// Export all member API functions
// ============================================================================

const memberApi = {
  getMembers,
  getMemberById,
  getMemberByNumber,
  createMember,
  updateMember,
  deleteMember,
  getMemberStatistics,
  setMemberPin,
  validateMemberPin,
};

export default memberApi;
