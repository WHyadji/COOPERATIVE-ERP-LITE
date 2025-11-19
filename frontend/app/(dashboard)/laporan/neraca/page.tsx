"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
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
  TextField,
  Grid,
  Card,
  CardContent,
  CircularProgress,
  Alert,
  Chip,
} from "@mui/material";
import {
  ArrowBack as ArrowBackIcon,
  Print as PrintIcon,
  Download as DownloadIcon,
} from "@mui/icons-material";
import { format } from "date-fns";
import * as reportsApi from "@/lib/api/reportsApi";
import type { LaporanPosisiKeuangan } from "@/types";

export default function NeracaPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");
  const [report, setReport] = useState<LaporanPosisiKeuangan | null>(null);
  const [tanggalPer, setTanggalPer] = useState(() =>
    format(new Date(), "yyyy-MM-dd")
  );

  const formatCurrency = (amount: number): string =>
    new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(amount);

  useEffect(() => {
    let ignore = false;
    const fetchReport = async () => {
      try {
        setLoading(true);
        setError("");
        const data = await reportsApi.getBalanceSheet(tanggalPer);
        if (!ignore) setReport(data);
      } catch (err) {
        if (!ignore)
          setError(
            err instanceof Error ? err.message : "Gagal memuat laporan neraca"
          );
      } finally {
        if (!ignore) setLoading(false);
      }
    };
    fetchReport();
    return () => {
      ignore = true;
    };
  }, [tanggalPer]);

  const isBalanced = report
    ? Math.abs(report.totalAset - (report.totalKewajiban + report.totalModal)) <
      0.01
    : false;

  return (
    <Box>
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
        className="no-print"
      >
        <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<ArrowBackIcon />}
            onClick={() => router.push("/laporan")}
          >
            Kembali
          </Button>
          <Box>
            <Typography variant="h4" fontWeight={700}>
              Neraca
            </Typography>
            <Typography color="text.secondary">
              Laporan Posisi Keuangan
            </Typography>
          </Box>
        </Box>
        <Box sx={{ display: "flex", gap: 2 }}>
          <TextField
            label="Tanggal Per"
            type="date"
            size="small"
            value={tanggalPer}
            onChange={(e) => setTanggalPer(e.target.value)}
            InputLabelProps={{ shrink: true }}
            sx={{ minWidth: 160 }}
          />
          <Button
            variant="outlined"
            startIcon={<PrintIcon />}
            onClick={() => window.print()}
            disabled={!report}
          >
            Cetak
          </Button>
          <Button
            variant="outlined"
            startIcon={<DownloadIcon />}
            onClick={() => alert("Fitur export PDF akan segera hadir")}
            disabled={!report}
          >
            Export PDF
          </Button>
        </Box>
      </Box>
      {error && (
        <Alert severity="error" sx={{ mb: 3 }} className="no-print">
          {error}
        </Alert>
      )}
      {loading && (
        <Box sx={{ display: "flex", justifyContent: "center", py: 8 }}>
          <CircularProgress />
        </Box>
      )}
      {!loading && report && (
        <>
          <Grid container spacing={2} sx={{ mb: 3 }} className="no-print">
            <Grid item xs={12} md={4}>
              <Card>
                <CardContent>
                  <Typography
                    color="text.secondary"
                    gutterBottom
                    variant="body2"
                  >
                    Total Aset
                  </Typography>
                  <Typography
                    variant="h5"
                    fontWeight={600}
                    color="primary.main"
                  >
                    {formatCurrency(report.totalAset)}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={4}>
              <Card>
                <CardContent>
                  <Typography
                    color="text.secondary"
                    gutterBottom
                    variant="body2"
                  >
                    Total Kewajiban
                  </Typography>
                  <Typography
                    variant="h5"
                    fontWeight={600}
                    color="warning.main"
                  >
                    {formatCurrency(report.totalKewajiban)}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={4}>
              <Card>
                <CardContent>
                  <Typography
                    color="text.secondary"
                    gutterBottom
                    variant="body2"
                  >
                    Total Modal
                  </Typography>
                  <Typography
                    variant="h5"
                    fontWeight={600}
                    color="secondary.main"
                  >
                    {formatCurrency(report.totalModal)}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
          <Paper>
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell width="120">Kode Akun</TableCell>
                    <TableCell>Nama Akun</TableCell>
                    <TableCell align="right" width="200">
                      Saldo (Rp)
                    </TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  <TableRow>
                    <TableCell
                      colSpan={3}
                      sx={{ bgcolor: "primary.lighter", fontWeight: 700 }}
                    >
                      ASET
                    </TableCell>
                  </TableRow>
                  {report.aset.map((item, i) => (
                    <TableRow key={`aset-${i}`} hover>
                      <TableCell>{item.kodeAkun}</TableCell>
                      <TableCell>{item.namaAkun}</TableCell>
                      <TableCell align="right">
                        {formatCurrency(item.saldo)}
                      </TableCell>
                    </TableRow>
                  ))}
                  <TableRow sx={{ bgcolor: "primary.light" }}>
                    <TableCell colSpan={2} sx={{ fontWeight: 700 }}>
                      TOTAL ASET
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 700 }}>
                      {formatCurrency(report.totalAset)}
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell colSpan={3} sx={{ py: 2 }}></TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell
                      colSpan={3}
                      sx={{ bgcolor: "warning.lighter", fontWeight: 700 }}
                    >
                      KEWAJIBAN
                    </TableCell>
                  </TableRow>
                  {report.kewajiban.map((item, i) => (
                    <TableRow key={`kewajiban-${i}`} hover>
                      <TableCell>{item.kodeAkun}</TableCell>
                      <TableCell>{item.namaAkun}</TableCell>
                      <TableCell align="right">
                        {formatCurrency(item.saldo)}
                      </TableCell>
                    </TableRow>
                  ))}
                  <TableRow sx={{ bgcolor: "warning.light" }}>
                    <TableCell colSpan={2} sx={{ fontWeight: 700 }}>
                      TOTAL KEWAJIBAN
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 700 }}>
                      {formatCurrency(report.totalKewajiban)}
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell colSpan={3} sx={{ py: 2 }}></TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell
                      colSpan={3}
                      sx={{ bgcolor: "secondary.lighter", fontWeight: 700 }}
                    >
                      MODAL
                    </TableCell>
                  </TableRow>
                  {report.modal.map((item, i) => (
                    <TableRow key={`modal-${i}`} hover>
                      <TableCell>{item.kodeAkun}</TableCell>
                      <TableCell>{item.namaAkun}</TableCell>
                      <TableCell align="right">
                        {formatCurrency(item.saldo)}
                      </TableCell>
                    </TableRow>
                  ))}
                  <TableRow sx={{ bgcolor: "secondary.light" }}>
                    <TableCell colSpan={2} sx={{ fontWeight: 700 }}>
                      TOTAL MODAL
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 700 }}>
                      {formatCurrency(report.totalModal)}
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </TableContainer>
            <Box sx={{ p: 2 }}>
              <Chip
                label={isBalanced ? "Balanced ✓" : "Unbalanced ✗"}
                color={isBalanced ? "success" : "error"}
                size="small"
              />
            </Box>
          </Paper>
        </>
      )}
      <style jsx global>{`
        @media print {
          .no-print {
            display: none !important;
          }
        }
      `}</style>
    </Box>
  );
}
