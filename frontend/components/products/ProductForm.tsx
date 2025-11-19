// ============================================================================
// Product Form Component - Create/Edit Products
// Material-UI form with validation for product creation and editing
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Grid,
  Alert,
  CircularProgress,
  InputAdornment,
  FormControlLabel,
  Switch,
  MenuItem,
} from "@mui/material";
import { useToast } from "@/lib/context/ToastContext";
import productApi from "@/lib/api/productApi";
import type { Produk, ProdukFormData } from "@/types";

// ============================================================================
// Component Props
// ============================================================================

interface ProductFormProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
  product?: Produk | null; // If provided, edit mode
}

// ============================================================================
// Product Form Component
// ============================================================================

export default function ProductForm({
  open,
  onClose,
  onSuccess,
  product,
}: ProductFormProps) {
  const isEditMode = !!product;
  const { showSuccess, showError } = useToast();

  // Form state
  const [formData, setFormData] = useState<ProdukFormData>({
    kodeProduk: "",
    namaProduk: "",
    kategori: "",
    deskripsi: "",
    harga: "",
    hargaBeli: "",
    stok: "0",
    stokMinimum: "0",
    satuan: "pcs",
    barcode: "",
    gambarUrl: "",
    statusAktif: true,
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>("");
  const [errors, setErrors] = useState<
    Partial<Record<keyof ProdukFormData, string>>
  >({});

  // Common unit options
  const satuanOptions = [
    { value: "pcs", label: "Pcs (Pieces)" },
    { value: "kg", label: "Kg (Kilogram)" },
    { value: "gram", label: "Gram" },
    { value: "liter", label: "Liter" },
    { value: "ml", label: "Ml (Mililiter)" },
    { value: "box", label: "Box" },
    { value: "pack", label: "Pack" },
    { value: "lusin", label: "Lusin" },
    { value: "karton", label: "Karton" },
    { value: "meter", label: "Meter" },
  ];

  // ============================================================================
  // Initialize form data when product changes
  // ============================================================================

  useEffect(() => {
    if (open) {
      if (product) {
        setFormData({
          kodeProduk: product.kodeProduk,
          namaProduk: product.namaProduk,
          kategori: product.kategori || "",
          deskripsi: product.deskripsi || "",
          harga: product.harga.toString(),
          hargaBeli: product.hargaBeli.toString(),
          stok: product.stok.toString(),
          stokMinimum: product.stokMinimum.toString(),
          satuan: product.satuan || "pcs",
          barcode: product.barcode || "",
          gambarUrl: product.gambarUrl || "",
          statusAktif: product.statusAktif,
        });
      } else {
        setFormData({
          kodeProduk: "",
          namaProduk: "",
          kategori: "",
          deskripsi: "",
          harga: "",
          hargaBeli: "",
          stok: "0",
          stokMinimum: "0",
          satuan: "pcs",
          barcode: "",
          gambarUrl: "",
          statusAktif: true,
        });
      }
      setError("");
      setErrors({});
    }
  }, [product, open]);

  // ============================================================================
  // Validation
  // ============================================================================

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof ProdukFormData, string>> = {};

    if (!formData.kodeProduk.trim()) {
      newErrors.kodeProduk = "Kode produk harus diisi";
    }

    if (!formData.namaProduk.trim()) {
      newErrors.namaProduk = "Nama produk harus diisi";
    }

    if (!formData.harga || parseFloat(formData.harga) < 0) {
      newErrors.harga = "Harga jual harus diisi dan >= 0";
    }

    if (formData.hargaBeli && parseFloat(formData.hargaBeli) < 0) {
      newErrors.hargaBeli = "Harga beli harus >= 0";
    }

    if (formData.stok && parseInt(formData.stok) < 0) {
      newErrors.stok = "Stok harus >= 0";
    }

    if (formData.stokMinimum && parseInt(formData.stokMinimum) < 0) {
      newErrors.stokMinimum = "Stok minimum harus >= 0";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleChange = (
    field: keyof ProdukFormData,
    value: string | boolean
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear error for this field
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    setLoading(true);
    setError("");

    try {
      const requestData = {
        kodeProduk: formData.kodeProduk.trim(),
        namaProduk: formData.namaProduk.trim(),
        kategori: formData.kategori.trim(),
        deskripsi: formData.deskripsi.trim(),
        harga: parseFloat(formData.harga),
        hargaBeli: formData.hargaBeli ? parseFloat(formData.hargaBeli) : 0,
        stok: formData.stok ? parseInt(formData.stok) : 0,
        stokMinimum: formData.stokMinimum ? parseInt(formData.stokMinimum) : 0,
        satuan: formData.satuan || "pcs",
        barcode: formData.barcode.trim(),
        gambarUrl: formData.gambarUrl.trim(),
        statusAktif: formData.statusAktif,
      };

      if (isEditMode && product) {
        await productApi.updateProduct(product.id, requestData);
        showSuccess(`Produk "${formData.namaProduk}" berhasil diperbarui`);
      } else {
        await productApi.createProduct(requestData);
        showSuccess(`Produk "${formData.namaProduk}" berhasil dibuat`);
      }

      onSuccess();
      onClose();
    } catch (err) {
      console.error("Failed to save product:", err);
      const errorMessage =
        err instanceof Error ? err.message : "Gagal menyimpan produk";
      setError(errorMessage);
      showError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    if (!loading) {
      onClose();
    }
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {isEditMode ? "Edit Produk" : "Tambah Produk Baru"}
      </DialogTitle>

      <form onSubmit={handleSubmit}>
        <DialogContent>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <Grid container spacing={2}>
            {/* Kode Produk */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                required
                label="Kode Produk"
                value={formData.kodeProduk}
                onChange={(e) => handleChange("kodeProduk", e.target.value)}
                error={!!errors.kodeProduk}
                helperText={errors.kodeProduk || "Contoh: PRD-001"}
                disabled={loading}
              />
            </Grid>

            {/* Barcode */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Barcode"
                value={formData.barcode}
                onChange={(e) => handleChange("barcode", e.target.value)}
                error={!!errors.barcode}
                helperText={errors.barcode || "Opsional"}
                disabled={loading}
              />
            </Grid>

            {/* Nama Produk */}
            <Grid item xs={12}>
              <TextField
                fullWidth
                required
                label="Nama Produk"
                value={formData.namaProduk}
                onChange={(e) => handleChange("namaProduk", e.target.value)}
                error={!!errors.namaProduk}
                helperText={errors.namaProduk}
                disabled={loading}
              />
            </Grid>

            {/* Kategori */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Kategori"
                value={formData.kategori}
                onChange={(e) => handleChange("kategori", e.target.value)}
                error={!!errors.kategori}
                helperText={errors.kategori || "Contoh: Makanan, Minuman"}
                disabled={loading}
              />
            </Grid>

            {/* Satuan */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                select
                label="Satuan"
                value={formData.satuan}
                onChange={(e) => handleChange("satuan", e.target.value)}
                disabled={loading}
              >
                {satuanOptions.map((option) => (
                  <MenuItem key={option.value} value={option.value}>
                    {option.label}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>

            {/* Harga Jual */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                required
                type="number"
                label="Harga Jual"
                value={formData.harga}
                onChange={(e) => handleChange("harga", e.target.value)}
                error={!!errors.harga}
                helperText={errors.harga}
                disabled={loading}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">Rp</InputAdornment>
                  ),
                }}
                inputProps={{ min: 0, step: 100 }}
              />
            </Grid>

            {/* Harga Beli / HPP */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Harga Beli (HPP)"
                value={formData.hargaBeli}
                onChange={(e) => handleChange("hargaBeli", e.target.value)}
                error={!!errors.hargaBeli}
                helperText={errors.hargaBeli || "Opsional"}
                disabled={loading}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">Rp</InputAdornment>
                  ),
                }}
                inputProps={{ min: 0, step: 100 }}
              />
            </Grid>

            {/* Stok */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Stok"
                value={formData.stok}
                onChange={(e) => handleChange("stok", e.target.value)}
                error={!!errors.stok}
                helperText={errors.stok}
                disabled={loading}
                inputProps={{ min: 0, step: 1 }}
              />
            </Grid>

            {/* Stok Minimum */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Stok Minimum"
                value={formData.stokMinimum}
                onChange={(e) => handleChange("stokMinimum", e.target.value)}
                error={!!errors.stokMinimum}
                helperText={
                  errors.stokMinimum ||
                  "Peringatan jika stok di bawah nilai ini"
                }
                disabled={loading}
                inputProps={{ min: 0, step: 1 }}
              />
            </Grid>

            {/* Deskripsi */}
            <Grid item xs={12}>
              <TextField
                fullWidth
                multiline
                rows={3}
                label="Deskripsi"
                value={formData.deskripsi}
                onChange={(e) => handleChange("deskripsi", e.target.value)}
                error={!!errors.deskripsi}
                helperText={errors.deskripsi || "Opsional"}
                disabled={loading}
              />
            </Grid>

            {/* Gambar URL */}
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="URL Gambar"
                value={formData.gambarUrl}
                onChange={(e) => handleChange("gambarUrl", e.target.value)}
                error={!!errors.gambarUrl}
                helperText={errors.gambarUrl || "Opsional - URL gambar produk"}
                disabled={loading}
              />
            </Grid>

            {/* Status Aktif */}
            {isEditMode && (
              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={formData.statusAktif}
                      onChange={(e) =>
                        handleChange("statusAktif", e.target.checked)
                      }
                      disabled={loading}
                    />
                  }
                  label="Produk Aktif"
                />
              </Grid>
            )}
          </Grid>
        </DialogContent>

        <DialogActions>
          <Button onClick={handleClose} disabled={loading}>
            Batal
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={loading}
            startIcon={loading && <CircularProgress size={20} />}
          >
            {loading ? "Menyimpan..." : isEditMode ? "Simpan" : "Tambah Produk"}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}
