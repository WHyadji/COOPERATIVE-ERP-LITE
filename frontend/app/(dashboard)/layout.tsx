// ============================================================================
// Dashboard Layout - Protected Layout with Sidebar and Header
// Material-UI responsive layout for authenticated users
// ============================================================================

'use client';

import React, { useState } from 'react';
import { Box, Toolbar } from '@mui/material';
import Sidebar, { DRAWER_WIDTH } from '@/components/layout/Sidebar';
import Header from '@/components/layout/Header';
import { ProtectedRoute } from '@/lib/context/AuthContext';
import { ToastProvider } from '@/lib/context/ToastContext';

// ============================================================================
// Dashboard Layout Component
// ============================================================================

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [mobileOpen, setMobileOpen] = useState(false);

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  // ============================================================================
  // Render
  // ============================================================================

  return (
    <ProtectedRoute>
      <ToastProvider>
        <Box sx={{ display: 'flex', minHeight: '100vh' }}>
          {/* Sidebar Navigation */}
          <Sidebar mobileOpen={mobileOpen} onDrawerToggle={handleDrawerToggle} />

          {/* Main Content Area */}
          <Box
            component="main"
            sx={{
              flexGrow: 1,
              width: { md: `calc(100% - ${DRAWER_WIDTH}px)` },
              minHeight: '100vh',
              backgroundColor: 'background.default',
            }}
          >
            {/* Header */}
            <Header onMenuClick={handleDrawerToggle} />

            {/* Toolbar spacer - creates space for fixed AppBar */}
            <Toolbar />

            {/* Page Content */}
            <Box sx={{ p: 3 }}>
              {children}
            </Box>
          </Box>
        </Box>
      </ToastProvider>
    </ProtectedRoute>
  );
}
