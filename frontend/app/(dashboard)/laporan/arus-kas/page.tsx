'use client';
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Box, Typography, Button, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, TextField, Grid, Card, CardContent, CircularProgress, Alert, Chip } from '@mui/material';
import { ArrowBack as ArrowBackIcon, Print as PrintIcon, Download as DownloadIcon } from '@mui/icons-material';
import { format, startOfMonth, endOfMonth } from 'date-fns';
import * as reportsApi from '@/lib/api/reportsApi';
import type { LaporanArusKas } from '@/types';

export default function ArusKasPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [report, setReport] = useState<LaporanArusKas | null>(null);
  const [tanggalMulai, setTanggalMulai] = useState(() => format(startOfMonth(new Date()), 'yyyy-MM-dd'));
  const [tanggalAkhir, setTanggalAkhir] = useState(() => format(endOfMonth(new Date()), 'yyyy-MM-dd'));
  const formatCurrency = (amount: number): string => new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
  
  useEffect(() => {
    let ignore = false;
    const fetchReport = async () => {
      if (!tanggalMulai || !tanggalAkhir || new Date(tanggalMulai) > new Date(tanggalAkhir)) {
        setError('Tanggal mulai dan akhir harus valid');
        setLoading(false);
        return;
      }
      try {
        setLoading(true);
        setError('');
        const data = await reportsApi.getCashFlow(tanggalMulai, tanggalAkhir);
        if (!ignore) setReport(data);
      } catch (err) {
        if (!ignore) setError(err instanceof Error ? err.message : 'Gagal memuat laporan arus kas');
      } finally {
        if (!ignore) setLoading(false);
      }
    };
    fetchReport();
    return () => { ignore = true; };
  }, [tanggalMulai, tanggalAkhir]);
  
  const isPositive = report ? report.kenaikanKasBersih >= 0 : true;
  
  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }} className="no-print">
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Button variant="outlined" startIcon={<ArrowBackIcon />} onClick={() => router.push('/laporan')}>Kembali</Button>
          <Box><Typography variant="h4" fontWeight={700}>Arus Kas</Typography><Typography color="text.secondary">Laporan Arus Kas</Typography></Box>
        </Box>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <TextField label="Tanggal Mulai" type="date" size="small" value={tanggalMulai} onChange={(e) => setTanggalMulai(e.target.value)} InputLabelProps={{ shrink: true }} sx={{ minWidth: 160 }} />
          <TextField label="Tanggal Akhir" type="date" size="small" value={tanggalAkhir} onChange={(e) => setTanggalAkhir(e.target.value)} InputLabelProps={{ shrink: true }} sx={{ minWidth: 160 }} />
          <Button variant="outlined" startIcon={<PrintIcon />} onClick={() => window.print()} disabled={!report}>Cetak</Button>
          <Button variant="outlined" startIcon={<DownloadIcon />} onClick={() => alert('Fitur export PDF akan segera hadir')} disabled={!report}>Export PDF</Button>
        </Box>
      </Box>
      {error && <Alert severity="error" sx={{ mb: 3 }} className="no-print">{error}</Alert>}
      {loading && <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}><CircularProgress /></Box>}
      {!loading && report && (
        <>
          <Grid container spacing={2} sx={{ mb: 3 }} className="no-print">
            <Grid item xs={12} md={3}><Card><CardContent><Typography color="text.secondary" gutterBottom variant="body2">Operasional</Typography><Typography variant="h5" fontWeight={600} color={report.totalOperasional >= 0 ? 'success.main' : 'error.main'}>{formatCurrency(report.totalOperasional)}</Typography></CardContent></Card></Grid>
            <Grid item xs={12} md={3}><Card><CardContent><Typography color="text.secondary" gutterBottom variant="body2">Investasi</Typography><Typography variant="h5" fontWeight={600} color={report.totalInvestasi >= 0 ? 'success.main' : 'error.main'}>{formatCurrency(report.totalInvestasi)}</Typography></CardContent></Card></Grid>
            <Grid item xs={12} md={3}><Card><CardContent><Typography color="text.secondary" gutterBottom variant="body2">Pendanaan</Typography><Typography variant="h5" fontWeight={600} color={report.totalPendanaan >= 0 ? 'success.main' : 'error.main'}>{formatCurrency(report.totalPendanaan)}</Typography></CardContent></Card></Grid>
            <Grid item xs={12} md={3}><Card sx={{ bgcolor: isPositive ? 'success.main' : 'error.main', color: 'white' }}><CardContent><Typography color="inherit" gutterBottom variant="body2">Net Cash Flow</Typography><Typography variant="h5" fontWeight={700} color="inherit">{formatCurrency(report.kenaikanKasBersih)}</Typography></CardContent></Card></Grid>
          </Grid>
          <Paper><TableContainer><Table>
            <TableHead><TableRow><TableCell width="120">Kode Akun</TableCell><TableCell>Keterangan</TableCell><TableCell align="right" width="200">Jumlah (Rp)</TableCell></TableRow></TableHead>
            <TableBody>
              <TableRow><TableCell colSpan={3} sx={{ bgcolor: 'primary.lighter', fontWeight: 700 }}>ARUS KAS DARI AKTIVITAS OPERASIONAL</TableCell></TableRow>
              {report.arusKasOperasional.map((item, i) => <TableRow key={`op-${i}`} hover><TableCell>{item.kodeAkun}</TableCell><TableCell>{item.namaAkun}</TableCell><TableCell align="right">{formatCurrency(item.saldo)}</TableCell></TableRow>)}
              <TableRow sx={{ bgcolor: 'primary.light' }}><TableCell colSpan={2} sx={{ fontWeight: 700 }}>Total Arus Kas dari Aktivitas Operasional</TableCell><TableCell align="right" sx={{ fontWeight: 700 }}>{formatCurrency(report.totalOperasional)}</TableCell></TableRow>
              <TableRow><TableCell colSpan={3} sx={{ py: 2 }}></TableCell></TableRow>
              <TableRow><TableCell colSpan={3} sx={{ bgcolor: 'warning.lighter', fontWeight: 700 }}>ARUS KAS DARI AKTIVITAS INVESTASI</TableCell></TableRow>
              {report.arusKasInvestasi.map((item, i) => <TableRow key={`inv-${i}`} hover><TableCell>{item.kodeAkun}</TableCell><TableCell>{item.namaAkun}</TableCell><TableCell align="right">{formatCurrency(item.saldo)}</TableCell></TableRow>)}
              <TableRow sx={{ bgcolor: 'warning.light' }}><TableCell colSpan={2} sx={{ fontWeight: 700 }}>Total Arus Kas dari Aktivitas Investasi</TableCell><TableCell align="right" sx={{ fontWeight: 700 }}>{formatCurrency(report.totalInvestasi)}</TableCell></TableRow>
              <TableRow><TableCell colSpan={3} sx={{ py: 2 }}></TableCell></TableRow>
              <TableRow><TableCell colSpan={3} sx={{ bgcolor: 'secondary.lighter', fontWeight: 700 }}>ARUS KAS DARI AKTIVITAS PENDANAAN</TableCell></TableRow>
              {report.arusKasPendanaan.map((item, i) => <TableRow key={`pend-${i}`} hover><TableCell>{item.kodeAkun}</TableCell><TableCell>{item.namaAkun}</TableCell><TableCell align="right">{formatCurrency(item.saldo)}</TableCell></TableRow>)}
              <TableRow sx={{ bgcolor: 'secondary.light' }}><TableCell colSpan={2} sx={{ fontWeight: 700 }}>Total Arus Kas dari Aktivitas Pendanaan</TableCell><TableCell align="right" sx={{ fontWeight: 700 }}>{formatCurrency(report.totalPendanaan)}</TableCell></TableRow>
              <TableRow><TableCell colSpan={3} sx={{ py: 2 }}></TableCell></TableRow>
              <TableRow><TableCell colSpan={3} sx={{ bgcolor: 'grey.100', fontWeight: 700 }}>REKONSILIASI KAS</TableCell></TableRow>
              <TableRow hover><TableCell colSpan={2} sx={{ pl: 4 }}>Kenaikan/(Penurunan) Kas Bersih</TableCell><TableCell align="right" sx={{ fontWeight: 600 }}>{formatCurrency(report.kenaikanKasBersih)}</TableCell></TableRow>
              <TableRow hover><TableCell colSpan={2} sx={{ pl: 4 }}>Saldo Kas Awal Periode</TableCell><TableCell align="right">{formatCurrency(report.saldoKasAwal)}</TableCell></TableRow>
              <TableRow sx={{ bgcolor: isPositive ? 'success.light' : 'warning.light' }}><TableCell colSpan={2} sx={{ fontWeight: 700, pl: 4 }}>Saldo Kas Akhir Periode</TableCell><TableCell align="right" sx={{ fontWeight: 700 }}>{formatCurrency(report.saldoKasAkhir)}</TableCell></TableRow>
            </TableBody>
          </Table></TableContainer>
          <Box sx={{ p: 2 }}><Chip label={isPositive ? 'Positive Cash Flow ✓' : 'Negative Cash Flow'} color={isPositive ? 'success' : 'warning'} size="small" /></Box>
          </Paper>
        </>
      )}
      <style jsx global>{`@media print { .no-print { display: none !important; } }`}</style>
    </Box>
  );
}
