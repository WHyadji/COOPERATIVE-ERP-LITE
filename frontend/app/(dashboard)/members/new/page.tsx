// ============================================================================
// Create Member Page - Add new member
// Multi-section form with React Hook Form and Zod validation
// ============================================================================

'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
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
} from '@mui/material';
import { Save as SaveIcon, Cancel as CancelIcon } from '@mui/icons-material';
import memberApi from '@/lib/api/memberApi';
import type { CreateMemberRequest, Gender } from '@/types';
import { format } from 'date-fns';

// ============================================================================
// Validation Schema
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
});

type MemberFormData = z.infer<typeof memberSchema>;

// ============================================================================
// Create Member Page Component
// ============================================================================

export default function CreateMemberPage() {
  const router = useRouter();
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');

  // ============================================================================
  // Form Setup
  // ============================================================================

  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isSubmitting },
  } = useForm<MemberFormData>({
    resolver: zodResolver(memberSchema),
    defaultValues: {
      namaLengkap: '',
      nik: '',
      tanggalLahir: '',
      tempatLahir: '',
      jenisKelamin: '',
      alamat: '',
      rt: '',
      rw: '',
      kelurahan: '',
      kecamatan: '',
      kotaKabupaten: '',
      provinsi: '',
      kodePos: '',
      noTelepon: '',
      email: '',
      pekerjaan: '',
      tanggalBergabung: format(new Date(), 'yyyy-MM-dd'),
    },
  });

  // ============================================================================
  // Submit Handler
  // ============================================================================

  const onSubmit = async (data: MemberFormData) => {
    try {
      setError('');
      setSuccess('');

      // Convert form data to API request
      const requestData: CreateMemberRequest = {
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
      };

      const newMember = await memberApi.createMember(requestData);

      setSuccess('Anggota berhasil ditambahkan!');

      // Redirect to member detail page after 1 second
      setTimeout(() => {
        router.push(`/dashboard/members/${newMember.id}`);
      }, 1000);
    } catch (err: unknown) {
      console.error('Failed to create member:', err);

      if (err && typeof err === 'object' && 'message' in err) {
        setError(err.message as string);
      } else {
        setError('Gagal menambahkan anggota. Silakan coba lagi.');
      }
    }
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h4" fontWeight={600} gutterBottom>
          Tambah Anggota Baru
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Lengkapi formulir di bawah ini untuk menambahkan anggota baru
        </Typography>
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
                disabled={isSubmitting}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                {...register('nik')}
                label="NIK (KTP)"
                fullWidth
                error={!!errors.nik}
                helperText={errors.nik?.message}
                disabled={isSubmitting}
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
                disabled={isSubmitting}
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
                disabled={isSubmitting}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <FormControl fullWidth error={!!errors.jenisKelamin} disabled={isSubmitting}>
                <InputLabel>Jenis Kelamin</InputLabel>
                <Controller
                  name="jenisKelamin"
                  control={control}
                  render={({ field }) => (
                    <Select {...field} label="Jenis Kelamin">
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

            <Grid item xs={12} md={6}>
              <TextField
                {...register('pekerjaan')}
                label="Pekerjaan"
                fullWidth
                error={!!errors.pekerjaan}
                helperText={errors.pekerjaan?.message}
                disabled={isSubmitting}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                {...register('tanggalBergabung')}
                label="Tanggal Bergabung *"
                type="date"
                fullWidth
                error={!!errors.tanggalBergabung}
                helperText={errors.tanggalBergabung?.message}
                disabled={isSubmitting}
                InputLabelProps={{ shrink: true }}
              />
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
                disabled={isSubmitting}
              />
            </Grid>

            <Grid item xs={12} md={3}>
              <TextField
                {...register('rt')}
                label="RT"
                fullWidth
                error={!!errors.rt}
                helperText={errors.rt?.message}
                disabled={isSubmitting}
              />
            </Grid>

            <Grid item xs={12} md={3}>
              <TextField
                {...register('rw')}
                label="RW"
                fullWidth
                error={!!errors.rw}
                helperText={errors.rw?.message}
                disabled={isSubmitting}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                {...register('kelurahan')}
                label="Kelurahan/Desa"
                fullWidth
                error={!!errors.kelurahan}
                helperText={errors.kelurahan?.message}
                disabled={isSubmitting}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('kecamatan')}
                label="Kecamatan"
                fullWidth
                error={!!errors.kecamatan}
                helperText={errors.kecamatan?.message}
                disabled={isSubmitting}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('kotaKabupaten')}
                label="Kota/Kabupaten"
                fullWidth
                error={!!errors.kotaKabupaten}
                helperText={errors.kotaKabupaten?.message}
                disabled={isSubmitting}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('provinsi')}
                label="Provinsi"
                fullWidth
                error={!!errors.provinsi}
                helperText={errors.provinsi?.message}
                disabled={isSubmitting}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <TextField
                {...register('kodePos')}
                label="Kode Pos"
                fullWidth
                error={!!errors.kodePos}
                helperText={errors.kodePos?.message}
                disabled={isSubmitting}
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
                disabled={isSubmitting}
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
                disabled={isSubmitting}
              />
            </Grid>
          </Grid>

          {/* Action Buttons */}
          <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 4 }}>
            <Button
              variant="outlined"
              startIcon={<CancelIcon />}
              onClick={() => router.back()}
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
              {isSubmitting ? 'Menyimpan...' : 'Simpan'}
            </Button>
          </Box>
        </form>
      </Paper>
    </Box>
  );
}
