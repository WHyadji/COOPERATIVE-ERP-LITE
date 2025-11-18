// ============================================================================
// Type Definitions for Cooperative ERP Lite Frontend
// Matching backend Go structs and API responses
// ============================================================================

// ----------------------------------------------------------------------------
// User & Authentication Types
// ----------------------------------------------------------------------------

export type UserRole = "admin" | "bendahara" | "kasir" | "anggota";

export interface User {
  id: string;
  idKoperasi: string;
  namaPengguna: string;
  namaLengkap: string;
  email: string;
  peran: UserRole;
  aktif: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface LoginRequest {
  namaPengguna: string;
  kataSandi: string;
}

export interface LoginResponse {
  token: string;
  pengguna: User;
}

export interface JWTPayload {
  idPengguna: string;
  idKoperasi: string;
  namaPengguna: string;
  namaLengkap: string;
  peran: UserRole;
  exp: number;
  iat: number;
  nbf: number;
  iss: string;
  sub: string;
}

// ----------------------------------------------------------------------------
// Member (Anggota) Types
// ----------------------------------------------------------------------------

export type MemberStatus = "aktif" | "nonaktif" | "diberhentikan";
export type Gender = "L" | "P"; // L = Laki-laki (Male), P = Perempuan (Female)

export interface Member {
  id: string;
  idKoperasi: string;
  nomorAnggota: string;
  namaLengkap: string;
  nik?: string;
  tanggalLahir?: string;
  tempatLahir?: string;
  jenisKelamin?: Gender;
  alamat?: string;
  rt?: string;
  rw?: string;
  kelurahan?: string;
  kecamatan?: string;
  kotaKabupaten?: string;
  provinsi?: string;
  kodePos?: string;
  noTelepon?: string;
  email?: string;
  pekerjaan?: string;
  tanggalBergabung: string;
  status: MemberStatus;
  fotoUrl?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateMemberRequest {
  namaLengkap: string;
  nik?: string;
  tanggalLahir?: string;
  tempatLahir?: string;
  jenisKelamin?: Gender;
  alamat?: string;
  rt?: string;
  rw?: string;
  kelurahan?: string;
  kecamatan?: string;
  kotaKabupaten?: string;
  provinsi?: string;
  kodePos?: string;
  noTelepon?: string;
  email?: string;
  pekerjaan?: string;
  tanggalBergabung: string;
}

export interface UpdateMemberRequest extends CreateMemberRequest {
  status?: MemberStatus;
}

export interface MemberStatistics {
  totalAnggota: number;
  anggotaAktif: number;
  anggotaNonaktif: number;
  anggotaBergabungBulanIni: number;
}

// ----------------------------------------------------------------------------
// API Response Types
// ----------------------------------------------------------------------------

export interface APIResponse<T = unknown> {
  success: boolean;
  message: string;
  data?: T;
}

export interface PaginationInfo {
  page: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
}

export interface PaginatedResponse<T> {
  success: boolean;
  message: string;
  data: T[];
  pagination: PaginationInfo;
}

export interface APIError {
  code: string;
  message: string;
  details?: unknown;
}

export interface APIErrorResponse {
  success: false;
  message: string;
  error: APIError;
}

// ----------------------------------------------------------------------------
// Form Types
// ----------------------------------------------------------------------------

export interface MemberFormData {
  namaLengkap: string;
  nik: string;
  tanggalLahir: Date | null;
  tempatLahir: string;
  jenisKelamin: Gender | "";
  alamat: string;
  rt: string;
  rw: string;
  kelurahan: string;
  kecamatan: string;
  kotaKabupaten: string;
  provinsi: string;
  kodePos: string;
  noTelepon: string;
  email: string;
  pekerjaan: string;
  tanggalBergabung: Date | null;
}

// ----------------------------------------------------------------------------
// Table/List Types
// ----------------------------------------------------------------------------

export interface MemberListFilters {
  search?: string;
  status?: MemberStatus | "all";
  page?: number;
  pageSize?: number;
}

// ----------------------------------------------------------------------------
// Cooperative Types (for multi-tenant context)
// ----------------------------------------------------------------------------

export interface Cooperative {
  id: string;
  nama: string;
  nomorBadanHukum?: string;
  alamat?: string;
  noTelepon?: string;
  email?: string;
  logoUrl?: string;
  aktif: boolean;
  createdAt?: string;
  updatedAt?: string;
}

// ----------------------------------------------------------------------------
// Utility Types
// ----------------------------------------------------------------------------

export type LoadingState = "idle" | "loading" | "success" | "error";

export interface AsyncState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
}
