package integration

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type KeycloakClientImpl struct {
	// publicKey, issuer, audience и др. можно добавить в конфиг
}

func NewKeycloakClient() *KeycloakClientImpl {
	return &KeycloakClientImpl{}
}

func (k *KeycloakClientImpl) ValidateToken(ctx context.Context, tokenStr string) (userID, email string, roles []string, err error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", "", nil, errors.New("invalid token format")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", nil, errors.New("invalid claims")
	}
	uid, _ := claims["sub"].(string)
	email, _ = claims["email"].(string)
	roles = extractRoles(claims)
	return uid, email, roles, nil
}

func extractRoles(claims jwt.MapClaims) []string {
	var roles []string
	if realm, ok := claims["realm_access"].(map[string]interface{}); ok {
		if rs, ok := realm["roles"].([]interface{}); ok {
			for _, r := range rs {
				if s, ok := r.(string); ok {
					roles = append(roles, s)
				}
			}
		}
	}
	return roles
}
