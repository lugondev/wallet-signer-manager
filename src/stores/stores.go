package stores

import (
	"context"

	auth "github.com/lugondev/wallet-signer-manager/src/auth/entities"
)

type Stores interface {
	// CreateWallet creates a wallet store
	CreateWallet(_ context.Context, name string, allowedTenants []string, userInfo *auth.UserInfo) error

	// ImportWallets import wallets from the vault into a wallet store
	ImportWallets(ctx context.Context, name string, userInfo *auth.UserInfo) error

	// Wallet get wallets store by name
	Wallet(ctx context.Context, storeName string, userInfo *auth.UserInfo) (WalletStore, error)

	// List stores
	List(ctx context.Context, storeType string, userInfo *auth.UserInfo) ([]string, error)

	// ListAllWallets list all wallets from all stores
	ListAllWallets(ctx context.Context, userInfo *auth.UserInfo) ([]string, error)
}
