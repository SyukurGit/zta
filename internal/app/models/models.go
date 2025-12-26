package models

import (
	// "time"

	"gorm.io/gorm"
)

// User merepresentasikan tabel pengguna di database
type User struct {
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt otomatis
	Username   string `gorm:"uniqueIndex;not null" json:"username"`
	Password   string `gorm:"not null" json:"-"` // Password disimpan ter-hash, tidak di-return JSON
	Role       string `gorm:"not null" json:"role"` // manager, auditor, staff
}

// FinancialReport merepresentasikan data sensitif
type FinancialReport struct {
	gorm.Model
	ReportID    string  `json:"report_id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Value       float64 `json:"value"`
	Status      string  `json:"status"` // draft, approved
	CreatedByID uint    `json:"created_by_id"`
}