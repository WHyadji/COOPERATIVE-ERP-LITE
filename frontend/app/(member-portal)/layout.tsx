// ============================================================================
// Member Portal Layout
// Layout wrapper for member portal pages with navigation
// ============================================================================

"use client";

import React, { useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import {
  Box,
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Container,
  Avatar,
  Menu,
  MenuItem,
  Divider,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import {
  Menu as MenuIcon,
  Dashboard,
  AccountBalance,
  History,
  Person,
  Logout,
  AccountCircle,
} from "@mui/icons-material";
import { useAuth, ProtectedRoute } from "@/lib/context/AuthContext";

// ============================================================================
// Navigation Items
// ============================================================================

const navigationItems = [
  { label: "Dashboard", icon: <Dashboard />, path: "/portal" },
  {
    label: "Saldo Simpanan",
    icon: <AccountBalance />,
    path: "/portal/balance",
  },
  {
    label: "Riwayat Transaksi",
    icon: <History />,
    path: "/portal/transactions",
  },
  { label: "Profil Saya", icon: <Person />, path: "/portal/profile" },
];

// ============================================================================
// Member Portal Layout Component
// ============================================================================

export default function MemberPortalLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const [mobileOpen, setMobileOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleProfileMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    handleProfileMenuClose();
    logout();
  };

  const handleNavigation = (path: string) => {
    router.push(path);
    if (isMobile) {
      setMobileOpen(false);
    }
  };

  // ============================================================================
  // Drawer Content
  // ============================================================================

  const drawer = (
    <Box sx={{ height: "100%", display: "flex", flexDirection: "column" }}>
      {/* Drawer Header */}
      <Box
        sx={{
          p: 3,
          textAlign: "center",
          bgcolor: "primary.main",
          color: "white",
        }}
      >
        <AccountCircle sx={{ fontSize: 60, mb: 1 }} />
        <Typography variant="h6" fontWeight={600}>
          Portal Anggota
        </Typography>
        {user && (
          <Typography variant="body2" sx={{ mt: 1, opacity: 0.9 }}>
            {user.namaLengkap}
          </Typography>
        )}
      </Box>

      <Divider />

      {/* Navigation Items */}
      <List sx={{ flexGrow: 1, pt: 2 }}>
        {navigationItems.map((item) => (
          <ListItem key={item.path} disablePadding>
            <ListItemButton
              selected={pathname === item.path}
              onClick={() => handleNavigation(item.path)}
              sx={{
                mx: 1,
                borderRadius: 1,
                "&.Mui-selected": {
                  bgcolor: "primary.light",
                  color: "primary.contrastText",
                  "&:hover": {
                    bgcolor: "primary.main",
                  },
                },
              }}
            >
              <ListItemIcon
                sx={{
                  color: pathname === item.path ? "inherit" : "text.secondary",
                }}
              >
                {item.icon}
              </ListItemIcon>
              <ListItemText primary={item.label} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>

      <Divider />

      {/* Logout Button */}
      <List>
        <ListItem disablePadding>
          <ListItemButton
            onClick={handleLogout}
            sx={{ mx: 1, borderRadius: 1 }}
          >
            <ListItemIcon>
              <Logout />
            </ListItemIcon>
            <ListItemText primary="Keluar" />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  );

  // ============================================================================
  // Render
  // ============================================================================

  // Don't show layout on login page
  if (pathname === "/portal/login") {
    return <>{children}</>;
  }

  return (
    <ProtectedRoute requiredRoles={["anggota"]}>
      <Box sx={{ display: "flex", minHeight: "100vh", bgcolor: "grey.50" }}>
        {/* App Bar */}
        <AppBar
          position="fixed"
          sx={{
            zIndex: theme.zIndex.drawer + 1,
            background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
          }}
        >
          <Toolbar>
            <IconButton
              color="inherit"
              edge="start"
              onClick={handleDrawerToggle}
              sx={{ mr: 2, display: { md: "none" } }}
            >
              <MenuIcon />
            </IconButton>

            <Typography
              variant="h6"
              noWrap
              component="div"
              sx={{ flexGrow: 1 }}
            >
              Portal Anggota Koperasi
            </Typography>

            {/* User Profile Menu */}
            <IconButton onClick={handleProfileMenuOpen} sx={{ color: "white" }}>
              <Avatar sx={{ width: 32, height: 32, bgcolor: "primary.dark" }}>
                {user?.namaLengkap?.charAt(0).toUpperCase()}
              </Avatar>
            </IconButton>

            <Menu
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleProfileMenuClose}
              transformOrigin={{ horizontal: "right", vertical: "top" }}
              anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
            >
              <Box sx={{ px: 2, py: 1 }}>
                <Typography variant="subtitle2">{user?.namaLengkap}</Typography>
                <Typography variant="caption" color="text.secondary">
                  {user?.email || "Anggota Koperasi"}
                </Typography>
              </Box>
              <Divider />
              <MenuItem
                onClick={() => {
                  handleProfileMenuClose();
                  router.push("/portal/profile");
                }}
              >
                <ListItemIcon>
                  <Person fontSize="small" />
                </ListItemIcon>
                Profil Saya
              </MenuItem>
              <MenuItem onClick={handleLogout}>
                <ListItemIcon>
                  <Logout fontSize="small" />
                </ListItemIcon>
                Keluar
              </MenuItem>
            </Menu>
          </Toolbar>
        </AppBar>

        {/* Drawer - Desktop */}
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: "none", md: "block" },
            width: 280,
            flexShrink: 0,
            "& .MuiDrawer-paper": {
              width: 280,
              boxSizing: "border-box",
            },
          }}
        >
          <Toolbar />
          {drawer}
        </Drawer>

        {/* Drawer - Mobile */}
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Better mobile performance
          }}
          sx={{
            display: { xs: "block", md: "none" },
            "& .MuiDrawer-paper": {
              width: 280,
              boxSizing: "border-box",
            },
          }}
        >
          <Toolbar />
          {drawer}
        </Drawer>

        {/* Main Content */}
        <Box
          component="main"
          sx={{
            flexGrow: 1,
            width: { md: `calc(100% - 280px)` },
          }}
        >
          <Toolbar />
          <Container maxWidth="lg" sx={{ py: 4 }}>
            {children}
          </Container>
        </Box>
      </Box>
    </ProtectedRoute>
  );
}
