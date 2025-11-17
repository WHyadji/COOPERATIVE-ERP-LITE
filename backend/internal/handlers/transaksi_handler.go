package handlers

import (
	"cooperative-erp-lite/internal/constants"
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"errors"
	"net/http"

	apperrors "cooperative-erp-lite/internal/errors"

	"github.com/gin-gonic/gin"
)

// TransaksiHandler menangani endpoint transaksi akuntansi
type TransaksiHandler struct {
	transaksiService *services.TransaksiService
}

// NewTransaksiHandler membuat instance baru TransaksiHandler
func NewTransaksiHandler(transaksiService *services.TransaksiService) *TransaksiHandler {
	return &TransaksiHandler{
		transaksiService: transaksiService,
	}
}

// Create menangani pembuatan transaksi akuntansi baru.
// Memvalidasi akses multi-tenant untuk memastikan transaksi dibuat dalam koperasi yang benar.
// Memvalidasi double-entry accounting (total debit harus sama dengan total kredit).
//
// Route: POST /api/v1/transaksi
// Request Body: BuatTransaksiRequest
// Response: TransaksiResponse dengan status 201 Created
//
// Error Responses:
//   - 400 Bad Request: Jika validasi debit-kredit gagal
//   - 500 Internal Server Error: Jika terjadi kesalahan server
func (h *TransaksiHandler) Create(c *gin.Context) {
	koperasiUUID, ok := AmbilIDKoperasiDariContext(c)
	if !ok {
		return
	}

	penggunaUUID, ok := AmbilIDPenggunaDariContext(c)
	if !ok {
		return
	}

	var req services.BuatTransaksiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	transaksi, err := h.transaksiService.BuatTransaksi(koperasiUUID, penggunaUUID, &req)
	if err != nil {
		// Check jika error validasi double-entry
		if errors.Is(err, apperrors.ErrDebitKreditTidakBalance) ||
			errors.Is(err, apperrors.ErrDebitKreditKeduanya) {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, constants.PesanTransaksiBerhasilDibuat, transaksi)
}

// List menangani pengambilan daftar transaksi dengan paginasi dan filter.
// Memvalidasi akses multi-tenant untuk memastikan hanya transaksi dalam koperasi yang benar yang diambil.
// Mendukung filter berdasarkan tanggal mulai, tanggal akhir, dan tipe transaksi.
//
// Route: GET /api/v1/transaksi
// Query Parameters:
//   - page: Nomor halaman (default: 1)
//   - pageSize: Jumlah item per halaman (default: 20, max: 100)
//   - tanggalMulai: Filter tanggal mulai (opsional)
//   - tanggalAkhir: Filter tanggal akhir (opsional)
//   - tipeTransaksi: Filter tipe transaksi (opsional)
// Response: Array TransaksiResponse dengan metadata paginasi dan status 200 OK
//
// Error Responses:
//   - 500 Internal Server Error: Jika terjadi kesalahan server
func (h *TransaksiHandler) List(c *gin.Context) {
	koperasiUUID, ok := AmbilIDKoperasiDariContext(c)
	if !ok {
		return
	}

	paginasi := AmbilParameterPaginasi(c)
	tanggalMulai := c.Query("tanggalMulai")
	tanggalAkhir := c.Query("tanggalAkhir")
	tipeTransaksi := c.Query("tipeTransaksi")

	transaksiList, total, err := h.transaksiService.DapatkanSemuaTransaksi(
		koperasiUUID, tanggalMulai, tanggalAkhir, tipeTransaksi, paginasi.Halaman, paginasi.UkuranHalaman,
	)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	metaPaginasi := utils.CalculatePaginationMeta(paginasi.Halaman, paginasi.UkuranHalaman, total)
	utils.PaginatedSuccessResponse(c, http.StatusOK, constants.PesanTransaksiBerhasilDiambil, transaksiList, metaPaginasi)
}

// GetByID menangani pengambilan detail transaksi berdasarkan ID.
// Memvalidasi akses multi-tenant untuk memastikan transaksi yang diambil milik koperasi yang benar.
//
// Route: GET /api/v1/transaksi/:id
// URL Parameters:
//   - id: UUID transaksi
// Response: TransaksiResponse dengan status 200 OK
//
// Error Responses:
//   - 400 Bad Request: Jika ID tidak valid
//   - 404 Not Found: Jika transaksi tidak ditemukan atau tidak milik koperasi pengguna
//   - 500 Internal Server Error: Jika terjadi kesalahan server
func (h *TransaksiHandler) GetByID(c *gin.Context) {
	koperasiUUID, ok := AmbilIDKoperasiDariContext(c)
	if !ok {
		return
	}

	id, ok := ParseUUIDDariParameter(c, "id")
	if !ok {
		return
	}

	transaksi, err := h.transaksiService.DapatkanTransaksi(koperasiUUID, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrTransaksiTidakDitemukan) {
			utils.NotFoundResponse(c, constants.PesanTransaksiTidakDitemukan)
		} else {
			utils.InternalServerErrorResponse(c, err.Error(), nil)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, constants.PesanTransaksiBerhasilDiambil, transaksi)
}

// Update menangani pembaruan transaksi akuntansi.
// Memvalidasi akses multi-tenant untuk memastikan transaksi yang diupdate milik koperasi yang benar.
// Memvalidasi double-entry accounting (total debit harus sama dengan total kredit).
// Transaksi yang sudah di-post tidak dapat diupdate.
//
// Route: PUT /api/v1/transaksi/:id
// URL Parameters:
//   - id: UUID transaksi
// Request Body: BuatTransaksiRequest
// Response: TransaksiResponse dengan status 200 OK
//
// Error Responses:
//   - 400 Bad Request: Jika ID tidak valid, validasi debit-kredit gagal, atau transaksi sudah di-post
//   - 404 Not Found: Jika transaksi tidak ditemukan atau tidak milik koperasi pengguna
//   - 500 Internal Server Error: Jika terjadi kesalahan server
func (h *TransaksiHandler) Update(c *gin.Context) {
	koperasiUUID, ok := AmbilIDKoperasiDariContext(c)
	if !ok {
		return
	}

	id, ok := ParseUUIDDariParameter(c, "id")
	if !ok {
		return
	}

	var req services.BuatTransaksiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	transaksi, err := h.transaksiService.PerbaruiTransaksi(koperasiUUID, id, &req)
	if err != nil {
		// Check if error is validation error from double-entry
		if errors.Is(err, apperrors.ErrDebitKreditTidakBalance) ||
			errors.Is(err, apperrors.ErrDebitKreditKeduanya) {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		// Check if transaction not found (could be because it doesn't exist or doesn't belong to cooperative)
		if errors.Is(err, apperrors.ErrTransaksiTidakDitemukan) {
			utils.NotFoundResponse(c, constants.PesanTransaksiTidakDitemukan)
			return
		}
		// Check for "posted" status error when the field is added
		if errors.Is(err, apperrors.ErrTransaksiSudahDiPost) {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, constants.PesanTransaksiBerhasilDiupdate, transaksi)
}

// Delete menangani penghapusan transaksi akuntansi.
// Memvalidasi akses multi-tenant untuk memastikan transaksi yang dihapus milik koperasi yang benar.
// Transaksi yang sudah di-post tidak dapat dihapus (harus di-reverse).
//
// Route: DELETE /api/v1/transaksi/:id
// URL Parameters:
//   - id: UUID transaksi
// Response: Pesan sukses dengan status 200 OK
//
// Error Responses:
//   - 400 Bad Request: Jika ID tidak valid atau transaksi sudah di-post
//   - 404 Not Found: Jika transaksi tidak ditemukan atau tidak milik koperasi pengguna
//   - 500 Internal Server Error: Jika terjadi kesalahan server
func (h *TransaksiHandler) Delete(c *gin.Context) {
	koperasiUUID, ok := AmbilIDKoperasiDariContext(c)
	if !ok {
		return
	}

	id, ok := ParseUUIDDariParameter(c, "id")
	if !ok {
		return
	}

	if err := h.transaksiService.HapusTransaksi(koperasiUUID, id); err != nil {
		if errors.Is(err, apperrors.ErrTransaksiSudahDiPostTidakBisaHapus) {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, constants.PesanTransaksiBerhasilDihapus, nil)
}

// Reverse menangani pembuatan jurnal pembalik (reversing entry) untuk transaksi.
// Memvalidasi akses multi-tenant untuk memastikan transaksi yang di-reverse milik koperasi yang benar.
// Membuat transaksi baru dengan debit dan kredit yang dibalik dari transaksi asli.
//
// Route: POST /api/v1/transaksi/:id/reverse
// URL Parameters:
//   - id: UUID transaksi yang akan di-reverse
// Request Body:
//   - keterangan: Alasan/keterangan reversal (required)
// Response: TransaksiResponse (transaksi reversal yang baru) dengan status 200 OK
//
// Error Responses:
//   - 400 Bad Request: Jika ID tidak valid, keterangan kosong, atau transaksi tidak memiliki baris
//   - 404 Not Found: Jika transaksi tidak ditemukan atau tidak milik koperasi pengguna
//   - 500 Internal Server Error: Jika terjadi kesalahan server
func (h *TransaksiHandler) Reverse(c *gin.Context) {
	koperasiUUID, ok := AmbilIDKoperasiDariContext(c)
	if !ok {
		return
	}

	penggunaUUID, ok := AmbilIDPenggunaDariContext(c)
	if !ok {
		return
	}

	id, ok := ParseUUIDDariParameter(c, "id")
	if !ok {
		return
	}

	var req struct {
		Keterangan string `json:"keterangan" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	reversedTransaksi, err := h.transaksiService.ReverseTransaksi(koperasiUUID, penggunaUUID, id, req.Keterangan)
	if err != nil {
		// Check if transaction not found (could be because it doesn't exist or doesn't belong to cooperative)
		if errors.Is(err, apperrors.ErrTransaksiTidakDitemukan) {
			utils.NotFoundResponse(c, constants.PesanTransaksiTidakDitemukan)
			return
		}
		// Check if transaction has no lines to reverse
		if errors.Is(err, apperrors.ErrTidakAdaBarisTransaksi) {
			utils.BadRequestResponse(c, err.Error())
			return
		}
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, constants.PesanTransaksiBerhasilDireverse, reversedTransaksi)
}
