package main

import (
    "backend-zerotrust-skripsi/config"
    "backend-zerotrust-skripsi/internal/app/handlers"
    "backend-zerotrust-skripsi/internal/security/audit"
    "backend-zerotrust-skripsi/internal/security/middleware"
    "backend-zerotrust-skripsi/internal/security/policies"
    "log"
    "time"

    "github.com/gin-contrib/cors" // Import library CORS
    "github.com/gin-gonic/gin"
)

func main() {
    // 1. Load Configuration
    cfg := config.LoadConfig()

    // 2. Init Components
    // Policy Engine sekarang butuh Config
    policyEngine := policies.NewEngine(cfg)
    auditLogger := audit.NewLogger()

    r := gin.Default()

    // 3. SETUP CORS (Sangat Penting untuk Frontend Next.js)
    // Tanpa ini, browser akan memblokir request dari localhost:3000
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // URL Frontend Next.js kamu
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-User-Role", "X-User-ID"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // 4. Setup Routes
    api := r.Group("/api/v1")
    
    // Attach Zero Trust Middleware
    api.Use(middleware.ZeroTrustPEP(policyEngine, auditLogger))
    
    {
        api.GET("/finance/reports", handlers.GetFinancialReport)
    }

    log.Printf("Sistem Audit Zero Trust berjalan di port %s", cfg.AppPort)
    if err := r.Run(cfg.AppPort); err != nil {
        log.Fatal(err)
    }
}