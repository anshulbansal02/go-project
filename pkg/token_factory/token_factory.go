package tokenfactory

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type TokenFactory struct {
	signingKey    []byte
	signingMethod jwt.SigningMethod
	parser        jwt.Parser
}

func New(signingMethod jwt.SigningMethod, secret []byte) *TokenFactory {
	return &TokenFactory{
		signingKey:    secret,
		signingMethod: signingMethod,
		parser:        *jwt.NewParser(),
	}
}

func (f *TokenFactory) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(f.signingMethod, claims)

	return token.SignedString(f.signingKey)
}

func (f *TokenFactory) IsTokenValid(token string, claims *jwt.Claims) (bool, error) {
	t, err := f.parser.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return t, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, jwt.ErrTokenExpired) {
			return false, nil
		}
		return false, err
	}

	return t.Valid, nil
}

func (f *TokenFactory) GetClaims(token string, dstClaims *jwt.Claims) error {
	_, err := f.parser.ParseWithClaims(token, *dstClaims, func(t *jwt.Token) (interface{}, error) {
		return t, nil
	})

	return err
}
