// ============================================================================
// Transaction Detail Page - View full journal entry details
// Shows transaction header, line items, totals, and balance status
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
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
  Chip,
  Alert,
  CircularProgress,
  Breadcrumbs,
  Link,
  Divider,
  Grid,
} from "@mui/material";
import {
  ArrowBack as ArrowBackIcon,
  Delete as DeleteIcon,
  Edit as EditIcon,
  Print as PrintIcon,
  CheckCircle as CheckCircleIcon,
  Warning as WarningIcon,
  Home as HomeIcon,
  Receipt as ReceiptIcon,
} from "@mui/icons-material";
import accountingApi from "@/lib/api/accountingApi";
import type { Transaksi } from "@/types";
import { format, parseISO } from "date-fns";
import { useToast } from "@/lib/context/ToastContext";
import TransactionForm from "@/components/accounting/TransactionForm";

// ============================================================================
// Transaction Detail Page Component
// ============================================================================

export default function TransactionDetailPage() {
  const params = useParams();
  const router = useRouter();
  const transactionId = params.id as string;
  const { showSuccess, showError } = useToast();

  const [transaction, setTransaction] = useState<Transaksi | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");
  const [formOpen, setFormOpen] = useState(false);
  const [refreshKey, setRefreshKey] = useState(0);

  // ============================================================================
  // Fetch Transaction Detail
  // ============================================================================

  useEffect(() => {
    let ignore = false;

    const fetchTransaction = async () => {
      try {
        setLoading(true);
        setError("");

        const data = await accountingApi.getTransactionById(transactionId);

        if (!ignore) {
          setTransaction(data);
        }
      } catch (err: unknown) {
        if (!ignore) {
          console.error("Failed to fetch transaction:", err);
          setError("Gagal memuat detail transaksi. Silakan coba lagi.");
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    fetchTransaction();

    return () => {
      ignore = true;
    };
  }, [transactionId, refreshKey]);

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleBack = () => {
    router.push("/akuntansi/jurnal");
  };

  const handleEdit = () => {
    setFormOpen(true);
  };

  const handleCloseForm = () => {
    setFormOpen(false);
  };

  const handleFormSuccess = () => {
    setFormOpen(false);
    setRefreshKey((prev) => prev + 1);
  };

  const handleDelete = async () => {
    if (!transaction) return;

    if (
      !confirm(
        `Apakah Anda yakin ingin menghapus transaksi "${transaction.nomorJurnal}"?`
      )
    ) {
      return;
    }

    try {
      await accountingApi.deleteTransaction(transaction.id);
      showSuccess(`Transaksi "${transaction.nomorJurnal}" berhasil dihapus`);
      router.push("/akuntansi/jurnal");
    } catch (err) {
      console.error("Failed to delete transaction:", err);
      showError("Gagal menghapus transaksi. Silakan coba lagi.");
    }
  };

  const handlePrint = () => {
    window.print();
  };

  // ============================================================================
  // Helper Functions
  // ============================================================================

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const formatDate = (dateString: string): string => {
    try {
      return format(parseISO(dateString), "dd/MM/yyyy");
    } catch {
      return dateString;
    }
  };

  // ============================================================================
  // Render Loading/Error States
  // ============================================================================

  if (loading) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          minHeight: "400px",
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  if (error || !transaction) {
    return (
      <Box>
        <Alert severity="error" sx={{ mb: 3 }}>
          {error || "Transaksi tidak ditemukan"}
        </Alert>
        <Button
          variant="outlined"
          startIcon={<ArrowBackIcon />}
          onClick={handleBack}
        >
          Kembali ke Daftar Jurnal
        </Button>
      </Box>
    );
  }

  // ============================================================================
  // Render Transaction Detail
  // ============================================================================

  return (
    <Box>
      {/* Breadcrumbs */}
      <Breadcrumbs sx={{ mb: 3 }}>
        <Link
          href="/dashboard"
          underline="hover"
          sx={{ display: "flex", alignItems: "center", cursor: "pointer" }}
          onClick={(e) => {
            e.preventDefault();
            router.push("/dashboard");
          }}
        >
          <HomeIcon sx={{ mr: 0.5 }} fontSize="small" />
          Dashboard
        </Link>
        <Link
          href="/akuntansi/jurnal"
          underline="hover"
          sx={{ display: "flex", alignItems: "center", cursor: "pointer" }}
          onClick={(e) => {
            e.preventDefault();
            router.push("/akuntansi/jurnal");
          }}
        >
          <ReceiptIcon sx={{ mr: 0.5 }} fontSize="small" />
          Jurnal Umum
        </Link>
        <Typography
          sx={{ display: "flex", alignItems: "center" }}
          color="text.primary"
        >
          Detail #{transaction.nomorJurnal}
        </Typography>
      </Breadcrumbs>

      {/* Header with Actions */}
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <Typography variant="h4" fontWeight={600}>
          Detail Transaksi
        </Typography>
        <Box sx={{ display: "flex", gap: 1 }} className="no-print">
          <Button
            variant="outlined"
            startIcon={<ArrowBackIcon />}
            onClick={handleBack}
          >
            Kembali
          </Button>
          <Button
            variant="outlined"
            color="primary"
            startIcon={<EditIcon />}
            onClick={handleEdit}
          >
            Edit
          </Button>
          <Button
            variant="outlined"
            color="error"
            startIcon={<DeleteIcon />}
            onClick={handleDelete}
          >
            Hapus
          </Button>
          <Button
            variant="contained"
            startIcon={<PrintIcon />}
            onClick={handlePrint}
          >
            Cetak
          </Button>
        </Box>
      </Box>

      {/* Transaction Header Information */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <Box sx={{ mb: 2 }}>
              <Typography
                variant="caption"
                color="text.secondary"
                display="block"
              >
                Nomor Jurnal
              </Typography>
              <Typography variant="h6" fontWeight={600} fontFamily="monospace">
                {transaction.nomorJurnal}
              </Typography>
            </Box>
            <Box sx={{ mb: 2 }}>
              <Typography
                variant="caption"
                color="text.secondary"
                display="block"
              >
                Tanggal Transaksi
              </Typography>
              <Typography variant="body1">
                {formatDate(transaction.tanggalTransaksi)}
              </Typography>
            </Box>
            <Box sx={{ mb: 2 }}>
              <Typography
                variant="caption"
                color="text.secondary"
                display="block"
              >
                Deskripsi
              </Typography>
              <Typography variant="body1">{transaction.deskripsi}</Typography>
            </Box>
          </Grid>
          <Grid item xs={12} md={6}>
            <Box sx={{ mb: 2 }}>
              <Typography
                variant="caption"
                color="text.secondary"
                display="block"
              >
                Nomor Referensi
              </Typography>
              <Typography variant="body1">
                {transaction.nomorReferensi || "-"}
              </Typography>
            </Box>
            {transaction.tipeTransaksi && (
              <Box sx={{ mb: 2 }}>
                <Typography
                  variant="caption"
                  color="text.secondary"
                  display="block"
                >
                  Tipe Transaksi
                </Typography>
                <Typography variant="body1">
                  {transaction.tipeTransaksi}
                </Typography>
              </Box>
            )}
            <Box sx={{ mb: 2 }}>
              <Typography
                variant="caption"
                color="text.secondary"
                display="block"
              >
                Status Balance
              </Typography>
              {transaction.statusBalanced ? (
                <Chip
                  icon={<CheckCircleIcon />}
                  label="Balanced (Debit = Kredit)"
                  color="success"
                  size="medium"
                />
              ) : (
                <Chip
                  icon={<WarningIcon />}
                  label="Unbalanced (Debit ≠ Kredit)"
                  color="error"
                  size="medium"
                />
              )}
            </Box>
          </Grid>
        </Grid>
      </Paper>

      {/* Line Items Table */}
      <Paper sx={{ mb: 3 }}>
        <Box sx={{ p: 2, borderBottom: "1px solid rgba(0, 0, 0, 0.12)" }}>
          <Typography variant="h6" fontWeight={600}>
            Baris Transaksi (Journal Entry Lines)
          </Typography>
        </Box>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Kode Akun</TableCell>
                <TableCell>Nama Akun</TableCell>
                <TableCell>Keterangan</TableCell>
                <TableCell align="right">Debit</TableCell>
                <TableCell align="right">Kredit</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {transaction.barisTransaksi &&
              transaction.barisTransaksi.length > 0 ? (
                transaction.barisTransaksi.map((line, index) => (
                  <TableRow key={line.id || index} hover>
                    <TableCell>
                      <Typography fontFamily="monospace" fontWeight={600}>
                        {line.kodeAkun || "-"}
                      </Typography>
                    </TableCell>
                    <TableCell>{line.namaAkun || "-"}</TableCell>
                    <TableCell>{line.keterangan || "-"}</TableCell>
                    <TableCell align="right">
                      {line.jumlahDebit > 0 ? (
                        <Typography fontWeight={600}>
                          {formatCurrency(line.jumlahDebit)}
                        </Typography>
                      ) : (
                        <Typography color="text.secondary">-</Typography>
                      )}
                    </TableCell>
                    <TableCell align="right">
                      {line.jumlahKredit > 0 ? (
                        <Typography fontWeight={600}>
                          {formatCurrency(line.jumlahKredit)}
                        </Typography>
                      ) : (
                        <Typography color="text.secondary">-</Typography>
                      )}
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={5} align="center" sx={{ py: 4 }}>
                    <Typography color="text.secondary">
                      Tidak ada baris transaksi
                    </Typography>
                  </TableCell>
                </TableRow>
              )}
              {/* Totals Row */}
              <TableRow sx={{ backgroundColor: "rgba(0, 0, 0, 0.04)" }}>
                <TableCell colSpan={3}>
                  <Typography variant="h6" fontWeight={700}>
                    TOTAL
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  <Typography variant="h6" fontWeight={700} color="primary">
                    {formatCurrency(transaction.totalDebit)}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  <Typography variant="h6" fontWeight={700} color="error">
                    {formatCurrency(transaction.totalKredit)}
                  </Typography>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* Summary Section */}
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" fontWeight={600} gutterBottom>
          Ringkasan Transaksi
        </Typography>
        <Divider sx={{ mb: 2 }} />
        <Grid container spacing={2}>
          <Grid item xs={12} md={4}>
            <Box
              sx={{
                p: 2,
                borderRadius: 1,
                backgroundColor: "primary.light",
                color: "primary.contrastText",
              }}
            >
              <Typography variant="caption" display="block">
                Total Debit
              </Typography>
              <Typography variant="h5" fontWeight={700}>
                {formatCurrency(transaction.totalDebit)}
              </Typography>
            </Box>
          </Grid>
          <Grid item xs={12} md={4}>
            <Box
              sx={{
                p: 2,
                borderRadius: 1,
                backgroundColor: "error.light",
                color: "error.contrastText",
              }}
            >
              <Typography variant="caption" display="block">
                Total Kredit
              </Typography>
              <Typography variant="h5" fontWeight={700}>
                {formatCurrency(transaction.totalKredit)}
              </Typography>
            </Box>
          </Grid>
          <Grid item xs={12} md={4}>
            <Box
              sx={{
                p: 2,
                borderRadius: 1,
                backgroundColor: transaction.statusBalanced
                  ? "success.light"
                  : "warning.light",
                color: transaction.statusBalanced
                  ? "success.contrastText"
                  : "warning.contrastText",
              }}
            >
              <Typography variant="caption" display="block">
                Selisih (Difference)
              </Typography>
              <Typography variant="h5" fontWeight={700}>
                {formatCurrency(
                  Math.abs(transaction.totalDebit - transaction.totalKredit)
                )}
              </Typography>
            </Box>
          </Grid>
        </Grid>

        {/* Balance Status Info */}
        {transaction.statusBalanced ? (
          <Alert severity="success" sx={{ mt: 3 }}>
            ✓ Transaksi ini <strong>balanced</strong> (seimbang). Total debit
            sama dengan total kredit.
          </Alert>
        ) : (
          <Alert severity="error" sx={{ mt: 3 }}>
            ✗ Transaksi ini <strong>unbalanced</strong> (tidak seimbang). Total
            debit tidak sama dengan total kredit. Perbedaan:{" "}
            {formatCurrency(
              Math.abs(transaction.totalDebit - transaction.totalKredit)
            )}
          </Alert>
        )}
      </Paper>

      {/* Metadata & Audit Trail */}
      <Box
        sx={{
          mt: 3,
          p: 3,
          backgroundColor: "rgba(0, 0, 0, 0.02)",
          borderRadius: 1,
        }}
      >
        <Typography variant="subtitle2" fontWeight={600} gutterBottom>
          Informasi Audit Trail
        </Typography>
        <Divider sx={{ mb: 2 }} />

        <Typography
          variant="caption"
          color="text.secondary"
          display="block"
          sx={{ mb: 1 }}
        >
          ID Transaksi: {transaction.id}
        </Typography>

        {transaction.namaDibuatOleh && transaction.tanggalDibuat && (
          <Typography
            variant="body2"
            color="text.secondary"
            display="block"
            sx={{ mb: 1 }}
          >
            Dibuat oleh <strong>{transaction.namaDibuatOleh}</strong> pada{" "}
            {formatDate(transaction.tanggalDibuat)}
          </Typography>
        )}

        {transaction.namaDiperbaruiOleh && transaction.tanggalDiperbarui && (
          <Typography variant="body2" color="text.secondary" display="block">
            Terakhir diperbarui oleh{" "}
            <strong>{transaction.namaDiperbaruiOleh}</strong> pada{" "}
            {formatDate(transaction.tanggalDiperbarui)}
          </Typography>
        )}

        {!transaction.namaDibuatOleh && transaction.tanggalDibuat && (
          <Typography variant="caption" color="text.secondary" display="block">
            Dibuat: {formatDate(transaction.tanggalDibuat)}
          </Typography>
        )}

        {!transaction.namaDiperbaruiOleh && transaction.tanggalDiperbarui && (
          <Typography variant="caption" color="text.secondary" display="block">
            Diperbarui: {formatDate(transaction.tanggalDiperbarui)}
          </Typography>
        )}
      </Box>

      {/* Print Styles */}
      <style jsx global>{`
        @media print {
          .no-print {
            display: none !important;
          }
          body {
            background: white;
          }
        }
      `}</style>

      {/* Transaction Form Dialog */}
      <TransactionForm
        open={formOpen}
        onClose={handleCloseForm}
        onSuccess={handleFormSuccess}
        transaction={transaction}
      />
    </Box>
  );
}
