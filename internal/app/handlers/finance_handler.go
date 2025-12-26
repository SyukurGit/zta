package handlers

import (
    "backend-zerotrust-skripsi/pkg/response" // Gunakan package response
    "net/http"
    "github.com/gin-gonic/gin"
)

func GetFinancialReport(c *gin.Context) {
    data := gin.H{
        "report_id": "REP-2024-001",
        "content":   "Laporan Audit Keuangan Internal - RAHASIA",
        "value":     1500000000,
    }
    
    // Gunakan wrapper response yang rapi
    response.Success(c, http.StatusOK, "Berhasil mengakses data sensitif", data)
}