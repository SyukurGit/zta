package policies

import "time"

// AccessRequest merepresentasikan 'Trust Context' atau data input untuk evaluasi
type AccessRequest struct {
    SubjectID   string    // ID User (Siapa?)
    SubjectRole string    // Role User (Jabatan?)
    Resource    string    // Endpoint yang dituju (Ke mana?)
    Action      string    // GET/POST (Mau ngapain?)
    DeviceIP    string    // IP Address (Dari mana?)
    RequestTime time.Time // Waktu request (Kapan?)
}

// AccessResult adalah output keputusan dari Policy Engine
type AccessResult struct {
    Allow  bool   // Boleh atau Tidak?
    Reason string // Alasan keputusan (Penting untuk Audit Log/Skripsi)
}