// ============================================================================
// Receipt Dialog Component - POS Receipt Display
// Shows sale receipt after successful transaction with print functionality
// ============================================================================

'use client';

import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableRow,
  Paper,
} from '@mui/material';
import {
  Print as PrintIcon,
  CheckCircle as CheckIcon,
  ShoppingCart as CartIcon,
} from '@mui/icons-material';
import type { PenjualanResponse } from '@/types';

// ============================================================================
// Component Props
// ============================================================================

interface ReceiptDialogProps {
  open: boolean;
  onClose: () => void;
  sale: PenjualanResponse | null;
  onNewSale: () => void;
}

// ============================================================================
// Receipt Dialog Component
// ============================================================================

export default function ReceiptDialog({
  open,
  onClose,
  sale,
  onNewSale,
}: ReceiptDialogProps) {
  // Currency formatter
  const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  // Format date
  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('id-ID', {
      day: '2-digit',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
  };

  // Handle print
  const handlePrint = () => {
    window.print();
  };

  // Handle new sale
  const handleNewSale = () => {
    onNewSale();
    onClose();
  };

  if (!sale) {
    return null;
  }

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
      fullWidth
      PaperProps={{
        sx: {
          '@media print': {
            boxShadow: 'none',
            margin: 0,
            maxWidth: '100%',
          },
        },
      }}
    >
      <DialogTitle sx={{ textAlign: 'center', bgcolor: 'success.main', color: 'white', '@media print': { bgcolor: 'white', color: 'black' } }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 1 }}>
          <CheckIcon sx={{ fontSize: 32 }} />
          <Typography variant="h5" component="span" fontWeight={600}>
            Transaksi Berhasil
          </Typography>
        </Box>
      </DialogTitle>

      <DialogContent sx={{ p: 3 }}>
        {/* Receipt Header */}
        <Box sx={{ textAlign: 'center', mb: 3 }}>
          <Typography variant="h6" fontWeight={700} gutterBottom>
            STRUK PEMBAYARAN
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {sale.nomorPenjualan}
          </Typography>
          <Typography variant="caption" color="text.secondary">
            {formatDate(sale.tanggalPenjualan)}
          </Typography>
        </Box>

        <Divider sx={{ mb: 2 }} />

        {/* Member Info */}
        {sale.namaAnggota && (
          <Box sx={{ mb: 2 }}>
            <Typography variant="caption" color="text.secondary">
              Anggota
            </Typography>
            <Typography variant="body2" fontWeight={600}>
              {sale.namaAnggota}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              No: {sale.nomorAnggota}
            </Typography>
          </Box>
        )}

        {/* Items */}
        <Paper variant="outlined" sx={{ mb: 2 }}>
          <Table size="small">
            <TableBody>
              {sale.itemPenjualan.map((item, index) => (
                <React.Fragment key={item.id}>
                  <TableRow>
                    <TableCell colSpan={3} sx={{ pb: 0.5 }}>
                      <Typography variant="body2" fontWeight={600}>
                        {item.namaProduk}
                      </Typography>
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell sx={{ pt: 0, pb: index === sale.itemPenjualan.length - 1 ? 1 : 2 }}>
                      <Typography variant="caption" color="text.secondary">
                        {item.kuantitas} × {formatCurrency(item.hargaSatuan)}
                      </Typography>
                    </TableCell>
                    <TableCell sx={{ pt: 0, pb: index === sale.itemPenjualan.length - 1 ? 1 : 2 }} />
                    <TableCell align="right" sx={{ pt: 0, pb: index === sale.itemPenjualan.length - 1 ? 1 : 2 }}>
                      <Typography variant="body2" fontWeight={600}>
                        {formatCurrency(item.subtotal)}
                      </Typography>
                    </TableCell>
                  </TableRow>
                </React.Fragment>
              ))}
            </TableBody>
          </Table>
        </Paper>

        <Divider sx={{ mb: 2 }} />

        {/* Totals */}
        <Box sx={{ mb: 2 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
            <Typography variant="body1" fontWeight={600}>
              TOTAL
            </Typography>
            <Typography variant="h6" fontWeight={700} color="primary">
              {formatCurrency(sale.totalBelanja)}
            </Typography>
          </Box>

          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
            <Typography variant="body2" color="text.secondary">
              Bayar
            </Typography>
            <Typography variant="body2">
              {formatCurrency(sale.jumlahBayar)}
            </Typography>
          </Box>

          <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
            <Typography variant="body2" color="text.secondary">
              Kembalian
            </Typography>
            <Typography variant="body1" fontWeight={600} color="success.main">
              {formatCurrency(sale.kembalian)}
            </Typography>
          </Box>
        </Box>

        <Divider sx={{ mb: 2 }} />

        {/* Footer */}
        <Box sx={{ textAlign: 'center' }}>
          <Typography variant="caption" color="text.secondary" display="block">
            Kasir: {sale.namaKasir}
          </Typography>
          {sale.catatan && (
            <Typography variant="caption" color="text.secondary" display="block" sx={{ mt: 1 }}>
              Catatan: {sale.catatan}
            </Typography>
          )}
          <Typography variant="caption" color="text.secondary" display="block" sx={{ mt: 2 }}>
            Terima kasih atas kunjungan Anda
          </Typography>
          <Typography variant="caption" color="text.secondary" display="block">
            --- Simpan struk ini sebagai bukti pembayaran ---
          </Typography>
        </Box>
      </DialogContent>

      <Divider />

      <DialogActions sx={{ p: 2, justifyContent: 'space-between', '@media print': { display: 'none' } }}>
        <Button
          variant="outlined"
          onClick={handlePrint}
          startIcon={<PrintIcon />}
        >
          Cetak
        </Button>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button onClick={onClose}>
            Tutup
          </Button>
          <Button
            variant="contained"
            onClick={handleNewSale}
            startIcon={<CartIcon />}
            size="large"
          >
            Transaksi Baru
          </Button>
        </Box>
      </DialogActions>
    </Dialog>
  );
}
