// ============================================================================
// Create Simpanan Page - Record new share capital deposit
// Form with React Hook Form and Zod validation
// ============================================================================

'use client';

import React, { useState, useEffect } from 'react';
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
  RadioGroup,
  FormControlLabel,
  Radio,
  FormLabel,
  Autocomplete,
  CircularProgress,
} from '@mui/material';
import { Save as SaveIcon, Cancel as CancelIcon } from '@mui/icons-material';
import simpananApi from '@/lib/api/simpananApi';
import memberApi from '@/lib/api/memberApi';
import type { CreateSimpananRequest, TipeSimpanan, Member } from '@/types';
import { format } from 'date-fns';

// ============================================================================
// Validation Schema
// ============================================================================

const simpananSchema = z.object({
  idAnggota: z.string().min(1, 'Anggota harus dipilih'),
  tipeSimpanan: z.enum(['pokok', 'wajib', 'sukarela'], {
    errorMap: () => ({ message: 'Tipe simpanan harus dipilih' }),
  }),
  tanggalTransaksi: z.string().min(1, 'Tanggal transaksi harus diisi'),
  jumlahSetoran: z
    .string()
    .min(1, 'Jumlah setoran harus diisi')
    .refine((val) => !isNaN(Number(val)) && Number(val) > 0, {
      message: 'Jumlah setoran harus lebih dari 0',
    }),
  keterangan: z.string().optional(),
});

type SimpananFormData = z.infer<typeof simpananSchema>;

// ============================================================================
// Create Simpanan Page Component
// ============================================================================

export default function CreateSimpananPage() {
  const router = useRouter();
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');
  const [members, setMembers] = useState<Member[]>([]);
  const [loadingMembers, setLoadingMembers] = useState(false);

  // ============================================================================
  // Form Setup
  // ============================================================================

  const {
    register,
    handleSubmit,
    control,
    formState: { errors, isSubmitting },
  } = useForm<SimpananFormData>({
    resolver: zodResolver(simpananSchema),
    defaultValues: {
      idAnggota: '',
      tipeSimpanan: 'wajib',
      tanggalTransaksi: format(new Date(), 'yyyy-MM-dd'),
      jumlahSetoran: '',
      keterangan: '',
    },
  });

  // ============================================================================
  // Fetch Members for Autocomplete
  // ============================================================================

  useEffect(() => {
    const fetchMembers = async () => {
      try {
        setLoadingMembers(true);
        const response = await memberApi.getMembers({
          status: 'aktif',
          pageSize: 1000, // Get all active members
        });
        setMembers(response.data);
      } catch (err) {
        console.error('Failed to fetch members:', err);
      } finally {
        setLoadingMembers(false);
      }
    };

    fetchMembers();
  }, []);

  // ============================================================================
  // Submit Handler
  // ============================================================================

  const onSubmit = async (data: SimpananFormData) => {
    try {
      setError('');
      setSuccess('');

      const requestData: CreateSimpananRequest = {
        idAnggota: data.idAnggota,
        tipeSimpanan: data.tipeSimpanan,
        tanggalTransaksi: data.tanggalTransaksi,
        jumlahSetoran: Number(data.jumlahSetoran),
        keterangan: data.keterangan || undefined,
      };

      await simpananApi.createSimpanan(requestData);

      setSuccess('Setoran simpanan berhasil dicatat!');

      // Redirect to simpanan list after 1 second
      setTimeout(() => {
        router.push('/dashboard/simpanan');
      }, 1000);
    } catch (err: unknown) {
      console.error('Failed to create simpanan:', err);

      if (err && typeof err === 'object' && 'message' in err) {
        setError(err.message as string);
      } else {
        setError('Gagal mencatat setoran. Silakan coba lagi.');
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
          Catat Setoran Simpanan
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Tambah catatan setoran simpanan anggota (Pokok, Wajib, atau Sukarela)
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
          <Grid container spacing={3}>
            {/* Pilih Anggota */}
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>
                Informasi Anggota
              </Typography>
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="idAnggota"
                control={control}
                render={({ field }) => (
                  <Autocomplete
                    {...field}
                    options={members}
                    getOptionLabel={(option) =>
                      typeof option === 'string'
                        ? option
                        : `${option.nomorAnggota} - ${option.namaLengkap}`
                    }
                    loading={loadingMembers}
                    onChange={(_, value) => field.onChange(value?.id || '')}
                    value={members.find((m) => m.id === field.value) || null}
                    renderInput={(params) => (
                      <TextField
                        {...params}
                        label="Pilih Anggota *"
                        error={!!errors.idAnggota}
                        helperText={errors.idAnggota?.message}
                        InputProps={{
                          ...params.InputProps,
                          endAdornment: (
                            <>
                              {loadingMembers ? (
                                <CircularProgress color="inherit" size={20} />
                              ) : null}
                              {params.InputProps.endAdornment}
                            </>
                          ),
                        }}
                      />
                    )}
                  />
                )}
              />
            </Grid>

            {/* Tipe Simpanan */}
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom sx={{ mt: 2 }}>
                Informasi Setoran
              </Typography>
            </Grid>

            <Grid item xs={12}>
              <Controller
                name="tipeSimpanan"
                control={control}
                render={({ field }) => (
                  <FormControl error={!!errors.tipeSimpanan} component="fieldset">
                    <FormLabel component="legend">Tipe Simpanan *</FormLabel>
                    <RadioGroup {...field} row>
                      <FormControlLabel
                        value="pokok"
                        control={<Radio />}
                        label="Simpanan Pokok (dibayar sekali)"
                      />
                      <FormControlLabel
                        value="wajib"
                        control={<Radio />}
                        label="Simpanan Wajib (bulanan)"
                      />
                      <FormControlLabel
                        value="sukarela"
                        control={<Radio />}
                        label="Simpanan Sukarela (opsional)"
                      />
                    </RadioGroup>
                    {errors.tipeSimpanan && (
                      <FormHelperText>{errors.tipeSimpanan.message}</FormHelperText>
                    )}
                  </FormControl>
                )}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Jumlah Setoran *"
                type="number"
                {...register('jumlahSetoran')}
                error={!!errors.jumlahSetoran}
                helperText={errors.jumlahSetoran?.message}
                InputProps={{
                  startAdornment: <Box sx={{ mr: 1 }}>Rp</Box>,
                }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Tanggal Transaksi *"
                type="date"
                {...register('tanggalTransaksi')}
                error={!!errors.tanggalTransaksi}
                helperText={errors.tanggalTransaksi?.message}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>

            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Keterangan"
                multiline
                rows={3}
                {...register('keterangan')}
                error={!!errors.keterangan}
                helperText={errors.keterangan?.message}
                placeholder="Catatan tambahan (opsional)"
              />
            </Grid>

            {/* Action Buttons */}
            <Grid item xs={12}>
              <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end', mt: 2 }}>
                <Button
                  variant="outlined"
                  startIcon={<CancelIcon />}
                  onClick={() => router.push('/dashboard/simpanan')}
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
                  {isSubmitting ? 'Menyimpan...' : 'Simpan Setoran'}
                </Button>
              </Box>
            </Grid>
          </Grid>
        </form>
      </Paper>
    </Box>
  );
}
