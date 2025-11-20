# Manual Testing Guide - Member Portal

This guide provides step-by-step instructions for manually testing all features of the Member Portal.

## Prerequisites

Before starting manual testing, ensure:

1. **Backend is running**:
   ```bash
   cd backend
   go run cmd/api/main.go
   ```
   Should see: `Server running on :8080`

2. **Frontend is running**:
   ```bash
   cd frontend
   npm run dev
   ```
   Should see: `Ready on http://localhost:3000`

3. **Test data is seeded**:
   ```bash
   cd backend
   go run cmd/seed-test-data/main.go
   ```
   Should see: `✓ Test data seeding completed successfully!`

4. **Test credentials ready**:
   - Nomor Anggota: `A001`
   - PIN: `123456`

## Testing Checklist

Use this checklist to track your testing progress:

### 1. Login Flow ✓
- [ ] Login page displays correctly
- [ ] Form validation works
- [ ] Invalid credentials show error
- [ ] Valid credentials redirect to dashboard
- [ ] PIN visibility toggle works

### 2. Dashboard ✓
- [ ] Welcome message displays
- [ ] Balance cards show correct amounts
- [ ] Recent transactions table loads
- [ ] Navigation buttons work

### 3. Balance Page ✓
- [ ] Member information displays
- [ ] All balance types shown
- [ ] Educational information visible
- [ ] Currency formatting correct

### 4. Transaction History ✓
- [ ] Transactions table loads
- [ ] Filtering by type works
- [ ] Date filtering works
- [ ] Reset button works
- [ ] Transaction count accurate

### 5. Profile Page ✓
- [ ] Profile information displays
- [ ] Edit mode enables correctly
- [ ] Save changes works (if implemented)
- [ ] Cancel button works

### 6. Navigation ✓
- [ ] Sidebar navigation works
- [ ] Mobile menu works
- [ ] Logout works
- [ ] URL routing correct

### 7. Responsive Design ✓
- [ ] Mobile view (375px)
- [ ] Tablet view (768px)
- [ ] Desktop view (1920px)

## Detailed Test Cases

---

## 1. Login Page Testing

### Test Case 1.1: Page Display
**Steps:**
1. Open browser to `http://localhost:3000/portal/login`
2. Observe the page layout

**Expected Results:**
- ✓ "Portal Anggota" header visible
- ✓ "Masuk ke portal anggota koperasi" subtitle visible
- ✓ "Nomor Anggota" input field visible
- ✓ "PIN (6 digit)" input field visible
- ✓ "Masuk" button visible
- ✓ Gradient background displays nicely

**Pass/Fail:** ___________

---

### Test Case 1.2: Form Validation - Empty Fields
**Steps:**
1. Click "Masuk" button without entering anything

**Expected Results:**
- ✓ "Nomor anggota harus diisi" error message appears
- ✓ "PIN harus 6 digit" error message appears
- ✓ No API call is made

**Pass/Fail:** ___________

---

### Test Case 1.3: Form Validation - Invalid PIN Format
**Steps:**
1. Enter "A001" in Nomor Anggota field
2. Enter "12345" (only 5 digits) in PIN field
3. Click "Masuk" button

**Expected Results:**
- ✓ "PIN harus 6 digit" error message appears
- ✓ No API call is made

**Pass/Fail:** ___________

---

### Test Case 1.4: Form Validation - Non-numeric PIN
**Steps:**
1. Enter "A001" in Nomor Anggota field
2. Enter "abcdef" in PIN field
3. Click "Masuk" button

**Expected Results:**
- ✓ "PIN harus berupa angka" error message appears
- ✓ No API call is made

**Pass/Fail:** ___________

---

### Test Case 1.5: Login with Invalid Credentials
**Steps:**
1. Enter "A001" in Nomor Anggota field
2. Enter "999999" (wrong PIN) in PIN field
3. Click "Masuk" button
4. Wait for response

**Expected Results:**
- ✓ Loading indicator appears briefly
- ✓ Error message "Login gagal" or similar appears
- ✓ User remains on login page

**Pass/Fail:** ___________

---

### Test Case 1.6: Login with Valid Credentials
**Steps:**
1. Enter "A001" in Nomor Anggota field
2. Enter "123456" in PIN field
3. Click "Masuk" button
4. Wait for response

**Expected Results:**
- ✓ Loading indicator appears briefly
- ✓ User is redirected to `/portal`
- ✓ Dashboard page loads
- ✓ Token is stored (check browser DevTools > Application > Cookies)

**Pass/Fail:** ___________

---

### Test Case 1.7: PIN Visibility Toggle
**Steps:**
1. Enter "123456" in PIN field
2. Observe that characters are hidden (••••••)
3. Click the eye icon button
4. Observe that PIN becomes visible
5. Click the eye icon again

**Expected Results:**
- ✓ Initially PIN is hidden (type="password")
- ✓ After clicking, PIN is visible (type="text")
- ✓ After clicking again, PIN is hidden again

**Pass/Fail:** ___________

---

## 2. Dashboard Testing

**Prerequisites:** Login with A001 / 123456

### Test Case 2.1: Welcome Message
**Steps:**
1. Observe the top of the dashboard page

**Expected Results:**
- ✓ "Selamat Datang, Test Member Portal" is visible
- ✓ Member name matches the logged-in user

**Pass/Fail:** ___________

---

### Test Case 2.2: Balance Cards Display
**Steps:**
1. Observe the balance cards section

**Expected Results:**
- ✓ "Total Simpanan" card shows Rp 4,000,000
- ✓ "Simpanan Pokok" card shows Rp 1,000,000
- ✓ "Simpanan Wajib" card shows Rp 2,500,000
- ✓ "Simpanan Sukarela" card shows Rp 500,000
- ✓ All amounts are formatted as Indonesian Rupiah (Rp X,XXX,XXX)
- ✓ Cards have gradient backgrounds
- ✓ Icons are visible

**Pass/Fail:** ___________

---

### Test Case 2.3: Recent Transactions Table
**Steps:**
1. Scroll to "Transaksi Terbaru" section
2. Observe the transactions table

**Expected Results:**
- ✓ Table shows recent transactions (5 most recent)
- ✓ Columns: Tanggal, Tipe Simpanan, Keterangan, Jumlah
- ✓ Dates are formatted in Indonesian (e.g., "15 Januari 2024")
- ✓ Type badges are color-coded (Pokok: blue, Wajib: green, Sukarela: purple)
- ✓ Amounts are formatted as Rupiah
- ✓ "Lihat Semua" button is visible

**Pass/Fail:** ___________

---

### Test Case 2.4: Navigation to Balance Page
**Steps:**
1. Click "Lihat Saldo Detail" button (or similar)
2. Wait for page load

**Expected Results:**
- ✓ URL changes to `/portal/balance`
- ✓ Balance detail page loads

**Pass/Fail:** ___________

---

### Test Case 2.5: Navigation to Transactions Page
**Steps:**
1. From dashboard, click "Riwayat Transaksi" in sidebar or "Lihat Semua" button
2. Wait for page load

**Expected Results:**
- ✓ URL changes to `/portal/transactions`
- ✓ Transaction history page loads

**Pass/Fail:** ___________

---

## 3. Balance Page Testing

**Prerequisites:** Navigate to `/portal/balance` while logged in

### Test Case 3.1: Member Information Display
**Steps:**
1. Observe the member information section

**Expected Results:**
- ✓ "Nomor Anggota: A001" is visible
- ✓ "Nama Anggota: Test Member Portal" is visible
- ✓ Information is clearly formatted

**Pass/Fail:** ___________

---

### Test Case 3.2: Total Balance Card
**Steps:**
1. Observe the large total balance card

**Expected Results:**
- ✓ "Total Simpanan" label visible
- ✓ Shows Rp 4,000,000
- ✓ Has gradient background
- ✓ Larger and more prominent than other cards

**Pass/Fail:** ___________

---

### Test Case 3.3: Individual Balance Cards
**Steps:**
1. Observe the three balance type cards

**Expected Results:**
- ✓ "Simpanan Pokok" card shows Rp 1,000,000
- ✓ "Simpanan Wajib" card shows Rp 2,500,000
- ✓ "Simpanan Sukarela" card shows Rp 500,000
- ✓ Each card has an icon
- ✓ Cards are laid out in a grid

**Pass/Fail:** ___________

---

### Test Case 3.4: Educational Information
**Steps:**
1. Scroll to the information sections
2. Read the content

**Expected Results:**
- ✓ "Tentang Simpanan Pokok:" section visible with explanation
- ✓ "Tentang Simpanan Wajib:" section visible with explanation
- ✓ "Tentang Simpanan Sukarela:" section visible with explanation
- ✓ "Informasi Penting" section visible
- ✓ Text is readable and informative

**Pass/Fail:** ___________

---

## 4. Transaction History Testing

**Prerequisites:** Navigate to `/portal/transactions` while logged in

### Test Case 4.1: Page Header
**Steps:**
1. Observe the page header

**Expected Results:**
- ✓ "Riwayat Transaksi" title visible
- ✓ "Semua transaksi simpanan Anda" subtitle visible

**Pass/Fail:** ___________

---

### Test Case 4.2: Filter Controls Display
**Steps:**
1. Observe the filter section

**Expected Results:**
- ✓ "Tipe Simpanan" dropdown visible (default: "Semua Tipe")
- ✓ "Tanggal Mulai" date picker visible
- ✓ "Tanggal Akhir" date picker visible
- ✓ "Terapkan" button visible
- ✓ "Reset" button visible

**Pass/Fail:** ___________

---

### Test Case 4.3: Transaction Summary
**Steps:**
1. Observe the summary card

**Expected Results:**
- ✓ "Total Transaksi: 8 Transaksi" displayed
- ✓ "Total Setoran" shows sum of all transactions (Rp 4,000,000)

**Pass/Fail:** ___________

---

### Test Case 4.4: Transactions Table Display
**Steps:**
1. Observe the transactions table

**Expected Results:**
- ✓ Headers: No. Referensi, Tanggal, Tipe Simpanan, Keterangan, Jumlah
- ✓ All 8 transactions are visible
- ✓ Transactions ordered by date (newest first or oldest first)
- ✓ Reference numbers are in monospace font
- ✓ Dates formatted correctly
- ✓ Type badges color-coded
- ✓ Amounts right-aligned and formatted as Rupiah

**Pass/Fail:** ___________

---

### Test Case 4.5: Filter by Transaction Type - Simpanan Pokok
**Steps:**
1. Click "Tipe Simpanan" dropdown
2. Select "Simpanan Pokok"
3. Click "Terapkan" button
4. Wait for results

**Expected Results:**
- ✓ Table updates to show only Simpanan Pokok transactions
- ✓ Should show 1 transaction (SP-2024-001)
- ✓ Summary updates: "Total Transaksi: 1 Transaksi"
- ✓ Total amount: Rp 1,000,000

**Pass/Fail:** ___________

---

### Test Case 4.6: Filter by Transaction Type - Simpanan Wajib
**Steps:**
1. Click "Tipe Simpanan" dropdown
2. Select "Simpanan Wajib"
3. Click "Terapkan" button
4. Wait for results

**Expected Results:**
- ✓ Table shows only Simpanan Wajib transactions
- ✓ Should show 5 transactions (SW-2024-001 to SW-2024-005)
- ✓ Summary: "Total Transaksi: 5 Transaksi"
- ✓ Total amount: Rp 2,500,000

**Pass/Fail:** ___________

---

### Test Case 4.7: Filter by Date Range
**Steps:**
1. Reset filters first (click "Reset")
2. Set "Tanggal Mulai" to "2024-02-01"
3. Set "Tanggal Akhir" to "2024-03-31"
4. Click "Terapkan" button
5. Wait for results

**Expected Results:**
- ✓ Table shows only transactions between Feb 1 and Mar 31
- ✓ Should show 4 transactions (SW-2024-002, SW-2024-003, SS-2024-001, SS-2024-002)
- ✓ Summary updates accordingly

**Pass/Fail:** ___________

---

### Test Case 4.8: Reset Filters
**Steps:**
1. Apply some filters (e.g., Simpanan Pokok only)
2. Click "Reset" button
3. Wait for results

**Expected Results:**
- ✓ "Tipe Simpanan" dropdown resets to "Semua Tipe"
- ✓ Date fields are cleared
- ✓ Table shows all 8 transactions again
- ✓ Summary shows total count and amount

**Pass/Fail:** ___________

---

### Test Case 4.9: Combined Filters
**Steps:**
1. Set "Tipe Simpanan" to "Simpanan Wajib"
2. Set "Tanggal Mulai" to "2024-03-01"
3. Set "Tanggal Akhir" to "2024-05-31"
4. Click "Terapkan"

**Expected Results:**
- ✓ Shows only Simpanan Wajib transactions from March to May
- ✓ Should show 3 transactions (SW-2024-003, SW-2024-004, SW-2024-005)

**Pass/Fail:** ___________

---

## 5. Profile Page Testing

**Prerequisites:** Navigate to `/portal/profile` while logged in

### Test Case 5.1: Profile Header
**Steps:**
1. Observe the profile header

**Expected Results:**
- ✓ "Profil Saya" title visible
- ✓ Avatar/icon displayed
- ✓ "Test Member Portal" name displayed
- ✓ "Nomor Anggota: A001" displayed
- ✓ Status badge (e.g., "Aktif") displayed

**Pass/Fail:** ___________

---

### Test Case 5.2: Personal Information Section
**Steps:**
1. Observe the "Informasi Pribadi" section

**Expected Results:**
- ✓ Section title "Informasi Pribadi" visible
- ✓ NIK field displayed (read-only)
- ✓ Jenis Kelamin field displayed (read-only)
- ✓ Tempat Lahir field displayed (read-only)
- ✓ Tanggal Lahir field displayed (read-only)
- ✓ Pekerjaan field displayed (read-only)
- ✓ All fields are disabled/read-only

**Pass/Fail:** ___________

---

### Test Case 5.3: Contact Information Section
**Steps:**
1. Observe the "Informasi Kontak" section

**Expected Results:**
- ✓ Section title "Informasi Kontak" visible
- ✓ Nomor Telepon field displayed
- ✓ Email field displayed
- ✓ Fields are initially disabled
- ✓ Data matches test member (081234567890, test.member@email.com)

**Pass/Fail:** ___________

---

### Test Case 5.4: Address Section
**Steps:**
1. Observe the "Alamat" section

**Expected Results:**
- ✓ Section title "Alamat" visible
- ✓ Alamat Lengkap field displayed
- ✓ RT field displayed
- ✓ RW field displayed
- ✓ Kelurahan field displayed
- ✓ Kecamatan field displayed
- ✓ Kota/Kabupaten field displayed
- ✓ Provinsi field displayed
- ✓ Kode Pos field displayed
- ✓ All fields initially disabled

**Pass/Fail:** ___________

---

### Test Case 5.5: Enable Edit Mode
**Steps:**
1. Click "Edit Profil" button
2. Observe the form changes

**Expected Results:**
- ✓ Contact information fields become enabled (phone, email)
- ✓ Address fields become enabled
- ✓ Personal information remains disabled (NIK, gender, birthdate)
- ✓ "Edit Profil" button is replaced by "Simpan Perubahan" and "Batal" buttons

**Pass/Fail:** ___________

---

### Test Case 5.6: Cancel Edit Mode
**Steps:**
1. Click "Edit Profil" button to enter edit mode
2. Modify some fields (e.g., change phone number)
3. Click "Batal" button

**Expected Results:**
- ✓ All fields return to disabled state
- ✓ Changes are discarded (original values restored)
- ✓ "Simpan Perubahan" and "Batal" buttons disappear
- ✓ "Edit Profil" button reappears

**Pass/Fail:** ___________

---

### Test Case 5.7: Save Changes (if implemented)
**Steps:**
1. Click "Edit Profil" button
2. Change phone number to "081999999999"
3. Change email to "new.email@test.com"
4. Click "Simpan Perubahan" button

**Expected Results:**
- ✓ Loading indicator appears
- ✓ Success message displayed
- ✓ Form returns to disabled state
- ✓ New values are saved and displayed
- ✓ Refresh page shows updated values

**Pass/Fail:** ___________ (If implemented)

---

## 6. Navigation Testing

**Prerequisites:** Logged in as A001

### Test Case 6.1: Sidebar Navigation - Dashboard
**Steps:**
1. From any page, click "Dashboard" in sidebar
2. Wait for page load

**Expected Results:**
- ✓ URL changes to `/portal`
- ✓ Dashboard page loads
- ✓ "Dashboard" menu item is highlighted/active

**Pass/Fail:** ___________

---

### Test Case 6.2: Sidebar Navigation - Saldo Simpanan
**Steps:**
1. From any page, click "Saldo Simpanan" in sidebar
2. Wait for page load

**Expected Results:**
- ✓ URL changes to `/portal/balance`
- ✓ Balance page loads
- ✓ "Saldo Simpanan" menu item is highlighted/active

**Pass/Fail:** ___________

---

### Test Case 6.3: Sidebar Navigation - Riwayat Transaksi
**Steps:**
1. From any page, click "Riwayat Transaksi" in sidebar
2. Wait for page load

**Expected Results:**
- ✓ URL changes to `/portal/transactions`
- ✓ Transactions page loads
- ✓ "Riwayat Transaksi" menu item is highlighted/active

**Pass/Fail:** ___________

---

### Test Case 6.4: Sidebar Navigation - Profil Saya
**Steps:**
1. From any page, click "Profil Saya" in sidebar
2. Wait for page load

**Expected Results:**
- ✓ URL changes to `/portal/profile`
- ✓ Profile page loads
- ✓ "Profil Saya" menu item is highlighted/active

**Pass/Fail:** ___________

---

### Test Case 6.5: User Menu and Logout
**Steps:**
1. From any page, click on the user avatar/menu in top-right
2. Observe the dropdown menu
3. Click "Keluar" (Logout)
4. Wait for redirect

**Expected Results:**
- ✓ Dropdown menu appears showing user name and "Keluar" option
- ✓ After clicking "Keluar", user is redirected to `/portal/login`
- ✓ Token is removed (check DevTools > Application > Cookies)
- ✓ Cannot access protected pages without login

**Pass/Fail:** ___________

---

### Test Case 6.6: Direct URL Access Without Login
**Steps:**
1. Logout if logged in
2. In browser address bar, navigate directly to `http://localhost:3000/portal`

**Expected Results:**
- ✓ User is redirected to `/portal/login`
- ✓ Protected route middleware works

**Pass/Fail:** ___________

---

### Test Case 6.7: Browser Back Button
**Steps:**
1. Navigate: Dashboard → Balance → Transactions
2. Click browser back button twice

**Expected Results:**
- ✓ First back: Returns to Balance page
- ✓ Second back: Returns to Dashboard
- ✓ Pages load correctly from browser cache

**Pass/Fail:** ___________

---

## 7. Responsive Design Testing

### Test Case 7.1: Mobile View (375px × 667px)
**Steps:**
1. Open browser DevTools (F12)
2. Toggle device toolbar (Ctrl+Shift+M)
3. Select "iPhone SE" or set viewport to 375 × 667
4. Navigate through all pages

**Expected Results:**
- ✓ Sidebar becomes a drawer (hidden by default)
- ✓ Hamburger menu icon appears in top-left
- ✓ Balance cards stack vertically
- ✓ Transaction table is scrollable horizontally
- ✓ Buttons are touch-friendly (adequate size)
- ✓ Text is readable without zooming
- ✓ No horizontal scrolling on page (except tables)

**Pass/Fail:** ___________

---

### Test Case 7.2: Tablet View (768px × 1024px)
**Steps:**
1. Set viewport to 768 × 1024 (iPad)
2. Navigate through all pages

**Expected Results:**
- ✓ Sidebar may be permanent or drawer (depending on design)
- ✓ Balance cards display 2-up grid
- ✓ Transaction table fully visible
- ✓ Forms use appropriate width (not full width)
- ✓ Layout is balanced and readable

**Pass/Fail:** ___________

---

### Test Case 7.3: Desktop View (1920px × 1080px)
**Steps:**
1. Set viewport to 1920 × 1080
2. Navigate through all pages

**Expected Results:**
- ✓ Sidebar is permanent and always visible
- ✓ Balance cards display in 3 or 4-column grid
- ✓ Content is centered with max-width (not stretched full width)
- ✓ Typography is comfortable for desktop reading
- ✓ Adequate whitespace

**Pass/Fail:** ___________

---

### Test Case 7.4: Mobile Drawer Menu
**Steps:**
1. Set viewport to mobile (375px)
2. Observe that sidebar is hidden
3. Click hamburger menu icon
4. Observe drawer opens
5. Click a menu item
6. Observe drawer closes

**Expected Results:**
- ✓ Hamburger icon visible and clickable
- ✓ Drawer slides in from left
- ✓ Overlay appears behind drawer
- ✓ Menu items visible and clickable
- ✓ Clicking item closes drawer
- ✓ Clicking outside drawer closes it

**Pass/Fail:** ___________

---

## 8. Error Handling Testing

### Test Case 8.1: Network Error Simulation
**Steps:**
1. Login successfully
2. Stop the backend server (Ctrl+C in backend terminal)
3. Try to navigate to Balance page or refresh

**Expected Results:**
- ✓ Loading indicator appears
- ✓ Error message displayed (e.g., "Gagal memuat data")
- ✓ User-friendly error message (not technical stack trace)
- ✓ No blank white screen

**Pass/Fail:** ___________

---

### Test Case 8.2: Token Expiration
**Steps:**
1. Login successfully
2. Manually delete token from cookies (DevTools > Application > Cookies)
3. Try to navigate to another page

**Expected Results:**
- ✓ User is redirected to login page
- ✓ Or shows "Session expired" message

**Pass/Fail:** ___________

---

## 9. Data Accuracy Testing

### Test Case 9.1: Balance Calculation Accuracy
**Steps:**
1. Login and view Balance page
2. Manually calculate: Pokok + Wajib + Sukarela
3. Compare with displayed total

**Expected:**
- Simpanan Pokok: Rp 1,000,000
- Simpanan Wajib: Rp 2,500,000
- Simpanan Sukarela: Rp 500,000
- **Total: Rp 4,000,000**

**Actual Total:** ___________

**Pass/Fail:** ___________

---

### Test Case 9.2: Transaction List Accuracy
**Steps:**
1. View Transaction History page
2. Count displayed transactions
3. Verify transaction details match test data

**Expected:**
- Total transactions: 8
- Pokok: 1 transaction
- Wajib: 5 transactions
- Sukarela: 2 transactions

**Actual Counts:** ___________

**Pass/Fail:** ___________

---

### Test Case 9.3: Transaction Sum Accuracy
**Steps:**
1. View Transaction History page
2. Manually sum all transaction amounts
3. Compare with displayed "Total Setoran"

**Expected:**
- 1,000,000 + 500,000×5 + 200,000 + 300,000 = **Rp 4,000,000**

**Actual Total:** ___________

**Pass/Fail:** ___________

---

## 10. Performance Testing (Basic)

### Test Case 10.1: Page Load Time
**Steps:**
1. Open browser DevTools > Network tab
2. Clear cache (Ctrl+Shift+Del)
3. Navigate to login page
4. Login and observe dashboard load time
5. Note the "Load" time in Network tab

**Expected:**
- ✓ Dashboard loads in < 3 seconds on good connection
- ✓ No excessive API calls (only necessary endpoints)

**Actual Load Time:** ___________ seconds

**Pass/Fail:** ___________

---

### Test Case 10.2: Data Refresh
**Steps:**
1. On Transaction History page, apply a filter
2. Observe loading time
3. Apply different filter
4. Observe loading time

**Expected:**
- ✓ Filter results appear in < 1 second
- ✓ Smooth transition without page refresh

**Pass/Fail:** ___________

---

## Test Summary

**Total Test Cases:** 50+

**Passed:** ___________

**Failed:** ___________

**Skipped:** ___________

**Pass Rate:** ___________%

---

## Issues Found

Use this section to document any bugs or issues discovered during testing:

### Issue #1
- **Test Case:** ___________
- **Severity:** Critical / Major / Minor
- **Description:** ___________
- **Steps to Reproduce:** ___________
- **Expected:** ___________
- **Actual:** ___________
- **Screenshot/Video:** ___________

### Issue #2
- **Test Case:** ___________
- **Severity:** Critical / Major / Minor
- **Description:** ___________
- **Steps to Reproduce:** ___________
- **Expected:** ___________
- **Actual:** ___________
- **Screenshot/Video:** ___________

*(Add more as needed)*

---

## Recommendations

Based on testing results, list any recommendations for improvements:

1. ___________
2. ___________
3. ___________

---

## Sign-off

**Tester Name:** ___________

**Date:** ___________

**Overall Assessment:** ___________

**Ready for Production:** Yes / No / With Fixes

**Notes:** ___________

---

## Appendix: Testing Tools

### Browser DevTools Shortcuts
- **F12**: Open DevTools
- **Ctrl+Shift+M**: Toggle device toolbar (responsive mode)
- **Ctrl+Shift+Del**: Clear cache and cookies
- **Ctrl+Shift+I**: Inspect element

### Useful Browser Extensions
- **React Developer Tools**: Inspect React components
- **Redux DevTools**: Monitor state changes
- **Lighthouse**: Performance and accessibility audits

### Network Tab Filters
- **XHR/Fetch**: See only API calls
- **Disable cache**: Force fresh loads
- **Throttling**: Simulate slow connection

---

**End of Manual Testing Guide**
