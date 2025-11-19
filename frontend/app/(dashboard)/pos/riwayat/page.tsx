// ============================================================================
// Sales History Page - POS Transaction History
// List of past sales with filters, search, and pagination
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
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
  TablePagination,
  TextField,
  Button,
  IconButton,
  Alert,
  CircularProgress,
  Card,
  CardContent,
  Grid,
  Chip,
} from "@mui/material";
import {
  Visibility as ViewIcon,
  PointOfSale as POSIcon,
  TrendingUp as TrendingUpIcon,
  Receipt as ReceiptIcon,
  AttachMoney as MoneyIcon,
} from "@mui/icons-material";
import { useRouter } from "next/navigation";
import { useToast } from "@/lib/context/ToastContext";
import posApi from "@/lib/api/posApi";
import type { Penjualan, RingkasanPenjualanHariIni } from "@/types";

// ============================================================================
// Sales History Page Component
// ============================================================================

export default function SalesHistoryPage() {
  const router = useRouter();
  const { showError } = useToast();

  // State
  const [sales, setSales] = useState<Penjualan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(20);
  const [totalItems, setTotalItems] = useState(0);
  const [search, setSearch] = useState("");
  const [tanggalMulai, setTanggalMulai] = useState("");
  const [tanggalAkhir, setTanggalAkhir] = useState("");
  const [summary, setSummary] = useState<RingkasanPenjualanHariIni | null>(
    null
  );

  // Currency formatter
  const formatCurrency = (value: number): string => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);
  };

  // Format date
  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("id-ID", {
      day: "2-digit",
      month: "short",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  };

  // Fetch today's summary
  useEffect(() => {
    const fetchSummary = async () => {
      try {
        const data = await posApi.getTodaySummary();
        setSummary(data);
      } catch {
        // Summary is optional, don't show error
      }
    };

    fetchSummary();
  }, []);

  // Fetch sales
  useEffect(() => {
    const fetchSales = async () => {
      setLoading(true);
      setError("");

      try {
        const response = await posApi.getSales({
          page: page + 1,
          pageSize: rowsPerPage,
          search: search || undefined,
          tanggalMulai: tanggalMulai || undefined,
          tanggalAkhir: tanggalAkhir || undefined,
        });

        setSales(response.data);
        setTotalItems(response.pagination.totalItems);
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Gagal memuat riwayat penjualan";
        setError(errorMessage);
        showError(errorMessage);
        setSales([]);
      } finally {
        setLoading(false);
      }
    };

    fetchSales();
  }, [page, rowsPerPage, search, tanggalMulai, tanggalAkhir, showError]);

  // Handle page change
  const handleChangePage = (_event: unknown, newPage: number) => {
    setPage(newPage);
  };

  // Handle rows per page change
  const handleChangeRowsPerPage = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  // Handle view sale
  const handleView = (id: string) => {
    router.push(`/pos/${id}`);
  };

  // Handle back to POS
  const handleBackToPOS = () => {
    router.push("/pos");
  };

  return (
    <Box>
      {/* Header */}
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <Typography variant="h4" fontWeight={700}>
          Riwayat Penjualan
        </Typography>
        <Button
          variant="contained"
          startIcon={<POSIcon />}
          onClick={handleBackToPOS}
        >
          Kembali ke Kasir
        </Button>
      </Box>

      {/* Today's Summary */}
      {summary && (
        <Grid container spacing={2} sx={{ mb: 3 }}>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
                  <ReceiptIcon sx={{ mr: 1, color: "primary.main" }} />
                  <Typography variant="body2" color="text.secondary">
                    Transaksi Hari Ini
                  </Typography>
                </Box>
                <Typography variant="h4" fontWeight={700}>
                  {summary.jumlahTransaksi}
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
                  <TrendingUpIcon sx={{ mr: 1, color: "success.main" }} />
                  <Typography variant="body2" color="text.secondary">
                    Total Penjualan
                  </Typography>
                </Box>
                <Typography variant="h5" fontWeight={700} color="success.main">
                  {formatCurrency(summary.totalPenjualan)}
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
                  <MoneyIcon sx={{ mr: 1, color: "info.main" }} />
                  <Typography variant="body2" color="text.secondary">
                    Total Penerimaan
                  </Typography>
                </Box>
                <Typography variant="h6" fontWeight={700} color="info.main">
                  {formatCurrency(summary.totalPenerimaan)}
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent>
                <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
                  <MoneyIcon sx={{ mr: 1, color: "warning.main" }} />
                  <Typography variant="body2" color="text.secondary">
                    Total Kembalian
                  </Typography>
                </Box>
                <Typography variant="h6" fontWeight={700} color="warning.main">
                  {formatCurrency(summary.totalKembalian)}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )}

      {/* Filters */}
      <Paper sx={{ p: 2, mb: 2 }}>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={4}>
            <TextField
              fullWidth
              label="Cari"
              placeholder="Nomor penjualan, anggota, atau kasir"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              size="small"
            />
          </Grid>
          <Grid item xs={12} sm={4}>
            <TextField
              fullWidth
              label="Tanggal Mulai"
              type="date"
              value={tanggalMulai}
              onChange={(e) => setTanggalMulai(e.target.value)}
              size="small"
              InputLabelProps={{ shrink: true }}
            />
          </Grid>
          <Grid item xs={12} sm={4}>
            <TextField
              fullWidth
              label="Tanggal Akhir"
              type="date"
              value={tanggalAkhir}
              onChange={(e) => setTanggalAkhir(e.target.value)}
              size="small"
              InputLabelProps={{ shrink: true }}
            />
          </Grid>
        </Grid>
      </Paper>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {/* Sales Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Nomor Penjualan</TableCell>
                <TableCell>Tanggal</TableCell>
                <TableCell>Anggota</TableCell>
                <TableCell align="right">Total</TableCell>
                <TableCell>Kasir</TableCell>
                <TableCell align="center">Aksi</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 8 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : sales.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 8 }}>
                    <Typography color="text.secondary">
                      {search || tanggalMulai || tanggalAkhir
                        ? "Tidak ada penjualan yang sesuai dengan filter"
                        : "Belum ada riwayat penjualan"}
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                sales.map((sale) => (
                  <TableRow key={sale.id} hover>
                    <TableCell>
                      <Typography variant="body2" fontWeight={600}>
                        {sale.nomorPenjualan}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {formatDate(sale.tanggalPenjualan)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      {sale.namaAnggota ? (
                        <Box>
                          <Typography variant="body2">
                            {sale.namaAnggota}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            {sale.nomorAnggota}
                          </Typography>
                        </Box>
                      ) : (
                        <Chip label="Guest" size="small" />
                      )}
                    </TableCell>
                    <TableCell align="right">
                      <Typography
                        variant="body2"
                        fontWeight={600}
                        color="primary"
                      >
                        {formatCurrency(sale.totalBelanja)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">{sale.namaKasir}</Typography>
                    </TableCell>
                    <TableCell align="center">
                      <IconButton
                        size="small"
                        color="primary"
                        onClick={() => handleView(sale.id)}
                        title="Lihat Detail"
                      >
                        <ViewIcon />
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
          component="div"
          count={totalItems}
          page={page}
          onPageChange={handleChangePage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          rowsPerPageOptions={[10, 20, 50, 100]}
          labelRowsPerPage="Baris per halaman:"
          labelDisplayedRows={({ from, to, count }) =>
            `${from}-${to} dari ${count !== -1 ? count : `lebih dari ${to}`}`
          }
        />
      </Paper>
    </Box>
  );
}
