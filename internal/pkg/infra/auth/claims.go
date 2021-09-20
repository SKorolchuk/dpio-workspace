package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents the authorization claims in JWT format.
type Claims struct {
	jwt.StandardClaims
	Roles []string `json:"roles"`
}

// AuthorizeCheck confirms existence of authorized roles in token.
func (token Claims) AuthorizeCheck(allowedRoles ...string) bool {
	for _, tokenRole := range token.Roles {
		for _, allowedRole := range allowedRoles {
			if tokenRole == allowedRole {
				return true
			}
		}
	}

	return false
}

// webContextKeyType represents the identifier of context key for user claims
type webContextKeyType int

// webContextKey is used to store and handle a Claims value in context.Context.
const webContextKey webContextKeyType = 356221

// SetClaims save user Claims in the web context.
func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, webContextKey, claims)
}

// GetClaims returns user Claims from the web context.
func GetClaims(ctx context.Context) (Claims, error) {
	value, found := ctx.Value(webContextKey).(Claims)
	if !found {
		return Claims{}, errors.New("user claims object missed from the web context")
	}

	return value, nil
}
