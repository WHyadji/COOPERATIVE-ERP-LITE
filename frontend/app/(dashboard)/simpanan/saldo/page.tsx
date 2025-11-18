// ============================================================================
// Simpanan Balance Report Page - View member balances
// Shows breakdown of Pokok, Wajib, Sukarela for all members
// ============================================================================

'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import {
  Box,
  Typography,
  Button,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Alert,
  CircularProgress,
  Card,
  CardContent,
  Grid,
  TextField,
  InputAdornment,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  Search as SearchIcon,
  Download as DownloadIcon,
} from '@mui/icons-material';
import simpananApi from '@/lib/api/simpananApi';
import type { SaldoSimpananAnggota } from '@/types';

// ============================================================================
// Balance Report Page Component
// ============================================================================

export default function SaldoSimpananPage() {
  const router = useRouter();
  const [balances, setBalances] = useState<SaldoSimpananAnggota[]>([]);
  const [filteredBalances, setFilteredBalances] = useState<SaldoSimpananAnggota[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState('');

  // ============================================================================
  // Fetch Balance Report
  // ============================================================================

  useEffect(() => {
    const fetchBalances = async () => {
      try {
        setLoading(true);
        setError('');

        const data = await simpananApi.getLaporanSaldo();
        setBalances(data);
        setFilteredBalances(data);
      } catch (err: unknown) {
        console.error('Failed to fetch balance report:', err);
        setError('Gagal memuat laporan saldo. Silakan coba lagi.');
      } finally {
        setLoading(false);
      }
    };

    fetchBalances();
  }, []);

  // ============================================================================
  // Search Filter
  // ============================================================================

  useEffect(() => {
    if (!searchQuery) {
      setFilteredBalances(balances);
      return;
    }

    const query = searchQuery.toLowerCase();
    const filtered = balances.filter(
      (balance) =>
        balance.namaAnggota.toLowerCase().includes(query) ||
        balance.nomorAnggota.toLowerCase().includes(query)
    );
    setFilteredBalances(filtered);
  }, [searchQuery, balances]);

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

  const calculateTotals = () => {
    return filteredBalances.reduce(
      (acc, balance) => ({
        pokok: acc.pokok + balance.simpananPokok,
        wajib: acc.wajib + balance.simpananWajib,
        sukarela: acc.sukarela + balance.simpananSukarela,
        total: acc.total + balance.totalSimpanan,
      }),
      { pokok: 0, wajib: 0, sukarela: 0, total: 0 }
    );
  };

  const totals = calculateTotals();

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Box>
      {/* Header */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box>
          <Button
            startIcon={<ArrowBackIcon />}
            onClick={() => router.push('/dashboard/simpanan')}
            sx={{ mb: 1 }}
          >
            Kembali
          </Button>
          <Typography variant="h4" fontWeight={600}>
            Laporan Saldo Simpanan
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Ringkasan saldo simpanan semua anggota
          </Typography>
        </Box>
        <Button
          variant="outlined"
          startIcon={<DownloadIcon />}
          onClick={() => alert('Fitur export akan segera hadir')}
        >
          Export Excel
        </Button>
      </Box>

      {/* Summary Cards */}
      <Grid container spacing={2} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom variant="body2">
                Total Simpanan Pokok
              </Typography>
              <Typography variant="h5" fontWeight={600}>
                {formatCurrency(totals.pokok)}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {filteredBalances.length} anggota
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom variant="body2">
                Total Simpanan Wajib
              </Typography>
              <Typography variant="h5" fontWeight={600}>
                {formatCurrency(totals.wajib)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="text.secondary" gutterBottom variant="body2">
                Total Simpanan Sukarela
              </Typography>
              <Typography variant="h5" fontWeight={600}>
                {formatCurrency(totals.sukarela)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ bgcolor: 'primary.main', color: 'white' }}>
            <CardContent>
              <Typography color="inherit" gutterBottom variant="body2">
                Grand Total
              </Typography>
              <Typography variant="h5" fontWeight={600} color="inherit">
                {formatCurrency(totals.total)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Search */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <TextField
          fullWidth
          label="Cari Anggota"
          variant="outlined"
          size="small"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Cari berdasarkan nama atau nomor anggota..."
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon />
              </InputAdornment>
            ),
          }}
        />
      </Paper>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Balance Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>No. Anggota</TableCell>
                <TableCell>Nama Anggota</TableCell>
                <TableCell align="right">Simpanan Pokok</TableCell>
                <TableCell align="right">Simpanan Wajib</TableCell>
                <TableCell align="right">Simpanan Sukarela</TableCell>
                <TableCell align="right" sx={{ fontWeight: 600 }}>
                  Total Simpanan
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : filteredBalances.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                    <Typography color="text.secondary">
                      {searchQuery
                        ? 'Tidak ada data yang sesuai dengan pencarian'
                        : 'Tidak ada data saldo simpanan'}
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                <>
                  {filteredBalances.map((balance) => (
                    <TableRow key={balance.idAnggota} hover>
                      <TableCell>{balance.nomorAnggota}</TableCell>
                      <TableCell>{balance.namaAnggota}</TableCell>
                      <TableCell align="right">
                        {formatCurrency(balance.simpananPokok)}
                      </TableCell>
                      <TableCell align="right">
                        {formatCurrency(balance.simpananWajib)}
                      </TableCell>
                      <TableCell align="right">
                        {formatCurrency(balance.simpananSukarela)}
                      </TableCell>
                      <TableCell
                        align="right"
                        sx={{ fontWeight: 600, bgcolor: 'action.hover' }}
                      >
                        {formatCurrency(balance.totalSimpanan)}
                      </TableCell>
                    </TableRow>
                  ))}
                  {/* Totals Row */}
                  <TableRow sx={{ bgcolor: 'primary.light' }}>
                    <TableCell colSpan={2} sx={{ fontWeight: 600 }}>
                      TOTAL
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 600 }}>
                      {formatCurrency(totals.pokok)}
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 600 }}>
                      {formatCurrency(totals.wajib)}
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 600 }}>
                      {formatCurrency(totals.sukarela)}
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 700, fontSize: '1.1rem' }}>
                      {formatCurrency(totals.total)}
                    </TableCell>
                  </TableRow>
                </>
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
    </Box>
  );
}
