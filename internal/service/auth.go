package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/sudeeya/avito-assignment/internal/config"
)

var _ Auth = (*AuthService)(nil)

type AuthService struct {
	tokenString string
}

var errWrongToken = errors.New("wrong token")

func newAuthService(cfg config.ServerConfig) (*AuthService, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err := token.SignedString([]byte(cfg.ServerSecretKey))
	if err != nil {
		return nil, fmt.Errorf("signing token: %w", err)
	}

	return &AuthService{
		tokenString: tokenString,
	}, nil
}

// IssueToken implements Auth.
func (a *AuthService) IssueToken(ctx context.Context) (string, error) {
	return a.tokenString, nil
}

// VerifyToken implements Auth.
func (a *AuthService) VerifyToken(ctx context.Context, tokenString string) error {
	if tokenString != a.tokenString {
		return errWrongToken
	}

	return nil
}
