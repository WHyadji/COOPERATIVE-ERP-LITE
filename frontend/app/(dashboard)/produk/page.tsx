// ============================================================================
// Product (Produk) Page - View and manage products
// Material-UI table with product list, filters, and CRUD operations
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import {
  Box,
  Typography,
  Button,
  TextField,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  IconButton,
  Chip,
  Alert,
  CircularProgress,
  MenuItem,
  Grid,
  InputAdornment,
} from "@mui/material";
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Visibility as VisibilityIcon,
  Search as SearchIcon,
  Warning as WarningIcon,
} from "@mui/icons-material";
import { useToast } from "@/lib/context/ToastContext";
import productApi from "@/lib/api/productApi";
import type { Produk, ProdukListFilters } from "@/types";
import ProductForm from "@/components/products/ProductForm";

// ============================================================================
// Product Page Component
// ============================================================================

export default function ProductPage() {
  const router = useRouter();
  const { showSuccess, showError } = useToast();

  const [products, setProducts] = useState<Produk[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");

  // Pagination & Filters
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(20);
  const [totalItems, setTotalItems] = useState(0);
  const [search, setSearch] = useState("");
  const [kategori, setKategori] = useState("");
  const [statusAktif, setStatusAktif] = useState<boolean | "all">("all");
  const [refreshKey, setRefreshKey] = useState(0);

  // Product form dialog
  const [formOpen, setFormOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<Produk | null>(null);

  // Categories (in a real app, this could be fetched from backend)
  const categories = [
    "Makanan",
    "Minuman",
    "Sembako",
    "Kebutuhan Rumah Tangga",
    "Alat Tulis",
    "Elektronik",
    "Pakaian",
    "Kesehatan",
    "Lainnya",
  ];

  // ============================================================================
  // Fetch Products
  // ============================================================================

  useEffect(() => {
    let ignore = false;

    const fetchProducts = async () => {
      try {
        setLoading(true);
        setError("");

        const filters: ProdukListFilters = {
          page: page + 1, // API uses 1-based pagination
          pageSize: rowsPerPage,
          search: search || undefined,
          kategori: kategori || undefined,
          statusAktif: statusAktif === "all" ? undefined : statusAktif,
        };

        const response = await productApi.getProducts(filters);

        if (!ignore) {
          setProducts(response.data);
          setTotalItems(response.pagination.totalItems);
        }
      } catch (err: unknown) {
        if (!ignore) {
          console.error("Failed to fetch products:", err);
          setError("Gagal memuat data produk. Silakan coba lagi.");
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    fetchProducts();

    return () => {
      ignore = true;
    };
  }, [page, rowsPerPage, search, kategori, statusAktif, refreshKey]);

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleCreate = () => {
    setSelectedProduct(null);
    setFormOpen(true);
  };

  const handleEdit = async (id: string) => {
    try {
      const product = await productApi.getProductById(id);
      setSelectedProduct(product);
      setFormOpen(true);
    } catch (err) {
      console.error("Failed to fetch product:", err);
      showError("Gagal memuat data produk. Silakan coba lagi.");
    }
  };

  const handleDelete = async (id: string, namaProduk: string) => {
    if (!confirm(`Apakah Anda yakin ingin menghapus produk "${namaProduk}"?`)) {
      return;
    }

    try {
      await productApi.deleteProduct(id);
      showSuccess(`Produk "${namaProduk}" berhasil dihapus`);
      setRefreshKey((prev) => prev + 1);
    } catch (err) {
      console.error("Failed to delete product:", err);
      showError("Gagal menghapus produk. Silakan coba lagi.");
    }
  };

  const handleView = (id: string) => {
    router.push(`/produk/${id}`);
  };

  const handleChangePage = (_event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleCloseForm = () => {
    setFormOpen(false);
    setSelectedProduct(null);
  };

  const handleFormSuccess = () => {
    setFormOpen(false);
    setSelectedProduct(null);
    setRefreshKey((prev) => prev + 1);
  };

  // ============================================================================
  // Helper Functions
  // ============================================================================

  const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(value);
  };

  const isLowStock = (product: Produk): boolean => {
    return product.stok <= product.stokMinimum && product.stokMinimum > 0;
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Box sx={{ p: 3 }}>
      {/* Header */}
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <Typography variant="h4" fontWeight={600}>
          Manajemen Produk
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleCreate}
        >
          Tambah Produk
        </Button>
      </Box>

      {/* Filters */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} sm={6} md={4}>
            <TextField
              fullWidth
              size="small"
              label="Cari produk"
              placeholder="Nama atau kode produk"
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
                setPage(0); // Reset to first page when searching
              }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            <TextField
              fullWidth
              select
              size="small"
              label="Kategori"
              value={kategori}
              onChange={(e) => {
                setKategori(e.target.value);
                setPage(0);
              }}
            >
              <MenuItem value="">Semua Kategori</MenuItem>
              {categories.map((cat) => (
                <MenuItem key={cat} value={cat}>
                  {cat}
                </MenuItem>
              ))}
            </TextField>
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            <TextField
              fullWidth
              select
              size="small"
              label="Status"
              value={statusAktif}
              onChange={(e) => {
                const value = e.target.value;
                setStatusAktif(value === "all" ? "all" : value === "true");
                setPage(0);
              }}
            >
              <MenuItem value="all">Semua Status</MenuItem>
              <MenuItem value="true">Aktif</MenuItem>
              <MenuItem value="false">Nonaktif</MenuItem>
            </TextField>
          </Grid>
        </Grid>
      </Paper>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Products Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Kode</TableCell>
                <TableCell>Nama Produk</TableCell>
                <TableCell>Kategori</TableCell>
                <TableCell align="right">Harga Jual</TableCell>
                <TableCell align="center">Stok</TableCell>
                <TableCell>Satuan</TableCell>
                <TableCell align="center">Status</TableCell>
                <TableCell align="center">Aksi</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={8} align="center" sx={{ py: 5 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : products.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} align="center" sx={{ py: 5 }}>
                    <Typography color="text.secondary">
                      {search || kategori || statusAktif !== "all"
                        ? "Tidak ada produk yang sesuai dengan filter"
                        : 'Belum ada produk. Klik "Tambah Produk" untuk mulai.'}
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                products.map((product) => (
                  <TableRow key={product.id} hover>
                    <TableCell>{product.kodeProduk}</TableCell>
                    <TableCell>
                      <Typography fontWeight={500}>
                        {product.namaProduk}
                      </Typography>
                      {isLowStock(product) && (
                        <Chip
                          icon={<WarningIcon />}
                          label="Stok Rendah"
                          color="warning"
                          size="small"
                          sx={{ mt: 0.5 }}
                        />
                      )}
                    </TableCell>
                    <TableCell>{product.kategori || "-"}</TableCell>
                    <TableCell align="right">
                      {formatCurrency(product.harga)}
                    </TableCell>
                    <TableCell align="center">
                      <Chip
                        label={product.stok}
                        color={isLowStock(product) ? "warning" : "default"}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>{product.satuan}</TableCell>
                    <TableCell align="center">
                      <Chip
                        label={product.statusAktif ? "Aktif" : "Nonaktif"}
                        color={product.statusAktif ? "success" : "default"}
                        size="small"
                      />
                    </TableCell>
                    <TableCell align="center">
                      <IconButton
                        size="small"
                        onClick={() => handleView(product.id)}
                        title="Lihat Detail"
                      >
                        <VisibilityIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() => handleEdit(product.id)}
                        title="Edit"
                        color="primary"
                      >
                        <EditIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() =>
                          handleDelete(product.id, product.namaProduk)
                        }
                        title="Hapus"
                        color="error"
                      >
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>

        <TablePagination
          component="div"
          count={totalItems}
          page={page}
          onPageChange={handleChangePage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          rowsPerPageOptions={[10, 20, 50, 100]}
          labelRowsPerPage="Baris per halaman:"
          labelDisplayedRows={({ from, to, count }) =>
            `${from}–${to} dari ${count !== -1 ? count : `lebih dari ${to}`}`
          }
        />
      </Paper>

      {/* Product Form Dialog */}
      <ProductForm
        open={formOpen}
        onClose={handleCloseForm}
        onSuccess={handleFormSuccess}
        product={selectedProduct}
      />
    </Box>
  );
}
