package policies

import (
	"backend-zerotrust-skripsi/config"
	"strings"
	// "time"
)

// Engine adalah Policy Decision Point (PDP)
type Engine struct {
	Cfg *config.AppConfig // Inject konfigurasi aplikasi
}

// NewEngine membuat instance PDP dengan config
func NewEngine(cfg *config.AppConfig) *Engine {
	return &Engine{Cfg: cfg}
}

// Evaluate adalah fungsi inti evaluasi kebijakan Zero Trust
func (e *Engine) Evaluate(req AccessRequest) AccessResult {

	// =========================================================
	// DEFAULT DENY — Zero Trust Principle
	// =========================================================

	// =========================================================
	// POLICY: Finance Reports
	// =========================================================
	if req.Resource == "/api/v1/finance/reports" {

		// -----------------------------------------------------
		// POLICY 1: Role-based Access
		// -----------------------------------------------------
		if req.SubjectRole != "manager" && req.SubjectRole != "auditor" {
			return AccessResult{
				Allow:  false,
				Reason: "Role tidak diizinkan mengakses data keuangan",
			}
		}

		// -----------------------------------------------------
		// POLICY 2: Time-based Access (Config-driven)
		// -----------------------------------------------------
		hour := req.RequestTime.Hour()
		if hour < e.Cfg.WorkStartHour || hour > e.Cfg.WorkEndHour {
			return AccessResult{
				Allow:  false,
				Reason: "Akses di luar jam operasional yang ditentukan",
			}
		}

		// -----------------------------------------------------
		// POLICY 3: Network / IP-based Access (Config-driven)
		// -----------------------------------------------------
		isSafeNetwork :=
			strings.HasPrefix(req.DeviceIP, e.Cfg.SafeNetwork) ||
				req.DeviceIP == "::1"

		if !isSafeNetwork {
			return AccessResult{
				Allow:  false,
				Reason: "Akses data sensitif harus dari jaringan kantor yang aman",
			}
		}

		// -----------------------------------------------------
		// ALL POLICIES PASSED
		// -----------------------------------------------------
		return AccessResult{
			Allow:  true,
			Reason: "Policy Check Passed",
		}
	}

	// =========================================================
	// UNKNOWN RESOURCE → DENY
	// =========================================================
	return AccessResult{
		Allow:  false,
		Reason: "Resource tidak didefinisikan dalam policy",
	}
}
