// ============================================================================
// Header Component - App Bar with User Info
// Material-UI AppBar with user menu and logout
// ============================================================================

'use client';

import React, { useState } from 'react';
import {
  AppBar,
  Toolbar,
  IconButton,
  Typography,
  Box,
  Avatar,
  Menu,
  MenuItem,
  ListItemIcon,
  Divider,
  Chip,
} from '@mui/material';
import {
  Menu as MenuIcon,
  AccountCircle as AccountCircleIcon,
  Logout as LogoutIcon,
  Settings as SettingsIcon,
} from '@mui/icons-material';
import { useAuth } from '@/lib/context/AuthContext';
import { DRAWER_WIDTH } from './Sidebar';

// ============================================================================
// Header Props
// ============================================================================

interface HeaderProps {
  onMenuClick: () => void;
  title?: string;
}

// ============================================================================
// Role Display Names (Indonesian)
// ============================================================================

const roleDisplayNames: Record<string, string> = {
  admin: 'Administrator',
  bendahara: 'Bendahara',
  kasir: 'Kasir',
  anggota: 'Anggota',
};

const roleColors: Record<string, 'error' | 'primary' | 'info' | 'default'> = {
  admin: 'error',
  bendahara: 'primary',
  kasir: 'info',
  anggota: 'default',
};

// ============================================================================
// Header Component
// ============================================================================

export default function Header({ onMenuClick, title }: HeaderProps) {
  const { user, logout } = useAuth();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  // ============================================================================
  // Menu Handlers
  // ============================================================================

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleMenuClose();
    logout();
  };

  // ============================================================================
  // Helper Functions
  // ============================================================================

  const getInitials = (name: string): string => {
    const names = name.trim().split(' ');
    if (names.length >= 2) {
      return `${names[0][0]}${names[names.length - 1][0]}`.toUpperCase();
    }
    return name.substring(0, 2).toUpperCase();
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <AppBar
      position="fixed"
      sx={{
        width: { md: `calc(100% - ${DRAWER_WIDTH}px)` },
        ml: { md: `${DRAWER_WIDTH}px` },
        backgroundColor: 'background.paper',
        color: 'text.primary',
        boxShadow: 1,
      }}
    >
      <Toolbar>
        {/* Mobile Menu Button */}
        <IconButton
          color="inherit"
          aria-label="open drawer"
          edge="start"
          onClick={onMenuClick}
          sx={{ mr: 2, display: { md: 'none' } }}
        >
          <MenuIcon />
        </IconButton>

        {/* Page Title */}
        <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
          {title || 'Dashboard'}
        </Typography>

        {/* User Info */}
        {user && (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            {/* User Role Badge */}
            <Chip
              label={roleDisplayNames[user.peran] || user.peran}
              color={roleColors[user.peran] || 'default'}
              size="small"
              sx={{ display: { xs: 'none', sm: 'flex' } }}
            />

            {/* User Menu Button */}
            <IconButton
              onClick={handleMenuOpen}
              size="small"
              sx={{ ml: 1 }}
              aria-controls={open ? 'user-menu' : undefined}
              aria-haspopup="true"
              aria-expanded={open ? 'true' : undefined}
            >
              <Avatar sx={{ width: 36, height: 36, bgcolor: 'primary.main' }}>
                {getInitials(user.namaLengkap)}
              </Avatar>
            </IconButton>
          </Box>
        )}

        {/* User Dropdown Menu */}
        <Menu
          anchorEl={anchorEl}
          id="user-menu"
          open={open}
          onClose={handleMenuClose}
          onClick={handleMenuClose}
          slotProps={{
            paper: {
              elevation: 3,
              sx: {
                overflow: 'visible',
                filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.15))',
                mt: 1.5,
                minWidth: 240,
                '& .MuiAvatar-root': {
                  width: 32,
                  height: 32,
                  ml: -0.5,
                  mr: 1,
                },
              },
            },
          }}
          transformOrigin={{ horizontal: 'right', vertical: 'top' }}
          anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
        >
          {/* User Info Header */}
          {user && (
            <Box sx={{ px: 2, py: 1.5 }}>
              <Typography variant="subtitle2" fontWeight={600}>
                {user.namaLengkap}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {user.email}
              </Typography>
              <Box sx={{ mt: 1 }}>
                <Chip
                  label={roleDisplayNames[user.peran] || user.peran}
                  color={roleColors[user.peran] || 'default'}
                  size="small"
                />
              </Box>
            </Box>
          )}

          <Divider />

          {/* Menu Items */}
          <MenuItem onClick={handleMenuClose}>
            <ListItemIcon>
              <AccountCircleIcon fontSize="small" />
            </ListItemIcon>
            Profil Saya
          </MenuItem>

          <MenuItem onClick={handleMenuClose}>
            <ListItemIcon>
              <SettingsIcon fontSize="small" />
            </ListItemIcon>
            Pengaturan
          </MenuItem>

          <Divider />

          <MenuItem onClick={handleLogout}>
            <ListItemIcon>
              <LogoutIcon fontSize="small" />
            </ListItemIcon>
            Keluar
          </MenuItem>
        </Menu>
      </Toolbar>
    </AppBar>
  );
}
