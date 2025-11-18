// ============================================================================
// Transaction Form Component - Create Journal Entries with Double-Entry
// Material-UI form with line items and automatic double-entry validation
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
  Grid,
  Alert,
  CircularProgress,
  Box,
  Typography,
  IconButton,
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
  TableContainer,
  Select,
  MenuItem,
  FormControl,
  Paper,
  Divider,
  Chip,
} from '@mui/material';
import {
  Add as AddIcon,
  Delete as DeleteIcon,
  CheckCircle as CheckCircleIcon,
  Warning as WarningIcon,
} from '@mui/icons-material';
import accountingApi from '@/lib/api/accountingApi';
import type { Akun, CreateTransaksiRequest, Transaksi } from '@/types';
import { format } from 'date-fns';
import { useToast } from '@/lib/context/ToastContext';

// ============================================================================
// Component Props
// ============================================================================

interface TransactionFormProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
  transaction?: Transaksi | null;
}

interface LineItem {
  idAkun: string;
  jumlahDebit: string;
  jumlahKredit: string;
  keterangan: string;
}

// ============================================================================
// Transaction Form Component
// ============================================================================

export default function TransactionForm({
  open,
  onClose,
  onSuccess,
  transaction = null,
}: TransactionFormProps) {
  const { showSuccess } = useToast();
  const isEditMode = !!transaction;

  // Form state
  const [nomorJurnal, setNomorJurnal] = useState('');
  const [tanggalTransaksi, setTanggalTransaksi] = useState(
    format(new Date(), 'yyyy-MM-dd')
  );
  const [deskripsi, setDeskripsi] = useState('');
  const [nomorReferensi, setNomorReferensi] = useState('');
  const [lineItems, setLineItems] = useState<LineItem[]>([
    { idAkun: '', jumlahDebit: '', jumlahKredit: '', keterangan: '' },
    { idAkun: '', jumlahDebit: '', jumlahKredit: '', keterangan: '' },
  ]);

  const [accounts, setAccounts] = useState<Akun[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingAccounts, setLoadingAccounts] = useState(false);
  const [error, setError] = useState<string>('');

  // ============================================================================
  // Fetch Accounts
  // ============================================================================

  useEffect(() => {
    if (open) {
      fetchAccounts();

      // Pre-populate form if editing, otherwise reset
      if (isEditMode && transaction) {
        setNomorJurnal(transaction.nomorJurnal);
        setTanggalTransaksi(transaction.tanggalTransaksi.split('T')[0]); // Extract date part
        setDeskripsi(transaction.deskripsi);
        setNomorReferensi(transaction.nomorReferensi || '');

        // Populate line items
        if (transaction.barisTransaksi && transaction.barisTransaksi.length > 0) {
          setLineItems(
            transaction.barisTransaksi.map((baris) => ({
              idAkun: baris.idAkun,
              jumlahDebit: baris.jumlahDebit.toString(),
              jumlahKredit: baris.jumlahKredit.toString(),
              keterangan: baris.keterangan || '',
            }))
          );
        }
      } else {
        // Reset form for create mode
        setNomorJurnal('');
        setTanggalTransaksi(format(new Date(), 'yyyy-MM-dd'));
        setDeskripsi('');
        setNomorReferensi('');
        setLineItems([
          { idAkun: '', jumlahDebit: '', jumlahKredit: '', keterangan: '' },
        { idAkun: '', jumlahDebit: '', jumlahKredit: '', keterangan: '' },
      ]);
      }
      setError('');
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, transaction]);

  const fetchAccounts = async () => {
    try {
      setLoadingAccounts(true);
      const data = await accountingApi.getAccounts('all', true);
      setAccounts(data);
    } catch (err) {
      console.error('Failed to fetch accounts:', err);
    } finally {
      setLoadingAccounts(false);
    }
  };

  // ============================================================================
  // Validation & Calculations
  // ============================================================================

  const calculateTotals = () => {
    let totalDebit = 0;
    let totalKredit = 0;

    lineItems.forEach((item) => {
      totalDebit += parseFloat(item.jumlahDebit) || 0;
      totalKredit += parseFloat(item.jumlahKredit) || 0;
    });

    return { totalDebit, totalKredit };
  };

  const isBalanced = () => {
    const { totalDebit, totalKredit } = calculateTotals();
    return Math.abs(totalDebit - totalKredit) < 0.01 && totalDebit > 0;
  };

  const validateForm = (): boolean => {
    if (!nomorJurnal.trim()) {
      setError('Nomor jurnal harus diisi');
      return false;
    }

    if (!deskripsi.trim()) {
      setError('Deskripsi harus diisi');
      return false;
    }

    // Check if at least 2 line items are filled
    const filledItems = lineItems.filter(
      (item) => item.idAkun && (item.jumlahDebit || item.jumlahKredit)
    );

    if (filledItems.length < 2) {
      setError('Minimal 2 baris transaksi harus diisi');
      return false;
    }

    // Validate each line item
    for (const item of filledItems) {
      if (!item.idAkun) {
        setError('Setiap baris harus memilih akun');
        return false;
      }

      const debit = parseFloat(item.jumlahDebit) || 0;
      const kredit = parseFloat(item.jumlahKredit) || 0;

      if (debit > 0 && kredit > 0) {
        setError('Satu baris tidak boleh memiliki debit dan kredit sekaligus');
        return false;
      }

      if (debit === 0 && kredit === 0) {
        setError('Setiap baris harus memiliki nilai debit atau kredit');
        return false;
      }
    }

    // Check if balanced
    if (!isBalanced()) {
      setError('Total debit harus sama dengan total kredit');
      return false;
    }

    return true;
  };

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleAddLine = () => {
    setLineItems([
      ...lineItems,
      { idAkun: '', jumlahDebit: '', jumlahKredit: '', keterangan: '' },
    ]);
  };

  const handleRemoveLine = (index: number) => {
    if (lineItems.length <= 2) {
      alert('Minimal 2 baris transaksi diperlukan');
      return;
    }
    setLineItems(lineItems.filter((_, i) => i !== index));
  };

  const handleLineItemChange = (
    index: number,
    field: keyof LineItem,
    value: string
  ) => {
    const newItems = [...lineItems];
    newItems[index] = { ...newItems[index], [field]: value };
    setLineItems(newItems);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      setLoading(true);
      setError('');

      // Filter only filled line items
      const filledItems = lineItems.filter(
        (item) => item.idAkun && (item.jumlahDebit || item.jumlahKredit)
      );

      const requestData: CreateTransaksiRequest = {
        nomorJurnal: nomorJurnal.trim(),
        tanggalTransaksi: tanggalTransaksi,
        deskripsi: deskripsi.trim(),
        nomorReferensi: nomorReferensi.trim() || undefined,
        tipeTransaksi: 'manual',
        barisTransaksi: filledItems.map((item) => ({
          idAkun: item.idAkun,
          jumlahDebit: parseFloat(item.jumlahDebit) || 0,
          jumlahKredit: parseFloat(item.jumlahKredit) || 0,
          keterangan: item.keterangan.trim() || undefined,
        })),
      };

      if (isEditMode && transaction) {
        await accountingApi.updateTransaction(transaction.id, requestData);
        showSuccess(`Transaksi "${nomorJurnal}" berhasil diperbarui`);
      } else {
        await accountingApi.createTransaction(requestData);
        showSuccess(`Transaksi "${nomorJurnal}" berhasil dibuat`);
      }
      onSuccess();
    } catch (err: unknown) {
      console.error('Failed to create transaction:', err);
      setError(
        err instanceof Error
          ? err.message
          : 'Gagal menyimpan transaksi. Silakan coba lagi.'
      );
    } finally {
      setLoading(false);
    }
  };

  // ============================================================================
  // Render Helpers
  // ============================================================================

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const { totalDebit, totalKredit } = calculateTotals();
  const balanced = isBalanced();

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>
          {isEditMode ? 'Edit Jurnal Umum' : 'Buat Jurnal Umum Baru'}
        </DialogTitle>

        <DialogContent>
          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          {/* Transaction Header */}
          <Grid container spacing={2} sx={{ mb: 3 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Nomor Jurnal"
                value={nomorJurnal}
                onChange={(e) => setNomorJurnal(e.target.value)}
                required
                disabled={loading}
                helperText="Contoh: JU-2025-001"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Tanggal Transaksi"
                type="date"
                value={tanggalTransaksi}
                onChange={(e) => setTanggalTransaksi(e.target.value)}
                required
                disabled={loading}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Deskripsi"
                value={deskripsi}
                onChange={(e) => setDeskripsi(e.target.value)}
                required
                disabled={loading}
                multiline
                rows={2}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Nomor Referensi (Opsional)"
                value={nomorReferensi}
                onChange={(e) => setNomorReferensi(e.target.value)}
                disabled={loading}
                helperText="Nomor bukti transaksi eksternal"
              />
            </Grid>
          </Grid>

          <Divider sx={{ my: 3 }} />

          {/* Line Items */}
          <Box sx={{ mb: 2 }}>
            <Box
              sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                mb: 2,
              }}
            >
              <Typography variant="h6">Baris Transaksi</Typography>
              <Button
                startIcon={<AddIcon />}
                onClick={handleAddLine}
                disabled={loading}
                size="small"
              >
                Tambah Baris
              </Button>
            </Box>

            <TableContainer component={Paper} variant="outlined">
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell width="30%">Akun</TableCell>
                    <TableCell width="18%">Debit</TableCell>
                    <TableCell width="18%">Kredit</TableCell>
                    <TableCell width="25%">Keterangan</TableCell>
                    <TableCell width="9%" align="center">
                      Aksi
                    </TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {lineItems.map((item, index) => (
                    <TableRow key={index}>
                      <TableCell>
                        <FormControl fullWidth size="small">
                          <Select
                            value={item.idAkun}
                            onChange={(e) =>
                              handleLineItemChange(index, 'idAkun', e.target.value)
                            }
                            disabled={loading || loadingAccounts}
                            displayEmpty
                          >
                            <MenuItem value="">
                              <em>Pilih Akun</em>
                            </MenuItem>
                            {loadingAccounts ? (
                              <MenuItem disabled>Loading accounts...</MenuItem>
                            ) : (
                              accounts?.map((account) => (
                                <MenuItem key={account.id} value={account.id}>
                                  {account.kodeAkun} - {account.namaAkun}
                                </MenuItem>
                              ))
                            )}
                          </Select>
                        </FormControl>
                      </TableCell>
                      <TableCell>
                        <TextField
                          fullWidth
                          size="small"
                          type="number"
                          value={item.jumlahDebit}
                          onChange={(e) =>
                            handleLineItemChange(index, 'jumlahDebit', e.target.value)
                          }
                          disabled={loading}
                          inputProps={{ min: 0, step: 0.01 }}
                        />
                      </TableCell>
                      <TableCell>
                        <TextField
                          fullWidth
                          size="small"
                          type="number"
                          value={item.jumlahKredit}
                          onChange={(e) =>
                            handleLineItemChange(
                              index,
                              'jumlahKredit',
                              e.target.value
                            )
                          }
                          disabled={loading}
                          inputProps={{ min: 0, step: 0.01 }}
                        />
                      </TableCell>
                      <TableCell>
                        <TextField
                          fullWidth
                          size="small"
                          value={item.keterangan}
                          onChange={(e) =>
                            handleLineItemChange(index, 'keterangan', e.target.value)
                          }
                          disabled={loading}
                        />
                      </TableCell>
                      <TableCell align="center">
                        <IconButton
                          size="small"
                          onClick={() => handleRemoveLine(index)}
                          disabled={loading || lineItems.length <= 2}
                          color="error"
                        >
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))}

                  {/* Totals Row */}
                  <TableRow sx={{ bgcolor: 'grey.50' }}>
                    <TableCell>
                      <Typography fontWeight={600}>TOTAL</Typography>
                    </TableCell>
                    <TableCell>
                      <Typography fontWeight={600} color="primary">
                        {formatCurrency(totalDebit)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography fontWeight={600} color="error">
                        {formatCurrency(totalKredit)}
                      </Typography>
                    </TableCell>
                    <TableCell colSpan={2} align="right">
                      {balanced ? (
                        <Chip
                          icon={<CheckCircleIcon />}
                          label="Balanced"
                          color="success"
                          size="small"
                        />
                      ) : (
                        <Chip
                          icon={<WarningIcon />}
                          label="Unbalanced"
                          color="error"
                          size="small"
                        />
                      )}
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </TableContainer>

            <Alert severity="info" sx={{ mt: 2 }}>
              <Typography variant="body2">
                <strong>Prinsip Double-Entry:</strong> Total Debit harus sama dengan
                Total Kredit. Setiap baris hanya boleh memiliki nilai Debit ATAU
                Kredit, tidak keduanya.
              </Typography>
            </Alert>
          </Box>
        </DialogContent>

        <DialogActions>
          <Button onClick={onClose} disabled={loading}>
            Batal
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={loading || !balanced}
            startIcon={loading && <CircularProgress size={20} />}
          >
            {loading ? 'Menyimpan...' : 'Simpan Jurnal'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}
