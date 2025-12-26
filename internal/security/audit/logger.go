package audit

import (
    "encoding/json"
    "os"
    "time"
)

// LogEntry mendefinisikan struktur log keamanan yang baku
type LogEntry struct {
    Timestamp   time.Time `json:"timestamp"`
    EventID     string    `json:"event_id"`      // UUID unik untuk setiap request (opsional, bagus untuk tracking)
    SubjectID   string    `json:"subject_id"`    // Who (Siapa)
    SubjectRole string    `json:"subject_role"`  // Role saat itu
    IPAddress   string    `json:"ip_address"`    // Where (Dari mana)
    Resource    string    `json:"resource"`      // What (Akses apa)
    Action      string    `json:"action"`        // Method (GET/POST)
    Decision    string    `json:"decision"`      // ALLOW / DENY
    Reason      string    `json:"reason"`        // Why (Alasan keputusan)
}

// Logger adalah service untuk mencatat audit trail
type Logger struct {
    // Di sistem nyata, ini bisa koneksi ke Elasticsearch, Splunk, atau file
    // Untuk skripsi, kita output ke Standard Output (Console) dengan format JSON
}

func NewLogger() *Logger {
    return &Logger{}
}

// LogAccess mencatat kejadian akses (baik sukses maupun gagal)
func (l *Logger) LogAccess(userID, role, ip, resource, action string, allowed bool, reason string) {
    
    decision := "DENY"
    if allowed {
        decision = "ALLOW"
    }

    entry := LogEntry{
        Timestamp:   time.Now(),
        SubjectID:   userID,
        SubjectRole: role,
        IPAddress:   ip,
        Resource:    resource,
        Action:      action,
        Decision:    decision,
        Reason:      reason,
    }

    // Encode struct ke JSON dan print ke terminal
    // JSON memudahkan log ini dibaca oleh mesin atau manusia
    _ = json.NewEncoder(os.Stdout).Encode(entry)
}