// ============================================================================
// Account Form Component - Create/Edit Chart of Accounts
// Material-UI form with validation for account creation and editing
// ============================================================================

'use client';

import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  FormHelperText,
  Grid,
  Alert,
  CircularProgress,
} from '@mui/material';
import accountingApi from '@/lib/api/accountingApi';
import type { Akun, TipeAkun, NormalSaldo, AkunFormData } from '@/types';

// ============================================================================
// Component Props
// ============================================================================

interface AccountFormProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
  account?: Akun | null; // If provided, edit mode
  parentAccounts?: Akun[]; // For hierarchical COA
}

// ============================================================================
// Account Form Component
// ============================================================================

export default function AccountForm({
  open,
  onClose,
  onSuccess,
  account,
  parentAccounts = [],
}: AccountFormProps) {
  const isEditMode = !!account;

  // Form state
  const [formData, setFormData] = useState<AkunFormData>({
    kodeAkun: '',
    namaAkun: '',
    tipeAkun: '',
    idInduk: '',
    normalSaldo: '',
    deskripsi: '',
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const [errors, setErrors] = useState<Partial<Record<keyof AkunFormData, string>>>(
    {}
  );

  // ============================================================================
  // Initialize form data when account changes
  // ============================================================================

  useEffect(() => {
    if (account) {
      setFormData({
        kodeAkun: account.kodeAkun,
        namaAkun: account.namaAkun,
        tipeAkun: account.tipeAkun,
        idInduk: account.idInduk || '',
        normalSaldo: account.normalSaldo,
        deskripsi: account.deskripsi || '',
      });
    } else {
      setFormData({
        kodeAkun: '',
        namaAkun: '',
        tipeAkun: '',
        idInduk: '',
        normalSaldo: '',
        deskripsi: '',
      });
    }
    setError('');
    setErrors({});
  }, [account, open]);

  // ============================================================================
  // Validation
  // ============================================================================

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof AkunFormData, string>> = {};

    if (!formData.kodeAkun.trim()) {
      newErrors.kodeAkun = 'Kode akun harus diisi';
    }

    if (!formData.namaAkun.trim()) {
      newErrors.namaAkun = 'Nama akun harus diisi';
    }

    if (!formData.tipeAkun) {
      newErrors.tipeAkun = 'Tipe akun harus dipilih';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // ============================================================================
  // Handlers
  // ============================================================================

  const handleChange = (field: keyof AkunFormData, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear error for this field
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: '' }));
    }

    // Auto-set normal saldo based on tipe akun
    if (field === 'tipeAkun' && !formData.normalSaldo) {
      const tipe = value as TipeAkun;
      let normalSaldo: NormalSaldo = 'debit';

      if (tipe === 'aset' || tipe === 'beban') {
        normalSaldo = 'debit';
      } else if (
        tipe === 'kewajiban' ||
        tipe === 'modal' ||
        tipe === 'pendapatan'
      ) {
        normalSaldo = 'kredit';
      }

      setFormData((prev) => ({ ...prev, normalSaldo }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      setLoading(true);
      setError('');

      const requestData = {
        kodeAkun: formData.kodeAkun.trim(),
        namaAkun: formData.namaAkun.trim(),
        tipeAkun: formData.tipeAkun as TipeAkun,
        idInduk: formData.idInduk || undefined,
        normalSaldo: formData.normalSaldo as NormalSaldo,
        deskripsi: formData.deskripsi?.trim() || undefined,
      };

      if (isEditMode && account) {
        await accountingApi.updateAccount(account.id, {
          ...requestData,
          statusAktif: account.statusAktif,
        });
      } else {
        await accountingApi.createAccount(requestData);
      }

      onSuccess();
    } catch (err: unknown) {
      console.error('Failed to save account:', err);
      setError(
        err instanceof Error
          ? err.message
          : 'Gagal menyimpan akun. Silakan coba lagi.'
      );
    } finally {
      setLoading(false);
    }
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{isEditMode ? 'Edit Akun' : 'Tambah Akun Baru'}</DialogTitle>

        <DialogContent>
          {error && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {error}
            </Alert>
          )}

          <Grid container spacing={3} sx={{ mt: 0.5 }}>
            {/* Kode Akun */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Kode Akun"
                value={formData.kodeAkun}
                onChange={(e) => handleChange('kodeAkun', e.target.value)}
                error={!!errors.kodeAkun}
                helperText={errors.kodeAkun || 'Contoh: 1101, 2201, 3101'}
                required
                disabled={loading}
              />
            </Grid>

            {/* Nama Akun */}
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Nama Akun"
                value={formData.namaAkun}
                onChange={(e) => handleChange('namaAkun', e.target.value)}
                error={!!errors.namaAkun}
                helperText={errors.namaAkun}
                required
                disabled={loading}
              />
            </Grid>

            {/* Tipe Akun */}
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth error={!!errors.tipeAkun} required>
                <InputLabel>Tipe Akun</InputLabel>
                <Select
                  value={formData.tipeAkun}
                  label="Tipe Akun"
                  onChange={(e) => handleChange('tipeAkun', e.target.value)}
                  disabled={loading}
                >
                  <MenuItem value="aset">Aset</MenuItem>
                  <MenuItem value="kewajiban">Kewajiban</MenuItem>
                  <MenuItem value="modal">Modal</MenuItem>
                  <MenuItem value="pendapatan">Pendapatan</MenuItem>
                  <MenuItem value="beban">Beban</MenuItem>
                </Select>
                {errors.tipeAkun && (
                  <FormHelperText>{errors.tipeAkun}</FormHelperText>
                )}
              </FormControl>
            </Grid>

            {/* Normal Saldo */}
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Normal Saldo</InputLabel>
                <Select
                  value={formData.normalSaldo}
                  label="Normal Saldo"
                  onChange={(e) => handleChange('normalSaldo', e.target.value)}
                  disabled={loading}
                >
                  <MenuItem value="debit">Debit</MenuItem>
                  <MenuItem value="kredit">Kredit</MenuItem>
                </Select>
                <FormHelperText>
                  Otomatis terisi berdasarkan tipe akun
                </FormHelperText>
              </FormControl>
            </Grid>

            {/* Akun Induk */}
            <Grid item xs={12}>
              <FormControl fullWidth>
                <InputLabel>Akun Induk (Opsional)</InputLabel>
                <Select
                  value={formData.idInduk}
                  label="Akun Induk (Opsional)"
                  onChange={(e) => handleChange('idInduk', e.target.value)}
                  disabled={loading}
                >
                  <MenuItem value="">
                    <em>Tidak ada (Akun level atas)</em>
                  </MenuItem>
                  {parentAccounts
                    .filter(
                      (acc) => !account || acc.id !== account.id // Don't show self as parent
                    )
                    .map((acc) => (
                      <MenuItem key={acc.id} value={acc.id}>
                        {acc.kodeAkun} - {acc.namaAkun}
                      </MenuItem>
                    ))}
                </Select>
                <FormHelperText>Untuk membuat struktur hierarki akun</FormHelperText>
              </FormControl>
            </Grid>

            {/* Deskripsi */}
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Deskripsi (Opsional)"
                value={formData.deskripsi}
                onChange={(e) => handleChange('deskripsi', e.target.value)}
                multiline
                rows={3}
                disabled={loading}
              />
            </Grid>
          </Grid>
        </DialogContent>

        <DialogActions>
          <Button onClick={onClose} disabled={loading}>
            Batal
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={loading}
            startIcon={loading && <CircularProgress size={20} />}
          >
            {loading ? 'Menyimpan...' : 'Simpan'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}
