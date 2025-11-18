// ============================================================================
// Member Lookup Component - Optional Member Selection for POS
// Autocomplete search for binding sales to members
// ============================================================================

'use client';

import React, { useState, useEffect } from 'react';
import {
  Autocomplete,
  TextField,
  Box,
  Typography,
  Chip,
  CircularProgress,
  InputAdornment,
  IconButton,
} from '@mui/material';
import {
  Person as PersonIcon,
  Clear as ClearIcon,
} from '@mui/icons-material';
import { useToast } from '@/lib/context/ToastContext';
import memberApi from '@/lib/api/memberApi';
import type { Member } from '@/types';

// ============================================================================
// Component Props
// ============================================================================

interface MemberLookupProps {
  selectedMember: Member | null;
  onMemberSelect: (member: Member | null) => void;
  disabled?: boolean;
}

// ============================================================================
// Member Lookup Component
// ============================================================================

export default function MemberLookup({
  selectedMember,
  onMemberSelect,
  disabled = false,
}: MemberLookupProps) {
  const { showError } = useToast();
  const [open, setOpen] = useState(false);
  const [options, setOptions] = useState<Member[]>([]);
  const [loading, setLoading] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const [searchTerm, setSearchTerm] = useState('');

  // Debounced search
  useEffect(() => {
    const timer = setTimeout(() => {
      setSearchTerm(inputValue);
    }, 300);

    return () => clearTimeout(timer);
  }, [inputValue]);

  // Fetch members when search term changes
  useEffect(() => {
    if (!searchTerm) {
      setOptions(selectedMember ? [selectedMember] : []);
      return;
    }

    const fetchMembers = async () => {
      setLoading(true);
      try {
        const response = await memberApi.getMembers({
          search: searchTerm,
          status: 'aktif',
          pageSize: 10,
        });

        setOptions(response.data || []);
      } catch {
        showError('Gagal mencari anggota');
        setOptions([]);
      } finally {
        setLoading(false);
      }
    };

    fetchMembers();
  }, [searchTerm, selectedMember, showError]);

  // Handle member selection
  const handleSelect = (_event: React.SyntheticEvent, value: Member | null) => {
    onMemberSelect(value);
  };

  // Handle clear
  const handleClear = () => {
    onMemberSelect(null);
    setInputValue('');
    setSearchTerm('');
  };

  // Get status color
  const getStatusColor = (status: string): 'success' | 'default' | 'error' => {
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

  return (
    <Box>
      <Autocomplete
        open={open}
        onOpen={() => setOpen(true)}
        onClose={() => setOpen(false)}
        disabled={disabled}
        value={selectedMember}
        inputValue={inputValue}
        onInputChange={(_event, newInputValue) => setInputValue(newInputValue)}
        onChange={handleSelect}
        isOptionEqualToValue={(option, value) => option.id === value.id}
        getOptionLabel={(option) => option.namaLengkap}
        options={options}
        loading={loading}
        noOptionsText={
          searchTerm
            ? 'Anggota tidak ditemukan'
            : 'Ketik untuk mencari anggota (nama atau nomor)'
        }
        renderInput={(params) => (
          <TextField
            {...params}
            label="Anggota (Opsional)"
            placeholder="Cari anggota atau transaksi sebagai guest"
            variant="outlined"
            fullWidth
            InputProps={{
              ...params.InputProps,
              startAdornment: (
                <InputAdornment position="start">
                  <PersonIcon />
                </InputAdornment>
              ),
              endAdornment: (
                <>
                  {loading ? <CircularProgress color="inherit" size={20} /> : null}
                  {selectedMember && (
                    <IconButton
                      size="small"
                      onClick={handleClear}
                      sx={{ mr: 1 }}
                    >
                      <ClearIcon fontSize="small" />
                    </IconButton>
                  )}
                  {params.InputProps.endAdornment}
                </>
              ),
            }}
          />
        )}
        renderOption={(props, option) => (
          <Box component="li" {...props} key={option.id}>
            <Box sx={{ flexGrow: 1 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                <Typography variant="body2" fontWeight={500}>
                  {option.namaLengkap}
                </Typography>
                <Chip
                  label={option.status}
                  size="small"
                  color={getStatusColor(option.status)}
                  sx={{ height: 18 }}
                />
              </Box>
              <Box sx={{ display: 'flex', gap: 2 }}>
                <Typography variant="caption" color="text.secondary">
                  No: {option.nomorAnggota}
                </Typography>
                {option.noTelepon && (
                  <Typography variant="caption" color="text.secondary">
                    Tel: {option.noTelepon}
                  </Typography>
                )}
              </Box>
            </Box>
          </Box>
        )}
      />

      {/* Selected Member Info */}
      {selectedMember && (
        <Box
          sx={{
            mt: 2,
            p: 2,
            bgcolor: 'primary.light',
            borderRadius: 2,
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <Box>
            <Typography variant="body2" fontWeight={600} color="primary.contrastText">
              {selectedMember.namaLengkap}
            </Typography>
            <Typography variant="caption" color="primary.contrastText">
              Nomor Anggota: {selectedMember.nomorAnggota}
            </Typography>
            {selectedMember.noTelepon && (
              <Typography variant="caption" color="primary.contrastText" display="block">
                Telepon: {selectedMember.noTelepon}
              </Typography>
            )}
          </Box>
          <IconButton
            size="small"
            onClick={handleClear}
            sx={{ color: 'primary.contrastText' }}
          >
            <ClearIcon />
          </IconButton>
        </Box>
      )}
    </Box>
  );
}
