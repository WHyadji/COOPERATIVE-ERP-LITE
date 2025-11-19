// ============================================================================
// POS Main Page - Point of Sale Interface
// Complete POS system with product selection, cart, and checkout
// ============================================================================

"use client";

import React, { useState } from "react";
import {
  Box,
  Grid,
  Paper,
  Typography,
  Button,
  Divider,
  Alert,
} from "@mui/material";
import {
  ShoppingCart as CartIcon,
  History as HistoryIcon,
} from "@mui/icons-material";
import { useRouter } from "next/navigation";
import { useToast } from "@/lib/context/ToastContext";
import ProductSearch from "@/components/pos/ProductSearch";
import ProductGrid from "@/components/pos/ProductGrid";
import ShoppingCart from "@/components/pos/ShoppingCart";
import MemberLookup from "@/components/pos/MemberLookup";
import CheckoutModal from "@/components/pos/CheckoutModal";
import ReceiptDialog from "@/components/pos/ReceiptDialog";
import type { CartItem, Produk, Member, PenjualanResponse } from "@/types";

// ============================================================================
// POS Main Page Component
// ============================================================================

export default function POSPage() {
  const router = useRouter();
  const { showSuccess, showError, showInfo } = useToast();

  // State
  const [cart, setCart] = useState<CartItem[]>([]);
  const [selectedMember, setSelectedMember] = useState<Member | null>(null);
  const [checkoutOpen, setCheckoutOpen] = useState(false);
  const [receiptOpen, setReceiptOpen] = useState(false);
  const [lastSale, setLastSale] = useState<PenjualanResponse | null>(null);

  // Add product to cart
  const handleAddToCart = (product: Produk) => {
    // Check if product already in cart
    const existingItemIndex = cart.findIndex(
      (item) => item.product.id === product.id
    );

    if (existingItemIndex >= 0) {
      // Update quantity if product already in cart
      const existingItem = cart[existingItemIndex];

      if (existingItem.quantity >= product.stok) {
        showError(`Stok produk "${product.namaProduk}" tidak mencukupi`);
        return;
      }

      const updatedCart = [...cart];
      updatedCart[existingItemIndex] = {
        ...existingItem,
        quantity: existingItem.quantity + 1,
        subtotal: (existingItem.quantity + 1) * product.harga,
      };

      setCart(updatedCart);
      showInfo(`Jumlah "${product.namaProduk}" ditambahkan ke keranjang`);
    } else {
      // Add new item to cart
      const newItem: CartItem = {
        product,
        quantity: 1,
        subtotal: product.harga,
      };

      setCart([...cart, newItem]);
      showSuccess(`"${product.namaProduk}" ditambahkan ke keranjang`);
    }
  };

  // Update quantity in cart
  const handleQuantityChange = (productId: string, quantity: number) => {
    const updatedCart = cart.map((item) => {
      if (item.product.id === productId) {
        return {
          ...item,
          quantity,
          subtotal: quantity * item.product.harga,
        };
      }
      return item;
    });

    setCart(updatedCart);
  };

  // Remove item from cart
  const handleRemoveItem = (productId: string) => {
    const item = cart.find((i) => i.product.id === productId);
    setCart(cart.filter((i) => i.product.id !== productId));

    if (item) {
      showInfo(`"${item.product.namaProduk}" dihapus dari keranjang`);
    }
  };

  // Clear cart
  const handleClearCart = () => {
    setCart([]);
    showInfo("Keranjang dikosongkan");
  };

  // Open checkout
  const handleCheckout = () => {
    if (cart.length === 0) {
      showError("Keranjang masih kosong");
      return;
    }

    setCheckoutOpen(true);
  };

  // Handle successful sale
  const handleSaleSuccess = (sale: PenjualanResponse) => {
    setLastSale(sale);
    setCart([]);
    setSelectedMember(null);
    setCheckoutOpen(false);
    setReceiptOpen(true);
    showSuccess(`Transaksi ${sale.nomorPenjualan} berhasil diproses`);
  };

  // Start new sale
  const handleNewSale = () => {
    setLastSale(null);
    setReceiptOpen(false);
    setCart([]);
    setSelectedMember(null);
  };

  // Navigate to history
  const handleViewHistory = () => {
    router.push("/pos/riwayat");
  };

  return (
    <Box>
      {/* Header */}
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <Typography variant="h4" fontWeight={700}>
          Kasir / POS
        </Typography>
        <Button
          variant="outlined"
          startIcon={<HistoryIcon />}
          onClick={handleViewHistory}
        >
          Riwayat Penjualan
        </Button>
      </Box>

      {/* Main Content */}
      <Grid container spacing={3}>
        {/* Left Column - Product Selection */}
        <Grid item xs={12} md={7}>
          <Paper sx={{ p: 3, height: "100%" }}>
            {/* Member Lookup */}
            <Box sx={{ mb: 3 }}>
              <Typography variant="h6" gutterBottom fontWeight={600}>
                Informasi Anggota
              </Typography>
              <MemberLookup
                selectedMember={selectedMember}
                onMemberSelect={setSelectedMember}
              />
            </Box>

            <Divider sx={{ my: 3 }} />

            {/* Product Search */}
            <Box sx={{ mb: 3 }}>
              <Typography variant="h6" gutterBottom fontWeight={600}>
                Cari Produk
              </Typography>
              <ProductSearch onProductSelect={handleAddToCart} />
            </Box>

            {/* Product Grid */}
            <Box>
              <Typography variant="h6" gutterBottom fontWeight={600}>
                Pilih Produk
              </Typography>
              <ProductGrid onProductSelect={handleAddToCart} />
            </Box>
          </Paper>
        </Grid>

        {/* Right Column - Shopping Cart */}
        <Grid item xs={12} md={5}>
          <Box sx={{ position: "sticky", top: 80 }}>
            {/* Shopping Cart */}
            <ShoppingCart
              items={cart}
              onQuantityChange={handleQuantityChange}
              onRemoveItem={handleRemoveItem}
              onClearCart={handleClearCart}
            />

            {/* Checkout Button */}
            <Paper sx={{ p: 2, mt: 2 }}>
              <Button
                fullWidth
                variant="contained"
                size="large"
                startIcon={<CartIcon />}
                onClick={handleCheckout}
                disabled={cart.length === 0}
                sx={{ py: 1.5 }}
              >
                Checkout
              </Button>

              {cart.length === 0 && (
                <Alert severity="info" sx={{ mt: 2 }}>
                  Tambahkan produk untuk memulai transaksi
                </Alert>
              )}
            </Paper>
          </Box>
        </Grid>
      </Grid>

      {/* Checkout Modal */}
      <CheckoutModal
        open={checkoutOpen}
        onClose={() => setCheckoutOpen(false)}
        items={cart}
        member={selectedMember}
        onSuccess={handleSaleSuccess}
      />

      {/* Receipt Dialog */}
      <ReceiptDialog
        open={receiptOpen}
        onClose={() => setReceiptOpen(false)}
        sale={lastSale}
        onNewSale={handleNewSale}
      />
    </Box>
  );
}
