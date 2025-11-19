// ============================================================================
// Member Transaction History Page
// View all transactions (simpanan deposits)
// ============================================================================

'use client';

import React, { useEffect, useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  CircularProgress,
  Alert,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  TextField,
  MenuItem,
  Grid,
  Button,
  Divider,
} from '@mui/material';
import { FilterList, Receipt } from '@mui/icons-material';
import { getMemberTransactions, type MemberTransactionFilters, type RiwayatTransaksiAnggota } from '@/lib/api/memberPortalApi';

// ============================================================================
// Helper Functions
// ============================================================================

const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
};

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
};

const getTipeSimpananLabel = (tipe: string): string => {
  const labels: Record<string, string> = {
    pokok: 'Simpanan Pokok',
    wajib: 'Simpanan Wajib',
    sukarela: 'Simpanan Sukarela',
  };
  return labels[tipe] || tipe;
};

const getTipeSimpananColor = (tipe: string): 'primary' | 'success' | 'info' => {
  const colors: Record<string, 'primary' | 'success' | 'info'> = {
    pokok: 'primary',
    wajib: 'success',
    sukarela: 'info',
  };
  return colors[tipe] || 'info';
};

// ============================================================================
// Member Transaction History Component
// ============================================================================

export default function MemberTransactionsPage() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [transactions, setTransactions] = useState<RiwayatTransaksiAnggota[]>([]);
  const [filters, setFilters] = useState<MemberTransactionFilters>({
    tipeSimpanan: 'all',
    tanggalMulai: '',
    tanggalAkhir: '',
  });

  // ============================================================================
  // Fetch Transactions
  // ============================================================================

  const fetchTransactions = async () => {
    try {
      setLoading(true);
      setError('');

      const filterParams: MemberTransactionFilters = {
        page: 1,
        pageSize: 100, // Get more transactions for client-side filtering
      };

      let data = await getMemberTransactions(filterParams);

      // Client-side filtering since backend doesn't support these filters
      if (filters.tipeSimpanan && filters.tipeSimpanan !== 'all') {
        data = data.filter(t => t.tipeSimpanan === filters.tipeSimpanan);
      }

      if (filters.tanggalMulai) {
        data = data.filter(t => t.tanggalTransaksi >= filters.tanggalMulai);
      }

      if (filters.tanggalAkhir) {
        data = data.filter(t => t.tanggalTransaksi <= filters.tanggalAkhir);
      }

      setTransactions(data);
    } catch (err: unknown) {
      console.error('Failed to fetch transactions:', err);
      if (err && typeof err === 'object' && 'message' in err) {
        setError(err.message as string);
      } else {
        setError('Gagal memuat riwayat transaksi');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTransactions();
  }, []);

  // ============================================================================
  // Handle Filter Change
  // ============================================================================

  const handleFilterChange = (field: keyof MemberTransactionFilters, value: string) => {
    setFilters((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleApplyFilters = () => {
    fetchTransactions();
  };

  const handleResetFilters = () => {
    setFilters({
      tipeSimpanan: 'all',
      tanggalMulai: '',
      tanggalAkhir: '',
    });
    setTimeout(() => {
      fetchTransactions();
    }, 0);
  };

  // ============================================================================
  // Calculate Total
  // ============================================================================

  const totalTransaksi = transactions.reduce((sum, t) => sum + t.jumlah, 0);

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Box>
      {/* Page Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" fontWeight={600} gutterBottom>
          Riwayat Transaksi
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Semua transaksi simpanan Anda
        </Typography>
      </Box>

      {/* Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
            <FilterList sx={{ mr: 1 }} />
            <Typography variant="h6">Filter</Typography>
          </Box>

          <Grid container spacing={2}>
            <Grid item xs={12} md={4}>
              <TextField
                select
                fullWidth
                label="Tipe Simpanan"
                value={filters.tipeSimpanan}
                onChange={(e) => handleFilterChange('tipeSimpanan', e.target.value)}
              >
                <MenuItem value="all">Semua Tipe</MenuItem>
                <MenuItem value="pokok">Simpanan Pokok</MenuItem>
                <MenuItem value="wajib">Simpanan Wajib</MenuItem>
                <MenuItem value="sukarela">Simpanan Sukarela</MenuItem>
              </TextField>
            </Grid>

            <Grid item xs={12} md={3}>
              <TextField
                fullWidth
                label="Tanggal Mulai"
                type="date"
                value={filters.tanggalMulai}
                onChange={(e) => handleFilterChange('tanggalMulai', e.target.value)}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>

            <Grid item xs={12} md={3}>
              <TextField
                fullWidth
                label="Tanggal Akhir"
                type="date"
                value={filters.tanggalAkhir}
                onChange={(e) => handleFilterChange('tanggalAkhir', e.target.value)}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>

            <Grid item xs={12} md={2}>
              <Box sx={{ display: 'flex', gap: 1, height: '100%' }}>
                <Button
                  fullWidth
                  variant="contained"
                  onClick={handleApplyFilters}
                  sx={{ height: '56px' }}
                >
                  Terapkan
                </Button>
                <Button
                  variant="outlined"
                  onClick={handleResetFilters}
                  sx={{ height: '56px' }}
                >
                  Reset
                </Button>
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Summary */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <Typography variant="caption" color="text.secondary">
                Total Transaksi
              </Typography>
              <Typography variant="h5" fontWeight={600}>
                {transactions.length} Transaksi
              </Typography>
            </Grid>
            <Grid item xs={12} md={6}>
              <Typography variant="caption" color="text.secondary">
                Total Setoran
              </Typography>
              <Typography variant="h5" fontWeight={600} color="success.main">
                {formatCurrency(totalTransaksi)}
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Transactions Table */}
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
            <Receipt sx={{ mr: 1 }} />
            <Typography variant="h6">Daftar Transaksi</Typography>
          </Box>

          <Divider sx={{ mb: 3 }} />

          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : error ? (
            <Alert severity="error">{error}</Alert>
          ) : transactions.length > 0 ? (
            <TableContainer component={Paper} variant="outlined">
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>No. Referensi</TableCell>
                    <TableCell>Tanggal</TableCell>
                    <TableCell>Tipe Simpanan</TableCell>
                    <TableCell>Keterangan</TableCell>
                    <TableCell align="right">Jumlah</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {transactions.map((transaksi) => (
                    <TableRow key={transaksi.id} hover>
                      <TableCell>
                        <Typography variant="body2" fontFamily="monospace">
                          {transaksi.nomorReferensi}
                        </Typography>
                      </TableCell>
                      <TableCell>{formatDate(transaksi.tanggalTransaksi)}</TableCell>
                      <TableCell>
                        <Chip
                          label={getTipeSimpananLabel(transaksi.tipeSimpanan)}
                          color={getTipeSimpananColor(transaksi.tipeSimpanan)}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>{transaksi.keterangan || '-'}</TableCell>
                      <TableCell align="right">
                        <Typography variant="body2" fontWeight={600} color="success.main">
                          {formatCurrency(transaksi.jumlah)}
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          ) : (
            <Alert severity="info">
              Tidak ada transaksi yang sesuai dengan filter
            </Alert>
          )}
        </CardContent>
      </Card>
    </Box>
  );
}
