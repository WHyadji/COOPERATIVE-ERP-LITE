'use client';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Box, Typography, Button, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, TextField, Grid, Card, CardContent, CircularProgress, Alert, Chip } from '@mui/material';
import { ArrowBack as ArrowBackIcon, Print as PrintIcon, Download as DownloadIcon } from '@mui/icons-material';
import { format } from 'date-fns';
import * as reportsApi from '@/lib/api/reportsApi';
import type { NeracaSaldo } from '@/types';

export default function NeracaSaldoPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [report, setReport] = useState<NeracaSaldo | null>(null);
  const [tanggalPer, setTanggalPer] = useState(() => format(new Date(), 'yyyy-MM-dd'));
  const formatCurrency = (amount: number): string => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
  
  useEffect(() => {
    let ignore = false;
    const fetchReport = async () => {
      try {
        setLoading(true);
        setError('');
        const data = await reportsApi.getTrialBalance(tanggalPer);
        if (!ignore) setReport(data);
      } catch (err) {
        if (!ignore) setError(err instanceof Error ? err.message : 'Gagal memuat neraca saldo');
      } finally {
        if (!ignore) setLoading(false);
      }
    };
    fetchReport();
    return () => { ignore = true; };
  }, [tanggalPer]);
  
  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }} className="no-print">
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Button variant="outlined" startIcon={<ArrowBackIcon />} onClick={() => router.push('/laporan')}>Kembali</Button>
          <Box><Typography variant="h4" fontWeight={700}>Neraca Saldo</Typography><Typography color="text.secondary">Trial Balance</Typography></Box>
        </Box>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <TextField label="Tanggal Per" type="date" size="small" value={tanggalPer} onChange={(e) => setTanggalPer(e.target.value)} InputLabelProps={{ shrink: true }} sx={{ minWidth: 160 }} />
          <Button variant="outlined" startIcon={<PrintIcon />} onClick={() => window.print()} disabled={!report}>Cetak</Button>
          <Button variant="outlined" startIcon={<DownloadIcon />} onClick={() => alert('Fitur export PDF akan segera hadir')} disabled={!report}>Export PDF</Button>
        </Box>
      </Box>
      {error && <Alert severity="error" sx={{ mb: 3 }} className="no-print">{error}</Alert>}
      {loading && <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}><CircularProgress /></Box>}
      {!loading && report && (
        <>
          <Grid container spacing={2} sx={{ mb: 3 }} className="no-print">
            <Grid item xs={12} md={4}><Card><CardContent><Typography color="text.secondary" gutterBottom variant="body2">Total Debit</Typography><Typography variant="h5" fontWeight={600} color="primary.main">{formatCurrency(report.totalDebit)}</Typography></CardContent></Card></Grid>
            <Grid item xs={12} md={4}><Card><CardContent><Typography color="text.secondary" gutterBottom variant="body2">Total Kredit</Typography><Typography variant="h5" fontWeight={600} color="warning.main">{formatCurrency(report.totalKredit)}</Typography></CardContent></Card></Grid>
            <Grid item xs={12} md={4}><Card sx={{ bgcolor: report.isBalanced ? 'success.main' : 'error.main', color: 'white' }}><CardContent><Typography color="inherit" gutterBottom variant="body2">Status</Typography><Typography variant="h5" fontWeight={700} color="inherit">{report.isBalanced ? 'Balanced' : 'Unbalanced'}</Typography></CardContent></Card></Grid>
          </Grid>
          <Paper><TableContainer><Table>
            <TableHead><TableRow><TableCell width="100">Kode Akun</TableCell><TableCell>Nama Akun</TableCell><TableCell width="120">Tipe</TableCell><TableCell align="right" width="180">Debit (Rp)</TableCell><TableCell align="right" width="180">Kredit (Rp)</TableCell></TableRow></TableHead>
            <TableBody>
              {report.items.map((item, i) => (
                <TableRow key={`item-${i}`} hover>
                  <TableCell>{item.kodeAkun}</TableCell>
                  <TableCell>{item.namaAkun}</TableCell>
                  <TableCell><Chip label={item.tipeAkun} size="small" color={item.tipeAkun === 'aset' ? 'primary' : item.tipeAkun === 'kewajiban' ? 'warning' : item.tipeAkun === 'modal' ? 'secondary' : item.tipeAkun === 'pendapatan' ? 'success' : 'error'} /></TableCell>
                  <TableCell align="right">{item.saldoDebit > 0 ? formatCurrency(item.saldoDebit) : '-'}</TableCell>
                  <TableCell align="right">{item.saldoKredit > 0 ? formatCurrency(item.saldoKredit) : '-'}</TableCell>
                </TableRow>
              ))}
              <TableRow sx={{ bgcolor: report.isBalanced ? 'success.lighter' : 'error.lighter' }}>
                <TableCell colSpan={3} sx={{ fontWeight: 700, fontSize: '1.1rem' }}>TOTAL</TableCell>
                <TableCell align="right" sx={{ fontWeight: 700, fontSize: '1.1rem' }}>{formatCurrency(report.totalDebit)}</TableCell>
                <TableCell align="right" sx={{ fontWeight: 700, fontSize: '1.1rem' }}>{formatCurrency(report.totalKredit)}</TableCell>
              </TableRow>
            </TableBody>
          </Table></TableContainer>
          <Box sx={{ p: 2 }}><Chip label={report.isBalanced ? 'Balanced ✓' : 'Unbalanced ✗'} color={report.isBalanced ? 'success' : 'error'} size="small" /></Box>
          </Paper>
        </>
      )}
      <style jsx global>{`@media print { .no-print { display: none !important; } }`}</style>
    </Box>
  );
}
