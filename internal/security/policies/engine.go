package policies

import (
	"backend-zerotrust-skripsi/config"
	"strings"
)

type Engine struct {
	Cfg *config.AppConfig
}

func NewEngine(cfg *config.AppConfig) *Engine {
	return &Engine{Cfg: cfg}
}

func (e *Engine) Evaluate(req AccessRequest) AccessResult {

	// =========================================================
	// POLICY: Finance Reports Resource
	// Menggunakan HasPrefix agar endpoint dengan ID (e.g., /reports/1) juga kena policy ini
	// =========================================================
	if strings.HasPrefix(req.Resource, "/api/v1/finance/reports") {

		// 1. GLOBAL CONTEXT CHECK (Wajib untuk semua method)
		// -----------------------------------------------------
		
		// A. Cek Jam Kerja
		hour := req.RequestTime.Hour()
		if hour < e.Cfg.WorkStartHour || hour > e.Cfg.WorkEndHour {
			return AccessResult{Allow: false, Reason: "Akses ditolak: Diluar jam operasional"}
		}

		// B. Cek Jaringan Aman
		isSafeNetwork := strings.HasPrefix(req.DeviceIP, e.Cfg.SafeNetwork) || req.DeviceIP == "::1"
		if !isSafeNetwork {
			return AccessResult{Allow: false, Reason: "Akses ditolak: Wajib menggunakan jaringan kantor aman"}
		}

		// 2. ACTION-BASED ROLE CHECK (Segregation of Duties)
		// -----------------------------------------------------
		
		// Skenario: MEMBUAT DATA (POST) -> Khusus Auditor
		if req.Action == "POST" {
			if req.SubjectRole != "auditor" {
				return AccessResult{
					Allow:  false, 
					Reason: "Pelanggaran SoD: Hanya Auditor yang berhak membuat laporan baru",
				}
			}
			return AccessResult{Allow: true, Reason: "Create Access Granted"}
		}

		// Skenario: UPDATE/APPROVAL (PUT) -> Khusus Manager
		if req.Action == "PUT" {
			if req.SubjectRole != "manager" {
				return AccessResult{
					Allow:  false,
					Reason: "Pelanggaran SoD: Hanya Manager yang berhak melakukan Approval",
				}
			}
			return AccessResult{Allow: true, Reason: "Approval Access Granted"}
		}

		// Skenario: MELIHAT DATA (GET) -> Manager & Auditor
		if req.Action == "GET" {
			if req.SubjectRole != "manager" && req.SubjectRole != "auditor" {
				return AccessResult{Allow: false, Reason: "Role tidak memiliki izin baca"}
			}
			return AccessResult{Allow: true, Reason: "Read Access Granted"}
		}
	}

	// Default Deny
	return AccessResult{Allow: false, Reason: "Akses ditolak secara default"}
}