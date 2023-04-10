package stores

import (
	"context"
	"github.com/lugondev/wallet-signer-manager/src/stores/entities"
)

type WalletStore interface {
	// Create creates a wallet
	Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.Wallet, error)

	// Import imports an externally created wallet

	Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.Wallet, error)

	// Get gets a wallet
	Get(ctx context.Context, pubkey string) (*entities.Wallet, error)

	// List lists all wallet addresses
	List(ctx context.Context, limit, offset uint64) ([]string, error)

	// Update updates wallet attributes
	Update(ctx context.Context, pubkey string, attr *entities.Attributes) (*entities.Wallet, error)

	// Delete deletes an account temporarily, by using Restore the account can be restored
	Delete(ctx context.Context, pubkey string) error

	// GetDeleted Gets a deleted wallets
	GetDeleted(ctx context.Context, pubkey string) (*entities.Wallet, error)

	// ListDeleted lists all deleted wallets
	ListDeleted(ctx context.Context, limit, offset uint64) ([]string, error)

	// Restore restores a previously deleted wallet
	Restore(ctx context.Context, pubkey string) error

	// Destroy destroys (purges) a wallet permanently
	Destroy(ctx context.Context, pubkey string) error

	// Sign signs data using the specified wallet (not exposed in the API)
	Sign(ctx context.Context, pubkey string, data []byte) ([]byte, error)
}
