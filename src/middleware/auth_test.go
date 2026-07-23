package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gandli/yctf/utils"
)

func TestAuthMiddleware(t *testing.T) {
	secret := "test-secret-min-32-bytes-long!!"
	token, _ := utils.GenerateToken("user-123", "player", secret)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/challenges", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := AuthMiddleware(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify context has claims
		claims, ok := r.Context().Value("claims").(*utils.Claims)
		if !ok {
			t.Error("claims not in context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if claims.UserID != "user-123" {
			t.Errorf("expected user-123, got %s", claims.UserID)
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthMiddlewareMissingToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/challenges", nil)
	w := httptest.NewRecorder()

	handler := AuthMiddleware("test-secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/challenges", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")
	w := httptest.NewRecorder()

	handler := AuthMiddleware("test-secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddlewareExpiredToken(t *testing.T) {
	secret := "test-secret-min-32-bytes-long!!"
	token, _ := utils.GenerateTokenWithExpiry("user-123", "player", secret, -1) // expired

	req := httptest.NewRequest(http.MethodGet, "/api/v1/challenges", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := AuthMiddleware(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 for expired token, got %d", w.Code)
	}
}

func TestAuthMiddlewareWrongSecret(t *testing.T) {
	token, _ := utils.GenerateToken("user-123", "player", "correct-secret-min-32-bytes!!")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/challenges", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := AuthMiddleware("wrong-secret-min-32-bytes!!")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 for wrong secret, got %d", w.Code)
	}
}

func TestRBACMiddleware(t *testing.T) {
	secret := "test-secret-min-32-bytes-long!!"

	adminToken, _ := utils.GenerateToken("admin-1", "admin", secret)
	playerToken, _ := utils.GenerateToken("player-1", "player", secret)

	// Admin should access admin endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w := httptest.NewRecorder()

	handler := RBACMiddleware(secret, "admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("admin should access admin endpoint, got %d", w.Code)
	}

	// Player should NOT access admin endpoint
	req = httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+playerToken)
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("player should not access admin endpoint, got %d", w.Code)
	}
}

func TestRBACMiddlewareAuthor(t *testing.T) {
	secret := "test-secret-min-32-bytes-long!!"
	authorToken, _ := utils.GenerateToken("author-1", "author", secret)

	// Author should access author endpoint
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", nil)
	req.Header.Set("Authorization", "Bearer "+authorToken)
	w := httptest.NewRecorder()

	handler := RBACMiddleware(secret, "admin", "author")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("author should access author endpoint, got %d", w.Code)
	}
}
