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
import { format, startOfMonth, endOfMonth } from "date-fns";
import * as reportsApi from "@/lib/api/reportsApi";
import type { LaporanLabaRugi } from "@/types";

export default function LabaRugiPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");
  const [report, setReport] = useState<LaporanLabaRugi | null>(null);
  const [tanggalMulai, setTanggalMulai] = useState(() =>
    format(startOfMonth(new Date()), "yyyy-MM-dd")
  );
  const [tanggalAkhir, setTanggalAkhir] = useState(() =>
    format(endOfMonth(new Date()), "yyyy-MM-dd")
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
      if (
        !tanggalMulai ||
        !tanggalAkhir ||
        new Date(tanggalMulai) > new Date(tanggalAkhir)
      ) {
        setError("Tanggal mulai dan akhir harus valid");
        setLoading(false);
        return;
      }
      try {
        setLoading(true);
        setError("");
        const data = await reportsApi.getIncomeStatement(
          tanggalMulai,
          tanggalAkhir
        );
        if (!ignore) setReport(data);
      } catch (err) {
        if (!ignore)
          setError(
            err instanceof Error
              ? err.message
              : "Gagal memuat laporan laba rugi"
          );
      } finally {
        if (!ignore) setLoading(false);
      }
    };
    fetchReport();
    return () => {
      ignore = true;
    };
  }, [tanggalMulai, tanggalAkhir]);

  const isProfit = report ? report.labaRugiBersih >= 0 : true;

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
              Laba Rugi
            </Typography>
            <Typography color="text.secondary">Laporan Laba Rugi</Typography>
          </Box>
        </Box>
        <Box sx={{ display: "flex", gap: 2 }}>
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
                    Total Pendapatan
                  </Typography>
                  <Typography
                    variant="h5"
                    fontWeight={600}
                    color="success.main"
                  >
                    {formatCurrency(report.totalPendapatan)}
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
                    Total Beban
                  </Typography>
                  <Typography variant="h5" fontWeight={600} color="error.main">
                    {formatCurrency(report.totalBeban)}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} md={4}>
              <Card
                sx={{
                  bgcolor: isProfit ? "success.main" : "error.main",
                  color: "white",
                }}
              >
                <CardContent>
                  <Typography color="inherit" gutterBottom variant="body2">
                    {isProfit ? "Laba Bersih" : "Rugi Bersih"}
                  </Typography>
                  <Typography variant="h5" fontWeight={700} color="inherit">
                    {formatCurrency(Math.abs(report.labaRugiBersih))}
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
                      Jumlah (Rp)
                    </TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  <TableRow>
                    <TableCell
                      colSpan={3}
                      sx={{ bgcolor: "success.lighter", fontWeight: 700 }}
                    >
                      PENDAPATAN
                    </TableCell>
                  </TableRow>
                  {report.pendapatan.map((item, i) => (
                    <TableRow key={`pendapatan-${i}`} hover>
                      <TableCell>{item.kodeAkun}</TableCell>
                      <TableCell>{item.namaAkun}</TableCell>
                      <TableCell align="right">
                        {formatCurrency(item.saldo)}
                      </TableCell>
                    </TableRow>
                  ))}
                  <TableRow sx={{ bgcolor: "success.light" }}>
                    <TableCell colSpan={2} sx={{ fontWeight: 700 }}>
                      TOTAL PENDAPATAN
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 700 }}>
                      {formatCurrency(report.totalPendapatan)}
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell colSpan={3} sx={{ py: 2 }}></TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell
                      colSpan={3}
                      sx={{ bgcolor: "error.lighter", fontWeight: 700 }}
                    >
                      BEBAN
                    </TableCell>
                  </TableRow>
                  {report.beban.map((item, i) => (
                    <TableRow key={`beban-${i}`} hover>
                      <TableCell>{item.kodeAkun}</TableCell>
                      <TableCell>{item.namaAkun}</TableCell>
                      <TableCell align="right">
                        {formatCurrency(item.saldo)}
                      </TableCell>
                    </TableRow>
                  ))}
                  <TableRow sx={{ bgcolor: "error.light" }}>
                    <TableCell colSpan={2} sx={{ fontWeight: 700 }}>
                      TOTAL BEBAN
                    </TableCell>
                    <TableCell align="right" sx={{ fontWeight: 700 }}>
                      {formatCurrency(report.totalBeban)}
                    </TableCell>
                  </TableRow>
                  <TableRow>
                    <TableCell colSpan={3} sx={{ py: 2 }}></TableCell>
                  </TableRow>
                  <TableRow
                    sx={{ bgcolor: isProfit ? "success.main" : "error.main" }}
                  >
                    <TableCell
                      colSpan={2}
                      sx={{
                        fontWeight: 700,
                        color: "white",
                        fontSize: "1.1rem",
                      }}
                    >
                      {isProfit ? "LABA BERSIH" : "RUGI BERSIH"}
                    </TableCell>
                    <TableCell
                      align="right"
                      sx={{
                        fontWeight: 700,
                        color: "white",
                        fontSize: "1.1rem",
                      }}
                    >
                      {formatCurrency(Math.abs(report.labaRugiBersih))}
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </TableContainer>
            <Box sx={{ p: 2 }}>
              <Chip
                label={isProfit ? "Profit ✓" : "Loss"}
                color={isProfit ? "success" : "error"}
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
