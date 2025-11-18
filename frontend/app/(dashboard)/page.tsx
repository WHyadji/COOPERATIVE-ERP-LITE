// ============================================================================
// Dashboard Home Page
// Main dashboard landing page
// ============================================================================

'use client';

import React from 'react';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  CardActions,
  Button,
} from '@mui/material';
import {
  People as PeopleIcon,
  AccountBalance as AccountBalanceIcon,
  PointOfSale as PointOfSaleIcon,
  Assessment as AssessmentIcon,
} from '@mui/icons-material';
import { useAuth } from '@/lib/context/AuthContext';
import { useRouter } from 'next/navigation';

// ============================================================================
// Dashboard Page Component
// ============================================================================

export default function DashboardPage() {
  const { user } = useAuth();
  const router = useRouter();

  // ============================================================================
  // Quick Actions based on role
  // ============================================================================

  const quickActions = [
    {
      title: 'Manajemen Anggota',
      description: 'Kelola data anggota koperasi',
      icon: <PeopleIcon sx={{ fontSize: 40, color: 'primary.main' }} />,
      path: '/dashboard/members',
      roles: ['admin', 'bendahara'],
    },
    {
      title: 'Simpanan',
      description: 'Kelola simpanan anggota',
      icon: <AccountBalanceIcon sx={{ fontSize: 40, color: 'success.main' }} />,
      path: '/dashboard/simpanan',
      roles: ['admin', 'bendahara'],
    },
    {
      title: 'Point of Sale',
      description: 'Transaksi penjualan kasir',
      icon: <PointOfSaleIcon sx={{ fontSize: 40, color: 'info.main' }} />,
      path: '/dashboard/pos',
      roles: ['admin', 'kasir'],
    },
    {
      title: 'Laporan',
      description: 'Lihat laporan keuangan',
      icon: <AssessmentIcon sx={{ fontSize: 40, color: 'warning.main' }} />,
      path: '/dashboard/reports',
      roles: ['admin', 'bendahara'],
    },
  ];

  const visibleActions = quickActions.filter((action) =>
    user ? action.roles.includes(user.peran) : false
  );

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Box>
      {/* Welcome Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" gutterBottom fontWeight={600}>
          Selamat Datang, {user?.namaLengkap}
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Dashboard Cooperative ERP Lite
        </Typography>
      </Box>

      {/* Quick Actions Grid */}
      <Grid container spacing={3}>
        {visibleActions.map((action) => (
          <Grid item xs={12} sm={6} md={4} key={action.path}>
            <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
              <CardContent sx={{ flexGrow: 1 }}>
                <Box sx={{ mb: 2 }}>{action.icon}</Box>
                <Typography variant="h6" gutterBottom>
                  {action.title}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {action.description}
                </Typography>
              </CardContent>
              <CardActions>
                <Button
                  size="small"
                  onClick={() => router.push(action.path)}
                  variant="text"
                >
                  Buka
                </Button>
              </CardActions>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Stats Section - Placeholder */}
      <Box sx={{ mt: 4 }}>
        <Typography variant="h6" gutterBottom>
          Ringkasan
        </Typography>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>
                  Total Anggota
                </Typography>
                <Typography variant="h4">-</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>
                  Transaksi Hari Ini
                </Typography>
                <Typography variant="h4">-</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>
                  Total Simpanan
                </Typography>
                <Typography variant="h4">-</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>
                  Penjualan Bulan Ini
                </Typography>
                <Typography variant="h4">-</Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Box>
    </Box>
  );
}
