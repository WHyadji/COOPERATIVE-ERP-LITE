"use client";

import { useRouter } from "next/navigation";
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  CardActionArea,
  Divider,
} from "@mui/material";
import {
  AccountBalance as AccountBalanceIcon,
  TrendingUp as TrendingUpIcon,
  AttachMoney as AttachMoneyIcon,
  Assignment as AssignmentIcon,
} from "@mui/icons-material";

export default function LaporanPage() {
  const router = useRouter();

  const reports = [
    {
      title: "Neraca",
      subtitle: "Laporan Posisi Keuangan",
      description:
        "Menampilkan aset, kewajiban, dan modal koperasi pada tanggal tertentu",
      icon: AccountBalanceIcon,
      path: "/laporan/neraca",
      color: "#1976d2",
    },
    {
      title: "Laba Rugi",
      subtitle: "Laporan Laba Rugi",
      description:
        "Menampilkan pendapatan, beban, dan laba/rugi bersih dalam periode tertentu",
      icon: TrendingUpIcon,
      path: "/laporan/laba-rugi",
      color: "#2e7d32",
    },
    {
      title: "Arus Kas",
      subtitle: "Laporan Arus Kas",
      description:
        "Menampilkan arus kas dari aktivitas operasional, investasi, dan pendanaan",
      icon: AttachMoneyIcon,
      path: "/laporan/arus-kas",
      color: "#ed6c02",
    },
    {
      title: "Neraca Saldo",
      subtitle: "Trial Balance",
      description: "Menampilkan seluruh akun dengan saldo debit dan kredit",
      icon: AssignmentIcon,
      path: "/laporan/neraca-saldo",
      color: "#9c27b0",
    },
  ];

  return (
    <Box>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" fontWeight={700} gutterBottom>
          Laporan Keuangan
        </Typography>
        <Typography color="text.secondary">
          Pilih jenis laporan yang ingin ditampilkan
        </Typography>
      </Box>

      <Grid container spacing={3}>
        {reports.map((report) => {
          const IconComponent = report.icon;
          return (
            <Grid item xs={12} sm={6} md={6} key={report.path}>
              <Card
                sx={{
                  height: "100%",
                  transition: "all 0.3s",
                  "&:hover": { transform: "translateY(-4px)", boxShadow: 4 },
                }}
              >
                <CardActionArea
                  onClick={() => router.push(report.path)}
                  sx={{ height: "100%", p: 1 }}
                >
                  <CardContent>
                    <Box
                      sx={{ display: "flex", alignItems: "flex-start", mb: 2 }}
                    >
                      <Box
                        sx={{
                          p: 1.5,
                          borderRadius: 2,
                          bgcolor: `${report.color}15`,
                          display: "flex",
                          alignItems: "center",
                          justifyContent: "center",
                          mr: 2,
                        }}
                      >
                        <IconComponent
                          sx={{ fontSize: 32, color: report.color }}
                        />
                      </Box>
                      <Box sx={{ flex: 1 }}>
                        <Typography variant="h6" fontWeight={600} gutterBottom>
                          {report.title}
                        </Typography>
                        <Typography
                          variant="body2"
                          color="text.secondary"
                          gutterBottom
                        >
                          {report.subtitle}
                        </Typography>
                      </Box>
                    </Box>
                    <Divider sx={{ mb: 2 }} />
                    <Typography variant="body2" color="text.secondary">
                      {report.description}
                    </Typography>
                  </CardContent>
                </CardActionArea>
              </Card>
            </Grid>
          );
        })}
      </Grid>

      <Box
        sx={{
          mt: 4,
          p: 3,
          bgcolor: "info.lighter",
          borderRadius: 2,
          border: "1px solid",
          borderColor: "info.light",
        }}
      >
        <Typography
          variant="subtitle2"
          color="info.dark"
          fontWeight={600}
          gutterBottom
        >
          Informasi
        </Typography>
        <Typography variant="body2" color="info.dark">
          Semua laporan dihasilkan secara real-time berdasarkan data transaksi
          terkini. Laporan keuangan mengikuti standar SAK ETAP untuk koperasi
          Indonesia.
        </Typography>
      </Box>
    </Box>
  );
}
