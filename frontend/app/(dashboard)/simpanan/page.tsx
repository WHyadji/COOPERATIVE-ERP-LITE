// ============================================================================
// Simpanan List Page - View and manage share capital transactions
// Material-UI table with filters, summary cards, and pagination
// ============================================================================

'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import {
  Box,
  Typography,
  Button,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  IconButton,
  Chip,
  Alert,
  CircularProgress,
  Grid,
  Card,
  CardContent,
} from '@mui/material';
import {
  Add as AddIcon,
  Visibility as VisibilityIcon,
  Assessment as AssessmentIcon,
} from '@mui/icons-material';
import simpananApi from '@/lib/api/simpananApi';
import type { Simpanan, TipeSimpanan, RingkasanSimpanan } from '@/types';
import { format } from 'date-fns';

// ============================================================================
// Simpanan List Page Component
// ============================================================================

export default function SimpananPage() {
  const router = useRouter();
  const [simpanan, setSimpanan] = useState<Simpanan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [ringkasan, setRingkasan] = useState<RingkasanSimpanan | null>(null);

  // Pagination & Filters
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(20);
  const [totalItems, setTotalItems] = useState(0);
  const [tipeFilter, setTipeFilter] = useState<TipeSimpanan | 'all'>('all');
  const [tanggalMulai, setTanggalMulai] = useState('');
  const [tanggalAkhir, setTanggalAkhir] = useState('');
  const [refreshKey, setRefreshKey] = useState(0);

  // ============================================================================
  // Fetch Simpanan with Race Condition Protection
  // ============================================================================

  useEffect(() => {
    let ignore = false;

    const fetchData = async () => {
      try {
        setLoading(true);
        setError('');

        // Fetch transactions
        const response = await simpananApi.getSimpananList({
          page: page + 1,
          pageSize: rowsPerPage,
          tipeSimpanan: tipeFilter,
          tanggalMulai: tanggalMulai || undefined,
          tanggalAkhir: tanggalAkhir || undefined,
        });

        // Fetch summary
        const ringkasanData = await simpananApi.getRingkasan();

        if (!ignore) {
          setSimpanan(response.data);
          setTotalItems(response.pagination.totalItems);
          setRingkasan(ringkasanData);
        }
      } catch (err: unknown) {
        if (!ignore) {
          console.error('Failed to fetch simpanan:', err);
          setError('Gagal memuat data simpanan. Silakan coba lagi.');
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    fetchData();

    return () => {
      ignore = true;
    };
  }, [page, rowsPerPage, tipeFilter, tanggalMulai, tanggalAkhir, refreshKey]);

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleChangePage = (_event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleView = (id: string) => {
    router.push(`/dashboard/simpanan/${id}`);
  };

  const handleFilterChange = () => {
    setPage(0);
    setRefreshKey((prev) => prev + 1);
  };

  // ============================================================================
  // Helper Functions
  // ============================================================================

  const getTipeLabel = (tipe: TipeSimpanan): string => {
    switch (tipe) {
      case 'pokok':
        return 'Simpanan Pokok';
      case 'wajib':
        return 'Simpanan Wajib';
      case 'sukarela':
        return 'Simpanan Sukarela';
      default:
        return tipe;
    }
  };

  const getTipeColor = (
    tipe: TipeSimpanan
  ): 'primary' | 'secondary' | 'success' => {
    switch (tipe) {
      case 'pokok':
        return 'primary';
      case 'wajib':
        return 'secondary';
      case 'sukarela':
        return 'success';
      default:
        return 'primary';
    }
  };

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Box>
      {/* Header */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" fontWeight={600}>
          Simpanan Anggota
        </Typography>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<AssessmentIcon />}
            onClick={() => router.push('/dashboard/simpanan/saldo')}
          >
            Laporan Saldo
          </Button>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => router.push('/dashboard/simpanan/new')}
          >
            Catat Setoran
          </Button>
        </Box>
      </Box>

      {/* Summary Cards */}
      {ringkasan && (
        <Grid container spacing={2} sx={{ mb: 3 }}>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom variant="body2">
                  Simpanan Pokok
                </Typography>
                <Typography variant="h5" fontWeight={600}>
                  {formatCurrency(ringkasan.totalSimpananPokok)}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom variant="body2">
                  Simpanan Wajib
                </Typography>
                <Typography variant="h5" fontWeight={600}>
                  {formatCurrency(ringkasan.totalSimpananWajib)}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom variant="body2">
                  Simpanan Sukarela
                </Typography>
                <Typography variant="h5" fontWeight={600}>
                  {formatCurrency(ringkasan.totalSimpananSukarela)}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card sx={{ bgcolor: 'primary.main', color: 'white' }}>
              <CardContent>
                <Typography color="inherit" gutterBottom variant="body2">
                  Total Simpanan
                </Typography>
                <Typography variant="h5" fontWeight={600} color="inherit">
                  {formatCurrency(ringkasan.totalSemuaSimpanan)}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )}

      {/* Filters */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
          {/* Type Filter */}
          <FormControl size="small" sx={{ minWidth: 180 }}>
            <InputLabel>Tipe Simpanan</InputLabel>
            <Select
              value={tipeFilter}
              label="Tipe Simpanan"
              onChange={(e) => {
                setTipeFilter(e.target.value as TipeSimpanan | 'all');
                setPage(0);
              }}
            >
              <MenuItem value="all">Semua Tipe</MenuItem>
              <MenuItem value="pokok">Simpanan Pokok</MenuItem>
              <MenuItem value="wajib">Simpanan Wajib</MenuItem>
              <MenuItem value="sukarela">Simpanan Sukarela</MenuItem>
            </Select>
          </FormControl>

          {/* Date Range */}
          <TextField
            label="Tanggal Mulai"
            type="date"
            size="small"
            value={tanggalMulai}
            onChange={(e) => setTanggalMulai(e.target.value)}
            InputLabelProps={{ shrink: true }}
            sx={{ minWidth: 160 }}
          />
          <TextField
            label="Tanggal Akhir"
            type="date"
            size="small"
            value={tanggalAkhir}
            onChange={(e) => setTanggalAkhir(e.target.value)}
            InputLabelProps={{ shrink: true }}
            sx={{ minWidth: 160 }}
          />

          <Button variant="outlined" onClick={handleFilterChange}>
            Terapkan Filter
          </Button>
        </Box>
      </Paper>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Simpanan Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Tanggal</TableCell>
                <TableCell>No. Referensi</TableCell>
                <TableCell>No. Anggota</TableCell>
                <TableCell>Nama Anggota</TableCell>
                <TableCell>Tipe Simpanan</TableCell>
                <TableCell align="right">Jumlah Setoran</TableCell>
                <TableCell>Keterangan</TableCell>
                <TableCell align="right">Aksi</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={8} align="center" sx={{ py: 4 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : simpanan.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} align="center" sx={{ py: 4 }}>
                    <Typography color="text.secondary">
                      Tidak ada data simpanan
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                simpanan.map((item) => (
                  <TableRow key={item.id} hover>
                    <TableCell>
                      {format(new Date(item.tanggalTransaksi), 'dd/MM/yyyy')}
                    </TableCell>
                    <TableCell>{item.nomorReferensi}</TableCell>
                    <TableCell>{item.nomorAnggota}</TableCell>
                    <TableCell>{item.namaAnggota}</TableCell>
                    <TableCell>
                      <Chip
                        label={getTipeLabel(item.tipeSimpanan)}
                        color={getTipeColor(item.tipeSimpanan)}
                        size="small"
                      />
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 600 }}>
                      {formatCurrency(item.jumlahSetoran)}
                    </TableCell>
                    <TableCell>{item.keterangan || '-'}</TableCell>
                    <TableCell align="right">
                      <IconButton
                        size="small"
                        onClick={() => handleView(item.id)}
                        title="Lihat Detail"
                      >
                        <VisibilityIcon fontSize="small" />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>

        {/* Pagination */}
        <TablePagination
          rowsPerPageOptions={[10, 20, 50, 100]}
          component="div"
          count={totalItems}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          labelRowsPerPage="Baris per halaman:"
          labelDisplayedRows={({ from, to, count }) =>
            `${from}–${to} dari ${count !== -1 ? count : `lebih dari ${to}`}`
          }
        />
      </Paper>
    </Box>
  );
}
