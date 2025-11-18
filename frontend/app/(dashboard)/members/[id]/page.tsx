// ============================================================================
// Member Detail Page - View and edit member information
// Displays member data with edit mode toggle
// ============================================================================

'use client';

import React, { useState, useEffect } from 'react';
import { useParams, useRouter, useSearchParams } from 'next/navigation';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
  Box,
  Typography,
  Paper,
  TextField,
  Button,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormHelperText,
  Alert,
  Divider,
  Chip,
  CircularProgress,
} from '@mui/material';
import {
  Edit as EditIcon,
  Save as SaveIcon,
  Cancel as CancelIcon,
  ArrowBack as ArrowBackIcon,
} from '@mui/icons-material';
import memberApi from '@/lib/api/memberApi';
import type { Member, UpdateMemberRequest, Gender, MemberStatus } from '@/types';
import { format } from 'date-fns';

// ============================================================================
// Validation Schema (same as create, but with status)
// ============================================================================

const memberSchema = z.object({
  namaLengkap: z.string().min(1, 'Nama lengkap harus diisi'),
  nik: z.string().optional(),
  tanggalLahir: z.string().optional(),
  tempatLahir: z.string().optional(),
  jenisKelamin: z.enum(['L', 'P', '']).optional(),
  alamat: z.string().optional(),
  rt: z.string().optional(),
  rw: z.string().optional(),
  kelurahan: z.string().optional(),
  kecamatan: z.string().optional(),
  kotaKabupaten: z.string().optional(),
  provinsi: z.string().optional(),
  kodePos: z.string().optional(),
  noTelepon: z.string().optional(),
  email: z.string().email('Email tidak valid').optional().or(z.literal('')),
  pekerjaan: z.string().optional(),
  tanggalBergabung: z.string().min(1, 'Tanggal bergabung harus diisi'),
  status: z.enum(['aktif', 'nonaktif', 'diberhentikan']),
});

type MemberFormData = z.infer<typeof memberSchema>;

// ============================================================================
// Member Detail Page Component
// ============================================================================

export default function MemberDetailPage() {
  const params = useParams();
  const router = useRouter();
  const searchParams = useSearchParams();
  const memberId = params.id as string;

  const [member, setMember] = useState<Member | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');
  const [editMode, setEditMode] = useState(searchParams.get('mode') === 'edit');

  // ============================================================================
  // Form Setup
  // ============================================================================

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<MemberFormData>({
    resolver: zodResolver(memberSchema),
  });

  // ============================================================================
  // Fetch Member Data with Race Condition Protection
  // ============================================================================

  useEffect(() => {
    if (!memberId) return;

    let ignore = false; // Cleanup flag to prevent race conditions

    const fetchMember = async () => {
      try {
        setLoading(true);
        const data = await memberApi.getMemberById(memberId);

        // Only update state if this effect is still current
        if (!ignore) {
          setMember(data);

          // Populate form
          reset({
            namaLengkap: data.namaLengkap,
            nik: data.nik || '',
            tanggalLahir: data.tanggalLahir
              ? format(new Date(data.tanggalLahir), 'yyyy-MM-dd')
              : '',
            tempatLahir: data.tempatLahir || '',
            jenisKelamin: data.jenisKelamin || '',
            alamat: data.alamat || '',
            rt: data.rt || '',
            rw: data.rw || '',
            kelurahan: data.kelurahan || '',
            kecamatan: data.kecamatan || '',
            kotaKabupaten: data.kotaKabupaten || '',
            provinsi: data.provinsi || '',
            kodePos: data.kodePos || '',
            noTelepon: data.noTelepon || '',
            email: data.email || '',
            pekerjaan: data.pekerjaan || '',
            tanggalBergabung: data.tanggalBergabung
              ? format(new Date(data.tanggalBergabung), 'yyyy-MM-dd')
              : '',
            status: data.status,
          });
        }
      } catch (err) {
        if (!ignore) {
          console.error('Failed to fetch member:', err);
          setError('Gagal memuat data anggota.');
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    fetchMember();

    // Cleanup function: mark results as stale if dependencies change
    return () => {
      ignore = true;
    };
  }, [memberId, reset]);

  // ============================================================================
  // Submit Handler
  // ============================================================================

  const onSubmit = async (data: MemberFormData) => {
    try {
      setError('');
      setSuccess('');

      const requestData: UpdateMemberRequest = {
        namaLengkap: data.namaLengkap,
        nik: data.nik || undefined,
        tanggalLahir: data.tanggalLahir || undefined,
        tempatLahir: data.tempatLahir || undefined,
        jenisKelamin: (data.jenisKelamin === '' ? undefined : data.jenisKelamin) as Gender | undefined,
        alamat: data.alamat || undefined,
        rt: data.rt || undefined,
        rw: data.rw || undefined,
        kelurahan: data.kelurahan || undefined,
        kecamatan: data.kecamatan || undefined,
        kotaKabupaten: data.kotaKabupaten || undefined,
        provinsi: data.provinsi || undefined,
        kodePos: data.kodePos || undefined,
        noTelepon: data.noTelepon || undefined,
        email: data.email || undefined,
        pekerjaan: data.pekerjaan || undefined,
        tanggalBergabung: data.tanggalBergabung,
        status: data.status as MemberStatus,
      };

      const updatedMember = await memberApi.updateMember(memberId, requestData);
      setMember(updatedMember);
      setSuccess('Data anggota berhasil diperbarui!');
      setEditMode(false);
    } catch (err: unknown) {
      console.error('Failed to update member:', err);

      if (err && typeof err === 'object' && 'message' in err) {
        setError(err.message as string);
      } else {
        setError('Gagal memperbarui data anggota. Silakan coba lagi.');
      }
    }
  };

  // ============================================================================
  // Helper Functions
  // ============================================================================

  const getStatusColor = (status: MemberStatus): 'success' | 'default' | 'error' => {
    switch (status) {
      case 'aktif':
        return 'success';
      case 'nonaktif':
        return 'default';
      case 'diberhentikan':
        return 'error';
      default:
        return 'default';
    }
  };

  // ============================================================================
  // Render
  // ============================================================================

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (!member) {
    return (
      <Box>
        <Alert severity="error">Anggota tidak ditemukan.</Alert>
        <Button onClick={() => router.back()} sx={{ mt: 2 }}>
          Kembali
        </Button>
      </Box>
    );
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box>
          <Button
            startIcon={<ArrowBackIcon />}
            onClick={() => router.back()}
            sx={{ mb: 1 }}
          >
            Kembali
          </Button>
          <Typography variant="h4" fontWeight={600}>
            Detail Anggota
          </Typography>
          <Typography variant="body2" color="text.secondary">
            No. Anggota: {member.nomorAnggota}
          </Typography>
        </Box>
        <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
          <Chip
            label={member.status}
            color={getStatusColor(member.status)}
          />
          {!editMode && (
            <Button
              variant="contained"
              startIcon={<EditIcon />}
              onClick={() => setEditMode(true)}
            >
              Edit
            </Button>
          )}
        </Box>
      </Box>

      {/* Success Alert */}
      {success && (
        <Alert severity="success" sx={{ mb: 3 }}>
          {success}
        </Alert>
      )}

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Form */}
      <Paper sx={{ p: 3 }}>
        <form onSubmit={handleSubmit(onSubmit)}>
          {/* Section 1: Data Pribadi */}
          <Typography variant="h6" gutterBottom>
            Data Pribadi
          </Typography>
          <Divider sx={{ mb: 3 }} />

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                {...register('namaLengkap')}
                label="Nama Lengkap *"
                fullWidth
                error={!!errors.namaLengkap}
                helperText={errors.namaLengkap?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                {...register('nik')}
                label="NIK (KTP)"
                fullWidth
                error={!!errors.nik}
                helperText={errors.nik?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
                inputProps={{ maxLength: 16 }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('tempatLahir')}
                label="Tempat Lahir"
                fullWidth
                error={!!errors.tempatLahir}
                helperText={errors.tempatLahir?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('tanggalLahir')}
                label="Tanggal Lahir"
                type="date"
                fullWidth
                error={!!errors.tanggalLahir}
                helperText={errors.tanggalLahir?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <FormControl fullWidth error={!!errors.jenisKelamin} disabled={!editMode || isSubmitting}>
                <InputLabel>Jenis Kelamin</InputLabel>
                <Controller
                  name="jenisKelamin"
                  control={control}
                  render={({ field }) => (
                    <Select {...field} label="Jenis Kelamin" readOnly={!editMode}>
                      <MenuItem value="">Pilih</MenuItem>
                      <MenuItem value="L">Laki-laki</MenuItem>
                      <MenuItem value="P">Perempuan</MenuItem>
                    </Select>
                  )}
                />
                {errors.jenisKelamin && (
                  <FormHelperText>{errors.jenisKelamin.message}</FormHelperText>
                )}
              </FormControl>
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('pekerjaan')}
                label="Pekerjaan"
                fullWidth
                error={!!errors.pekerjaan}
                helperText={errors.pekerjaan?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('tanggalBergabung')}
                label="Tanggal Bergabung *"
                type="date"
                fullWidth
                error={!!errors.tanggalBergabung}
                helperText={errors.tanggalBergabung?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <FormControl fullWidth error={!!errors.status} disabled={!editMode || isSubmitting}>
                <InputLabel>Status</InputLabel>
                <Controller
                  name="status"
                  control={control}
                  render={({ field }) => (
                    <Select {...field} label="Status" readOnly={!editMode}>
                      <MenuItem value="aktif">Aktif</MenuItem>
                      <MenuItem value="nonaktif">Non-aktif</MenuItem>
                      <MenuItem value="diberhentikan">Diberhentikan</MenuItem>
                    </Select>
                  )}
                />
                {errors.status && (
                  <FormHelperText>{errors.status.message}</FormHelperText>
                )}
              </FormControl>
            </Grid>
          </Grid>

          {/* Section 2: Alamat */}
          <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
            Alamat
          </Typography>
          <Divider sx={{ mb: 3 }} />

          <Grid container spacing={3}>
            <Grid item xs={12}>
              <TextField
                {...register('alamat')}
                label="Alamat Lengkap"
                fullWidth
                multiline
                rows={2}
                error={!!errors.alamat}
                helperText={errors.alamat?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={3}>
              <TextField
                {...register('rt')}
                label="RT"
                fullWidth
                error={!!errors.rt}
                helperText={errors.rt?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={3}>
              <TextField
                {...register('rw')}
                label="RW"
                fullWidth
                error={!!errors.rw}
                helperText={errors.rw?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                {...register('kelurahan')}
                label="Kelurahan/Desa"
                fullWidth
                error={!!errors.kelurahan}
                helperText={errors.kelurahan?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('kecamatan')}
                label="Kecamatan"
                fullWidth
                error={!!errors.kecamatan}
                helperText={errors.kecamatan?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('kotaKabupaten')}
                label="Kota/Kabupaten"
                fullWidth
                error={!!errors.kotaKabupaten}
                helperText={errors.kotaKabupaten?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('provinsi')}
                label="Provinsi"
                fullWidth
                error={!!errors.provinsi}
                helperText={errors.provinsi?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('kodePos')}
                label="Kode Pos"
                fullWidth
                error={!!errors.kodePos}
                helperText={errors.kodePos?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
                inputProps={{ maxLength: 5 }}
              />
            </Grid>
          </Grid>

          {/* Section 3: Kontak */}
          <Typography variant="h6" gutterBottom sx={{ mt: 4 }}>
            Informasi Kontak
          </Typography>
          <Divider sx={{ mb: 3 }} />

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                {...register('noTelepon')}
                label="No. Telepon"
                fullWidth
                error={!!errors.noTelepon}
                helperText={errors.noTelepon?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                {...register('email')}
                label="Email"
                type="email"
                fullWidth
                error={!!errors.email}
                helperText={errors.email?.message}
                disabled={!editMode || isSubmitting}
                InputProps={{ readOnly: !editMode }}
              />
            </Grid>
          </Grid>

          {/* Action Buttons - Only show in edit mode */}
          {editMode && (
            <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 4 }}>
              <Button
                variant="outlined"
                startIcon={<CancelIcon />}
                onClick={() => {
                  setEditMode(false);
                  setError('');
                  setSuccess('');
                  // Reset form to original values
                  if (member) {
                    reset({
                      namaLengkap: member.namaLengkap,
                      nik: member.nik || '',
                      tanggalLahir: member.tanggalLahir
                        ? format(new Date(member.tanggalLahir), 'yyyy-MM-dd')
                        : '',
                      tempatLahir: member.tempatLahir || '',
                      jenisKelamin: member.jenisKelamin || '',
                      alamat: member.alamat || '',
                      rt: member.rt || '',
                      rw: member.rw || '',
                      kelurahan: member.kelurahan || '',
                      kecamatan: member.kecamatan || '',
                      kotaKabupaten: member.kotaKabupaten || '',
                      provinsi: member.provinsi || '',
                      kodePos: member.kodePos || '',
                      noTelepon: member.noTelepon || '',
                      email: member.email || '',
                      pekerjaan: member.pekerjaan || '',
                      tanggalBergabung: member.tanggalBergabung
                        ? format(new Date(member.tanggalBergabung), 'yyyy-MM-dd')
                        : '',
                      status: member.status,
                    });
                  }
                }}
                disabled={isSubmitting}
              >
                Batal
              </Button>
              <Button
                type="submit"
                variant="contained"
                startIcon={<SaveIcon />}
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Menyimpan...' : 'Simpan Perubahan'}
              </Button>
            </Box>
          )}
        </form>
      </Paper>
    </Box>
  );
}
