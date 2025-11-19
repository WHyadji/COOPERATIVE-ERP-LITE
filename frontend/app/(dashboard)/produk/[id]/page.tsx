// ============================================================================
// Product Detail Page - View product details and manage stock
// Comprehensive product information with stock adjustment functionality
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import {
  Box,
  Typography,
  Button,
  Paper,
  Grid,
  Chip,
  Divider,
  Alert,
  CircularProgress,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Card,
  CardContent,
} from "@mui/material";
import {
  ArrowBack as ArrowBackIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Warning as WarningIcon,
  CheckCircle as CheckCircleIcon,
  Inventory as InventoryIcon,
} from "@mui/icons-material";
import { useToast } from "@/lib/context/ToastContext";
import productApi from "@/lib/api/productApi";
import type { Produk } from "@/types";
import ProductForm from "@/components/products/ProductForm";

// ============================================================================
// Product Detail Page Component
// ============================================================================

export default function ProductDetailPage() {
  const router = useRouter();
  const params = useParams();
  const productId = params.id as string;
  const { showSuccess, showError } = useToast();

  const [product, setProduct] = useState<Produk | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");
  const [refreshKey, setRefreshKey] = useState(0);

  // Form dialog states
  const [formOpen, setFormOpen] = useState(false);
  const [stockDialogOpen, setStockDialogOpen] = useState(false);
  const [newStock, setNewStock] = useState("");
  const [stockAdjusting, setStockAdjusting] = useState(false);

  // ============================================================================
  // Fetch Product Details
  // ============================================================================

  useEffect(() => {
    let ignore = false;

    const fetchProduct = async () => {
      try {
        setLoading(true);
        setError("");

        const data = await productApi.getProductById(productId);

        if (!ignore) {
          setProduct(data);
          setNewStock(data.stok.toString());
        }
      } catch (err: unknown) {
        if (!ignore) {
          console.error("Failed to fetch product:", err);
          setError("Gagal memuat data produk. Silakan coba lagi.");
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    fetchProduct();

    return () => {
      ignore = true;
    };
  }, [productId, refreshKey]);

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleBack = () => {
    router.push("/produk");
  };

  const handleEdit = () => {
    setFormOpen(true);
  };

  const handleDelete = async () => {
    if (!product) return;

    if (
      !confirm(
        `Apakah Anda yakin ingin menghapus produk "${product.namaProduk}"?`
      )
    ) {
      return;
    }

    try {
      await productApi.deleteProduct(product.id);
      showSuccess(`Produk "${product.namaProduk}" berhasil dihapus`);
      router.push("/produk");
    } catch (err) {
      console.error("Failed to delete product:", err);
      showError("Gagal menghapus produk. Silakan coba lagi.");
    }
  };

  const handleCloseForm = () => {
    setFormOpen(false);
  };

  const handleFormSuccess = () => {
    setFormOpen(false);
    setRefreshKey((prev) => prev + 1);
  };

  const handleOpenStockDialog = () => {
    if (product) {
      setNewStock(product.stok.toString());
    }
    setStockDialogOpen(true);
  };

  const handleCloseStockDialog = () => {
    setStockDialogOpen(false);
  };

  const handleStockAdjustment = async () => {
    if (!product) return;

    const stockValue = parseInt(newStock);
    if (isNaN(stockValue) || stockValue < 0) {
      showError("Stok harus berupa angka >= 0");
      return;
    }

    setStockAdjusting(true);

    try {
      await productApi.updateProductStock(product.id, stockValue);
      showSuccess(
        `Stok produk "${product.namaProduk}" berhasil diperbarui menjadi ${stockValue}`
      );
      setStockDialogOpen(false);
      setRefreshKey((prev) => prev + 1);
    } catch (err) {
      console.error("Failed to update stock:", err);
      showError("Gagal memperbarui stok. Silakan coba lagi.");
    } finally {
      setStockAdjusting(false);
    }
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

  const isLowStock = (): boolean => {
    if (!product) return false;
    return product.stok <= product.stokMinimum && product.stokMinimum > 0;
  };

  const calculateMargin = (): number => {
    if (!product || !product.hargaBeli || product.hargaBeli === 0) return 0;
    return ((product.harga - product.hargaBeli) / product.hargaBeli) * 100;
  };

  // ============================================================================
  // Render
  // ============================================================================

  if (loading) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          minHeight: "60vh",
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  if (error || !product) {
    return (
      <Box sx={{ p: 3 }}>
        <Alert severity="error">{error || "Produk tidak ditemukan"}</Alert>
        <Button
          variant="outlined"
          startIcon={<ArrowBackIcon />}
          onClick={handleBack}
          sx={{ mt: 2 }}
        >
          Kembali ke Daftar Produk
        </Button>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 3 }}>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={handleBack}
          sx={{ mb: 2 }}
        >
          Kembali
        </Button>

        <Box
          sx={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "flex-start",
          }}
        >
          <Box>
            <Typography variant="h4" fontWeight={600} gutterBottom>
              {product.namaProduk}
            </Typography>
            <Box sx={{ display: "flex", gap: 1, flexWrap: "wrap" }}>
              <Chip
                label={product.statusAktif ? "Aktif" : "Nonaktif"}
                color={product.statusAktif ? "success" : "default"}
                size="small"
              />
              {isLowStock() && (
                <Chip
                  icon={<WarningIcon />}
                  label="Stok Rendah"
                  color="warning"
                  size="small"
                />
              )}
            </Box>
          </Box>

          <Box sx={{ display: "flex", gap: 1 }}>
            <Button
              variant="outlined"
              color="primary"
              startIcon={<EditIcon />}
              onClick={handleEdit}
            >
              Edit
            </Button>
            <Button
              variant="outlined"
              color="error"
              startIcon={<DeleteIcon />}
              onClick={handleDelete}
            >
              Hapus
            </Button>
          </Box>
        </Box>
      </Box>

      <Grid container spacing={3}>
        {/* Product Information */}
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" fontWeight={600} gutterBottom>
              Informasi Produk
            </Typography>
            <Divider sx={{ mb: 2 }} />

            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="body2" color="text.secondary">
                  Kode Produk
                </Typography>
                <Typography variant="body1" fontWeight={500}>
                  {product.kodeProduk}
                </Typography>
              </Grid>

              <Grid item xs={6}>
                <Typography variant="body2" color="text.secondary">
                  Barcode
                </Typography>
                <Typography variant="body1" fontWeight={500}>
                  {product.barcode || "-"}
                </Typography>
              </Grid>

              <Grid item xs={6}>
                <Typography variant="body2" color="text.secondary">
                  Kategori
                </Typography>
                <Typography variant="body1" fontWeight={500}>
                  {product.kategori || "-"}
                </Typography>
              </Grid>

              <Grid item xs={6}>
                <Typography variant="body2" color="text.secondary">
                  Satuan
                </Typography>
                <Typography variant="body1" fontWeight={500}>
                  {product.satuan}
                </Typography>
              </Grid>

              {product.deskripsi && (
                <Grid item xs={12}>
                  <Typography variant="body2" color="text.secondary">
                    Deskripsi
                  </Typography>
                  <Typography variant="body1">{product.deskripsi}</Typography>
                </Grid>
              )}
            </Grid>
          </Paper>

          {/* Pricing Information */}
          <Paper sx={{ p: 3, mt: 3 }}>
            <Typography variant="h6" fontWeight={600} gutterBottom>
              Informasi Harga
            </Typography>
            <Divider sx={{ mb: 2 }} />

            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="body2" color="text.secondary">
                  Harga Jual
                </Typography>
                <Typography variant="h6" color="primary" fontWeight={600}>
                  {formatCurrency(product.harga)}
                </Typography>
              </Grid>

              <Grid item xs={6}>
                <Typography variant="body2" color="text.secondary">
                  Harga Beli (HPP)
                </Typography>
                <Typography variant="body1" fontWeight={500}>
                  {product.hargaBeli ? formatCurrency(product.hargaBeli) : "-"}
                </Typography>
              </Grid>

              {product.hargaBeli > 0 && (
                <>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Margin Keuntungan
                    </Typography>
                    <Typography variant="body1" fontWeight={500}>
                      {formatCurrency(product.harga - product.hargaBeli)}
                    </Typography>
                  </Grid>

                  <Grid item xs={6}>
                    <Typography variant="body2" color="text.secondary">
                      Persentase Margin
                    </Typography>
                    <Typography
                      variant="body1"
                      fontWeight={500}
                      color={
                        calculateMargin() > 0 ? "success.main" : "error.main"
                      }
                    >
                      {calculateMargin().toFixed(2)}%
                    </Typography>
                  </Grid>
                </>
              )}
            </Grid>
          </Paper>
        </Grid>

        {/* Stock Management */}
        <Grid item xs={12} md={4}>
          <Card
            sx={{
              bgcolor: isLowStock() ? "warning.light" : "primary.light",
              color: "primary.contrastText",
            }}
          >
            <CardContent>
              <Box sx={{ display: "flex", alignItems: "center", mb: 2 }}>
                <InventoryIcon sx={{ mr: 1, fontSize: 40 }} />
                <Box>
                  <Typography variant="body2" sx={{ opacity: 0.9 }}>
                    Stok Saat Ini
                  </Typography>
                  <Typography variant="h3" fontWeight={600}>
                    {product.stok}
                  </Typography>
                </Box>
              </Box>

              <Divider sx={{ my: 2, bgcolor: "rgba(255,255,255,0.3)" }} />

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ opacity: 0.9 }}>
                  Stok Minimum: {product.stokMinimum}
                </Typography>
                <Typography variant="body2" sx={{ opacity: 0.9 }}>
                  Satuan: {product.satuan}
                </Typography>
              </Box>

              {isLowStock() && (
                <Alert severity="warning" icon={<WarningIcon />} sx={{ mb: 2 }}>
                  Stok produk sudah mencapai batas minimum!
                </Alert>
              )}

              <Button
                fullWidth
                variant="contained"
                color="inherit"
                onClick={handleOpenStockDialog}
                sx={{
                  bgcolor: "rgba(255,255,255,0.9)",
                  color: "primary.main",
                  "&:hover": {
                    bgcolor: "rgba(255,255,255,1)",
                  },
                }}
              >
                Sesuaikan Stok
              </Button>
            </CardContent>
          </Card>

          {/* Stock Status */}
          <Paper sx={{ p: 2, mt: 2 }}>
            <Typography variant="subtitle2" fontWeight={600} gutterBottom>
              Status Stok
            </Typography>
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
                gap: 1,
                color: isLowStock() ? "warning.main" : "success.main",
              }}
            >
              {isLowStock() ? (
                <>
                  <WarningIcon />
                  <Typography>Stok Rendah</Typography>
                </>
              ) : (
                <>
                  <CheckCircleIcon />
                  <Typography>Stok Aman</Typography>
                </>
              )}
            </Box>
          </Paper>
        </Grid>
      </Grid>

      {/* Product Form Dialog */}
      <ProductForm
        open={formOpen}
        onClose={handleCloseForm}
        onSuccess={handleFormSuccess}
        product={product}
      />

      {/* Stock Adjustment Dialog */}
      <Dialog
        open={stockDialogOpen}
        onClose={handleCloseStockDialog}
        maxWidth="xs"
        fullWidth
      >
        <DialogTitle>Sesuaikan Stok Produk</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Masukkan jumlah stok baru untuk produk{" "}
            <strong>{product.namaProduk}</strong>
          </Typography>

          <Alert severity="info" sx={{ mb: 2 }}>
            Stok saat ini: <strong>{product.stok}</strong> {product.satuan}
          </Alert>

          <TextField
            fullWidth
            type="number"
            label="Stok Baru"
            value={newStock}
            onChange={(e) => setNewStock(e.target.value)}
            inputProps={{ min: 0, step: 1 }}
            helperText={`Satuan: ${product.satuan}`}
            disabled={stockAdjusting}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseStockDialog} disabled={stockAdjusting}>
            Batal
          </Button>
          <Button
            variant="contained"
            onClick={handleStockAdjustment}
            disabled={stockAdjusting}
            startIcon={stockAdjusting && <CircularProgress size={20} />}
          >
            {stockAdjusting ? "Memperbarui..." : "Simpan"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
