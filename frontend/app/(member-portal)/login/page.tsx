// ============================================================================
// Member Portal Login Page
// Login page specifically for cooperative members
// ============================================================================

'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
  Box,
  Card,
  CardContent,
  TextField,
  Button,
  Typography,
  Alert,
  Container,
  InputAdornment,
  IconButton,
} from '@mui/material';
import { Visibility, VisibilityOff, AccountCircle } from '@mui/icons-material';
import { useAuth } from '@/lib/context/AuthContext';
import type { LoginRequest } from '@/types';

// ============================================================================
// Validation Schema
// ============================================================================

const loginSchema = z.object({
  namaPengguna: z.string().min(1, 'Nama pengguna harus diisi'),
  kataSandi: z.string().min(1, 'Kata sandi harus diisi'),
});

type LoginFormData = z.infer<typeof loginSchema>;

// ============================================================================
// Member Portal Login Page Component
// ============================================================================

export default function MemberLoginPage() {
  const router = useRouter();
  const { login, isAuthenticated, user } = useAuth();
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string>('');

  // Redirect if already authenticated
  React.useEffect(() => {
    if (isAuthenticated && user) {
      if (user.peran === 'anggota') {
        router.push('/portal');
      } else {
        router.push('/dashboard');
      }
    }
  }, [isAuthenticated, user, router]);

  // ============================================================================
  // Form Setup
  // ============================================================================

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      namaPengguna: '',
      kataSandi: '',
    },
  });

  // ============================================================================
  // Submit Handler
  // ============================================================================

  const onSubmit = async (data: LoginFormData) => {
    try {
      setError('');
      await login(data as LoginRequest);

      // Note: Redirect is handled in useEffect after user state is updated
    } catch (err: unknown) {
      console.error('Login failed:', err);

      // Extract error message
      if (err && typeof err === 'object' && 'message' in err) {
        setError(err.message as string);
      } else {
        setError('Login gagal. Silakan periksa kembali nama pengguna dan kata sandi Anda.');
      }
    }
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Container maxWidth="sm">
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          py: 4,
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        }}
      >
        <Card
          sx={{
            width: '100%',
            maxWidth: 500,
            boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
          }}
        >
          <CardContent sx={{ p: 4 }}>
            {/* Header */}
            <Box sx={{ textAlign: 'center', mb: 4 }}>
              <AccountCircle sx={{ fontSize: 60, color: 'primary.main', mb: 2 }} />
              <Typography variant="h4" component="h1" gutterBottom fontWeight={600}>
                Portal Anggota
              </Typography>
              <Typography variant="body1" color="text.secondary">
                Masuk ke portal anggota koperasi
              </Typography>
            </Box>

            {/* Error Alert */}
            {error && (
              <Alert severity="error" sx={{ mb: 3 }}>
                {error}
              </Alert>
            )}

            {/* Login Form */}
            <form onSubmit={handleSubmit(onSubmit)}>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                {/* Username Field */}
                <TextField
                  {...register('namaPengguna')}
                  label="Nama Pengguna / Nomor Anggota"
                  variant="outlined"
                  fullWidth
                  error={!!errors.namaPengguna}
                  helperText={errors.namaPengguna?.message}
                  autoComplete="username"
                  autoFocus
                  disabled={isSubmitting}
                />

                {/* Password Field */}
                <TextField
                  {...register('kataSandi')}
                  label="Kata Sandi"
                  type={showPassword ? 'text' : 'password'}
                  variant="outlined"
                  fullWidth
                  error={!!errors.kataSandi}
                  helperText={errors.kataSandi?.message}
                  autoComplete="current-password"
                  disabled={isSubmitting}
                  InputProps={{
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton
                          aria-label="toggle password visibility"
                          onClick={() => setShowPassword(!showPassword)}
                          edge="end"
                        >
                          {showPassword ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                      </InputAdornment>
                    ),
                  }}
                />

                {/* Submit Button */}
                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  fullWidth
                  disabled={isSubmitting}
                  sx={{
                    mt: 1,
                    py: 1.5,
                    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                    '&:hover': {
                      background: 'linear-gradient(135deg, #5568d3 0%, #63408b 100%)',
                    }
                  }}
                >
                  {isSubmitting ? 'Memproses...' : 'Masuk'}
                </Button>
              </Box>
            </form>

            {/* Footer Info */}
            <Box sx={{ mt: 4, textAlign: 'center' }}>
              <Typography variant="caption" color="text.secondary">
                Portal Anggota Koperasi
              </Typography>
              <Typography variant="caption" display="block" color="text.secondary" sx={{ mt: 1 }}>
                Butuh bantuan? Hubungi pengurus koperasi Anda
              </Typography>
            </Box>
          </CardContent>
        </Card>
      </Box>
    </Container>
  );
}
