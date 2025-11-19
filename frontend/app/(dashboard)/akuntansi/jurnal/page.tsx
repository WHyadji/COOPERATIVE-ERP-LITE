// ============================================================================
// Journal Entry (Transaksi) Page - View and manage journal entries
// Material-UI table with transaction list, filters, and double-entry forms
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import {
  Box,
  Typography,
  Button,
  TextField,
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
} from "@mui/material";
import {
  Add as AddIcon,
  Visibility as VisibilityIcon,
  Delete as DeleteIcon,
  Edit as EditIcon,
  CheckCircle as CheckCircleIcon,
  Warning as WarningIcon,
} from "@mui/icons-material";
import accountingApi from "@/lib/api/accountingApi";
import type { Transaksi, TransaksiListFilters } from "@/types";
import { format, parseISO } from "date-fns";
import TransactionForm from "@/components/accounting/TransactionForm";
import { useToast } from "@/lib/context/ToastContext";

// ============================================================================
// Journal Entry Page Component
// ============================================================================

export default function JournalEntryPage() {
  const router = useRouter();
  const { showSuccess, showError } = useToast();
  const [transactions, setTransactions] = useState<Transaksi[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");

  // Pagination & Filters
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(20);
  const [totalItems, setTotalItems] = useState(0);
  const [tanggalMulai, setTanggalMulai] = useState("");
  const [tanggalAkhir, setTanggalAkhir] = useState("");
  const [refreshKey, setRefreshKey] = useState(0);

  // Transaction form dialog
  const [formOpen, setFormOpen] = useState(false);
  const [selectedTransaction, setSelectedTransaction] =
    useState<Transaksi | null>(null);

  // ============================================================================
  // Fetch Transactions
  // ============================================================================

  useEffect(() => {
    let ignore = false;

    const fetchTransactions = async () => {
      try {
        setLoading(true);
        setError("");

        const filters: TransaksiListFilters = {
          page: page + 1, // API uses 1-based pagination
          pageSize: rowsPerPage,
          tanggalMulai: tanggalMulai || undefined,
          tanggalAkhir: tanggalAkhir || undefined,
        };

        const response = await accountingApi.getTransactions(filters);

        if (!ignore) {
          setTransactions(response.data);
          setTotalItems(response.pagination.totalItems);
        }
      } catch (err: unknown) {
        if (!ignore) {
          console.error("Failed to fetch transactions:", err);
          setError("Gagal memuat data transaksi. Silakan coba lagi.");
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    fetchTransactions();

    return () => {
      ignore = true;
    };
  }, [page, rowsPerPage, tanggalMulai, tanggalAkhir, refreshKey]);

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleCreate = () => {
    setSelectedTransaction(null);
    setFormOpen(true);
  };

  const handleEdit = async (id: string) => {
    try {
      const transaction = await accountingApi.getTransactionById(id);
      setSelectedTransaction(transaction);
      setFormOpen(true);
    } catch (err) {
      console.error("Failed to fetch transaction:", err);
      showError("Gagal memuat data transaksi. Silakan coba lagi.");
    }
  };

  const handleView = (id: string) => {
    router.push(`/akuntansi/jurnal/${id}`);
  };

  const handleDelete = async (id: string, nomorJurnal: string) => {
    if (
      !confirm(`Apakah Anda yakin ingin menghapus transaksi "${nomorJurnal}"?`)
    ) {
      return;
    }

    try {
      await accountingApi.deleteTransaction(id);
      showSuccess(`Transaksi "${nomorJurnal}" berhasil dihapus`);
      setRefreshKey((prev) => prev + 1);
    } catch (err) {
      console.error("Failed to delete transaction:", err);
      showError("Gagal menghapus transaksi. Silakan coba lagi.");
    }
  };

  const handleChangePage = (_event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleCloseForm = () => {
    setFormOpen(false);
  };

  const handleFormSuccess = () => {
    setFormOpen(false);
    setRefreshKey((prev) => prev + 1);
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
  // Render
  // ============================================================================

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
        <Typography variant="h4" fontWeight={600}>
          Jurnal Umum (Journal Entries)
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={handleCreate}
        >
          Tambah Jurnal
        </Button>
      </Box>

      {/* Filters */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap" }}>
          {/* Date Range */}
          <TextField
            label="Tanggal Mulai"
            type="date"
            size="small"
            value={tanggalMulai}
            onChange={(e) => setTanggalMulai(e.target.value)}
            InputLabelProps={{ shrink: true }}
            sx={{ minWidth: 180 }}
          />
          <TextField
            label="Tanggal Akhir"
            type="date"
            size="small"
            value={tanggalAkhir}
            onChange={(e) => setTanggalAkhir(e.target.value)}
            InputLabelProps={{ shrink: true }}
            sx={{ minWidth: 180 }}
          />
          <Button
            variant="outlined"
            onClick={() => {
              setTanggalMulai("");
              setTanggalAkhir("");
            }}
          >
            Reset Filter
          </Button>
        </Box>
      </Paper>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Transactions Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>No. Jurnal</TableCell>
                <TableCell>Tanggal</TableCell>
                <TableCell>Deskripsi</TableCell>
                <TableCell>No. Referensi</TableCell>
                <TableCell align="right">Total Debit</TableCell>
                <TableCell align="right">Total Kredit</TableCell>
                <TableCell align="center">Status</TableCell>
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
              ) : transactions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={8} align="center" sx={{ py: 4 }}>
                    <Typography color="text.secondary">
                      Tidak ada data transaksi
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                transactions.map((transaction) => (
                  <TableRow key={transaction.id} hover>
                    <TableCell>
                      <Typography fontFamily="monospace" fontWeight={600}>
                        {transaction.nomorJurnal}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      {formatDate(transaction.tanggalTransaksi)}
                    </TableCell>
                    <TableCell>{transaction.deskripsi}</TableCell>
                    <TableCell>{transaction.nomorReferensi || "-"}</TableCell>
                    <TableCell align="right">
                      {formatCurrency(transaction.totalDebit)}
                    </TableCell>
                    <TableCell align="right">
                      {formatCurrency(transaction.totalKredit)}
                    </TableCell>
                    <TableCell align="center">
                      {transaction.statusBalanced ? (
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
                    <TableCell align="right">
                      <IconButton
                        size="small"
                        onClick={() => handleView(transaction.id)}
                        title="Lihat Detail"
                      >
                        <VisibilityIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() => handleEdit(transaction.id)}
                        title="Edit"
                        color="primary"
                      >
                        <EditIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() =>
                          handleDelete(transaction.id, transaction.nomorJurnal)
                        }
                        title="Hapus"
                        color="error"
                      >
                        <DeleteIcon fontSize="small" />
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

      {/* Transaction Form Dialog */}
      <TransactionForm
        open={formOpen}
        onClose={handleCloseForm}
        onSuccess={handleFormSuccess}
        transaction={selectedTransaction}
      />
    </Box>
  );
}
