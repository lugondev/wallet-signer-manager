package vaults

import (
	"context"

	auth "github.com/lugondev/wallet-signer-manager/src/auth/entities"
	"github.com/lugondev/wallet-signer-manager/src/entities"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock

type Vaults interface {
	// CreateHashicorp creates a Hashicorp Vault client
	CreateHashicorp(ctx context.Context, name string, config *entities.HashicorpConfig, allowedTenants []string, userInfo *auth.UserInfo) error

	// Get gets a valut by name
	Get(ctx context.Context, name string, userInfo *auth.UserInfo) (*entities.Vault, error)
}
