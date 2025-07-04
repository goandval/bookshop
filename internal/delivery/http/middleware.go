package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/yourorg/bookshop/internal/integration"
	"golang.org/x/exp/slog"
)

type AuthMiddleware struct {
	keycloak integration.KeycloakClient
	Logger   *slog.Logger
}

func NewAuthMiddleware(keycloak integration.KeycloakClient, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{keycloak: keycloak, Logger: logger}
}

func (a *AuthMiddleware) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			a.Logger.Warn("missing token in request")
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		userID, email, roles, err := a.keycloak.ValidateToken(r.Context(), token)
		if err != nil {
			a.Logger.Warn("invalid token", "err", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "roles", roles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := r.Context().Value("roles").([]string)
			if !ok {
				a.Logger.Warn("no roles in context")
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			for _, roleVal := range roles {
				if roleVal == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			a.Logger.Warn("forbidden: required role missing", "required", role, "roles", roles)
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

func extractBearerToken(header string) string {
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}
