package handlers

import (
	"backend-zerotrust-skripsi/internal/app/models"
	"backend-zerotrust-skripsi/internal/database"
	"backend-zerotrust-skripsi/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//
// =======================
// DTO INPUT
// =======================

// Input DTO untuk Create Report
type CreateReportInput struct {
	ReportID string  `json:"report_id" binding:"required"`
	Title    string  `json:"title" binding:"required"`
	Content  string  `json:"content" binding:"required"`
	Value    float64 `json:"value" binding:"required"`
}

// Input DTO khusus Update Status
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"` // Draft / Approved / Rejected
}

//
// =======================
// HANDLERS
// =======================

// GetFinancialReport (READ)
func GetFinancialReport(c *gin.Context) {
	var reports []models.FinancialReport
	if err := database.DB.Find(&reports).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil data", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Berhasil memuat data laporan", reports)
}

// CreateFinancialReport (CREATE)
func CreateFinancialReport(c *gin.Context) {
	// 1. Validasi Input JSON
	var input CreateReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Input data tidak valid", err.Error())
		return
	}

	// 2. Ambil User ID dari Context (dari middleware ZeroTrustPEP)
	userID := c.GetUint("user_id")

	// 3. Simpan ke Database
	newReport := models.FinancialReport{
		ReportID:    input.ReportID,
		Title:       input.Title,
		Content:     input.Content,
		Value:       input.Value,
		Status:      "Draft",
		CreatedByID: userID,
		Model:       gorm.Model{CreatedAt: time.Now()},
	}

	if err := database.DB.Create(&newReport).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal menyimpan laporan", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Laporan berhasil dibuat", newReport)
}

// UpdateFinancialReport (UPDATE STATUS - APPROVAL)
func UpdateFinancialReport(c *gin.Context) {
	// 1. Ambil ID dari URL
	id := c.Param("id")

	// 2. Cari laporan
	var report models.FinancialReport
	if err := database.DB.First(&report, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Laporan tidak ditemukan", nil)
		return
	}

	// 3. Validasi input JSON
	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Input status tidak valid", err.Error())
		return
	}

	// 4. Update status
	report.Status = input.Status
	if err := database.DB.Save(&report).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengupdate status", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Status laporan berhasil diperbarui", report)
}
