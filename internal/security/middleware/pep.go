package middleware

import (
	"backend-zerotrust-skripsi/internal/security/audit"
	"backend-zerotrust-skripsi/internal/security/policies"
	"backend-zerotrust-skripsi/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ZeroTrustPEP adalah Policy Enforcement Point (PEP)
func ZeroTrustPEP(engine *policies.Engine, logger *audit.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// =====================================================
		// 1. CONTEXT COLLECTION
		// =====================================================
		simulatedRole := c.GetHeader("X-User-Role")
		simulatedID := c.GetHeader("X-User-ID")

		resourcePath := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		requestContext := policies.AccessRequest{
			SubjectID:   simulatedID,
			SubjectRole: simulatedRole,
			Resource:    resourcePath,
			Action:      method,
			DeviceIP:    clientIP,
			RequestTime: time.Now(),
		}

		// =====================================================
		// 2. POLICY EVALUATION (PDP)
		// =====================================================
		result := engine.Evaluate(requestContext)

		// =====================================================
		// 3. AUDIT LOGGING
		// =====================================================
		logger.LogAccess(
			simulatedID,
			simulatedRole,
			clientIP,
			resourcePath,
			method,
			result.Allow,
			result.Reason,
		)

		// =====================================================
		// 4. ENFORCEMENT (STANDARD RESPONSE)
		// =====================================================
		if !result.Allow {
			response.Error(
				c,
				http.StatusForbidden,
				"Akses Ditolak oleh Zero Trust Policy",
				gin.H{
					"reason": result.Reason,
					"context": gin.H{
						"role": simulatedRole,
						"ip":   clientIP,
					},
				},
			)
			c.Abort()
			return
		}

		// =====================================================
		// 5. ALLOW → NEXT HANDLER
		// =====================================================
		c.Next()
	}
}
