// ============================================================================
// Member Profile Page
// View and edit member profile information
// ============================================================================

'use client';

import React, { useEffect, useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  CircularProgress,
  Alert,
  Grid,
  TextField,
  Button,
  Divider,
  Avatar,
  Chip,
} from '@mui/material';
import {
  Person,
  Edit,
  Save,
  Cancel,
  Phone,
  Email,
  Home,
  Work,
  Cake,
  LocationOn,
} from '@mui/icons-material';
import { getMemberProfile, updateMemberProfile } from '@/lib/api/memberPortalApi';
import type { Member } from '@/types';

// ============================================================================
// Helper Functions
// ============================================================================

const formatDate = (dateString: string | undefined): string => {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
};

const getStatusColor = (status: string): 'success' | 'warning' | 'error' => {
  const colors: Record<string, 'success' | 'warning' | 'error'> = {
    aktif: 'success',
    nonaktif: 'warning',
    diberhentikan: 'error',
  };
  return colors[status] || 'warning';
};

const getStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    aktif: 'Aktif',
    nonaktif: 'Non-Aktif',
    diberhentikan: 'Diberhentikan',
  };
  return labels[status] || status;
};

// ============================================================================
// Member Profile Component
// ============================================================================

export default function MemberProfilePage() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');
  const [isEditing, setIsEditing] = useState(false);
  const [profile, setProfile] = useState<Member | null>(null);
  const [editedProfile, setEditedProfile] = useState<Partial<Member>>({});

  // ============================================================================
  // Fetch Profile
  // ============================================================================

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        setLoading(true);
        setError('');
        const data = await getMemberProfile();
        setProfile(data);
        setEditedProfile(data);
      } catch (err: unknown) {
        console.error('Failed to fetch profile:', err);
        if (err && typeof err === 'object' && 'message' in err) {
          setError(err.message as string);
        } else {
          setError('Gagal memuat data profil');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  // ============================================================================
  // Handle Edit
  // ============================================================================

  const handleEdit = () => {
    setIsEditing(true);
    setSuccess('');
    setError('');
  };

  const handleCancel = () => {
    setIsEditing(false);
    setEditedProfile(profile || {});
    setError('');
  };

  const handleChange = (field: keyof Member, value: string) => {
    setEditedProfile((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleSave = async () => {
    try {
      setSaving(true);
      setError('');
      setSuccess('');

      // Only send editable fields
      const updateData = {
        noTelepon: editedProfile.noTelepon,
        email: editedProfile.email,
        alamat: editedProfile.alamat,
        rt: editedProfile.rt,
        rw: editedProfile.rw,
        kelurahan: editedProfile.kelurahan,
        kecamatan: editedProfile.kecamatan,
        kotaKabupaten: editedProfile.kotaKabupaten,
        provinsi: editedProfile.provinsi,
        kodePos: editedProfile.kodePos,
      };

      const updatedProfile = await updateMemberProfile(updateData);
      setProfile(updatedProfile);
      setEditedProfile(updatedProfile);
      setIsEditing(false);
      setSuccess('Profil berhasil diperbarui');
    } catch (err: unknown) {
      console.error('Failed to update profile:', err);
      if (err && typeof err === 'object' && 'message' in err) {
        setError(err.message as string);
      } else {
        setError('Gagal memperbarui profil');
      }
    } finally {
      setSaving(false);
    }
  };

  // ============================================================================
  // Loading State
  // ============================================================================

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 400 }}>
        <CircularProgress />
      </Box>
    );
  }

  // ============================================================================
  // Error State
  // ============================================================================

  if (!profile && error) {
    return (
      <Alert severity="error" sx={{ mb: 3 }}>
        {error}
      </Alert>
    );
  }

  // ============================================================================
  // Render Profile
  // ============================================================================

  return (
    <Box>
      {/* Page Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" fontWeight={600} gutterBottom>
          Profil Saya
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Informasi data anggota koperasi
        </Typography>
      </Box>

      {/* Success/Error Messages */}
      {success && (
        <Alert severity="success" sx={{ mb: 3 }} onClose={() => setSuccess('')}>
          {success}
        </Alert>
      )}
      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError('')}>
          {error}
        </Alert>
      )}

      {/* Profile Header Card */}
      {profile && (
        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Box sx={{ display: 'flex', alignItems: 'center', flexWrap: 'wrap', gap: 3 }}>
              <Avatar
                sx={{
                  width: 100,
                  height: 100,
                  bgcolor: 'primary.main',
                  fontSize: 40,
                }}
              >
                {profile.namaLengkap?.charAt(0).toUpperCase()}
              </Avatar>

              <Box sx={{ flexGrow: 1 }}>
                <Typography variant="h5" fontWeight={600} gutterBottom>
                  {profile.namaLengkap}
                </Typography>
                <Typography variant="body1" color="text.secondary" gutterBottom>
                  Nomor Anggota: <strong>{profile.nomorAnggota}</strong>
                </Typography>
                <Box sx={{ mt: 1 }}>
                  <Chip
                    label={getStatusLabel(profile.status)}
                    color={getStatusColor(profile.status)}
                    size="small"
                  />
                </Box>
              </Box>

              {!isEditing && (
                <Button
                  variant="contained"
                  startIcon={<Edit />}
                  onClick={handleEdit}
                >
                  Edit Profil
                </Button>
              )}
            </Box>
          </CardContent>
        </Card>
      )}

      {/* Personal Information */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
            <Person sx={{ mr: 1 }} />
            <Typography variant="h6">Informasi Pribadi</Typography>
          </Box>
          <Divider sx={{ mb: 3 }} />

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="NIK"
                value={profile?.nik || '-'}
                disabled
                InputProps={{
                  startAdornment: <Person sx={{ mr: 1, color: 'text.secondary' }} />,
                }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Jenis Kelamin"
                value={profile?.jenisKelamin === 'L' ? 'Laki-laki' : profile?.jenisKelamin === 'P' ? 'Perempuan' : '-'}
                disabled
                InputProps={{
                  startAdornment: <Person sx={{ mr: 1, color: 'text.secondary' }} />,
                }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Tempat Lahir"
                value={profile?.tempatLahir || '-'}
                disabled
                InputProps={{
                  startAdornment: <LocationOn sx={{ mr: 1, color: 'text.secondary' }} />,
                }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Tanggal Lahir"
                value={formatDate(profile?.tanggalLahir)}
                disabled
                InputProps={{
                  startAdornment: <Cake sx={{ mr: 1, color: 'text.secondary' }} />,
                }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Pekerjaan"
                value={profile?.pekerjaan || '-'}
                disabled
                InputProps={{
                  startAdornment: <Work sx={{ mr: 1, color: 'text.secondary' }} />,
                }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Tanggal Bergabung"
                value={formatDate(profile?.tanggalBergabung)}
                disabled
                InputProps={{
                  startAdornment: <Cake sx={{ mr: 1, color: 'text.secondary' }} />,
                }}
              />
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Contact Information */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
            <Phone sx={{ mr: 1 }} />
            <Typography variant="h6">Informasi Kontak</Typography>
          </Box>
          <Divider sx={{ mb: 3 }} />

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Nomor Telepon"
                value={isEditing ? editedProfile.noTelepon || '' : profile?.noTelepon || '-'}
                onChange={(e) => handleChange('noTelepon', e.target.value)}
                disabled={!isEditing}
                InputProps={{
                  startAdornment: <Phone sx={{ mr: 1, color: 'text.secondary' }} />,
                }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Email"
                type="email"
                value={isEditing ? editedProfile.email || '' : profile?.email || '-'}
                onChange={(e) => handleChange('email', e.target.value)}
                disabled={!isEditing}
                InputProps={{
                  startAdornment: <Email sx={{ mr: 1, color: 'text.secondary' }} />,
                }}
              />
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Address Information */}
      <Card>
        <CardContent>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
            <Home sx={{ mr: 1 }} />
            <Typography variant="h6">Alamat</Typography>
          </Box>
          <Divider sx={{ mb: 3 }} />

          <Grid container spacing={3}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Alamat Lengkap"
                value={isEditing ? editedProfile.alamat || '' : profile?.alamat || '-'}
                onChange={(e) => handleChange('alamat', e.target.value)}
                disabled={!isEditing}
                multiline
                rows={2}
                InputProps={{
                  startAdornment: <Home sx={{ mr: 1, color: 'text.secondary', alignSelf: 'flex-start', mt: 1 }} />,
                }}
              />
            </Grid>

            <Grid item xs={6} md={3}>
              <TextField
                fullWidth
                label="RT"
                value={isEditing ? editedProfile.rt || '' : profile?.rt || '-'}
                onChange={(e) => handleChange('rt', e.target.value)}
                disabled={!isEditing}
              />
            </Grid>

            <Grid item xs={6} md={3}>
              <TextField
                fullWidth
                label="RW"
                value={isEditing ? editedProfile.rw || '' : profile?.rw || '-'}
                onChange={(e) => handleChange('rw', e.target.value)}
                disabled={!isEditing}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Kelurahan"
                value={isEditing ? editedProfile.kelurahan || '' : profile?.kelurahan || '-'}
                onChange={(e) => handleChange('kelurahan', e.target.value)}
                disabled={!isEditing}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Kecamatan"
                value={isEditing ? editedProfile.kecamatan || '' : profile?.kecamatan || '-'}
                onChange={(e) => handleChange('kecamatan', e.target.value)}
                disabled={!isEditing}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Kota/Kabupaten"
                value={isEditing ? editedProfile.kotaKabupaten || '' : profile?.kotaKabupaten || '-'}
                onChange={(e) => handleChange('kotaKabupaten', e.target.value)}
                disabled={!isEditing}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Provinsi"
                value={isEditing ? editedProfile.provinsi || '' : profile?.provinsi || '-'}
                onChange={(e) => handleChange('provinsi', e.target.value)}
                disabled={!isEditing}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Kode Pos"
                value={isEditing ? editedProfile.kodePos || '' : profile?.kodePos || '-'}
                onChange={(e) => handleChange('kodePos', e.target.value)}
                disabled={!isEditing}
              />
            </Grid>
          </Grid>

          {/* Edit Actions */}
          {isEditing && (
            <Box sx={{ mt: 3, display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
              <Button
                variant="outlined"
                startIcon={<Cancel />}
                onClick={handleCancel}
                disabled={saving}
              >
                Batal
              </Button>
              <Button
                variant="contained"
                startIcon={<Save />}
                onClick={handleSave}
                disabled={saving}
              >
                {saving ? 'Menyimpan...' : 'Simpan Perubahan'}
              </Button>
            </Box>
          )}
        </CardContent>
      </Card>

      {/* Info Notice */}
      <Alert severity="info" sx={{ mt: 3 }}>
        <Typography variant="body2">
          <strong>Catatan:</strong> Hanya informasi kontak dan alamat yang dapat diubah.
          Untuk mengubah data pribadi lainnya, silakan hubungi pengurus koperasi.
        </Typography>
      </Alert>
    </Box>
  );
}
