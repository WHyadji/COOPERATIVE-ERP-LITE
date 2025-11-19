// ============================================================================
// Member Portal Login Page
// Login page specifically for cooperative members
// ============================================================================

"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
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
} from "@mui/material";
import { Visibility, VisibilityOff, AccountCircle } from "@mui/icons-material";
import apiClient, { tokenManager } from "@/lib/api/client";
import type { APIResponse } from "@/types";

// Member Portal Login Response
interface MemberLoginResponse {
  token: string;
  anggota: {
    id: string;
    idKoperasi: string;
    nomorAnggota: string;
    namaLengkap: string;
    email?: string;
    status: string;
  };
}

// ============================================================================
// Validation Schema
// ============================================================================

const loginSchema = z.object({
  nomorAnggota: z.string().min(1, "Nomor anggota harus diisi"),
  pin: z
    .string()
    .length(6, "PIN harus 6 digit")
    .regex(/^\d+$/, "PIN harus berupa angka"),
});

type LoginFormData = z.infer<typeof loginSchema>;

// ============================================================================
// Member Portal Login Page Component
// ============================================================================

export default function MemberLoginPage() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string>("");

  // Check if already logged in
  React.useEffect(() => {
    const token = tokenManager.getToken();
    if (token) {
      router.push("/portal");
    }
  }, [router]);

  // Member Portal Login Function
  const loginMemberPortal = async (nomorAnggota: string, pin: string) => {
    try {
      // Note: Backend requires idKoperasi - for MVP we can use a default or env variable
      // In production, this should be determined by the domain or user selection
      const idKoperasi =
        process.env.NEXT_PUBLIC_DEFAULT_KOPERASI_ID || "default-koperasi-id";

      const response = await apiClient.post<APIResponse<MemberLoginResponse>>(
        "/portal/login",
        { nomorAnggota, pin },
        { params: { idKoperasi } }
      );

      if (response.data.success && response.data.data) {
        const { token } = response.data.data;

        // Store token
        tokenManager.setToken(token);

        // Redirect to portal
        router.push("/portal");
      } else {
        throw new Error("Login failed: Invalid response");
      }
    } catch (err: unknown) {
      console.error("Member portal login error:", err);
      throw err;
    }
  };

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
      nomorAnggota: "",
      pin: "",
    },
  });

  // ============================================================================
  // Submit Handler
  // ============================================================================

  const onSubmit = async (data: LoginFormData) => {
    try {
      setError("");
      // Member portal login uses separate endpoint with nomor anggota + PIN
      await loginMemberPortal(data.nomorAnggota, data.pin);

      // Note: Redirect is handled in useEffect after user state is updated
    } catch (err: unknown) {
      console.error("Login failed:", err);

      // Extract error message
      if (err && typeof err === "object" && "message" in err) {
        setError(err.message as string);
      } else {
        setError(
          "Login gagal. Silakan periksa kembali nomor anggota dan PIN Anda."
        );
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
          minHeight: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          py: 4,
          background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
        }}
      >
        <Card
          sx={{
            width: "100%",
            maxWidth: 500,
            boxShadow: "0 8px 32px rgba(0, 0, 0, 0.1)",
          }}
        >
          <CardContent sx={{ p: 4 }}>
            {/* Header */}
            <Box sx={{ textAlign: "center", mb: 4 }}>
              <AccountCircle
                sx={{ fontSize: 60, color: "primary.main", mb: 2 }}
              />
              <Typography
                variant="h4"
                component="h1"
                gutterBottom
                fontWeight={600}
              >
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
              <Box sx={{ display: "flex", flexDirection: "column", gap: 3 }}>
                {/* Member Number Field */}
                <TextField
                  {...register("nomorAnggota")}
                  label="Nomor Anggota"
                  variant="outlined"
                  fullWidth
                  error={!!errors.nomorAnggota}
                  helperText={
                    errors.nomorAnggota?.message ||
                    "Masukkan nomor anggota Anda"
                  }
                  autoComplete="off"
                  autoFocus
                  disabled={isSubmitting}
                  placeholder="Contoh: A001"
                />

                {/* PIN Field */}
                <TextField
                  {...register("pin")}
                  label="PIN (6 digit)"
                  type={showPassword ? "text" : "password"}
                  variant="outlined"
                  fullWidth
                  error={!!errors.pin}
                  helperText={
                    errors.pin?.message || "Masukkan PIN 6 digit Anda"
                  }
                  autoComplete="off"
                  disabled={isSubmitting}
                  inputProps={{
                    maxLength: 6,
                    inputMode: "numeric",
                    pattern: "[0-9]*",
                  }}
                  InputProps={{
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton
                          aria-label="toggle PIN visibility"
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
                    background:
                      "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
                    "&:hover": {
                      background:
                        "linear-gradient(135deg, #5568d3 0%, #63408b 100%)",
                    },
                  }}
                >
                  {isSubmitting ? "Memproses..." : "Masuk"}
                </Button>
              </Box>
            </form>

            {/* Footer Info */}
            <Box sx={{ mt: 4, textAlign: "center" }}>
              <Typography variant="caption" color="text.secondary">
                Portal Anggota Koperasi
              </Typography>
              <Typography
                variant="caption"
                display="block"
                color="text.secondary"
                sx={{ mt: 1 }}
              >
                Butuh bantuan? Hubungi pengurus koperasi Anda
              </Typography>
            </Box>
          </CardContent>
        </Card>
      </Box>
    </Container>
  );
}
