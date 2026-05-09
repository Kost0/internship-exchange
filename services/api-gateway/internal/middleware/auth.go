package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type contextKey string

const (
	ContextUserID   contextKey = "user_id"
	ContextUserRole contextKey = "user_role"
)

type AuthMiddleware struct {
	jwtSecret string
	redis     *redis.Client
}

func NewAuthMiddleware(jwtSecret string, redis *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret, redis: redis}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.extractAndValidate(r)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)

			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)

			return
		}

		userID, _ := claims["sub"].(string)
		role, _ := claims["role"].(string)

		ctx := context.WithValue(r.Context(), ContextUserID, userID)
		ctx = context.WithValue(ctx, ContextUserRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole, ok := r.Context().Value(ContextUserRole).(string)
		if !ok || userRole != role {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) extractAndValidate(r *http.Request) (*jwt.Token, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return nil, jwt.ErrTokenMalformed
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, jwt.ErrTokenMalformed
	}

	tokenStr := parts[1]

	cached, err := m.redis.Get(r.Context(), "jwt:"+tokenStr).Result()
	if err == nil && cached == "valid" {
		return jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.jwtSecret), nil
		})
	}

	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(m.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	m.redis.Set(r.Context(), "jwt:"+tokenStr, "valid", 5*time.Minute)

	return token, nil
}

func GetUserID(ctx context.Context) string {
	v, _ := ctx.Value(ContextUserID).(string)

	return v
}

func GetUserRole(ctx context.Context) string {
	v, _ := ctx.Value(ContextUserRole).(string)

	return v
}
