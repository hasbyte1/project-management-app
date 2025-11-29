package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hasbyte1/project-management-app/internal/models"
	"github.com/hasbyte1/project-management-app/pkg/jwt"
)

type contextKey string

const UserIDKey contextKey = "userID"

type AuthMiddleware struct {
	tokenManager *jwt.TokenManager
}

func NewAuthMiddleware(tokenManager *jwt.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenManager}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondJSON(w, http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Error:   "Missing authorization header",
			})
			return
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondJSON(w, http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Error:   "Invalid authorization header format",
			})
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.tokenManager.ValidateAccessToken(token)
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, models.ErrorResponse{
				Success: false,
				Error:   "Invalid or expired token",
			})
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}
