// ============================================================================
// Checkout Modal Component - POS Payment Processing
// Handles payment input, change calculation, and sale confirmation
// ============================================================================

'use client';

import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Typography,
  Box,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Alert,
  Grid,
  InputAdornment,
} from '@mui/material';
import {
  Payments as PaymentsIcon,
  CheckCircle as CheckIcon,
} from '@mui/icons-material';
import { useToast } from '@/lib/context/ToastContext';
import posApi from '@/lib/api/posApi';
import type { CartItem, Member, PenjualanResponse } from '@/types';

// ============================================================================
// Component Props
// ============================================================================

interface CheckoutModalProps {
  open: boolean;
  onClose: () => void;
  items: CartItem[];
  member: Member | null;
  onSuccess: (sale: PenjualanResponse) => void;
}

// ============================================================================
// Checkout Modal Component
// ============================================================================

export default function CheckoutModal({
  open,
  onClose,
  items,
  member,
  onSuccess,
}: CheckoutModalProps) {
  const { showError } = useToast();
  const [jumlahBayar, setJumlahBayar] = useState('');
  const [catatan, setCatatan] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Calculate totals
  const totalBelanja = items.reduce((sum, item) => sum + item.subtotal, 0);
  const jumlahBayarNum = parseFloat(jumlahBayar) || 0;
  const kembalian = jumlahBayarNum - totalBelanja;

  // Currency formatter
  const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  // Format number input
  const formatNumberInput = (value: string): string => {
    // Remove non-numeric characters
    const numericValue = value.replace(/[^\d]/g, '');
    return numericValue;
  };

  // Handle payment input change
  const handlePaymentChange = (value: string) => {
    const formatted = formatNumberInput(value);
    setJumlahBayar(formatted);
    setError('');
  };

  // Quick amount buttons
  const quickAmounts = [
    totalBelanja,
    Math.ceil(totalBelanja / 50000) * 50000, // Round up to nearest 50k
    Math.ceil(totalBelanja / 100000) * 100000, // Round up to nearest 100k
  ].filter((amount, index, self) => self.indexOf(amount) === index); // Remove duplicates

  // Handle quick amount click
  const handleQuickAmount = (amount: number) => {
    setJumlahBayar(amount.toString());
    setError('');
  };

  // Validate payment
  const validatePayment = (): boolean => {
    if (!jumlahBayar || jumlahBayarNum <= 0) {
      setError('Jumlah bayar harus diisi');
      return false;
    }

    if (jumlahBayarNum < totalBelanja) {
      setError('Jumlah bayar kurang dari total belanja');
      return false;
    }

    return true;
  };

  // Handle submit
  const handleSubmit = async () => {
    if (!validatePayment()) {
      return;
    }

    setLoading(true);
    setError('');

    try {
      const saleData = {
        idAnggota: member?.id,
        items: items.map((item) => ({
          idProduk: item.product.id,
          kuantitas: item.quantity,
          hargaSatuan: item.product.harga,
        })),
        jumlahBayar: jumlahBayarNum,
        catatan: catatan.trim() || undefined,
      };

      const sale = await posApi.createSale(saleData);
      onSuccess(sale);
      onClose();

      // Reset form
      setJumlahBayar('');
      setCatatan('');
      setError('');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Gagal memproses transaksi';
      setError(errorMessage);
      showError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // Reset form when modal opens
  useEffect(() => {
    if (open) {
      setJumlahBayar('');
      setCatatan('');
      setError('');
    }
  }, [open]);

  return (
    <Dialog
      open={open}
      onClose={loading ? undefined : onClose}
      maxWidth="md"
      fullWidth
      disableEscapeKeyDown={loading}
    >
      <DialogTitle>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <PaymentsIcon />
          <Typography variant="h6" component="span">
            Pembayaran
          </Typography>
        </Box>
      </DialogTitle>

      <DialogContent>
        {/* Error Alert */}
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {/* Member Info */}
        {member && (
          <Box sx={{ mb: 2, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
            <Typography variant="body2" color="text.secondary">
              Anggota
            </Typography>
            <Typography variant="body1" fontWeight={600}>
              {member.namaLengkap}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              No: {member.nomorAnggota}
            </Typography>
          </Box>
        )}

        {/* Items Summary */}
        <TableContainer component={Paper} variant="outlined" sx={{ mb: 2 }}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Produk</TableCell>
                <TableCell align="right">Harga</TableCell>
                <TableCell align="center">Qty</TableCell>
                <TableCell align="right">Subtotal</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {items.map((item) => (
                <TableRow key={item.product.id}>
                  <TableCell>
                    <Typography variant="body2">{item.product.namaProduk}</Typography>
                  </TableCell>
                  <TableCell align="right">
                    {formatCurrency(item.product.harga)}
                  </TableCell>
                  <TableCell align="center">
                    {item.quantity} {item.product.satuan}
                  </TableCell>
                  <TableCell align="right" sx={{ fontWeight: 600 }}>
                    {formatCurrency(item.subtotal)}
                  </TableCell>
                </TableRow>
              ))}
              <TableRow>
                <TableCell colSpan={3} sx={{ fontWeight: 600, fontSize: '1.1rem' }}>
                  TOTAL
                </TableCell>
                <TableCell align="right" sx={{ fontWeight: 700, fontSize: '1.1rem', color: 'primary.main' }}>
                  {formatCurrency(totalBelanja)}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </TableContainer>

        {/* Payment Input */}
        <Box sx={{ mb: 2 }}>
          <TextField
            fullWidth
            label="Jumlah Bayar"
            value={jumlahBayar ? formatCurrency(jumlahBayarNum) : ''}
            onChange={(e) => handlePaymentChange(e.target.value)}
            placeholder="0"
            disabled={loading}
            autoFocus
            InputProps={{
              startAdornment: <InputAdornment position="start">Rp</InputAdornment>,
            }}
            sx={{
              '& input': {
                fontSize: '1.5rem',
                fontWeight: 600,
              },
            }}
          />

          {/* Quick Amount Buttons */}
          <Grid container spacing={1} sx={{ mt: 1 }}>
            {quickAmounts.map((amount) => (
              <Grid item xs={4} key={amount}>
                <Button
                  fullWidth
                  variant="outlined"
                  onClick={() => handleQuickAmount(amount)}
                  disabled={loading}
                  size="small"
                >
                  {formatCurrency(amount)}
                </Button>
              </Grid>
            ))}
          </Grid>
        </Box>

        {/* Change Calculation */}
        <Box
          sx={{
            p: 2,
            bgcolor: kembalian >= 0 ? 'success.light' : 'error.light',
            borderRadius: 1,
            mb: 2,
          }}
        >
          <Typography variant="body2" color="text.secondary">
            Kembalian
          </Typography>
          <Typography variant="h4" fontWeight={700} color={kembalian >= 0 ? 'success.dark' : 'error.dark'}>
            {formatCurrency(Math.max(0, kembalian))}
          </Typography>
        </Box>

        {/* Notes */}
        <TextField
          fullWidth
          label="Catatan (Opsional)"
          value={catatan}
          onChange={(e) => setCatatan(e.target.value)}
          disabled={loading}
          multiline
          rows={2}
          placeholder="Tambahkan catatan untuk transaksi ini..."
        />
      </DialogContent>

      <Divider />

      <DialogActions sx={{ p: 2 }}>
        <Button onClick={onClose} disabled={loading}>
          Batal
        </Button>
        <Button
          variant="contained"
          onClick={handleSubmit}
          disabled={loading || jumlahBayarNum < totalBelanja}
          startIcon={loading ? null : <CheckIcon />}
          size="large"
        >
          {loading ? 'Memproses...' : 'Proses Pembayaran'}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
