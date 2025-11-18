// ============================================================================
// Shopping Cart Component - POS Cart with Quantity Controls
// Displays cart items, quantities, and totals for POS transactions
// ============================================================================

'use client';

import React from 'react';
import {
  Box,
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  TextField,
  Button,
  Divider,
  Alert,
} from '@mui/material';
import {
  Add as AddIcon,
  Remove as RemoveIcon,
  Delete as DeleteIcon,
  ShoppingCart as CartIcon,
} from '@mui/icons-material';
import type { CartItem } from '@/types';

// ============================================================================
// Component Props
// ============================================================================

interface ShoppingCartProps {
  items: CartItem[];
  onQuantityChange: (productId: string, quantity: number) => void;
  onRemoveItem: (productId: string) => void;
  onClearCart: () => void;
}

// ============================================================================
// Shopping Cart Component
// ============================================================================

export default function ShoppingCart({
  items,
  onQuantityChange,
  onRemoveItem,
  onClearCart,
}: ShoppingCartProps) {
  // Currency formatter
  const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  // Calculate total
  const total = items.reduce((sum, item) => sum + item.subtotal, 0);

  // Handle quantity increment
  const handleIncrement = (item: CartItem) => {
    if (item.quantity < item.product.stok) {
      onQuantityChange(item.product.id, item.quantity + 1);
    }
  };

  // Handle quantity decrement
  const handleDecrement = (item: CartItem) => {
    if (item.quantity > 1) {
      onQuantityChange(item.product.id, item.quantity - 1);
    }
  };

  // Handle direct quantity input
  const handleQuantityInput = (item: CartItem, value: string) => {
    const quantity = parseInt(value);
    if (!isNaN(quantity) && quantity >= 1 && quantity <= item.product.stok) {
      onQuantityChange(item.product.id, quantity);
    }
  };

  // Empty cart state
  if (items.length === 0) {
    return (
      <Paper sx={{ p: 3, height: '100%', display: 'flex', flexDirection: 'column' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2 }}>
          <Typography variant="h6" fontWeight={600}>
            Keranjang Belanja
          </Typography>
        </Box>

        <Divider sx={{ mb: 3 }} />

        <Box
          sx={{
            flexGrow: 1,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            py: 8,
          }}
        >
          <CartIcon sx={{ fontSize: 80, color: 'action.disabled', mb: 2 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            Keranjang Kosong
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Tambahkan produk untuk memulai transaksi
          </Typography>
        </Box>
      </Paper>
    );
  }

  // Cart with items
  return (
    <Paper sx={{ p: 3, height: '100%', display: 'flex', flexDirection: 'column' }}>
      {/* Header */}
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2 }}>
        <Typography variant="h6" fontWeight={600}>
          Keranjang Belanja ({items.length} item)
        </Typography>
        <Button
          variant="outlined"
          color="error"
          size="small"
          startIcon={<DeleteIcon />}
          onClick={onClearCart}
        >
          Kosongkan
        </Button>
      </Box>

      <Divider sx={{ mb: 2 }} />

      {/* Cart Items Table */}
      <TableContainer sx={{ flexGrow: 1, mb: 2, maxHeight: 'calc(100vh - 400px)', overflow: 'auto' }}>
        <Table stickyHeader size="small">
          <TableHead>
            <TableRow>
              <TableCell>Produk</TableCell>
              <TableCell align="right">Harga</TableCell>
              <TableCell align="center">Jumlah</TableCell>
              <TableCell align="right">Subtotal</TableCell>
              <TableCell align="center">Aksi</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {items.map((item) => (
              <TableRow key={item.product.id} hover>
                {/* Product Name & Code */}
                <TableCell>
                  <Typography variant="body2" fontWeight={500}>
                    {item.product.namaProduk}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    {item.product.kodeProduk}
                  </Typography>
                </TableCell>

                {/* Price */}
                <TableCell align="right">
                  <Typography variant="body2">
                    {formatCurrency(item.product.harga)}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    /{item.product.satuan}
                  </Typography>
                </TableCell>

                {/* Quantity Controls */}
                <TableCell align="center">
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 0.5 }}>
                    <IconButton
                      size="small"
                      onClick={() => handleDecrement(item)}
                      disabled={item.quantity <= 1}
                      sx={{ p: 0.5 }}
                    >
                      <RemoveIcon fontSize="small" />
                    </IconButton>

                    <TextField
                      size="small"
                      type="number"
                      value={item.quantity}
                      onChange={(e) => handleQuantityInput(item, e.target.value)}
                      inputProps={{
                        min: 1,
                        max: item.product.stok,
                        style: { textAlign: 'center', width: '50px' },
                      }}
                      sx={{ mx: 0.5 }}
                    />

                    <IconButton
                      size="small"
                      onClick={() => handleIncrement(item)}
                      disabled={item.quantity >= item.product.stok}
                      sx={{ p: 0.5 }}
                    >
                      <AddIcon fontSize="small" />
                    </IconButton>
                  </Box>

                  {/* Stock Warning */}
                  {item.quantity >= item.product.stok && (
                    <Typography variant="caption" color="error" display="block" sx={{ mt: 0.5 }}>
                      Stok maksimal: {item.product.stok}
                    </Typography>
                  )}
                </TableCell>

                {/* Subtotal */}
                <TableCell align="right">
                  <Typography variant="body2" fontWeight={600}>
                    {formatCurrency(item.subtotal)}
                  </Typography>
                </TableCell>

                {/* Remove Button */}
                <TableCell align="center">
                  <IconButton
                    size="small"
                    color="error"
                    onClick={() => onRemoveItem(item.product.id)}
                  >
                    <DeleteIcon fontSize="small" />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Low Stock Warning */}
      {items.some((item) => item.quantity >= item.product.stok) && (
        <Alert severity="warning" sx={{ mb: 2 }}>
          Beberapa produk mencapai stok maksimal yang tersedia
        </Alert>
      )}

      {/* Total */}
      <Divider sx={{ mb: 2 }} />
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h6" fontWeight={600}>
          Total
        </Typography>
        <Typography variant="h5" fontWeight={700} color="primary">
          {formatCurrency(total)}
        </Typography>
      </Box>
    </Paper>
  );
}
