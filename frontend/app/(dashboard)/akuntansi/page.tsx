// ============================================================================
// Chart of Accounts Page - View and manage chart of accounts
// Material-UI table with account hierarchy, filters, and CRUD operations
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import {
  Box,
  Typography,
  Button,
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
  IconButton,
  Chip,
  Alert,
  CircularProgress,
} from "@mui/material";
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  PlayArrowOutlined as SeedIcon,
  AccountBalance as AccountBalanceIcon,
} from "@mui/icons-material";
import accountingApi from "@/lib/api/accountingApi";
import type { Akun, TipeAkun } from "@/types";
import AccountForm from "@/components/accounting/AccountForm";

// ============================================================================
// Chart of Accounts Page Component
// ============================================================================

export default function ChartOfAccountsPage() {
  const router = useRouter();
  const [accounts, setAccounts] = useState<Akun[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");

  // Filters
  const [tipeFilter, setTipeFilter] = useState<TipeAkun | "all">("all");
  const [statusFilter, setStatusFilter] = useState<boolean | undefined>(true);
  const [refreshKey, setRefreshKey] = useState(0);

  // Account form dialog
  const [formOpen, setFormOpen] = useState(false);
  const [currentAccount, setCurrentAccount] = useState<Akun | null>(null);

  // ============================================================================
  // Fetch Accounts
  // ============================================================================

  useEffect(() => {
    let ignore = false;

    const fetchAccounts = async () => {
      try {
        setLoading(true);
        setError("");

        const data = await accountingApi.getAccounts(tipeFilter, statusFilter);

        if (!ignore) {
          setAccounts(data);
        }
      } catch (err: unknown) {
        if (!ignore) {
          console.error("Failed to fetch accounts:", err);
          setError("Gagal memuat data akun. Silakan coba lagi.");
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    fetchAccounts();

    return () => {
      ignore = true;
    };
  }, [tipeFilter, statusFilter, refreshKey]);

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleCreate = () => {
    setCurrentAccount(null);
    setFormOpen(true);
  };

  const handleEdit = (account: Akun) => {
    setCurrentAccount(account);
    setFormOpen(true);
  };

  const handleDelete = async (id: string, namaAkun: string) => {
    if (!confirm(`Apakah Anda yakin ingin menghapus akun "${namaAkun}"?`)) {
      return;
    }

    try {
      await accountingApi.deleteAccount(id);
      setRefreshKey((prev) => prev + 1);
    } catch (err) {
      console.error("Failed to delete account:", err);
      alert("Gagal menghapus akun. Silakan coba lagi.");
    }
  };

  const handleSeedCOA = async () => {
    if (
      !confirm(
        "Apakah Anda yakin ingin membuat Chart of Accounts default? Ini akan menambahkan akun-akun standar koperasi."
      )
    ) {
      return;
    }

    try {
      await accountingApi.seedDefaultCOA();
      setRefreshKey((prev) => prev + 1);
      alert("Chart of Accounts default berhasil dibuat!");
    } catch (err) {
      console.error("Failed to seed COA:", err);
      alert("Gagal membuat Chart of Accounts. Silakan coba lagi.");
    }
  };

  const handleViewLedger = (id: string) => {
    router.push(`/akuntansi/ledger/${id}`);
  };

  const handleCloseForm = () => {
    setFormOpen(false);
    setCurrentAccount(null);
  };

  const handleFormSuccess = () => {
    setFormOpen(false);
    setCurrentAccount(null);
    setRefreshKey((prev) => prev + 1);
  };

  // ============================================================================
  // Helper Functions
  // ============================================================================

  const getTipeAkunLabel = (tipe: TipeAkun): string => {
    const labels: Record<TipeAkun, string> = {
      aset: "Aset",
      kewajiban: "Kewajiban",
      modal: "Modal",
      pendapatan: "Pendapatan",
      beban: "Beban",
    };
    return labels[tipe];
  };

  const getTipeAkunColor = (
    tipe: TipeAkun
  ): "primary" | "error" | "success" | "warning" | "info" => {
    const colors: Record<
      TipeAkun,
      "primary" | "error" | "success" | "warning" | "info"
    > = {
      aset: "primary",
      kewajiban: "error",
      modal: "success",
      pendapatan: "info",
      beban: "warning",
    };
    return colors[tipe];
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
          Bagan Akun (Chart of Accounts)
        </Typography>
        <Box sx={{ display: "flex", gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<SeedIcon />}
            onClick={handleSeedCOA}
            disabled={accounts.length > 0}
          >
            Buat COA Default
          </Button>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={handleCreate}
          >
            Tambah Akun
          </Button>
        </Box>
      </Box>

      {/* Filters */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap" }}>
          {/* Tipe Akun Filter */}
          <FormControl size="small" sx={{ minWidth: 200 }}>
            <InputLabel>Tipe Akun</InputLabel>
            <Select
              value={tipeFilter}
              label="Tipe Akun"
              onChange={(e) =>
                setTipeFilter(e.target.value as TipeAkun | "all")
              }
            >
              <MenuItem value="all">Semua Tipe</MenuItem>
              <MenuItem value="aset">Aset</MenuItem>
              <MenuItem value="kewajiban">Kewajiban</MenuItem>
              <MenuItem value="modal">Modal</MenuItem>
              <MenuItem value="pendapatan">Pendapatan</MenuItem>
              <MenuItem value="beban">Beban</MenuItem>
            </Select>
          </FormControl>

          {/* Status Filter */}
          <FormControl size="small" sx={{ minWidth: 150 }}>
            <InputLabel>Status</InputLabel>
            <Select
              value={
                statusFilter === undefined ? "all" : statusFilter.toString()
              }
              label="Status"
              onChange={(e) =>
                setStatusFilter(
                  e.target.value === "all"
                    ? undefined
                    : e.target.value === "true"
                )
              }
            >
              <MenuItem value="all">Semua</MenuItem>
              <MenuItem value="true">Aktif</MenuItem>
              <MenuItem value="false">Non-aktif</MenuItem>
            </Select>
          </FormControl>
        </Box>
      </Paper>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Accounts Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Kode Akun</TableCell>
                <TableCell>Nama Akun</TableCell>
                <TableCell>Tipe</TableCell>
                <TableCell>Normal Saldo</TableCell>
                <TableCell>Akun Induk</TableCell>
                <TableCell>Status</TableCell>
                <TableCell align="right">Aksi</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={7} align="center" sx={{ py: 4 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : accounts.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center" sx={{ py: 4 }}>
                    <Typography color="text.secondary">
                      Tidak ada data akun. Klik &quot;Buat COA Default&quot;
                      untuk memulai.
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                accounts.map((account) => (
                  <TableRow key={account.id} hover>
                    <TableCell>
                      <Typography fontFamily="monospace" fontWeight={600}>
                        {account.kodeAkun}
                      </Typography>
                    </TableCell>
                    <TableCell>{account.namaAkun}</TableCell>
                    <TableCell>
                      <Chip
                        label={getTipeAkunLabel(account.tipeAkun)}
                        color={getTipeAkunColor(account.tipeAkun)}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={account.normalSaldo.toUpperCase()}
                        variant="outlined"
                        size="small"
                      />
                    </TableCell>
                    <TableCell>{account.namaInduk || "-"}</TableCell>
                    <TableCell>
                      <Chip
                        label={account.statusAktif ? "Aktif" : "Non-aktif"}
                        color={account.statusAktif ? "success" : "default"}
                        size="small"
                      />
                    </TableCell>
                    <TableCell align="right">
                      <IconButton
                        size="small"
                        onClick={() => handleViewLedger(account.id)}
                        title="Lihat Ledger"
                      >
                        <AccountBalanceIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() => handleEdit(account)}
                        title="Edit"
                      >
                        <EditIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() =>
                          handleDelete(account.id, account.namaAkun)
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
      </Paper>

      {/* Account Form Dialog */}
      <AccountForm
        open={formOpen}
        onClose={handleCloseForm}
        onSuccess={handleFormSuccess}
        account={currentAccount}
        parentAccounts={accounts}
      />
    </Box>
  );
}
