// ============================================================================
// Member Balance Detail Page
// Detailed view of member's share capital (Pokok, Wajib, Sukarela)
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
  Grid,
  Divider,
  Paper,
  List,
  ListItem,
  ListItemText,
} from '@mui/material';
import {
  AccountBalance,
  Savings,
  TrendingUp,
  Receipt,
} from '@mui/icons-material';
import { getMemberBalance } from '@/lib/api/memberPortalApi';
import type { SaldoSimpananAnggota } from '@/types';

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

// ============================================================================
// Member Balance Page Component
// ============================================================================

export default function MemberBalancePage() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [balance, setBalance] = useState<SaldoSimpananAnggota | null>(null);

  // ============================================================================
  // Fetch Balance Data
  // ============================================================================

  useEffect(() => {
    const fetchBalance = async () => {
      try {
        setLoading(true);
        setError('');
        const data = await getMemberBalance();
        setBalance(data);
      } catch (err: unknown) {
        console.error('Failed to fetch balance:', err);
        if (err && typeof err === 'object' && 'message' in err) {
          setError(err.message as string);
        } else {
          setError('Gagal memuat data saldo');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchBalance();
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
  // Render Balance
  // ============================================================================

  return (
    <Box>
      {/* Page Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" fontWeight={600} gutterBottom>
          Saldo Simpanan
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Detail saldo simpanan anggota
        </Typography>
      </Box>

      {/* Member Info */}
      {balance && (
        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Grid container spacing={2}>
              <Grid item xs={12} md={4}>
                <Typography variant="caption" color="text.secondary">
                  Nomor Anggota
                </Typography>
                <Typography variant="h6">{balance.nomorAnggota}</Typography>
              </Grid>
              <Grid item xs={12} md={8}>
                <Typography variant="caption" color="text.secondary">
                  Nama Anggota
                </Typography>
                <Typography variant="h6">{balance.namaAnggota}</Typography>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      )}

      {/* Total Balance */}
      <Card
        sx={{
          mb: 3,
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          color: 'white',
        }}
      >
        <CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
            <AccountBalance sx={{ fontSize: 50, mr: 2 }} />
            <Box>
              <Typography variant="h6">Total Simpanan</Typography>
              <Typography variant="h3" fontWeight={700} sx={{ mt: 1 }}>
                {formatCurrency(balance?.totalSimpanan || 0)}
              </Typography>
            </Box>
          </Box>
          <Divider sx={{ bgcolor: 'rgba(255,255,255,0.3)', my: 2 }} />
          <Typography variant="body2" sx={{ opacity: 0.9 }}>
            Akumulasi dari Simpanan Pokok, Wajib, dan Sukarela
          </Typography>
        </CardContent>
      </Card>

      {/* Balance Breakdown */}
      <Grid container spacing={3}>
        {/* Simpanan Pokok */}
        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                <Savings sx={{ fontSize: 40, mr: 2, color: 'primary.main' }} />
                <Typography variant="h6">Simpanan Pokok</Typography>
              </Box>

              <Typography variant="h4" fontWeight={700} color="primary.main" gutterBottom>
                {formatCurrency(balance?.simpananPokok || 0)}
              </Typography>

              <Divider sx={{ my: 2 }} />

              <Paper variant="outlined" sx={{ p: 2, bgcolor: 'grey.50' }}>
                <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                  Tentang Simpanan Pokok:
                </Typography>
                <List dense>
                  <ListItem disablePadding>
                    <ListItemText
                      primary="• Setoran satu kali saat bergabung"
                      primaryTypographyProps={{ variant: 'caption' }}
                    />
                  </ListItem>
                  <ListItem disablePadding>
                    <ListItemText
                      primary="• Tidak dapat ditarik selama menjadi anggota"
                      primaryTypographyProps={{ variant: 'caption' }}
                    />
                  </ListItem>
                  <ListItem disablePadding>
                    <ListItemText
                      primary="• Dikembalikan saat berhenti"
                      primaryTypographyProps={{ variant: 'caption' }}
                    />
                  </ListItem>
                </List>
              </Paper>
            </CardContent>
          </Card>
        </Grid>

        {/* Simpanan Wajib */}
        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                <TrendingUp sx={{ fontSize: 40, mr: 2, color: 'success.main' }} />
                <Typography variant="h6">Simpanan Wajib</Typography>
              </Box>

              <Typography variant="h4" fontWeight={700} color="success.main" gutterBottom>
                {formatCurrency(balance?.simpananWajib || 0)}
              </Typography>

              <Divider sx={{ my: 2 }} />

              <Paper variant="outlined" sx={{ p: 2, bgcolor: 'grey.50' }}>
                <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                  Tentang Simpanan Wajib:
                </Typography>
                <List dense>
                  <ListItem disablePadding>
                    <ListItemText
                      primary="• Setoran rutin setiap bulan"
                      primaryTypographyProps={{ variant: 'caption' }}
                    />
                  </ListItem>
                  <ListItem disablePadding>
                    <ListItemText
                      primary="• Jumlah sesuai aturan koperasi"
                      primaryTypographyProps={{ variant: 'caption' }}
                    />
                  </ListItem>
                  <ListItem disablePadding>
                    <ListItemText
                      primary="• Dapat ditarik sesuai ketentuan"
                      primaryTypographyProps={{ variant: 'caption' }}
                    />
                  </ListItem>
                </List>
              </Paper>
            </CardContent>
          </Card>
        </Grid>

        {/* Simpanan Sukarela */}
        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
                <Receipt sx={{ fontSize: 40, mr: 2, color: 'info.main' }} />
                <Typography variant="h6">Simpanan Sukarela</Typography>
              </Box>

              <Typography variant="h4" fontWeight={700} color="info.main" gutterBottom>
                {formatCurrency(balance?.simpananSukarela || 0)}
              </Typography>

              <Divider sx={{ my: 2 }} />

              <Paper variant="outlined" sx={{ p: 2, bgcolor: 'grey.50' }}>
                <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                  Tentang Simpanan Sukarela:
                </Typography>
                <List dense>
                  <ListItem disablePadding>
                    <ListItemText
                      primary="• Setoran kapan saja, tidak wajib"
                      primaryTypographyProps={{ variant: 'caption' }}
                    />
                  </ListItem>
                  <ListItem disablePadding>
                    <ListItemText
                      primary="• Jumlah sesuai keinginan anggota"
                      primaryTypographyProps={{ variant: 'caption' }}
                    />
                  </ListItem>
                  <ListItem disablePadding>
                    <ListItemText
                      primary="• Dapat ditarik sewaktu-waktu"
                      primaryTypographyProps={{ variant: 'caption' }}
                    />
                  </ListItem>
                </List>
              </Paper>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Additional Info */}
      <Card sx={{ mt: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Informasi Penting
          </Typography>
          <Divider sx={{ mb: 2 }} />
          <Typography variant="body2" color="text.secondary" paragraph>
            • Saldo simpanan diperbarui secara otomatis setiap kali ada transaksi setoran
          </Typography>
          <Typography variant="body2" color="text.secondary" paragraph>
            • Untuk melakukan setoran, silakan hubungi pengurus koperasi atau datang langsung ke kantor koperasi
          </Typography>
          <Typography variant="body2" color="text.secondary">
            • Jika ada perbedaan saldo, segera laporkan ke pengurus koperasi
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
}
