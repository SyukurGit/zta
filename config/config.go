package config

// AppConfig menyimpan semua setting aplikasi
type AppConfig struct {
    AppPort      string
    WorkStartHour int
    WorkEndHour   int
    SafeNetwork   string // Network prefix (misal 127.0.0.1)
}

// LoadConfig memuat konfigurasi (bisa dari .env nanti, sekarang kita set default)
func LoadConfig() *AppConfig {
    return &AppConfig{
        AppPort:       ":8080",
        WorkStartHour: 8,  // Jam 8 Pagi
        WorkEndHour:   17, // Jam 5 Sore
        SafeNetwork:   "127.0.0.1",
    }
}