// ============================================================================
// Product Grid Component - Visual Product Selection for POS
// Grid layout with category filters and quick add to cart
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
import {
  Box,
  Grid,
  Card,
  CardContent,
  CardMedia,
  CardActionArea,
  Typography,
  Chip,
  Button,
  ButtonGroup,
  CircularProgress,
  Alert,
  Badge,
} from "@mui/material";
import {
  Add as AddIcon,
  Inventory as InventoryIcon,
  Warning as WarningIcon,
} from "@mui/icons-material";
import { useToast } from "@/lib/context/ToastContext";
import productApi from "@/lib/api/productApi";
import type { Produk } from "@/types";

// ============================================================================
// Component Props
// ============================================================================

interface ProductGridProps {
  onProductSelect: (product: Produk) => void;
  disabled?: boolean;
}

// ============================================================================
// Categories
// ============================================================================

const categories = [
  "Semua",
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
// Product Grid Component
// ============================================================================

export default function ProductGrid({
  onProductSelect,
  disabled = false,
}: ProductGridProps) {
  const { showError } = useToast();
  const [products, setProducts] = useState<Produk[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState("Semua");

  // Currency formatter
  const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  // Check if product has low stock
  const isLowStock = (product: Produk): boolean => {
    return product.stok <= product.stokMinimum && product.stokMinimum > 0;
  };

  // Fetch products
  useEffect(() => {
    const fetchProducts = async () => {
      setLoading(true);
      try {
        const response = await productApi.getProducts({
          statusAktif: true,
          kategori: selectedCategory === "Semua" ? undefined : selectedCategory,
          pageSize: 100, // Show more products in POS
        });

        setProducts(response.data || []);
      } catch {
        showError("Gagal memuat produk");
        setProducts([]);
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, [selectedCategory, showError]);

  // Handle product click
  const handleProductClick = (product: Produk) => {
    if (product.stok === 0) {
      showError(`Produk "${product.namaProduk}" habis stok`);
      return;
    }

    onProductSelect(product);
  };

  return (
    <Box>
      {/* Category Filter */}
      <Box sx={{ mb: 3, overflowX: "auto" }}>
        <ButtonGroup
          variant="outlined"
          size="small"
          sx={{ flexWrap: "wrap", gap: 1 }}
        >
          {categories.map((category) => (
            <Button
              key={category}
              variant={selectedCategory === category ? "contained" : "outlined"}
              onClick={() => setSelectedCategory(category)}
              sx={{ borderRadius: 2 }}
            >
              {category}
            </Button>
          ))}
        </ButtonGroup>
      </Box>

      {/* Loading State */}
      {loading && (
        <Box sx={{ display: "flex", justifyContent: "center", py: 8 }}>
          <CircularProgress />
        </Box>
      )}

      {/* Empty State */}
      {!loading && products.length === 0 && (
        <Alert severity="info" sx={{ mt: 2 }}>
          Tidak ada produk aktif
          {selectedCategory !== "Semua"
            ? ` dalam kategori "${selectedCategory}"`
            : ""}
        </Alert>
      )}

      {/* Product Grid */}
      {!loading && products.length > 0 && (
        <Grid container spacing={2}>
          {products.map((product) => (
            <Grid item xs={6} sm={4} md={3} lg={2} key={product.id}>
              <Card
                sx={{
                  height: "100%",
                  display: "flex",
                  flexDirection: "column",
                  opacity: disabled || product.stok === 0 ? 0.6 : 1,
                  position: "relative",
                }}
              >
                {/* Low Stock Badge */}
                {isLowStock(product) && product.stok > 0 && (
                  <Badge
                    badgeContent={<WarningIcon sx={{ fontSize: 16 }} />}
                    color="warning"
                    sx={{ position: "absolute", top: 8, right: 8, zIndex: 1 }}
                  />
                )}

                {/* Out of Stock Badge */}
                {product.stok === 0 && (
                  <Chip
                    label="HABIS"
                    color="error"
                    size="small"
                    sx={{ position: "absolute", top: 8, right: 8, zIndex: 1 }}
                  />
                )}

                <CardActionArea
                  onClick={() => !disabled && handleProductClick(product)}
                  disabled={disabled || product.stok === 0}
                  sx={{
                    flexGrow: 1,
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "stretch",
                  }}
                >
                  {/* Product Image */}
                  {product.gambarUrl ? (
                    <CardMedia
                      component="img"
                      height="120"
                      image={product.gambarUrl}
                      alt={product.namaProduk}
                      sx={{ objectFit: "cover" }}
                    />
                  ) : (
                    <Box
                      sx={{
                        height: 120,
                        display: "flex",
                        alignItems: "center",
                        justifyContent: "center",
                        bgcolor: "grey.200",
                      }}
                    >
                      <InventoryIcon sx={{ fontSize: 48, color: "grey.400" }} />
                    </Box>
                  )}

                  {/* Product Info */}
                  <CardContent sx={{ flexGrow: 1, p: 1.5 }}>
                    <Typography
                      variant="body2"
                      fontWeight={600}
                      gutterBottom
                      sx={{
                        overflow: "hidden",
                        textOverflow: "ellipsis",
                        display: "-webkit-box",
                        WebkitLineClamp: 2,
                        WebkitBoxOrient: "vertical",
                        minHeight: "2.5em",
                      }}
                    >
                      {product.namaProduk}
                    </Typography>

                    <Typography
                      variant="caption"
                      color="text.secondary"
                      display="block"
                      gutterBottom
                    >
                      {product.kodeProduk}
                    </Typography>

                    <Typography
                      variant="h6"
                      color="primary"
                      fontWeight={700}
                      gutterBottom
                    >
                      {formatCurrency(product.harga)}
                    </Typography>

                    <Box
                      sx={{
                        display: "flex",
                        justifyContent: "space-between",
                        alignItems: "center",
                      }}
                    >
                      <Typography
                        variant="caption"
                        color={
                          product.stok === 0
                            ? "error"
                            : isLowStock(product)
                              ? "warning.main"
                              : "text.secondary"
                        }
                      >
                        Stok: {product.stok} {product.satuan}
                      </Typography>

                      <AddIcon fontSize="small" color="primary" />
                    </Box>
                  </CardContent>
                </CardActionArea>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}
    </Box>
  );
}
