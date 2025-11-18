package services

import (
	"cooperative-erp-lite/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LaporanService menangani logika bisnis laporan dan analytics
type LaporanService struct {
	db               *gorm.DB
	akunService      *AkunService
	simpananService  *SimpananService
	penjualanService *PenjualanService
}

// NewLaporanService membuat instance baru LaporanService
func NewLaporanService(db *gorm.DB, akunService *AkunService, simpananService *SimpananService, penjualanService *PenjualanService) *LaporanService {
	return &LaporanService{
		db:               db,
		akunService:      akunService,
		simpananService:  simpananService,
		penjualanService: penjualanService,
	}
}

// LaporanPosisiKeuangan adalah struktur untuk Balance Sheet
type LaporanPosisiKeuangan struct {
	TanggalLaporan time.Time             `json:"tanggalLaporan"`
	Aset           []ItemLaporanKeuangan `json:"aset"`
	TotalAset      float64               `json:"totalAset"`
	Kewajiban      []ItemLaporanKeuangan `json:"kewajiban"`
	TotalKewajiban float64               `json:"totalKewajiban"`
	Modal          []ItemLaporanKeuangan `json:"modal"`
	TotalModal     float64               `json:"totalModal"`
}

// ItemLaporanKeuangan adalah struktur untuk item dalam laporan
type ItemLaporanKeuangan struct {
	KodeAkun string  `json:"kodeAkun"`
	NamaAkun string  `json:"namaAkun"`
	Saldo    float64 `json:"saldo"`
}

// GenerateLaporanPosisiKeuangan membuat laporan neraca/balance sheet
// Optimized version using single query with aggregation to eliminate N+1 query problem
func (s *LaporanService) GenerateLaporanPosisiKeuangan(idKoperasi uuid.UUID, tanggalPer string) (*LaporanPosisiKeuangan, error) {
	// Parse tanggal
	var tanggalLaporan time.Time
	if tanggalPer == "" {
		tanggalLaporan = time.Now()
	} else {
		var err error
		tanggalLaporan, err = time.Parse("2006-01-02", tanggalPer)
		if err != nil {
			return nil, errors.New("format tanggal tidak valid")
		}
	}

	laporan := &LaporanPosisiKeuangan{
		TanggalLaporan: tanggalLaporan,
		Aset:           []ItemLaporanKeuangan{},
		Kewajiban:      []ItemLaporanKeuangan{},
		Modal:          []ItemLaporanKeuangan{},
	}

	// Structure to hold aggregated account balances
	type AccountBalance struct {
		KodeAkun    string
		NamaAkun    string
		TipeAkun    models.TipeAkun
		NormalSaldo string
		TotalDebit  float64
		TotalKredit float64
	}

	var balances []AccountBalance

	// Single optimized query with JOIN and aggregation
	// This replaces the N+1 pattern (1 query for accounts + N queries for balances)
	query := s.db.Table("akun").
		Select(`
			akun.kode_akun,
			akun.nama_akun,
			akun.tipe_akun,
			akun.normal_saldo,
			COALESCE(SUM(baris_transaksi.jumlah_debit), 0) as total_debit,
			COALESCE(SUM(baris_transaksi.jumlah_kredit), 0) as total_kredit
		`).
		Joins("LEFT JOIN baris_transaksi ON baris_transaksi.id_akun = akun.id").
		Joins("LEFT JOIN transaksi ON transaksi.id = baris_transaksi.id_transaksi").
		Where("akun.id_koperasi = ?", idKoperasi)

	// Apply date filter if provided
	if tanggalPer != "" {
		query = query.Where("transaksi.tanggal_transaksi <= ? OR transaksi.id IS NULL", tanggalPer)
	}

	query = query.Group("akun.id, akun.kode_akun, akun.nama_akun, akun.tipe_akun, akun.normal_saldo").
		Order("akun.kode_akun ASC")

	err := query.Scan(&balances).Error
	if err != nil {
		return nil, errors.New("gagal mengambil data laporan posisi keuangan")
	}

	// Process balances and categorize by account type
	for _, balance := range balances {
		// Calculate balance based on normal balance
		var saldo float64
		if balance.NormalSaldo == "debit" {
			saldo = balance.TotalDebit - balance.TotalKredit
		} else {
			saldo = balance.TotalKredit - balance.TotalDebit
		}

		// Skip accounts with zero balance
		if saldo == 0 {
			continue
		}

		item := ItemLaporanKeuangan{
			KodeAkun: balance.KodeAkun,
			NamaAkun: balance.NamaAkun,
			Saldo:    saldo,
		}

		// Categorize by account type (Aset, Kewajiban, Modal only for Balance Sheet)
		switch balance.TipeAkun {
		case models.AkunAset:
			laporan.Aset = append(laporan.Aset, item)
			laporan.TotalAset += saldo
		case models.AkunKewajiban:
			laporan.Kewajiban = append(laporan.Kewajiban, item)
			laporan.TotalKewajiban += saldo
		case models.AkunModal:
			laporan.Modal = append(laporan.Modal, item)
			laporan.TotalModal += saldo
		}
	}

	return laporan, nil
}

// LaporanLabaRugi adalah struktur untuk Income Statement
type LaporanLabaRugi struct {
	PeriodeMulai    time.Time             `json:"periodeMulai"`
	PeriodeAkhir    time.Time             `json:"periodeAkhir"`
	Pendapatan      []ItemLaporanKeuangan `json:"pendapatan"`
	TotalPendapatan float64               `json:"totalPendapatan"`
	Beban           []ItemLaporanKeuangan `json:"beban"`
	TotalBeban      float64               `json:"totalBeban"`
	LabaRugiBersih  float64               `json:"labaRugiBersih"`
}

// GenerateLaporanLabaRugi membuat laporan laba rugi
// Optimized version using single query with date filtering to eliminate N+1 query problem
func (s *LaporanService) GenerateLaporanLabaRugi(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir string) (*LaporanLabaRugi, error) {
	// Parse tanggal
	periodeMulai, err := time.Parse("2006-01-02", tanggalMulai)
	if err != nil {
		return nil, errors.New("format tanggal mulai tidak valid")
	}

	periodeAkhir, err := time.Parse("2006-01-02", tanggalAkhir)
	if err != nil {
		return nil, errors.New("format tanggal akhir tidak valid")
	}

	laporan := &LaporanLabaRugi{
		PeriodeMulai: periodeMulai,
		PeriodeAkhir: periodeAkhir,
		Pendapatan:   []ItemLaporanKeuangan{},
		Beban:        []ItemLaporanKeuangan{},
	}

	// Structure to hold aggregated income and expense data
	type IncomeExpenseBalance struct {
		KodeAkun    string
		NamaAkun    string
		TipeAkun    models.TipeAkun
		NormalSaldo string
		TotalDebit  float64
		TotalKredit float64
	}

	var balances []IncomeExpenseBalance

	// Single optimized query for period-specific balances
	// This replaces the N+1 pattern (fetching accounts + calculating balance for each)
	err = s.db.Table("akun").
		Select(`
			akun.kode_akun,
			akun.nama_akun,
			akun.tipe_akun,
			akun.normal_saldo,
			COALESCE(SUM(baris_transaksi.jumlah_debit), 0) as total_debit,
			COALESCE(SUM(baris_transaksi.jumlah_kredit), 0) as total_kredit
		`).
		Joins("LEFT JOIN baris_transaksi ON baris_transaksi.id_akun = akun.id").
		Joins("LEFT JOIN transaksi ON transaksi.id = baris_transaksi.id_transaksi AND transaksi.tanggal_transaksi BETWEEN ? AND ?", tanggalMulai, tanggalAkhir).
		Where("akun.id_koperasi = ? AND akun.tipe_akun IN (?)", idKoperasi, []models.TipeAkun{models.AkunPendapatan, models.AkunBeban}).
		Group("akun.id, akun.kode_akun, akun.nama_akun, akun.tipe_akun, akun.normal_saldo").
		Order("akun.kode_akun ASC").
		Scan(&balances).Error

	if err != nil {
		return nil, errors.New("gagal mengambil data laporan laba rugi")
	}

	// Process balances and categorize
	for _, balance := range balances {
		// Calculate balance based on normal balance
		var saldoPeriode float64
		if balance.NormalSaldo == "debit" {
			saldoPeriode = balance.TotalDebit - balance.TotalKredit
		} else {
			saldoPeriode = balance.TotalKredit - balance.TotalDebit
		}

		// Skip accounts with zero balance
		if saldoPeriode == 0 {
			continue
		}

		item := ItemLaporanKeuangan{
			KodeAkun: balance.KodeAkun,
			NamaAkun: balance.NamaAkun,
			Saldo:    saldoPeriode,
		}

		// Categorize by account type
		switch balance.TipeAkun {
		case models.AkunPendapatan:
			laporan.Pendapatan = append(laporan.Pendapatan, item)
			laporan.TotalPendapatan += saldoPeriode
		case models.AkunBeban:
			laporan.Beban = append(laporan.Beban, item)
			laporan.TotalBeban += saldoPeriode
		}
	}

	// Hitung laba/rugi bersih
	laporan.LabaRugiBersih = laporan.TotalPendapatan - laporan.TotalBeban

	return laporan, nil
}

// LaporanArusKas adalah struktur untuk Cash Flow Statement
type LaporanArusKas struct {
	PeriodeMulai       time.Time             `json:"periodeMulai"`
	PeriodeAkhir       time.Time             `json:"periodeAkhir"`
	ArusKasOperasional []ItemLaporanKeuangan `json:"arusKasOperasional"`
	TotalOperasional   float64               `json:"totalOperasional"`
	ArusKasInvestasi   []ItemLaporanKeuangan `json:"arusKasInvestasi"`
	TotalInvestasi     float64               `json:"totalInvestasi"`
	ArusKasPendanaan   []ItemLaporanKeuangan `json:"arusKasPendanaan"`
	TotalPendanaan     float64               `json:"totalPendanaan"`
	KenaikanKasBersih  float64               `json:"kenaikanKasBersih"`
	SaldoKasAwal       float64               `json:"saldoKasAwal"`
	SaldoKasAkhir      float64               `json:"saldoKasAkhir"`
}

// GenerateLaporanArusKas membuat laporan arus kas
func (s *LaporanService) GenerateLaporanArusKas(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir string) (*LaporanArusKas, error) {
	// Parse tanggal
	periodeMulai, err := time.Parse("2006-01-02", tanggalMulai)
	if err != nil {
		return nil, errors.New("format tanggal mulai tidak valid")
	}

	periodeAkhir, err := time.Parse("2006-01-02", tanggalAkhir)
	if err != nil {
		return nil, errors.New("format tanggal akhir tidak valid")
	}

	// Dapatkan akun kas (1101)
	akunKas, err := s.akunService.DapatkanAkunByKode(idKoperasi, "1101")
	if err != nil {
		return nil, errors.New("akun kas tidak ditemukan")
	}

	// Saldo kas awal periode
	tanggalSebelum := periodeMulai.AddDate(0, 0, -1).Format("2006-01-02")
	saldoKasAwal, _ := s.akunService.HitungSaldoAkun(akunKas.ID, tanggalSebelum)

	// Saldo kas akhir periode
	saldoKasAkhir, _ := s.akunService.HitungSaldoAkun(akunKas.ID, tanggalAkhir)

	// Untuk MVP, simplified cash flow (langsung dari mutasi kas)
	laporan := &LaporanArusKas{
		PeriodeMulai:       periodeMulai,
		PeriodeAkhir:       periodeAkhir,
		SaldoKasAwal:       saldoKasAwal,
		SaldoKasAkhir:      saldoKasAkhir,
		KenaikanKasBersih:  saldoKasAkhir - saldoKasAwal,
		ArusKasOperasional: []ItemLaporanKeuangan{},
		ArusKasInvestasi:   []ItemLaporanKeuangan{},
		ArusKasPendanaan:   []ItemLaporanKeuangan{},
	}

	// Simplified: kenaikan = total operasional + investasi + pendanaan
	laporan.TotalOperasional = laporan.KenaikanKasBersih

	return laporan, nil
}

// GenerateLaporanSaldoAnggota membuat laporan saldo simpanan anggota
func (s *LaporanService) GenerateLaporanSaldoAnggota(idKoperasi uuid.UUID) ([]models.SaldoSimpananAnggota, error) {
	return s.simpananService.DapatkanLaporanSaldoAnggota(idKoperasi)
}

// LaporanPenjualan adalah struktur untuk sales report
type LaporanPenjualan struct {
	PeriodeMulai      time.Time                `json:"periodeMulai"`
	PeriodeAkhir      time.Time                `json:"periodeAkhir"`
	TotalPenjualan    float64                  `json:"totalPenjualan"`
	JumlahTransaksi   int64                    `json:"jumlahTransaksi"`
	RataRataTransaksi float64                  `json:"rataRataTransaksi"`
	TopProduk         []map[string]interface{} `json:"topProduk"`
}

// GenerateLaporanPenjualan membuat laporan penjualan
func (s *LaporanService) GenerateLaporanPenjualan(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir string) (*LaporanPenjualan, error) {
	periodeMulai, err := time.Parse("2006-01-02", tanggalMulai)
	if err != nil {
		return nil, errors.New("format tanggal mulai tidak valid")
	}

	periodeAkhir, err := time.Parse("2006-01-02", tanggalAkhir)
	if err != nil {
		return nil, errors.New("format tanggal akhir tidak valid")
	}

	// Dapatkan summary penjualan
	summary, err := s.penjualanService.HitungTotalPenjualan(idKoperasi, tanggalMulai, tanggalAkhir)
	if err != nil {
		return nil, err
	}

	// Dapatkan top produk
	topProduk, err := s.penjualanService.DapatkanTopProduk(idKoperasi, 10)
	if err != nil {
		return nil, err
	}

	laporan := &LaporanPenjualan{
		PeriodeMulai:      periodeMulai,
		PeriodeAkhir:      periodeAkhir,
		TotalPenjualan:    summary["totalPenjualan"].(float64),
		JumlahTransaksi:   summary["jumlahTransaksi"].(int64),
		RataRataTransaksi: summary["rataRata"].(float64),
		TopProduk:         topProduk,
	}

	return laporan, nil
}

// LaporanTransaksiHarian adalah struktur untuk daily transaction report
type LaporanTransaksiHarian struct {
	Tanggal         time.Time `json:"tanggal"`
	TotalKasMasuk   float64   `json:"totalKasMasuk"`
	TotalKasKeluar  float64   `json:"totalKasKeluar"`
	SaldoKasAkhir   float64   `json:"saldoKasAkhir"`
	JumlahPenjualan int64     `json:"jumlahPenjualan"`
	JumlahSimpanan  int64     `json:"jumlahSimpanan"`
}

// GenerateLaporanTransaksiHarian membuat laporan transaksi harian
func (s *LaporanService) GenerateLaporanTransaksiHarian(idKoperasi uuid.UUID, tanggal string) (*LaporanTransaksiHarian, error) {
	tgl, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid")
	}

	// Dapatkan akun kas
	akunKas, err := s.akunService.DapatkanAkunByKode(idKoperasi, "1101")
	if err != nil {
		return nil, errors.New("akun kas tidak ditemukan")
	}

	// Hitung saldo kas akhir hari
	saldoKas, _ := s.akunService.HitungSaldoAkun(akunKas.ID, tanggal)

	// Hitung total kas masuk (debit ke kas)
	type KasResult struct {
		Total float64
	}
	var kasMasuk KasResult
	s.db.Model(&models.BarisTransaksi{}).
		Select("COALESCE(SUM(jumlah_debit), 0) as total").
		Joins("JOIN transaksi ON transaksi.id = baris_transaksi.id_transaksi").
		Where("baris_transaksi.id_akun = ? AND DATE(transaksi.tanggal_transaksi) = ?", akunKas.ID, tanggal).
		Scan(&kasMasuk)

	// Hitung total kas keluar (kredit dari kas)
	var kasKeluar KasResult
	s.db.Model(&models.BarisTransaksi{}).
		Select("COALESCE(SUM(jumlah_kredit), 0) as total").
		Joins("JOIN transaksi ON transaksi.id = baris_transaksi.id_transaksi").
		Where("baris_transaksi.id_akun = ? AND DATE(transaksi.tanggal_transaksi) = ?", akunKas.ID, tanggal).
		Scan(&kasKeluar)

	// Hitung jumlah transaksi
	var jumlahPenjualan, jumlahSimpanan int64
	s.db.Model(&models.Penjualan{}).
		Where("id_koperasi = ? AND DATE(tanggal_penjualan) = ?", idKoperasi, tanggal).
		Count(&jumlahPenjualan)

	s.db.Model(&models.Simpanan{}).
		Where("id_koperasi = ? AND DATE(tanggal_transaksi) = ?", idKoperasi, tanggal).
		Count(&jumlahSimpanan)

	laporan := &LaporanTransaksiHarian{
		Tanggal:         tgl,
		TotalKasMasuk:   kasMasuk.Total,
		TotalKasKeluar:  kasKeluar.Total,
		SaldoKasAkhir:   saldoKas,
		JumlahPenjualan: jumlahPenjualan,
		JumlahSimpanan:  jumlahSimpanan,
	}

	return laporan, nil
}

// GetDashboardStats mengambil statistik untuk dashboard
func (s *LaporanService) GetDashboardStats(idKoperasi uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total anggota aktif
	var totalAnggota int64
	s.db.Model(&models.Anggota{}).
		Where("id_koperasi = ? AND status = ?", idKoperasi, models.StatusAktif).
		Count(&totalAnggota)
	stats["totalAnggota"] = totalAnggota

	// Ringkasan simpanan
	ringkasanSimpanan, err := s.simpananService.DapatkanRingkasanSimpanan(idKoperasi)
	if err == nil {
		stats["totalSimpanan"] = ringkasanSimpanan.TotalSemuaSimpanan
	}

	// Penjualan hari ini
	penjualanHariIni, err := s.penjualanService.DapatkanPenjualanHariIni(idKoperasi)
	if err == nil {
		stats["penjualanHariIni"] = penjualanHariIni["totalPenjualan"]
		stats["transaksiHariIni"] = penjualanHariIni["jumlahTransaksi"]
	}

	// Saldo kas
	akunKas, err := s.akunService.DapatkanAkunByKode(idKoperasi, "1101")
	if err == nil {
		saldoKas, _ := s.akunService.HitungSaldoAkun(akunKas.ID, "")
		stats["saldoKas"] = saldoKas
	}

	return stats, nil
}

// GenerateLaporanPerubahanModal generates statement of changes in equity
func (s *LaporanService) GenerateLaporanPerubahanModal(idKoperasi uuid.UUID, tanggalMulai, tanggalAkhir string) (map[string]interface{}, error) {
	// This is a placeholder implementation
	// TODO: Implement full statement of changes in equity
	return map[string]interface{}{
		"message": "Laporan Perubahan Modal - Not yet implemented",
	}, nil
}

// GenerateBukuBesar generates general ledger for an account
func (s *LaporanService) GenerateBukuBesar(idKoperasi, idAkun uuid.UUID, tanggalMulai, tanggalAkhir string) (map[string]interface{}, error) {
	// Get account information
	akun, err := s.akunService.DapatkanAkun(idAkun)
	if err != nil {
		return nil, err
	}

	// Validate multi-tenancy
	var akunModel models.Akun
	err = s.db.Where("id = ? AND id_koperasi = ?", idAkun, idKoperasi).First(&akunModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("akun tidak ditemukan atau tidak memiliki akses")
		}
		return nil, err
	}

	// Get all transaction lines for this account within date range
	type TransactionDetail struct {
		Tanggal    string  `json:"tanggal"`
		NoJurnal   string  `json:"noJurnal"`
		Keterangan string  `json:"keterangan"`
		Debit      float64 `json:"debit"`
		Kredit     float64 `json:"kredit"`
		Saldo      float64 `json:"saldo"`
	}

	var details []TransactionDetail
	var runningBalance float64

	// Get starting balance (before tanggalMulai)
	if tanggalMulai != "" {
		runningBalance, _ = s.akunService.HitungSaldoAkun(idAkun, tanggalMulai)
	}

	// Query transaction lines
	query := s.db.Table("jurnal_detail").
		Select("jurnal.tanggal, jurnal.no_jurnal, jurnal.keterangan, jurnal_detail.debit, jurnal_detail.kredit").
		Joins("JOIN jurnal ON jurnal_detail.id_jurnal = jurnal.id").
		Where("jurnal_detail.id_akun = ? AND jurnal.id_koperasi = ?", idAkun, idKoperasi)

	if tanggalMulai != "" {
		query = query.Where("jurnal.tanggal >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("jurnal.tanggal <= ?", tanggalAkhir)
	}

	query = query.Order("jurnal.tanggal ASC, jurnal.created_at ASC")

	var rows []struct {
		Tanggal    string
		NoJurnal   string
		Keterangan string
		Debit      float64
		Kredit     float64
	}

	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	// Process each transaction
	for _, row := range rows {
		// Update running balance
		if akun.NormalSaldo == "debit" {
			runningBalance += row.Debit - row.Kredit
		} else {
			runningBalance += row.Kredit - row.Debit
		}

		details = append(details, TransactionDetail{
			Tanggal:    row.Tanggal,
			NoJurnal:   row.NoJurnal,
			Keterangan: row.Keterangan,
			Debit:      row.Debit,
			Kredit:     row.Kredit,
			Saldo:      runningBalance,
		})
	}

	return map[string]interface{}{
		"akun": map[string]interface{}{
			"kode": akun.KodeAkun,
			"nama": akun.NamaAkun,
			"tipe": akun.TipeAkun,
		},
		"periode": map[string]string{
			"tanggalMulai": tanggalMulai,
			"tanggalAkhir": tanggalAkhir,
		},
		"transaksi":  details,
		"saldoAkhir": runningBalance,
	}, nil
}

// GenerateNeracaSaldo generates trial balance
func (s *LaporanService) GenerateNeracaSaldo(idKoperasi uuid.UUID, tanggalPer string) (map[string]interface{}, error) {
	// Get all active accounts
	akunList, err := s.akunService.DapatkanSemuaAkun(idKoperasi, "", nil)
	if err != nil {
		return nil, err
	}

	type NeracaSaldoItem struct {
		KodeAkun    string  `json:"kodeAkun"`
		NamaAkun    string  `json:"namaAkun"`
		TipeAkun    string  `json:"tipeAkun"`
		SaldoDebit  float64 `json:"saldoDebit"`
		SaldoKredit float64 `json:"saldoKredit"`
	}

	var items []NeracaSaldoItem
	var totalDebit, totalKredit float64

	for _, akun := range akunList {
		saldo, _ := s.akunService.HitungSaldoAkun(akun.ID, tanggalPer)

		item := NeracaSaldoItem{
			KodeAkun: akun.KodeAkun,
			NamaAkun: akun.NamaAkun,
			TipeAkun: string(akun.TipeAkun),
		}

		if saldo > 0 {
			if akun.NormalSaldo == "debit" {
				item.SaldoDebit = saldo
				totalDebit += saldo
			} else {
				item.SaldoKredit = saldo
				totalKredit += saldo
			}
		}

		items = append(items, item)
	}

	return map[string]interface{}{
		"tanggalPer":  tanggalPer,
		"items":       items,
		"totalDebit":  totalDebit,
		"totalKredit": totalKredit,
		"isBalanced":  totalDebit == totalKredit,
	}, nil
}
