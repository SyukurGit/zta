package middleware

import (
	"backend-zerotrust-skripsi/internal/security/audit"
	"backend-zerotrust-skripsi/internal/security/policies"
	"backend-zerotrust-skripsi/internal/security/token" // Import token
	"backend-zerotrust-skripsi/pkg/response"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Tambahkan parameter jwtService *token.JWTService
func ZeroTrustPEP(engine *policies.Engine, logger *audit.Logger, jwtService *token.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Missing Authorization Header", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Panggil ValidateToken dari struct service
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			logger.LogAccess("unknown", "unknown", c.ClientIP(), c.Request.URL.Path, c.Request.Method, false, "Invalid or Expired Token")
			response.Error(c, http.StatusUnauthorized, "Token Invalid atau Expired", err.Error())
			c.Abort()
			return
		}

		currentUserID := strconv.Itoa(int(claims.UserID))
		currentUserRole := claims.Role

		// ... (Sisa kode policy evaluation & logging biarkan sama seperti sebelumnya) ...
		
		// Setup Context untuk Policy Engine
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		
		resourcePath := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		requestContext := policies.AccessRequest{
			SubjectID:   currentUserID,
			SubjectRole: currentUserRole,
			Resource:    resourcePath,
			Action:      method,
			DeviceIP:    clientIP,
			RequestTime: time.Now(),
		}

		result := engine.Evaluate(requestContext)

		logger.LogAccess(
			currentUserID,
			currentUserRole,
			clientIP,
			resourcePath,
			method,
			result.Allow,
			result.Reason,
		)

		if !result.Allow {
			response.Error(
				c,
				http.StatusForbidden,
				"Akses Ditolak oleh Zero Trust Policy",
				gin.H{
					"reason": result.Reason,
					"context": gin.H{
						"role": currentUserRole,
						"ip":   clientIP,
					},
				},
			)
			c.Abort()
			return
		}

		c.Next()
	}
}