// ============================================================================
// Account Ledger Page - View transaction history for specific account
// Shows running balance with debit/credit entries
// ============================================================================

'use client';

import React, { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
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
  Alert,
  CircularProgress,
  Breadcrumbs,
  Link,
  Chip,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  Print as PrintIcon,
  Download as DownloadIcon,
} from '@mui/icons-material';
import accountingApi from '@/lib/api/accountingApi';
import type { Akun } from '@/types';
import { format, parseISO } from 'date-fns';

// ============================================================================
// Ledger Entry Interface
// ============================================================================

interface LedgerEntry {
  tanggal: string;
  nomorJurnal: string;
  deskripsi: string;
  debit: number;
  kredit: number;
  saldo: number;
}

// ============================================================================
// Account Ledger Page Component
// ============================================================================

export default function AccountLedgerPage() {
  const params = useParams();
  const router = useRouter();
  const accountId = params.id as string;

  const [account, setAccount] = useState<Akun | null>(null);
  const [ledgerEntries, setLedgerEntries] = useState<LedgerEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  // Filters
  const [tanggalMulai, setTanggalMulai] = useState('');
  const [tanggalAkhir, setTanggalAkhir] = useState('');
  const [refreshKey, setRefreshKey] = useState(0);

  // ============================================================================
  // Fetch Account Details
  // ============================================================================

  useEffect(() => {
    const fetchAccount = async () => {
      try {
        const data = await accountingApi.getAccountById(accountId);
        setAccount(data);
      } catch (err) {
        console.error('Failed to fetch account:', err);
      }
    };

    if (accountId) {
      fetchAccount();
    }
  }, [accountId]);

  // ============================================================================
  // Fetch Ledger Entries
  // ============================================================================

  useEffect(() => {
    let ignore = false;

    const fetchLedger = async () => {
      try {
        setLoading(true);
        setError('');

        const data = await accountingApi.getAccountLedger(
          accountId,
          tanggalMulai || undefined,
          tanggalAkhir || undefined
        );

        if (!ignore) {
          setLedgerEntries(data);
        }
      } catch (err: unknown) {
        if (!ignore) {
          console.error('Failed to fetch ledger:', err);
          setError('Gagal memuat data buku besar. Silakan coba lagi.');
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    if (accountId) {
      fetchLedger();
    }

    return () => {
      ignore = true;
    };
  }, [accountId, tanggalMulai, tanggalAkhir, refreshKey]);

  // ============================================================================
  // Helper Functions
  // ============================================================================

  const formatCurrency = (amount: number): string => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const formatDate = (dateString: string): string => {
    try {
      return format(parseISO(dateString), 'dd/MM/yyyy');
    } catch {
      return dateString;
    }
  };

  const getTotalDebit = (): number => {
    return ledgerEntries.reduce((sum, entry) => sum + entry.debit, 0);
  };

  const getTotalKredit = (): number => {
    return ledgerEntries.reduce((sum, entry) => sum + entry.kredit, 0);
  };

  const getFinalBalance = (): number => {
    if (ledgerEntries.length === 0) return 0;
    return ledgerEntries[ledgerEntries.length - 1].saldo;
  };

  // ============================================================================
  // Render
  // ============================================================================

  if (!account && !loading) {
    return (
      <Box>
        <Alert severity="error">Akun tidak ditemukan</Alert>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={() => router.push('/akuntansi')}
          sx={{ mt: 2 }}
        >
          Kembali ke Bagan Akun
        </Button>
      </Box>
    );
  }

  return (
    <Box>
      {/* Breadcrumbs */}
      <Breadcrumbs sx={{ mb: 3 }}>
        <Link
          underline="hover"
          color="inherit"
          onClick={() => router.push('/akuntansi')}
          sx={{ cursor: 'pointer' }}
        >
          Bagan Akun
        </Link>
        <Typography color="text.primary">Buku Besar</Typography>
      </Breadcrumbs>

      {/* Header */}
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'flex-start',
          mb: 3,
        }}
      >
        <Box>
          <Typography variant="h4" fontWeight={600} gutterBottom>
            Buku Besar (General Ledger)
          </Typography>
          {account && (
            <Box sx={{ display: 'flex', gap: 2, mt: 1 }}>
              <Chip
                label={`${account.kodeAkun} - ${account.namaAkun}`}
                color="primary"
              />
              <Chip
                label={`Normal: ${account.normalSaldo.toUpperCase()}`}
                variant="outlined"
              />
            </Box>
          )}
        </Box>
        <Button
          variant="outlined"
          startIcon={<ArrowBackIcon />}
          onClick={() => router.push('/akuntansi')}
        >
          Kembali
        </Button>
      </Box>

      {/* Filters */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap', alignItems: 'center' }}>
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
              setTanggalMulai('');
              setTanggalAkhir('');
            }}
          >
            Reset Filter
          </Button>
          <Box sx={{ flexGrow: 1 }} />
          <Button startIcon={<PrintIcon />} onClick={() => window.print()}>
            Cetak
          </Button>
        </Box>
      </Paper>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Ledger Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Tanggal</TableCell>
                <TableCell>No. Jurnal</TableCell>
                <TableCell>Deskripsi</TableCell>
                <TableCell align="right">Debit</TableCell>
                <TableCell align="right">Kredit</TableCell>
                <TableCell align="right">Saldo</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : ledgerEntries.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                    <Typography color="text.secondary">
                      Tidak ada transaksi untuk akun ini
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                <>
                  {ledgerEntries.map((entry, index) => (
                    <TableRow key={index} hover>
                      <TableCell>{formatDate(entry.tanggal)}</TableCell>
                      <TableCell>
                        <Typography fontFamily="monospace">
                          {entry.nomorJurnal}
                        </Typography>
                      </TableCell>
                      <TableCell>{entry.deskripsi}</TableCell>
                      <TableCell align="right">
                        {entry.debit > 0 ? formatCurrency(entry.debit) : '-'}
                      </TableCell>
                      <TableCell align="right">
                        {entry.kredit > 0 ? formatCurrency(entry.kredit) : '-'}
                      </TableCell>
                      <TableCell align="right">
                        <Typography
                          fontWeight={600}
                          color={entry.saldo >= 0 ? 'primary' : 'error'}
                        >
                          {formatCurrency(Math.abs(entry.saldo))}
                          {entry.saldo < 0 && ' (CR)'}
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ))}

                  {/* Totals Row */}
                  <TableRow sx={{ bgcolor: 'grey.100' }}>
                    <TableCell colSpan={3}>
                      <Typography fontWeight={700}>TOTAL</Typography>
                    </TableCell>
                    <TableCell align="right">
                      <Typography fontWeight={700} color="primary">
                        {formatCurrency(getTotalDebit())}
                      </Typography>
                    </TableCell>
                    <TableCell align="right">
                      <Typography fontWeight={700} color="error">
                        {formatCurrency(getTotalKredit())}
                      </Typography>
                    </TableCell>
                    <TableCell align="right">
                      <Typography fontWeight={700}>
                        {formatCurrency(Math.abs(getFinalBalance()))}
                        {getFinalBalance() < 0 && ' (CR)'}
                      </Typography>
                    </TableCell>
                  </TableRow>
                </>
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      {/* Summary */}
      {!loading && ledgerEntries.length > 0 && account && (
        <Paper sx={{ p: 3, mt: 3 }}>
          <Typography variant="h6" gutterBottom>
            Ringkasan
          </Typography>
          <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 3 }}>
            <Box>
              <Typography variant="body2" color="text.secondary">
                Total Debit
              </Typography>
              <Typography variant="h6" color="primary">
                {formatCurrency(getTotalDebit())}
              </Typography>
            </Box>
            <Box>
              <Typography variant="body2" color="text.secondary">
                Total Kredit
              </Typography>
              <Typography variant="h6" color="error">
                {formatCurrency(getTotalKredit())}
              </Typography>
            </Box>
            <Box>
              <Typography variant="body2" color="text.secondary">
                Saldo Akhir
              </Typography>
              <Typography variant="h6">
                {formatCurrency(Math.abs(getFinalBalance()))}
                {getFinalBalance() < 0 && ' (CR)'}
              </Typography>
            </Box>
          </Box>
        </Paper>
      )}
    </Box>
  );
}
