package auth

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Auth0Validator struct {
	domain   string
	audience string
	jwks     *JWKSClient
}

func NewAuth0Validator(domain, audience string) *Auth0Validator {
	return &Auth0Validator{
		domain:   domain,
		audience: audience,
		jwks:     NewJWKSClient(domain),
	}
}

func (v *Auth0Validator) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid header not found")
		}

		return v.jwks.GetKey(ctx, kid)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	expectedIssuer := fmt.Sprintf("https://%s/", v.domain)
	iss, _ := claims.GetIssuer()
	if iss != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %s", expectedIssuer, iss)
	}

	if !v.hasValidAudience(claims) {
		return nil, fmt.Errorf("invalid audience")
	}

	return claims, nil
}

func (v *Auth0Validator) hasValidAudience(claims *Claims) bool {
	aud, _ := claims.GetAudience()
	return slices.Contains(aud, v.audience)
}
