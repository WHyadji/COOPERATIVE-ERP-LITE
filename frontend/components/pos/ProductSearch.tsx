// ============================================================================
// Product Search Component - Barcode/Name Search for POS
// Autocomplete search with barcode scanner support
// ============================================================================

"use client";

import React, { useState, useEffect, useCallback } from "react";
import {
  Autocomplete,
  TextField,
  Box,
  Typography,
  CircularProgress,
  InputAdornment,
} from "@mui/material";
import {
  Search as SearchIcon,
  QrCodeScanner as BarcodeIcon,
} from "@mui/icons-material";
import { useToast } from "@/lib/context/ToastContext";
import productApi from "@/lib/api/productApi";
import type { Produk } from "@/types";

// ============================================================================
// Component Props
// ============================================================================

interface ProductSearchProps {
  onProductSelect: (product: Produk) => void;
  disabled?: boolean;
}

// ============================================================================
// Product Search Component
// ============================================================================

export default function ProductSearch({
  onProductSelect,
  disabled = false,
}: ProductSearchProps) {
  const { showError } = useToast();
  const [open, setOpen] = useState(false);
  const [options, setOptions] = useState<Produk[]>([]);
  const [loading, setLoading] = useState(false);
  const [inputValue, setInputValue] = useState("");
  const [searchTerm, setSearchTerm] = useState("");

  // Currency formatter
  const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  // Debounced search
  useEffect(() => {
    const timer = setTimeout(() => {
      setSearchTerm(inputValue);
    }, 300);

    return () => clearTimeout(timer);
  }, [inputValue]);

  // Fetch products when search term changes
  useEffect(() => {
    if (!searchTerm) {
      setOptions([]);
      return;
    }

    const fetchProducts = async () => {
      setLoading(true);
      try {
        // Try barcode lookup first if input looks like a barcode (numbers only)
        if (/^\d+$/.test(searchTerm)) {
          try {
            const product = await productApi.getProductByBarcode(searchTerm);
            if (product && product.statusAktif) {
              setOptions([product]);
              setLoading(false);
              return;
            }
          } catch {
            // If barcode not found, fall through to regular search
          }
        }

        // Regular search by name or code
        const response = await productApi.getProducts({
          search: searchTerm,
          statusAktif: true,
          pageSize: 10,
        });

        setOptions(response.data || []);
      } catch {
        showError("Gagal mencari produk");
        setOptions([]);
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, [searchTerm, showError]);

  // Handle product selection
  const handleSelect = useCallback(
    (_event: React.SyntheticEvent, value: Produk | null) => {
      if (value) {
        onProductSelect(value);
        setInputValue("");
        setSearchTerm("");
        setOptions([]);
      }
    },
    [onProductSelect]
  );

  return (
    <Autocomplete
      open={open}
      onOpen={() => setOpen(true)}
      onClose={() => setOpen(false)}
      disabled={disabled}
      value={null}
      inputValue={inputValue}
      onInputChange={(_event, newInputValue) => setInputValue(newInputValue)}
      onChange={handleSelect}
      isOptionEqualToValue={(option, value) => option.id === value.id}
      getOptionLabel={(option) => option.namaProduk}
      options={options}
      loading={loading}
      noOptionsText={
        searchTerm
          ? "Produk tidak ditemukan"
          : "Ketik untuk mencari produk (nama, kode, atau barcode)"
      }
      renderInput={(params) => (
        <TextField
          {...params}
          label="Cari Produk"
          placeholder="Scan barcode atau ketik nama/kode produk"
          variant="outlined"
          fullWidth
          InputProps={{
            ...params.InputProps,
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon />
              </InputAdornment>
            ),
            endAdornment: (
              <>
                {loading ? (
                  <CircularProgress color="inherit" size={20} />
                ) : null}
                {params.InputProps.endAdornment}
                <InputAdornment position="end">
                  <BarcodeIcon color="action" />
                </InputAdornment>
              </>
            ),
          }}
        />
      )}
      renderOption={(props, option) => (
        <Box component="li" {...props} key={option.id}>
          <Box sx={{ flexGrow: 1 }}>
            <Typography variant="body2" fontWeight={500}>
              {option.namaProduk}
            </Typography>
            <Box sx={{ display: "flex", gap: 2, mt: 0.5 }}>
              <Typography variant="caption" color="text.secondary">
                Kode: {option.kodeProduk}
              </Typography>
              {option.barcode && (
                <Typography variant="caption" color="text.secondary">
                  Barcode: {option.barcode}
                </Typography>
              )}
              <Typography variant="caption" color="primary">
                Stok: {option.stok} {option.satuan}
              </Typography>
            </Box>
          </Box>
          <Box sx={{ textAlign: "right", ml: 2 }}>
            <Typography variant="body2" fontWeight={600} color="primary">
              {formatCurrency(option.harga)}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              /{option.satuan}
            </Typography>
          </Box>
        </Box>
      )}
      sx={{ mb: 2 }}
    />
  );
}
