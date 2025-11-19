// ============================================================================
// Product API - CRUD Operations for Products
// Type-safe API calls for product management
// ============================================================================

import apiClient from "./client";
import type {
  Produk,
  CreateProdukRequest,
  UpdateProdukRequest,
  ProdukListFilters,
  APIResponse,
  PaginatedResponse,
} from "@/types";

// ============================================================================
// Product (Produk) API Functions
// ============================================================================

/**
 * Get all products with pagination and filters
 */
export const getProducts = async (
  filters?: ProdukListFilters
): Promise<PaginatedResponse<Produk>> => {
  const params: Record<string, string> = {};

  if (filters?.search) {
    params.search = filters.search;
  }

  if (filters?.kategori) {
    params.kategori = filters.kategori;
  }

  if (filters?.statusAktif !== undefined && filters.statusAktif !== "all") {
    params.statusAktif = filters.statusAktif.toString();
  }

  if (filters?.page) {
    params.page = filters.page.toString();
  }

  if (filters?.pageSize) {
    params.pageSize = filters.pageSize.toString();
  }

  const response = await apiClient.get<PaginatedResponse<Produk>>("/produk", {
    params,
  });

  if (!response.data.success) {
    throw new Error("Failed to fetch products");
  }

  return response.data;
};

/**
 * Get product by ID
 */
export const getProductById = async (id: string): Promise<Produk> => {
  const response = await apiClient.get<APIResponse<Produk>>(`/produk/${id}`);

  if (!response.data.success || !response.data.data) {
    throw new Error("Failed to fetch product");
  }

  return response.data.data;
};

/**
 * Get product by barcode
 */
export const getProductByBarcode = async (barcode: string): Promise<Produk> => {
  const response = await apiClient.get<APIResponse<Produk>>(
    `/produk/barcode/${barcode}`
  );

  if (!response.data.success || !response.data.data) {
    throw new Error("Failed to fetch product");
  }

  return response.data.data;
};

/**
 * Get low stock products
 */
export const getLowStockProducts = async (): Promise<Produk[]> => {
  const response = await apiClient.get<APIResponse<Produk[]>>(
    "/produk/stok-rendah"
  );

  if (!response.data.success || !response.data.data) {
    throw new Error("Failed to fetch low stock products");
  }

  return response.data.data;
};

/**
 * Create a new product
 */
export const createProduct = async (
  data: CreateProdukRequest
): Promise<Produk> => {
  const response = await apiClient.post<APIResponse<Produk>>("/produk", data);

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || "Failed to create product");
  }

  return response.data.data;
};

/**
 * Update an existing product
 */
export const updateProduct = async (
  id: string,
  data: UpdateProdukRequest
): Promise<Produk> => {
  const response = await apiClient.put<APIResponse<Produk>>(
    `/produk/${id}`,
    data
  );

  if (!response.data.success || !response.data.data) {
    throw new Error(response.data.message || "Failed to update product");
  }

  return response.data.data;
};

/**
 * Delete a product
 */
export const deleteProduct = async (id: string): Promise<void> => {
  const response = await apiClient.delete<APIResponse<null>>(`/produk/${id}`);

  if (!response.data.success) {
    throw new Error(response.data.message || "Failed to delete product");
  }
};

/**
 * Update product stock (for stock adjustments)
 */
export const updateProductStock = async (
  id: string,
  stok: number
): Promise<Produk> => {
  const product = await getProductById(id);

  return updateProduct(id, {
    ...product,
    stok,
  });
};

// Export all functions as a single object (optional, for convenience)
const productApi = {
  getProducts,
  getProductById,
  getProductByBarcode,
  getLowStockProducts,
  createProduct,
  updateProduct,
  deleteProduct,
  updateProductStock,
};

export default productApi;
