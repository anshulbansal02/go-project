package tokenfactory

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type TokenFactory[T jwt.Claims] struct {
	signingKey    []byte
	signingMethod jwt.SigningMethod
	parser        jwt.Parser
}

func New[T jwt.Claims](signingMethod jwt.SigningMethod, secret []byte) *TokenFactory[T] {
	return &TokenFactory[T]{
		signingKey:    secret,
		signingMethod: signingMethod,
		parser:        *jwt.NewParser(),
	}
}

func (f *TokenFactory[T]) GenerateToken(claims T) (string, error) {

	token := jwt.NewWithClaims(f.signingMethod, claims)

	return token.SignedString(f.signingKey)
}

func (f *TokenFactory[T]) IsTokenValid(token string, claims *T) (bool, error) {
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

func (f *TokenFactory[T]) GetClaims(token string, dstClaims jwt.Claims) error {
	_, err := f.parser.ParseWithClaims(token, dstClaims, func(t *jwt.Token) (interface{}, error) {
		return f.signingKey, nil
	})

	return err
}
