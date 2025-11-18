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
// Share Capital (Simpanan) Types
// ----------------------------------------------------------------------------

export type TipeSimpanan = "pokok" | "wajib" | "sukarela";

export interface Simpanan {
  id: string;
  idKoperasi: string;
  idAnggota: string;
  namaAnggota: string;
  nomorAnggota: string;
  tipeSimpanan: TipeSimpanan;
  tanggalTransaksi: string;
  jumlahSetoran: number;
  keterangan: string;
  nomorReferensi: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateSimpananRequest {
  idAnggota: string;
  tipeSimpanan: TipeSimpanan;
  tanggalTransaksi: string;
  jumlahSetoran: number;
  keterangan?: string;
}

export interface SaldoSimpananAnggota {
  idAnggota: string;
  nomorAnggota: string;
  namaAnggota: string;
  simpananPokok: number;
  simpananWajib: number;
  simpananSukarela: number;
  totalSimpanan: number;
}

export interface RingkasanSimpanan {
  totalSimpananPokok: number;
  totalSimpananWajib: number;
  totalSimpananSukarela: number;
  totalSemuaSimpanan: number;
  jumlahAnggota: number;
}

export interface SimpananListFilters {
  tipeSimpanan?: TipeSimpanan | "all";
  idAnggota?: string;
  tanggalMulai?: string;
  tanggalAkhir?: string;
  page?: number;
  pageSize?: number;
}

export interface SimpananFormData {
  idAnggota: string;
  tipeSimpanan: TipeSimpanan | "";
  tanggalTransaksi: Date | null;
  jumlahSetoran: string;
  keterangan: string;
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
// Accounting (Akun & Transaksi) Types
// ----------------------------------------------------------------------------

export type TipeAkun = "aset" | "kewajiban" | "modal" | "pendapatan" | "beban";
export type NormalSaldo = "debit" | "kredit";

export interface Akun {
  id: string;
  idKoperasi: string;
  kodeAkun: string;
  namaAkun: string;
  tipeAkun: TipeAkun;
  idInduk?: string;
  namaInduk?: string;
  normalSaldo: NormalSaldo;
  deskripsi?: string;
  statusAktif: boolean;
  saldo?: number;
  tanggalDibuat?: string;
  tanggalDiperbarui?: string;
}

export interface CreateAkunRequest {
  kodeAkun: string;
  namaAkun: string;
  tipeAkun: TipeAkun;
  idInduk?: string;
  normalSaldo?: NormalSaldo;
  deskripsi?: string;
}

export interface UpdateAkunRequest extends CreateAkunRequest {
  statusAktif?: boolean;
}

export interface BarisTransaksi {
  id?: string;
  idAkun: string;
  kodeAkun?: string;
  namaAkun?: string;
  jumlahDebit: number;
  jumlahKredit: number;
  keterangan?: string;
}

export interface Transaksi {
  id: string;
  idKoperasi: string;
  nomorJurnal: string;
  tanggalTransaksi: string;
  deskripsi: string;
  nomorReferensi?: string;
  tipeTransaksi?: string;
  totalDebit: number;
  totalKredit: number;
  statusBalanced: boolean;
  barisTransaksi?: BarisTransaksi[];
  tanggalDibuat?: string;
  tanggalDiperbarui?: string;
}

export interface CreateTransaksiRequest {
  nomorJurnal: string;
  tanggalTransaksi: string;
  deskripsi: string;
  nomorReferensi?: string;
  tipeTransaksi?: string;
  barisTransaksi: {
    idAkun: string;
    jumlahDebit: number;
    jumlahKredit: number;
    keterangan?: string;
  }[];
}

export interface TransaksiListFilters {
  tanggalMulai?: string;
  tanggalAkhir?: string;
  tipeTransaksi?: string;
  page?: number;
  pageSize?: number;
}

export interface TransaksiFormData {
  nomorJurnal: string;
  tanggalTransaksi: Date | null;
  deskripsi: string;
  nomorReferensi: string;
  tipeTransaksi: string;
  barisTransaksi: {
    idAkun: string;
    jumlahDebit: string;
    jumlahKredit: string;
    keterangan: string;
  }[];
}

export interface AkunFormData {
  kodeAkun: string;
  namaAkun: string;
  tipeAkun: TipeAkun | "";
  idInduk: string;
  normalSaldo: NormalSaldo | "";
  deskripsi: string;
}

export interface LedgerEntry {
  tanggal: string;
  nomorJurnal: string;
  deskripsi: string;
  debit: number;
  kredit: number;
  saldo: number;
}

// ----------------------------------------------------------------------------
// Product (Produk) Types
// ----------------------------------------------------------------------------

export interface Produk {
  id: string;
  kodeProduk: string;
  namaProduk: string;
  kategori: string;
  deskripsi?: string;
  harga: number;           // Selling price
  hargaBeli: number;       // Cost price / HPP
  stok: number;
  stokMinimum: number;
  satuan: string;          // pcs, kg, liter, etc.
  barcode?: string;
  gambarUrl?: string;
  statusAktif: boolean;
}

export interface CreateProdukRequest {
  kodeProduk: string;
  namaProduk: string;
  kategori?: string;
  deskripsi?: string;
  harga: number;
  hargaBeli?: number;
  stok?: number;
  stokMinimum?: number;
  satuan?: string;
  barcode?: string;
  gambarUrl?: string;
}

export interface UpdateProdukRequest extends CreateProdukRequest {
  statusAktif?: boolean;
}

export interface ProdukListFilters {
  search?: string;
  kategori?: string;
  statusAktif?: boolean | "all";
  page?: number;
  pageSize?: number;
}

export interface ProdukFormData {
  kodeProduk: string;
  namaProduk: string;
  kategori: string;
  deskripsi: string;
  harga: string;
  hargaBeli: string;
  stok: string;
  stokMinimum: string;
  satuan: string;
  barcode: string;
  gambarUrl: string;
  statusAktif: boolean;
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
