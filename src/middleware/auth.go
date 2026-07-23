package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gandli/yctf/utils"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"unauthorized","message":"Authorization header required"}`, http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				http.Error(w, `{"error":"unauthorized","message":"Invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			claims, err := utils.ValidateToken(tokenString, secret)
			if err != nil {
				http.Error(w, `{"error":"unauthorized","message":"Invalid token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RBACMiddleware(secret string, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"unauthorized","message":"Authorization header required"}`, http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := utils.ValidateToken(tokenString, secret)
			if err != nil {
				http.Error(w, `{"error":"unauthorized","message":"Invalid token"}`, http.StatusUnauthorized)
				return
			}

			roleAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				http.Error(w, `{"error":"forbidden","message":"Insufficient permissions"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
