package jwt

import (
	"context"

	"github.com/lugondev/wallet-signer-manager/src/auth/entities"
)

//go:generate mockgen -source=validator.go -destination=mock/validator.go -package=mock

type Validator interface {
	ValidateToken(ctx context.Context, token string) (interface{}, error)
	ParseClaims(interface{}) (*entities.UserClaims, error)
}
