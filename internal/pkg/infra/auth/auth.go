package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

const DefaultJWTSignInMethod = "RS256"

var (
	ErrorKeyNotFound          = errors.New("action key not found in key store")
	ErrorSignInMethodNotFound = errors.New("default JWT sign in method not found")
	ErrorKIDNotFound          = errors.New("kid Header not found in token")
	ErrorKIDInvalidString     = errors.New("kid Header value must be string")
	ErrorTokenAuthority       = errors.New("token authority cannot be confirmed")
)

// KeyStore declares interface for a set of methods to retrieve
// private and public keys.
type KeyStore interface {
	GetPrivateKey(keyId string) (*rsa.PrivateKey, error)
	GetPublicKey(keyId string) (*rsa.PublicKey, error)
}

// AuthenticationContext is used in operations to generate token with user Claims.
type AuthenticationContext struct {
	currentKeyId string
	keyStore     KeyStore
	method       jwt.SigningMethod
	getKey       func(token *jwt.Token) (interface{}, error)
	parser       jwt.Parser
}

// NewAuthenticationContext constructs an AuthenticationContext for authentication and authorization processes.
func NewAuthenticationContext(activeKeyId string, keyStore KeyStore) (*AuthenticationContext, error) {
	_, err := keyStore.GetPrivateKey(activeKeyId)
	if err != nil {
		return nil, ErrorKeyNotFound
	}

	method := jwt.GetSigningMethod(DefaultJWTSignInMethod)
	if method == nil {
		return nil, ErrorSignInMethodNotFound
	}

	getKey := func(token *jwt.Token) (interface{}, error) {
		kid, found := token.Header["kid"]
		if !found {
			return nil, ErrorKIDNotFound
		}

		keyId, ok := kid.(string)
		if !ok {
			return nil, ErrorKIDInvalidString
		}

		return keyStore.GetPublicKey(keyId)
	}

	parser := jwt.Parser{
		ValidMethods: []string{DefaultJWTSignInMethod},
	}

	authContext := AuthenticationContext{
		currentKeyId: activeKeyId,
		keyStore:     keyStore,
		method:       method,
		getKey:       getKey,
		parser:       parser,
	}

	return &authContext, nil
}

// GenerateToken returns a signed JWT string with user Claims.
func (ctx *AuthenticationContext) GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(ctx.method, claims)
	token.Header["kid"] = ctx.currentKeyId

	pvKey, err := ctx.keyStore.GetPrivateKey(ctx.currentKeyId)
	if err != nil {
		return "", ErrorKIDNotFound
	}

	signedToken, err := token.SignedString(pvKey)
	if err != nil {
		return "", fmt.Errorf("signing token error: %w", err)
	}

	return signedToken, nil
}

// ReadClaimsFromToken returns user Claims that stored in a target JWT token.
// Token should be signed with a correct key.
func (ctx *AuthenticationContext) ReadClaimsFromToken(signedToken string) (Claims, error) {
	var claims Claims

	token, err := ctx.parser.ParseWithClaims(signedToken, &claims, ctx.getKey)
	if err != nil {
		return Claims{}, fmt.Errorf("error during parse of signed token: %w", err)
	}

	if !token.Valid {
		return Claims{}, ErrorTokenAuthority
	}

	return claims, nil
}
