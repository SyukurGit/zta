package database

import (
	"backend-zerotrust-skripsi/config"
	"backend-zerotrust-skripsi/internal/app/models"
	"log"
	"golang.org/x/crypto/bcrypt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Pastikan menerima parameter *config.AppConfig
func ConnectDB(cfg *config.AppConfig) {
	var err error
	// Gunakan DBName dari config
	DB, err = gorm.Open(sqlite.Open(cfg.DBName), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.FinancialReport{})
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}

	log.Printf("✅ Database %s terhubung!", cfg.DBName)
	seedUsers()
	seedReports()
}

// seedUsers mengisi data user dummy jika tabel kosong
func seedUsers() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return // Sudah ada data, skip
	}

	// Password default untuk semua user: "rahasia"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("rahasia"), bcrypt.DefaultCost)
	passStr := string(passwordHash)

	users := []models.User{
		{Username: "manager", Password: passStr, Role: "manager"},
		{Username: "auditor", Password: passStr, Role: "auditor"},
		{Username: "staff", Password: passStr, Role: "staff"}, // Role yang akan ditolak aksesnya
	}

	DB.Create(&users)
	log.Println("🌱 Seed Data: User dummy berhasil dibuat (Password: 'rahasia')")
}

func seedReports() {
	var count int64
	DB.Model(&models.FinancialReport{}).Count(&count)
	if count > 0 {
		return
	}

	reports := []models.FinancialReport{
		{ReportID: "REP-001", Title: "Q1 Audit", Content: "Rahasia Perusahaan A", Value: 1500000000, Status: "Approved"},
		{ReportID: "REP-002", Title: "Q2 Audit", Content: "Rahasia Perusahaan B", Value: 2500000000, Status: "Draft"},
	}
	DB.Create(&reports)
	log.Println("🌱 Seed Data: Laporan dummy berhasil dibuat")
}