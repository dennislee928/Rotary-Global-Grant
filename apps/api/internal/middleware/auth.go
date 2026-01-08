package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/service"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, vo.ErrorVO{
				Code:    "UNAUTHORIZED",
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check for Bearer token (JWT)
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authSvc.ValidateToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, vo.ErrorVO{
					Code:    "UNAUTHORIZED",
					Message: "Invalid or expired token",
				})
				c.Abort()
				return
			}

			// Set user info in context
			userID, _ := uuid.Parse(claims.UserID)
			c.Set("userID", userID)
			c.Set("userEmail", claims.Email)
			c.Set("userRole", claims.Role)
			c.Next()
			return
		}

		// Check for API key
		if strings.HasPrefix(authHeader, "ApiKey ") {
			apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
			user, key, err := authSvc.ValidateAPIKey(c.Request.Context(), apiKey)
			if err != nil {
				c.JSON(http.StatusUnauthorized, vo.ErrorVO{
					Code:    "UNAUTHORIZED",
					Message: "Invalid API key",
				})
				c.Abort()
				return
			}

			// Set user info in context
			c.Set("userID", user.ID)
			c.Set("userEmail", user.Email)
			c.Set("userRole", user.Role)
			c.Set("apiKeyScopes", key.Scopes)
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, vo.ErrorVO{
			Code:    "UNAUTHORIZED",
			Message: "Invalid authorization header format",
		})
		c.Abort()
	}
}

// OptionalAuthMiddleware creates an optional authentication middleware
// It sets user info if auth is present, but doesn't require it
func OptionalAuthMiddleware(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check for Bearer token (JWT)
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authSvc.ValidateToken(tokenString)
			if err == nil {
				userID, _ := uuid.Parse(claims.UserID)
				c.Set("userID", userID)
				c.Set("userEmail", claims.Email)
				c.Set("userRole", claims.Role)
			}
		}

		// Check for API key
		if strings.HasPrefix(authHeader, "ApiKey ") {
			apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
			user, key, err := authSvc.ValidateAPIKey(c.Request.Context(), apiKey)
			if err == nil {
				c.Set("userID", user.ID)
				c.Set("userEmail", user.Email)
				c.Set("userRole", user.Role)
				c.Set("apiKeyScopes", key.Scopes)
			}
		}

		c.Next()
	}
}

// RoleMiddleware creates a role-based authorization middleware
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, vo.ErrorVO{
				Code:    "UNAUTHORIZED",
				Message: "Authentication required",
			})
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, vo.ErrorVO{
			Code:    "FORBIDDEN",
			Message: "Insufficient permissions",
		})
		c.Abort()
	}
}

// ScopeMiddleware creates a scope-based authorization middleware for API keys
func ScopeMiddleware(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopesRaw, exists := c.Get("apiKeyScopes")
		if !exists {
			// JWT token, no scopes needed
			c.Next()
			return
		}

		scopes := scopesRaw.([]string)
		for _, scope := range scopes {
			if scope == requiredScope || scope == "admin" {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, vo.ErrorVO{
			Code:    "FORBIDDEN",
			Message: "API key missing required scope: " + requiredScope,
		})
		c.Abort()
	}
}
