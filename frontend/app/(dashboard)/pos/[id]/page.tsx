// ============================================================================
// Sale Detail Page - View Individual Sale/Receipt
// Displays complete sale details with print functionality
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
import {
  Box,
  Paper,
  Typography,
  Button,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableRow,
  Breadcrumbs,
  Link,
  Alert,
  CircularProgress,
  Chip,
} from "@mui/material";
import {
  ArrowBack as BackIcon,
  Print as PrintIcon,
  Receipt as ReceiptIcon,
} from "@mui/icons-material";
import { useRouter, useParams } from "next/navigation";
import { useToast } from "@/lib/context/ToastContext";
import posApi from "@/lib/api/posApi";
import type { Penjualan } from "@/types";

// ============================================================================
// Sale Detail Page Component
// ============================================================================

export default function SaleDetailPage() {
  const router = useRouter();
  const params = useParams();
  const { showError } = useToast();
  const saleId = params.id as string;

  // State
  const [sale, setSale] = useState<Penjualan | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

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
      month: "long",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    }).format(date);
  };

  // Fetch sale details
  useEffect(() => {
    const fetchSale = async () => {
      setLoading(true);
      setError("");

      try {
        const data = await posApi.getSaleById(saleId);
        setSale(data);
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Gagal memuat detail penjualan";
        setError(errorMessage);
        showError(errorMessage);
      } finally {
        setLoading(false);
      }
    };

    if (saleId) {
      fetchSale();
    }
  }, [saleId, showError]);

  // Handle print
  const handlePrint = () => {
    window.print();
  };

  // Handle back
  const handleBack = () => {
    router.push("/pos/riwayat");
  };

  // Loading state
  if (loading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  // Error state
  if (error || !sale) {
    return (
      <Box>
        <Alert severity="error" sx={{ mb: 2 }}>
          {error || "Data penjualan tidak ditemukan"}
        </Alert>
        <Button
          variant="outlined"
          startIcon={<BackIcon />}
          onClick={handleBack}
        >
          Kembali ke Riwayat
        </Button>
      </Box>
    );
  }

  return (
    <Box>
      {/* Breadcrumbs */}
      <Breadcrumbs sx={{ mb: 2, "@media print": { display: "none" } }}>
        <Link
          underline="hover"
          color="inherit"
          href="/pos"
          onClick={(e) => {
            e.preventDefault();
            router.push("/pos");
          }}
        >
          Kasir
        </Link>
        <Link
          underline="hover"
          color="inherit"
          href="/pos/riwayat"
          onClick={(e) => {
            e.preventDefault();
            router.push("/pos/riwayat");
          }}
        >
          Riwayat Penjualan
        </Link>
        <Typography color="text.primary">{sale.nomorPenjualan}</Typography>
      </Breadcrumbs>

      {/* Header */}
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
          "@media print": { display: "none" },
        }}
      >
        <Typography variant="h4" fontWeight={700}>
          Detail Penjualan
        </Typography>
        <Box sx={{ display: "flex", gap: 1 }}>
          <Button
            variant="outlined"
            startIcon={<BackIcon />}
            onClick={handleBack}
          >
            Kembali
          </Button>
          <Button
            variant="contained"
            startIcon={<PrintIcon />}
            onClick={handlePrint}
          >
            Cetak Struk
          </Button>
        </Box>
      </Box>

      {/* Receipt Content */}
      <Paper
        sx={{
          maxWidth: 600,
          mx: "auto",
          p: 4,
          "@media print": { boxShadow: "none", p: 2 },
        }}
      >
        {/* Receipt Header */}
        <Box sx={{ textAlign: "center", mb: 4 }}>
          <ReceiptIcon sx={{ fontSize: 48, color: "primary.main", mb: 1 }} />
          <Typography variant="h5" fontWeight={700} gutterBottom>
            STRUK PENJUALAN
          </Typography>
          <Typography variant="h6" color="primary" gutterBottom>
            {sale.nomorPenjualan}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {formatDate(sale.tanggalPenjualan)}
          </Typography>
        </Box>

        <Divider sx={{ mb: 3 }} />

        {/* Sale Info */}
        <Box sx={{ mb: 3 }}>
          <Table size="small">
            <TableBody>
              {sale.namaAnggota && (
                <>
                  <TableRow>
                    <TableCell sx={{ borderBottom: "none", py: 0.5 }}>
                      <Typography variant="body2" color="text.secondary">
                        Anggota
                      </Typography>
                    </TableCell>
                    <TableCell sx={{ borderBottom: "none", py: 0.5 }}>
                      <Typography variant="body2" fontWeight={600}>
                        {sale.namaAnggota}
                      </Typography>
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell sx={{ borderBottom: "none", py: 0.5 }}>
                      <Typography variant="body2" color="text.secondary">
                        Nomor Anggota
                      </Typography>
                    </TableCell>
                    <TableCell sx={{ borderBottom: "none", py: 0.5 }}>
                      <Typography variant="body2">
                        {sale.nomorAnggota}
                      </Typography>
                    </TableCell>
                  </TableRow>
                </>
              )}
              <TableRow>
                <TableCell sx={{ borderBottom: "none", py: 0.5 }}>
                  <Typography variant="body2" color="text.secondary">
                    Kasir
                  </Typography>
                </TableCell>
                <TableCell sx={{ borderBottom: "none", py: 0.5 }}>
                  <Typography variant="body2">{sale.namaKasir}</Typography>
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell sx={{ borderBottom: "none", py: 0.5 }}>
                  <Typography variant="body2" color="text.secondary">
                    Metode Bayar
                  </Typography>
                </TableCell>
                <TableCell sx={{ borderBottom: "none", py: 0.5 }}>
                  <Chip label="Tunai" size="small" color="success" />
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </Box>

        <Divider sx={{ mb: 3 }} />

        {/* Items */}
        <Box sx={{ mb: 3 }}>
          <Typography variant="subtitle2" fontWeight={600} gutterBottom>
            Daftar Produk
          </Typography>
          <Table size="small">
            <TableBody>
              {sale.itemPenjualan.map((item, index) => (
                <React.Fragment key={item.id}>
                  <TableRow>
                    <TableCell
                      colSpan={3}
                      sx={{ pb: 0.5, borderBottom: "none" }}
                    >
                      <Typography variant="body2" fontWeight={600}>
                        {item.namaProduk}
                      </Typography>
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell
                      sx={{
                        pt: 0,
                        pb: index === sale.itemPenjualan.length - 1 ? 1 : 2,
                        borderBottom: "none",
                      }}
                    >
                      <Typography variant="caption" color="text.secondary">
                        {item.kuantitas} × {formatCurrency(item.hargaSatuan)}
                      </Typography>
                    </TableCell>
                    <TableCell
                      sx={{
                        pt: 0,
                        pb: index === sale.itemPenjualan.length - 1 ? 1 : 2,
                        borderBottom: "none",
                      }}
                    />
                    <TableCell
                      align="right"
                      sx={{
                        pt: 0,
                        pb: index === sale.itemPenjualan.length - 1 ? 1 : 2,
                        borderBottom: "none",
                      }}
                    >
                      <Typography variant="body2" fontWeight={600}>
                        {formatCurrency(item.subtotal)}
                      </Typography>
                    </TableCell>
                  </TableRow>
                </React.Fragment>
              ))}
            </TableBody>
          </Table>
        </Box>

        <Divider sx={{ mb: 3 }} />

        {/* Totals */}
        <Box sx={{ mb: 3 }}>
          <Box sx={{ display: "flex", justifyContent: "space-between", mb: 2 }}>
            <Typography variant="h6" fontWeight={700}>
              TOTAL
            </Typography>
            <Typography variant="h5" fontWeight={700} color="primary">
              {formatCurrency(sale.totalBelanja)}
            </Typography>
          </Box>

          <Box sx={{ display: "flex", justifyContent: "space-between", mb: 1 }}>
            <Typography variant="body2" color="text.secondary">
              Bayar
            </Typography>
            <Typography variant="body2">
              {formatCurrency(sale.jumlahBayar)}
            </Typography>
          </Box>

          <Box sx={{ display: "flex", justifyContent: "space-between" }}>
            <Typography variant="body2" color="text.secondary">
              Kembalian
            </Typography>
            <Typography variant="body1" fontWeight={600} color="success.main">
              {formatCurrency(sale.kembalian)}
            </Typography>
          </Box>
        </Box>

        {/* Notes */}
        {sale.catatan && (
          <>
            <Divider sx={{ mb: 2 }} />
            <Box sx={{ mb: 3 }}>
              <Typography variant="caption" color="text.secondary">
                Catatan:
              </Typography>
              <Typography variant="body2">{sale.catatan}</Typography>
            </Box>
          </>
        )}

        <Divider sx={{ mb: 2 }} />

        {/* Footer */}
        <Box sx={{ textAlign: "center" }}>
          <Typography
            variant="caption"
            color="text.secondary"
            display="block"
            sx={{ mb: 1 }}
          >
            Terima kasih atas kunjungan Anda
          </Typography>
          <Typography variant="caption" color="text.secondary" display="block">
            --- Simpan struk ini sebagai bukti pembayaran ---
          </Typography>
        </Box>
      </Paper>
    </Box>
  );
}
