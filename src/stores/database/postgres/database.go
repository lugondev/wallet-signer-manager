package postgres

import (
	"context"

	"github.com/lugondev/wallet-signer-manager/pkg/errors"

	"github.com/lugondev/wallet-signer-manager/src/infra/log"
	"github.com/lugondev/wallet-signer-manager/src/infra/postgres"
	"github.com/lugondev/wallet-signer-manager/src/stores/database"
)

type Database struct {
	logger log.Logger
	client postgres.Client
}

var _ database.Database = &Database{}

func New(logger log.Logger, client postgres.Client) *Database {
	return &Database{
		logger: logger,
		client: client,
	}
}

func (db *Database) Ping(ctx context.Context) error {
	err := db.client.Ping(ctx)
	if err != nil {
		errMessage := "database connection error"
		db.logger.WithError(err).Error(errMessage)
		return errors.DependencyFailureError(errMessage)
	}

	return nil
}

func (db *Database) Wallets(storeID string) database.Wallets {
	return NewWallets(storeID, db.client, db.logger.With("store_id", storeID))
}
