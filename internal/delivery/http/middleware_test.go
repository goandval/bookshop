package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourorg/bookshop/internal/mocks"
)

func TestAuthMiddleware_JWTAuth_Success(t *testing.T) {
	keycloak := new(mocks.KeycloakClient)
	keycloak.On("ValidateToken", mock.Anything, "valid-token").Return("user-1", "user@ex.com", []string{"user"}, nil)
	mw := NewAuthMiddleware(keycloak)

	called := false
	h := mw.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		ctx := r.Context()
		assert.Equal(t, "user-1", ctx.Value("userID"))
		assert.Equal(t, "user@ex.com", ctx.Value("email"))
		assert.Contains(t, ctx.Value("roles").([]string), "user")
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	assert.True(t, called)
	assert.Equal(t, 200, rw.Code)
}

func TestAuthMiddleware_RequireRole_Forbidden(t *testing.T) {
	keycloak := new(mocks.KeycloakClient)
	mw := NewAuthMiddleware(keycloak)

	h := mw.RequireRole("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), "roles", []string{"user"})
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req.WithContext(ctx))
	assert.Equal(t, 403, rw.Code)
}
