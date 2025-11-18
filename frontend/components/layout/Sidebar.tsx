// ============================================================================
// Sidebar Component - Navigation Menu
// Material-UI Drawer with role-based menu items
// ============================================================================

'use client';

import React from 'react';
import { usePathname, useRouter } from 'next/navigation';
import {
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  Divider,
} from '@mui/material';
import {
  Dashboard as DashboardIcon,
  People as PeopleIcon,
  AccountBalance as AccountBalanceIcon,
  Receipt as ReceiptIcon,
  Inventory as InventoryIcon,
  PointOfSale as PointOfSaleIcon,
  Assessment as AssessmentIcon,
  Settings as SettingsIcon,
} from '@mui/icons-material';
import { useAuth } from '@/lib/context/AuthContext';
import type { UserRole } from '@/types';

// ============================================================================
// Constants
// ============================================================================

const DRAWER_WIDTH = 260;

// ============================================================================
// Menu Item Type
// ============================================================================

interface MenuItem {
  label: string;
  icon: React.ReactNode;
  path: string;
  roles: UserRole[]; // Roles allowed to see this menu item
}

// ============================================================================
// Menu Items Configuration
// ============================================================================

const menuItems: MenuItem[] = [
  {
    label: 'Dashboard',
    icon: <DashboardIcon />,
    path: '/dashboard',
    roles: ['admin', 'bendahara', 'kasir', 'anggota'],
  },
  {
    label: 'Anggota',
    icon: <PeopleIcon />,
    path: '/dashboard/members',
    roles: ['admin', 'bendahara'],
  },
  {
    label: 'Simpanan',
    icon: <AccountBalanceIcon />,
    path: '/dashboard/simpanan',
    roles: ['admin', 'bendahara'],
  },
  {
    label: 'POS / Kasir',
    icon: <PointOfSaleIcon />,
    path: '/pos',
    roles: ['admin', 'kasir'],
  },
  {
    label: 'Produk',
    icon: <InventoryIcon />,
    path: '/produk',
    roles: ['admin', 'bendahara', 'kasir'],
  },
  {
    label: 'Bagan Akun',
    icon: <ReceiptIcon />,
    path: '/akuntansi',
    roles: ['admin', 'bendahara'],
  },
  {
    label: 'Jurnal Umum',
    icon: <ReceiptIcon />,
    path: '/akuntansi/jurnal',
    roles: ['admin', 'bendahara'],
  },
  {
    label: 'Laporan',
    icon: <AssessmentIcon />,
    path: '/dashboard/reports',
    roles: ['admin', 'bendahara'],
  },
  {
    label: 'Pengaturan',
    icon: <SettingsIcon />,
    path: '/dashboard/settings',
    roles: ['admin'],
  },
];

// ============================================================================
// Sidebar Props
// ============================================================================

interface SidebarProps {
  mobileOpen: boolean;
  onDrawerToggle: () => void;
}

// ============================================================================
// Sidebar Component
// ============================================================================

export default function Sidebar({ mobileOpen, onDrawerToggle }: SidebarProps) {
  const pathname = usePathname();
  const router = useRouter();
  const { user } = useAuth();

  // Filter menu items based on user role
  const visibleMenuItems = menuItems.filter((item) =>
    user ? item.roles.includes(user.peran) : false
  );

  const handleNavigate = (path: string) => {
    router.push(path);
    // Close mobile drawer after navigation
    if (mobileOpen) {
      onDrawerToggle();
    }
  };

  // ============================================================================
  // Drawer Content
  // ============================================================================

  const drawerContent = (
    <Box>
      {/* Logo/Header */}
      <Toolbar>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <AccountBalanceIcon color="primary" sx={{ fontSize: 32 }} />
          <Box>
            <Typography variant="h6" noWrap fontWeight={600}>
              Koperasi ERP
            </Typography>
            <Typography variant="caption" color="text.secondary">
              Sistem Koperasi
            </Typography>
          </Box>
        </Box>
      </Toolbar>

      <Divider />

      {/* Navigation Menu */}
      <List sx={{ px: 1, py: 2 }}>
        {visibleMenuItems.map((item) => {
          const isActive = pathname === item.path || pathname?.startsWith(`${item.path}/`);

          return (
            <ListItem key={item.path} disablePadding sx={{ mb: 0.5 }}>
              <ListItemButton
                onClick={() => handleNavigate(item.path)}
                selected={isActive}
                sx={{
                  borderRadius: 2,
                  '&.Mui-selected': {
                    backgroundColor: 'primary.main',
                    color: 'white',
                    '&:hover': {
                      backgroundColor: 'primary.dark',
                    },
                    '& .MuiListItemIcon-root': {
                      color: 'white',
                    },
                  },
                }}
              >
                <ListItemIcon
                  sx={{
                    color: isActive ? 'white' : 'action.active',
                    minWidth: 40,
                  }}
                >
                  {item.icon}
                </ListItemIcon>
                <ListItemText
                  primary={item.label}
                  primaryTypographyProps={{
                    fontSize: '0.95rem',
                    fontWeight: isActive ? 600 : 400,
                  }}
                />
              </ListItemButton>
            </ListItem>
          );
        })}
      </List>
    </Box>
  );

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <Box component="nav" sx={{ width: { md: DRAWER_WIDTH }, flexShrink: { md: 0 } }}>
      {/* Mobile Drawer */}
      <Drawer
        variant="temporary"
        open={mobileOpen}
        onClose={onDrawerToggle}
        ModalProps={{
          keepMounted: true, // Better mobile performance
        }}
        sx={{
          display: { xs: 'block', md: 'none' },
          '& .MuiDrawer-paper': {
            boxSizing: 'border-box',
            width: DRAWER_WIDTH,
          },
        }}
      >
        {drawerContent}
      </Drawer>

      {/* Desktop Drawer */}
      <Drawer
        variant="permanent"
        sx={{
          display: { xs: 'none', md: 'block' },
          '& .MuiDrawer-paper': {
            boxSizing: 'border-box',
            width: DRAWER_WIDTH,
          },
        }}
        open
      >
        {drawerContent}
      </Drawer>
    </Box>
  );
}

export { DRAWER_WIDTH };
