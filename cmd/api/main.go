package main

import (
	"backend-zerotrust-skripsi/config"
	"backend-zerotrust-skripsi/internal/app/handlers"
	"backend-zerotrust-skripsi/internal/database"
	"backend-zerotrust-skripsi/internal/security/audit"
	"backend-zerotrust-skripsi/internal/security/middleware"
	"backend-zerotrust-skripsi/internal/security/policies"
	"backend-zerotrust-skripsi/internal/security/token"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Config (Dari .env)
	cfg := config.LoadConfig()

	// 2. Connect Database (Inject Config)
	database.ConnectDB(cfg)

	// 3. Init Services (Inject Config)
	jwtService := token.NewJWTService(cfg.JWTSecret) // Service Auth
	policyEngine := policies.NewEngine(cfg)          // Service Policy
	auditLogger := audit.NewLogger()                 // Service Logger

    // 4. Init Handlers
    authHandler := handlers.NewAuthHandler(jwtService)

	r := gin.Default()

	// 5. CORS Setup
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 6. Setup Routes
	api := r.Group("/api/v1")

	// A. PUBLIC ROUTES
	// Sekarang memanggil method Login milik struct authHandler
	api.POST("/auth/login", authHandler.Login)

	// B. PROTECTED ROUTES
	protected := api.Group("/")
	// Inject jwtService ke Middleware
	protected.Use(middleware.ZeroTrustPEP(policyEngine, auditLogger, jwtService))
	{
		protected.GET("/finance/reports", handlers.GetFinancialReport)
		// Route GET (Read)
		// protected.GET("/finance/reports", handlers.GetFinancialReport)
		
		// Route POST (Create) - NEW!
		// ZeroTrustPEP akan otomatis memblokir jika user bukan Auditor (sesuai engine.go tadi)
		protected.POST("/finance/reports", handlers.CreateFinancialReport)


		protected.PUT("/finance/reports/:id", handlers.UpdateFinancialReport)
		
	}

	log.Printf("🛡️  Zero Trust System Running (Env: %s) on port %s", "Production-Like", cfg.AppPort)
	if err := r.Run(cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}