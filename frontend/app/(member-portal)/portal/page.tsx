// ============================================================================
// Member Portal Dashboard
// Main dashboard showing balance and recent transactions
// ============================================================================

'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  CircularProgress,
  Alert,
  Button,
  Divider,
  Chip,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from '@mui/material';
import {
  AccountBalance,
  TrendingUp,
  Receipt,
  ArrowForward,
  Savings,
} from '@mui/icons-material';
import { getMemberDashboard } from '@/lib/api/memberPortalApi';
import type { MemberDashboardSummary, Simpanan } from '@/lib/api/memberPortalApi';
import { useAuth } from '@/lib/context/AuthContext';

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
// Member Dashboard Component
// ============================================================================

export default function MemberPortalDashboard() {
  const router = useRouter();
  const { user } = useAuth();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [dashboardData, setDashboardData] = useState<MemberDashboardSummary | null>(null);

  // ============================================================================
  // Fetch Dashboard Data
  // ============================================================================

  useEffect(() => {
    const fetchDashboard = async () => {
      try {
        setLoading(true);
        setError('');
        const data = await getMemberDashboard();
        setDashboardData(data);
      } catch (err: unknown) {
        console.error('Failed to fetch dashboard:', err);
        if (err && typeof err === 'object' && 'message' in err) {
          setError(err.message as string);
        } else {
          setError('Gagal memuat data dashboard');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchDashboard();
  }, []);

  // ============================================================================
  // Loading State
  // ============================================================================

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 400 }}>
        <CircularProgress />
      </Box>
    );
  }

  // ============================================================================
  // Error State
  // ============================================================================

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 3 }}>
        {error}
      </Alert>
    );
  }

  // ============================================================================
  // Render Dashboard
  // ============================================================================

  return (
    <Box>
      {/* Welcome Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" fontWeight={600} gutterBottom>
          Selamat Datang, {user?.namaLengkap}
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Berikut adalah ringkasan simpanan dan transaksi Anda
        </Typography>
      </Box>

      {/* Balance Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {/* Total Simpanan */}
        <Grid item xs={12} md={6} lg={3}>
          <Card
            sx={{
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              color: 'white',
              height: '100%',
            }}
          >
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <AccountBalance sx={{ fontSize: 40, mr: 2 }} />
                <Typography variant="h6">Total Simpanan</Typography>
              </Box>
              <Typography variant="h4" fontWeight={700}>
                {formatCurrency(dashboardData?.saldoSimpanan.totalSimpanan || 0)}
              </Typography>
              <Typography variant="body2" sx={{ mt: 1, opacity: 0.9 }}>
                Akumulasi semua simpanan
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Simpanan Pokok */}
        <Grid item xs={12} md={6} lg={3}>
          <Card sx={{ height: '100%' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Savings sx={{ fontSize: 40, mr: 2, color: 'primary.main' }} />
                <Typography variant="h6">Simpanan Pokok</Typography>
              </Box>
              <Typography variant="h4" fontWeight={700} color="primary.main">
                {formatCurrency(dashboardData?.saldoSimpanan.simpananPokok || 0)}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                Setoran pokok
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Simpanan Wajib */}
        <Grid item xs={12} md={6} lg={3}>
          <Card sx={{ height: '100%' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <TrendingUp sx={{ fontSize: 40, mr: 2, color: 'success.main' }} />
                <Typography variant="h6">Simpanan Wajib</Typography>
              </Box>
              <Typography variant="h4" fontWeight={700} color="success.main">
                {formatCurrency(dashboardData?.saldoSimpanan.simpananWajib || 0)}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                Setoran wajib bulanan
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Simpanan Sukarela */}
        <Grid item xs={12} md={6} lg={3}>
          <Card sx={{ height: '100%' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Receipt sx={{ fontSize: 40, mr: 2, color: 'info.main' }} />
                <Typography variant="h6">Simpanan Sukarela</Typography>
              </Box>
              <Typography variant="h4" fontWeight={700} color="info.main">
                {formatCurrency(dashboardData?.saldoSimpanan.simpananSukarela || 0)}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                Setoran sukarela
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Recent Transactions */}
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
            <Typography variant="h6" fontWeight={600}>
              Transaksi Terbaru
            </Typography>
            <Button
              endIcon={<ArrowForward />}
              onClick={() => router.push('/portal/transactions')}
            >
              Lihat Semua
            </Button>
          </Box>

          <Divider sx={{ mb: 3 }} />

          {dashboardData?.transaksiTerbaru && dashboardData.transaksiTerbaru.length > 0 ? (
            <TableContainer component={Paper} variant="outlined">
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Tanggal</TableCell>
                    <TableCell>Tipe Simpanan</TableCell>
                    <TableCell>Keterangan</TableCell>
                    <TableCell align="right">Jumlah</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {dashboardData.transaksiTerbaru.map((transaksi: Simpanan) => (
                    <TableRow key={transaksi.id} hover>
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
                          {formatCurrency(transaksi.jumlahSetoran)}
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          ) : (
            <Alert severity="info">
              Belum ada transaksi
            </Alert>
          )}
        </CardContent>
      </Card>

      {/* Quick Actions */}
      <Grid container spacing={2} sx={{ mt: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Button
            fullWidth
            variant="outlined"
            startIcon={<AccountBalance />}
            onClick={() => router.push('/portal/balance')}
            sx={{ py: 1.5 }}
          >
            Lihat Saldo Detail
          </Button>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Button
            fullWidth
            variant="outlined"
            startIcon={<Receipt />}
            onClick={() => router.push('/portal/transactions')}
            sx={{ py: 1.5 }}
          >
            Riwayat Transaksi
          </Button>
        </Grid>
      </Grid>
    </Box>
  );
}
