package main

import (
	"cooperative-erp-lite/internal/config"
	"cooperative-erp-lite/internal/models"
	"cooperative-erp-lite/internal/services"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("🌱 Starting database seeding...")

	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("❌ Gagal memuat konfigurasi:", err)
	}

	// 2. Initialize database
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatal("❌ Gagal menghubungkan database:", err)
	}
	defer config.CloseDatabase()

	db := config.GetDB()

	log.Println("✅ Koneksi database berhasil")

	// 3. Clear existing data (optional - hati-hati!)
	log.Println("🗑️  Membersihkan data lama...")
	clearExistingData(db)

	// 4. Seed Koperasi
	log.Println("📋 Seeding Koperasi...")
	koperasi := seedKoperasi(db)

	// 5. Seed Users
	log.Println("👥 Seeding Users...")
	admin, bendahara, kasir := seedUsers(db, koperasi.ID)

	// 6. Seed Chart of Accounts
	log.Println("💰 Seeding Chart of Accounts...")
	seedChartOfAccounts(db, koperasi.ID)

	// 7. Seed Members
	log.Println("🧑‍🤝‍🧑 Seeding Members...")
	members := seedMembers(db, koperasi.ID)

	// 8. Seed Products
	log.Println("📦 Seeding Products...")
	products := seedProducts(db, koperasi.ID)

	// 9. Seed Simpanan (Share Capital)
	log.Println("💵 Seeding Simpanan...")
	transaksiService := services.NewTransaksiService(db)
	simpananService := services.NewSimpananService(db, transaksiService)
	seedSimpanan(simpananService, koperasi.ID, members)

	// 10. Seed Sample Sales
	log.Println("🛒 Seeding Sample Sales...")
	produkService := services.NewProdukService(db)
	penjualanService := services.NewPenjualanService(db, produkService, transaksiService)
	seedSampleSales(penjualanService, koperasi.ID, kasir.ID, members, products)

	log.Println("✅ Seeding completed successfully!")
	log.Println("")
	log.Println("📊 Summary:")
	log.Println("   - Koperasi: 1")
	log.Println(fmt.Sprintf("   - Users: 3 (Admin: %s, Bendahara: %s, Kasir: %s)", admin.NamaPengguna, bendahara.NamaPengguna, kasir.NamaPengguna))
	log.Println("   - Chart of Accounts: 31")
	log.Println(fmt.Sprintf("   - Members: %d", len(members)))
	log.Println(fmt.Sprintf("   - Products: %d", len(products)))
	log.Println("   - Simpanan Transactions: Multiple")
	log.Println("   - Sales Transactions: Multiple")
	log.Println("")
	log.Println("🔑 Default Login Credentials:")
	log.Println("   Admin     - Username: admin     | Password: admin123")
	log.Println("   Bendahara - Username: bendahara | Password: bendahara123")
	log.Println("   Kasir     - Username: kasir     | Password: kasir123")
	log.Println("")
	log.Println("🚀 You can now start the server and test the API!")
}

func clearExistingData(db *config.DB) {
	// Clear in order to respect foreign key constraints
	db.Exec("DELETE FROM baris_transaksi")
	db.Exec("DELETE FROM transaksi")
	db.Exec("DELETE FROM item_penjualan")
	db.Exec("DELETE FROM penjualan")
	db.Exec("DELETE FROM simpanan")
	db.Exec("DELETE FROM produk")
	db.Exec("DELETE FROM anggota")
	db.Exec("DELETE FROM akun")
	db.Exec("DELETE FROM pengguna")
	db.Exec("DELETE FROM koperasi")
}

func seedKoperasi(db *config.DB) *models.Koperasi {
	koperasi := &models.Koperasi{
		NamaKoperasi: "Koperasi Maju Bersama",
		Alamat:       "Jl. Merdeka No. 123, Jakarta Pusat",
		Telepon:      "021-12345678",
		Email:        "info@majubersama.coop",
		Npwp:         "01.234.567.8-901.000",
		TanggalBerdiri: parseDate("2020-01-15"),
		NomorBadanHukum: "123/BH/KWK.11/I/2020",
		NamaKetuaKoperasi: "Budi Santoso",
		StatusAktif:    true,
	}

	if err := db.Create(koperasi).Error; err != nil {
		log.Fatal("❌ Error seeding koperasi:", err)
	}

	log.Printf("   ✓ Koperasi '%s' created", koperasi.NamaKoperasi)
	return koperasi
}

func seedUsers(db *config.DB, idKoperasi uuid.UUID) (*models.Pengguna, *models.Pengguna, *models.Pengguna) {
	// Admin
	admin := &models.Pengguna{
		IDKoperasi:   idKoperasi,
		NamaPengguna: "admin",
		KataSandi:    hashPassword("admin123"),
		NamaLengkap:  "Administrator",
		Email:        "admin@majubersama.coop",
		Peran:        models.PeranAdmin,
		StatusAktif:  true,
	}

	// Bendahara
	bendahara := &models.Pengguna{
		IDKoperasi:   idKoperasi,
		NamaPengguna: "bendahara",
		KataSandi:    hashPassword("bendahara123"),
		NamaLengkap:  "Siti Nurhaliza",
		Email:        "bendahara@majubersama.coop",
		Peran:        models.PeranBendahara,
		StatusAktif:  true,
	}

	// Kasir
	kasir := &models.Pengguna{
		IDKoperasi:   idKoperasi,
		NamaPengguna: "kasir",
		KataSandi:    hashPassword("kasir123"),
		NamaLengkap:  "Ahmad Yani",
		Email:        "kasir@majubersama.coop",
		Peran:        models.PeranKasir,
		StatusAktif:  true,
	}

	if err := db.Create(admin).Error; err != nil {
		log.Fatal("❌ Error seeding admin:", err)
	}
	if err := db.Create(bendahara).Error; err != nil {
		log.Fatal("❌ Error seeding bendahara:", err)
	}
	if err := db.Create(kasir).Error; err != nil {
		log.Fatal("❌ Error seeding kasir:", err)
	}

	log.Printf("   ✓ User 'admin' created")
	log.Printf("   ✓ User 'bendahara' created")
	log.Printf("   ✓ User 'kasir' created")

	return admin, bendahara, kasir
}

func seedChartOfAccounts(db *config.DB, idKoperasi uuid.UUID) {
	akunService := services.NewAkunService(db)

	if err := akunService.InisialisasiCOADefault(idKoperasi); err != nil {
		log.Fatal("❌ Error seeding COA:", err)
	}

	log.Printf("   ✓ 31 Chart of Accounts created")
}

func seedMembers(db *config.DB, idKoperasi uuid.UUID) []*models.Anggota {
	anggotaService := services.NewAnggotaService(db)

	members := []services.BuatAnggotaRequest{
		{
			NamaLengkap:     "Andi Wijaya",
			TempatLahir:     "Jakarta",
			TanggalLahir:    "1985-05-15",
			JenisKelamin:    "L",
			Alamat:          "Jl. Sudirman No. 45, Jakarta",
			Telepon:         "081234567890",
			Email:           "andi.wijaya@email.com",
			Pekerjaan:       "Wiraswasta",
			TanggalBergabung: parseDate("2023-01-10"),
		},
		{
			NamaLengkap:     "Budi Santoso",
			TempatLahir:     "Bandung",
			TanggalLahir:    "1990-08-22",
			JenisKelamin:    "L",
			Alamat:          "Jl. Asia Afrika No. 12, Bandung",
			Telepon:         "082345678901",
			Email:           "budi.santoso@email.com",
			Pekerjaan:       "Pegawai Swasta",
			TanggalBergabung: parseDate("2023-02-15"),
		},
		{
			NamaLengkap:     "Citra Dewi",
			TempatLahir:     "Surabaya",
			TanggalLahir:    "1988-03-30",
			JenisKelamin:    "P",
			Alamat:          "Jl. Pemuda No. 78, Surabaya",
			Telepon:         "083456789012",
			Email:           "citra.dewi@email.com",
			Pekerjaan:       "Guru",
			TanggalBergabung: parseDate("2023-03-20"),
		},
		{
			NamaLengkap:     "Dedi Kurniawan",
			TempatLahir:     "Yogyakarta",
			TanggalLahir:    "1992-11-05",
			JenisKelamin:    "L",
			Alamat:          "Jl. Malioboro No. 99, Yogyakarta",
			Telepon:         "084567890123",
			Email:           "dedi.kurniawan@email.com",
			Pekerjaan:       "Desainer Grafis",
			TanggalBergabung: parseDate("2023-04-10"),
		},
		{
			NamaLengkap:     "Eka Putri",
			TempatLahir:     "Semarang",
			TanggalLahir:    "1987-07-18",
			JenisKelamin:    "P",
			Alamat:          "Jl. Pahlawan No. 56, Semarang",
			Telepon:         "085678901234",
			Email:           "eka.putri@email.com",
			Pekerjaan:       "Dokter",
			TanggalBergabung: parseDate("2023-05-25"),
		},
		{
			NamaLengkap:     "Faisal Rahman",
			TempatLahir:     "Medan",
			TanggalLahir:    "1991-12-08",
			JenisKelamin:    "L",
			Alamat:          "Jl. Sisingamangaraja No. 34, Medan",
			Telepon:         "086789012345",
			Email:           "faisal.rahman@email.com",
			Pekerjaan:       "Pengusaha",
			TanggalBergabung: parseDate("2023-06-12"),
		},
		{
			NamaLengkap:     "Gita Savitri",
			TempatLahir:     "Malang",
			TanggalLahir:    "1989-04-25",
			JenisKelamin:    "P",
			Alamat:          "Jl. Ijen No. 67, Malang",
			Telepon:         "087890123456",
			Email:           "gita.savitri@email.com",
			Pekerjaan:       "Akuntan",
			TanggalBergabung: parseDate("2023-07-08"),
		},
		{
			NamaLengkap:     "Hendra Gunawan",
			TempatLahir:     "Palembang",
			TanggalLahir:    "1986-09-14",
			JenisKelamin:    "L",
			Alamat:          "Jl. Sudirman No. 89, Palembang",
			Telepon:         "088901234567",
			Email:           "hendra.gunawan@email.com",
			Pekerjaan:       "Insinyur",
			TanggalBergabung: parseDate("2023-08-22"),
		},
	}

	var createdMembers []*models.Anggota

	for i, req := range members {
		member, err := anggotaService.BuatAnggota(idKoperasi, &req)
		if err != nil {
			log.Printf("   ⚠️  Warning: Error creating member %d: %v", i+1, err)
			continue
		}
		createdMembers = append(createdMembers, member)
		log.Printf("   ✓ Member '%s' created with number %s", member.NamaLengkap, member.NomorAnggota)
	}

	return createdMembers
}

func seedProducts(db *config.DB, idKoperasi uuid.UUID) []*models.Produk {
	produkService := services.NewProdukService(db)

	products := []services.BuatProdukRequest{
		{
			NamaProduk:   "Beras Premium 5kg",
			Kategori:     "Sembako",
			Barcode:      "8991234567890",
			HargaBeli:    55000,
			HargaJual:    65000,
			Stok:         100,
			StokMinimum:  20,
			Satuan:       "kg",
			Deskripsi:    "Beras premium kualitas terbaik",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Minyak Goreng 2L",
			Kategori:     "Sembako",
			Barcode:      "8991234567891",
			HargaBeli:    28000,
			HargaJual:    32000,
			Stok:         80,
			StokMinimum:  15,
			Satuan:       "botol",
			Deskripsi:    "Minyak goreng kemasan 2 liter",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Gula Pasir 1kg",
			Kategori:     "Sembako",
			Barcode:      "8991234567892",
			HargaBeli:    12000,
			HargaJual:    14000,
			Stok:         120,
			StokMinimum:  25,
			Satuan:       "kg",
			Deskripsi:    "Gula pasir putih kemasan 1kg",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Telur Ayam 1kg",
			Kategori:     "Sembako",
			Barcode:      "8991234567893",
			HargaBeli:    25000,
			HargaJual:    28000,
			Stok:         60,
			StokMinimum:  10,
			Satuan:       "kg",
			Deskripsi:    "Telur ayam segar",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Tepung Terigu 1kg",
			Kategori:     "Sembako",
			Barcode:      "8991234567894",
			HargaBeli:    10000,
			HargaJual:    12000,
			Stok:         90,
			StokMinimum:  20,
			Satuan:       "kg",
			Deskripsi:    "Tepung terigu serbaguna",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Kopi Bubuk 200g",
			Kategori:     "Minuman",
			Barcode:      "8991234567895",
			HargaBeli:    18000,
			HargaJual:    22000,
			Stok:         50,
			StokMinimum:  10,
			Satuan:       "bungkus",
			Deskripsi:    "Kopi bubuk robusta",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Teh Celup 25 sachet",
			Kategori:     "Minuman",
			Barcode:      "8991234567896",
			HargaBeli:    8000,
			HargaJual:    10000,
			Stok:         70,
			StokMinimum:  15,
			Satuan:       "kotak",
			Deskripsi:    "Teh celup isi 25 sachet",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Susu UHT 1L",
			Kategori:     "Minuman",
			Barcode:      "8991234567897",
			HargaBeli:    15000,
			HargaJual:    18000,
			Stok:         40,
			StokMinimum:  10,
			Satuan:       "kotak",
			Deskripsi:    "Susu UHT full cream 1 liter",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Sabun Mandi 85g",
			Kategori:     "Toiletries",
			Barcode:      "8991234567898",
			HargaBeli:    3500,
			HargaJual:    5000,
			Stok:         100,
			StokMinimum:  20,
			Satuan:       "pcs",
			Deskripsi:    "Sabun mandi batang",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Pasta Gigi 150g",
			Kategori:     "Toiletries",
			Barcode:      "8991234567899",
			HargaBeli:    7000,
			HargaJual:    9000,
			Stok:         80,
			StokMinimum:  15,
			Satuan:       "tube",
			Deskripsi:    "Pasta gigi keluarga",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Shampo Sachet 12ml",
			Kategori:     "Toiletries",
			Barcode:      "8991234567800",
			HargaBeli:    1500,
			HargaJual:    2000,
			Stok:         150,
			StokMinimum:  30,
			Satuan:       "sachet",
			Deskripsi:    "Shampo sachet 12ml",
			StatusAktif:  true,
		},
		{
			NamaProduk:   "Detergen 1kg",
			Kategori:     "Toiletries",
			Barcode:      "8991234567801",
			HargaBeli:    15000,
			HargaJual:    18000,
			Stok:         60,
			StokMinimum:  12,
			Satuan:       "kg",
			Deskripsi:    "Detergen bubuk 1kg",
			StatusAktif:  true,
		},
	}

	var createdProducts []*models.Produk

	for i, req := range products {
		product, err := produkService.BuatProduk(idKoperasi, &req)
		if err != nil {
			log.Printf("   ⚠️  Warning: Error creating product %d: %v", i+1, err)
			continue
		}
		createdProducts = append(createdProducts, product)
		log.Printf("   ✓ Product '%s' created", product.NamaProduk)
	}

	return createdProducts
}

func seedSimpanan(simpananService *services.SimpananService, idKoperasi uuid.UUID, members []*models.Anggota) {
	// Seed Simpanan Pokok (one-time deposit)
	for _, member := range members {
		req := &services.CatatSetoranRequest{
			IDAnggota:    member.ID,
			JenisSimpanan: models.JenisSimpananPokok,
			Jumlah:       100000, // Rp 100,000 per member
			TanggalSetor: time.Now(),
			Keterangan:   "Simpanan Pokok - setoran awal",
		}

		_, err := simpananService.CatatSetoran(idKoperasi, req)
		if err != nil {
			log.Printf("   ⚠️  Warning: Error creating simpanan pokok for %s: %v", member.NamaLengkap, err)
		}
	}
	log.Printf("   ✓ Simpanan Pokok for %d members created", len(members))

	// Seed Simpanan Wajib (monthly deposits)
	for _, member := range members[:4] { // Only first 4 members
		for i := 0; i < 3; i++ { // 3 months of deposits
			req := &services.CatatSetoranRequest{
				IDAnggota:    member.ID,
				JenisSimpanan: models.JenisSimpananWajib,
				Jumlah:       50000, // Rp 50,000 per month
				TanggalSetor: time.Now().AddDate(0, -i, 0),
				Keterangan:   fmt.Sprintf("Simpanan Wajib bulan %s", time.Now().AddDate(0, -i, 0).Format("January 2006")),
			}

			_, err := simpananService.CatatSetoran(idKoperasi, req)
			if err != nil {
				log.Printf("   ⚠️  Warning: Error creating simpanan wajib: %v", err)
			}
		}
	}
	log.Printf("   ✓ Simpanan Wajib (3 months) for 4 members created")

	// Seed Simpanan Sukarela (voluntary deposits)
	sukarela := []struct {
		memberIdx int
		amount    float64
	}{
		{0, 200000},
		{1, 150000},
		{2, 300000},
		{4, 100000},
	}

	for _, s := range sukarela {
		if s.memberIdx < len(members) {
			req := &services.CatatSetoranRequest{
				IDAnggota:    members[s.memberIdx].ID,
				JenisSimpanan: models.JenisSimpananSukarela,
				Jumlah:       s.amount,
				TanggalSetor: time.Now(),
				Keterangan:   "Simpanan Sukarela",
			}

			_, err := simpananService.CatatSetoran(idKoperasi, req)
			if err != nil {
				log.Printf("   ⚠️  Warning: Error creating simpanan sukarela: %v", err)
			}
		}
	}
	log.Printf("   ✓ Simpanan Sukarela for 4 members created")
}

func seedSampleSales(penjualanService *services.PenjualanService, idKoperasi, idKasir uuid.UUID, members []*models.Anggota, products []*models.Produk) {
	// Create sample sales transactions
	sales := []struct {
		memberIdx int
		items     []struct {
			productIdx int
			quantity   float64
		}
		metodePembayaran string
		keterangan       string
	}{
		{
			memberIdx: 0,
			items: []struct {
				productIdx int
				quantity   float64
			}{
				{0, 2}, // 2x Beras 5kg
				{1, 1}, // 1x Minyak Goreng
				{2, 1}, // 1x Gula Pasir
			},
			metodePembayaran: "tunai",
			keterangan:       "Pembelian bulanan",
		},
		{
			memberIdx: 1,
			items: []struct {
				productIdx int
				quantity   float64
			}{
				{5, 2}, // 2x Kopi Bubuk
				{6, 3}, // 3x Teh Celup
				{7, 1}, // 1x Susu UHT
			},
			metodePembayaran: "tunai",
			keterangan:       "Belanja minuman",
		},
		{
			memberIdx: 2,
			items: []struct {
				productIdx int
				quantity   float64
			}{
				{8, 3},  // 3x Sabun Mandi
				{9, 2},  // 2x Pasta Gigi
				{10, 5}, // 5x Shampo Sachet
			},
			metodePembayaran: "tunai",
			keterangan:       "Belanja toiletries",
		},
	}

	for idx, sale := range sales {
		if sale.memberIdx >= len(members) {
			continue
		}

		var items []services.ItemPenjualanRequest
		for _, item := range sale.items {
			if item.productIdx < len(products) {
				items = append(items, services.ItemPenjualanRequest{
					IDProduk: products[item.productIdx].ID,
					Jumlah:   item.quantity,
				})
			}
		}

		req := &services.ProsesPenjualanRequest{
			IDAnggota:        &members[sale.memberIdx].ID,
			Items:            items,
			MetodePembayaran: sale.metodePembayaran,
			Keterangan:       sale.keterangan,
		}

		_, err := penjualanService.ProsesPenjualan(idKoperasi, idKasir, req)
		if err != nil {
			log.Printf("   ⚠️  Warning: Error creating sale %d: %v", idx+1, err)
		} else {
			log.Printf("   ✓ Sale transaction %d created", idx+1)
		}
	}
}

// Helper functions

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}
	return string(hash)
}

func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Fatal("Error parsing date:", err)
	}
	return date
}
