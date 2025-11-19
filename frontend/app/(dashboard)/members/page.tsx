// ============================================================================
// Member List Page - View and manage members
// Material-UI table with search, filters, and pagination
// ============================================================================

"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import {
  Box,
  Typography,
  Button,
  TextField,
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
  TablePagination,
  IconButton,
  Chip,
  Alert,
  CircularProgress,
  InputAdornment,
} from "@mui/material";
import {
  Add as AddIcon,
  Visibility as VisibilityIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
} from "@mui/icons-material";
import memberApi from "@/lib/api/memberApi";
import type { Member, MemberStatus } from "@/types";
import { format } from "date-fns";

// ============================================================================
// Member List Page Component
// ============================================================================

export default function MembersPage() {
  const router = useRouter();
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>("");

  // Pagination & Filters
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(20);
  const [totalItems, setTotalItems] = useState(0);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<MemberStatus | "all">("all");
  const [refreshKey, setRefreshKey] = useState(0); // For manual refresh triggers

  // ============================================================================
  // Fetch Members with Race Condition Protection
  // ============================================================================

  useEffect(() => {
    let ignore = false; // Cleanup flag to prevent race conditions

    const fetchMembers = async () => {
      try {
        setLoading(true);
        setError("");

        const response = await memberApi.getMembers({
          page: page + 1, // API uses 1-based pagination
          pageSize: rowsPerPage,
          search: searchQuery || undefined,
          status: statusFilter,
        });

        // Only update state if this effect is still current
        if (!ignore) {
          setMembers(response.data);
          setTotalItems(response.pagination.totalItems);
        }
      } catch (err: unknown) {
        if (!ignore) {
          console.error("Failed to fetch members:", err);
          setError("Gagal memuat data anggota. Silakan coba lagi.");
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    fetchMembers();

    // Cleanup function: mark results as stale if dependencies change
    return () => {
      ignore = true;
    };
  }, [page, rowsPerPage, statusFilter, searchQuery, refreshKey]);

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleSearch = () => {
    setPage(0); // Reset to first page - will trigger useEffect to refetch
  };

  const handleSearchKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSearch();
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

  const handleView = (id: string) => {
    router.push(`/dashboard/members/${id}`);
  };

  const handleEdit = (id: string) => {
    router.push(`/dashboard/members/${id}?mode=edit`);
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Apakah Anda yakin ingin menghapus anggota ini?")) {
      return;
    }

    try {
      await memberApi.deleteMember(id);
      // Trigger re-fetch by incrementing refreshKey
      setRefreshKey((prev) => prev + 1);
    } catch (err) {
      console.error("Failed to delete member:", err);
      alert("Gagal menghapus anggota. Silakan coba lagi.");
    }
  };

  // ============================================================================
  // Helper Functions
  // ============================================================================

  const getStatusColor = (
    status: MemberStatus
  ): "success" | "default" | "error" => {
    switch (status) {
      case "aktif":
        return "success";
      case "nonaktif":
        return "default";
      case "diberhentikan":
        return "error";
      default:
        return "default";
    }
  };

  const getStatusLabel = (status: MemberStatus): string => {
    switch (status) {
      case "aktif":
        return "Aktif";
      case "nonaktif":
        return "Non-aktif";
      case "diberhentikan":
        return "Diberhentikan";
      default:
        return status;
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
          Manajemen Anggota
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => router.push("/dashboard/members/new")}
        >
          Tambah Anggota
        </Button>
      </Box>

      {/* Filters */}
      <Paper sx={{ p: 2, mb: 3 }}>
        <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap" }}>
          {/* Search */}
          <TextField
            label="Cari Anggota"
            variant="outlined"
            size="small"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            onKeyPress={handleSearchKeyPress}
            sx={{ flexGrow: 1, minWidth: 200 }}
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton onClick={handleSearch} edge="end">
                    <SearchIcon />
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />

          {/* Status Filter */}
          <FormControl size="small" sx={{ minWidth: 150 }}>
            <InputLabel>Status</InputLabel>
            <Select
              value={statusFilter}
              label="Status"
              onChange={(e) =>
                setStatusFilter(e.target.value as MemberStatus | "all")
              }
            >
              <MenuItem value="all">Semua</MenuItem>
              <MenuItem value="aktif">Aktif</MenuItem>
              <MenuItem value="nonaktif">Non-aktif</MenuItem>
              <MenuItem value="diberhentikan">Diberhentikan</MenuItem>
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

      {/* Members Table */}
      <Paper>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>No. Anggota</TableCell>
                <TableCell>Nama Lengkap</TableCell>
                <TableCell>No. Telepon</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Tgl Bergabung</TableCell>
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
              ) : members.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center" sx={{ py: 4 }}>
                    <Typography color="text.secondary">
                      Tidak ada data anggota
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                members.map((member) => (
                  <TableRow key={member.id} hover>
                    <TableCell>{member.nomorAnggota}</TableCell>
                    <TableCell>{member.namaLengkap}</TableCell>
                    <TableCell>{member.noTelepon || "-"}</TableCell>
                    <TableCell>{member.email || "-"}</TableCell>
                    <TableCell>
                      {member.tanggalBergabung
                        ? format(
                            new Date(member.tanggalBergabung),
                            "dd/MM/yyyy"
                          )
                        : "-"}
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={getStatusLabel(member.status)}
                        color={getStatusColor(member.status)}
                        size="small"
                      />
                    </TableCell>
                    <TableCell align="right">
                      <IconButton
                        size="small"
                        onClick={() => handleView(member.id)}
                        title="Lihat Detail"
                      >
                        <VisibilityIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() => handleEdit(member.id)}
                        title="Edit"
                      >
                        <EditIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() => handleDelete(member.id)}
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
    </Box>
  );
}
