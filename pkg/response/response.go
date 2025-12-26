package response

import (
    "github.com/gin-gonic/gin"
    // "net/http"
)

// APIResponse adalah kontrak baku komunikasi antara Backend dan Frontend
type APIResponse struct {
    Success bool        `json:"success"` // true/false
    Message string      `json:"message"` // Pesan untuk user/dev
    Data    interface{} `json:"data,omitempty"` // Data payload (kosong jika error)
    Error   interface{} `json:"error,omitempty"` // Detail error (jika ada)
}

// Success mengirim respon sukses standar
func Success(c *gin.Context, code int, message string, data interface{}) {
    c.JSON(code, APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    })
}

// Error mengirim respon error standar
func Error(c *gin.Context, code int, message string, errDetail interface{}) {
    c.JSON(code, APIResponse{
        Success: false,
        Message: message,
        Error:   errDetail,
    })
}