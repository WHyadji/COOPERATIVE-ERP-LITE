package handlers

import (
	"cooperative-erp-lite/internal/services"
	"cooperative-erp-lite/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProdukHandler menangani endpoint manajemen produk
type ProdukHandler struct {
	produkService *services.ProdukService
}

// NewProdukHandler membuat instance baru ProdukHandler
func NewProdukHandler(produkService *services.ProdukService) *ProdukHandler {
	return &ProdukHandler{
		produkService: produkService,
	}
}

// Create handles POST /api/v1/produk
func (h *ProdukHandler) Create(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	var req services.BuatProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	produk, err := h.produkService.BuatProduk(koperasiUUID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Produk berhasil dibuat", produk)
}

// List handles GET /api/v1/produk
func (h *ProdukHandler) List(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	kategori := c.Query("kategori")
	search := c.Query("search")
	statusAktif := c.Query("statusAktif")

	// Parse status aktif
	var statusAktifPtr *bool
	if statusAktif != "" {
		aktif := statusAktif == "true"
		statusAktifPtr = &aktif
	}

	produkList, total, err := h.produkService.DapatkanSemuaProduk(koperasiUUID, kategori, search, statusAktifPtr, page, pageSize)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	pagination := utils.CalculatePaginationMeta(page, pageSize, total)
	utils.PaginatedSuccessResponse(c, http.StatusOK, "Data produk berhasil diambil", produkList, pagination)
}

// GetByID handles GET /api/v1/produk/:id
func (h *ProdukHandler) GetByID(c *gin.Context) {
	koperasiUUID, ok := AmbilIDKoperasiDariContext(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID produk tidak valid")
		return
	}

	produk, err := h.produkService.DapatkanProduk(koperasiUUID, id)
	if err != nil {
		utils.NotFoundResponse(c, "Produk tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data produk berhasil diambil", produk)
}

// GetByBarcode handles GET /api/v1/produk/barcode/:barcode
func (h *ProdukHandler) GetByBarcode(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	barcode := c.Param("barcode")

	produk, err := h.produkService.DapatkanProdukByBarcode(koperasiUUID, barcode)
	if err != nil {
		utils.NotFoundResponse(c, "Produk tidak ditemukan")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data produk berhasil diambil", produk)
}

// Update handles PUT /api/v1/produk/:id
func (h *ProdukHandler) Update(c *gin.Context) {
	koperasiUUID, ok := AmbilIDKoperasiDariContext(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID produk tidak valid")
		return
	}

	var req services.PerbaruiProdukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	produk, err := h.produkService.PerbaruiProduk(koperasiUUID, id, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Produk berhasil diupdate", produk)
}

// Delete handles DELETE /api/v1/produk/:id
func (h *ProdukHandler) Delete(c *gin.Context) {
	koperasiUUID, ok := AmbilIDKoperasiDariContext(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID produk tidak valid")
		return
	}

	if err := h.produkService.HapusProduk(koperasiUUID, id); err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Produk berhasil dihapus", nil)
}

// AdjustStok handles POST /api/v1/produk/:id/adjust-stok
func (h *ProdukHandler) AdjustStok(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "ID produk tidak valid")
		return
	}

	var req struct {
		Jumlah     int    `json:"jumlah" binding:"required"`
		Keterangan string `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	produk, err := h.produkService.AdjustStok(koperasiUUID, id, req.Jumlah, req.Keterangan)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stok berhasil disesuaikan", produk)
}

// GetStokRendah handles GET /api/v1/produk/stok-rendah
func (h *ProdukHandler) GetStokRendah(c *gin.Context) {
	idKoperasi, _ := c.Get("idKoperasi")
	koperasiUUID := idKoperasi.(uuid.UUID)

	produkList, err := h.produkService.DapatkanProdukStokRendah(koperasiUUID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data produk stok rendah berhasil diambil", produkList)
}
