package database

import (
	"context"

	"github.com/lugondev/wallet-signer-manager/src/stores/entities"
)

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	Ping(ctx context.Context) error
	Wallets(storeID string) Wallets
}

type Wallets interface {
	RunInTransaction(ctx context.Context, persistFunc func(dbtx Wallets) error) error
	Get(ctx context.Context, pubkey string) (*entities.Wallet, error)
	GetDeleted(ctx context.Context, pubkey string) (*entities.Wallet, error)
	GetAll(ctx context.Context) ([]*entities.Wallet, error)
	GetAllDeleted(ctx context.Context) ([]*entities.Wallet, error)
	SearchAddresses(ctx context.Context, isDeleted bool, limit, offset uint64) ([]string, error)
	Add(ctx context.Context, wallet *entities.Wallet) (*entities.Wallet, error)
	Update(ctx context.Context, wallet *entities.Wallet) (*entities.Wallet, error)
	Delete(ctx context.Context, pubkey string) error
	Restore(ctx context.Context, pubkey string) error
	Purge(ctx context.Context, pubkey string) error
}
