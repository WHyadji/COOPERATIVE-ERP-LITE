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

	// Dapatkan semua akun
	akunList, err := s.akunService.DapatkanSemuaAkun(idKoperasi, "", nil)
	if err != nil {
		return nil, err
	}

	// Hitung saldo untuk setiap akun
	for _, akunResp := range akunList {
		saldo, err := s.akunService.HitungSaldoAkun(akunResp.ID, tanggalPer)
		if err != nil {
			continue
		}

		// Skip akun dengan saldo 0
		if saldo == 0 {
			continue
		}

		item := ItemLaporanKeuangan{
			KodeAkun: akunResp.KodeAkun,
			NamaAkun: akunResp.NamaAkun,
			Saldo:    saldo,
		}

		// Kategorikan berdasarkan tipe akun
		switch akunResp.TipeAkun {
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

	// Dapatkan akun pendapatan dan beban
	akunList, err := s.akunService.DapatkanSemuaAkun(idKoperasi, "", nil)
	if err != nil {
		return nil, err
	}

	// Hitung saldo untuk periode
	for _, akunResp := range akunList {
		if akunResp.TipeAkun != models.AkunPendapatan && akunResp.TipeAkun != models.AkunBeban {
			continue
		}

		// Hitung saldo akhir periode
		saldoAkhir, _ := s.akunService.HitungSaldoAkun(akunResp.ID, tanggalAkhir)

		// Hitung saldo awal periode (sehari sebelum tanggal mulai)
		tanggalSebelum := periodeMulai.AddDate(0, 0, -1).Format("2006-01-02")
		saldoAwal, _ := s.akunService.HitungSaldoAkun(akunResp.ID, tanggalSebelum)

		// Saldo periode = saldo akhir - saldo awal
		saldoPeriode := saldoAkhir - saldoAwal

		if saldoPeriode == 0 {
			continue
		}

		item := ItemLaporanKeuangan{
			KodeAkun: akunResp.KodeAkun,
			NamaAkun: akunResp.NamaAkun,
			Saldo:    saldoPeriode,
		}

		// Kategorikan
		if akunResp.TipeAkun == models.AkunPendapatan {
			laporan.Pendapatan = append(laporan.Pendapatan, item)
			laporan.TotalPendapatan += saldoPeriode
		} else if akunResp.TipeAkun == models.AkunBeban {
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
	// Delegate to AkunService
	result, err := s.akunService.GetBukuBesar(idKoperasi, idAkun, tanggalMulai, tanggalAkhir)
	if err != nil {
		return nil, err
	}

	// Type assert to map[string]interface{}
	if bukuBesar, ok := result.(map[string]interface{}); ok {
		return bukuBesar, nil
	}

	return nil, errors.New("format buku besar tidak valid")
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
