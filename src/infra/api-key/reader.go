package apikey

import (
	"context"

	"github.com/lugondev/wallet-signer-manager/src/auth/entities"
)

//go:generate mockgen -source=reader.go -destination=mock/reader.go -package=mock

// Reader reads manifests from filesystem
type Reader interface {
	Load(ctx context.Context) (map[string]*entities.UserClaims, error)
}
