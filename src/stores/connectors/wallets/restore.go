package wallets

import (
	"context"

	"github.com/lugondev/signer-key-manager/src/auth/entities"

	"github.com/lugondev/signer-key-manager/pkg/errors"
	"github.com/lugondev/signer-key-manager/src/stores/database"
)

func (c Connector) Restore(ctx context.Context, pubkey string) error {
	logger := c.logger.With("pubkey", pubkey)
	logger.Debug("restoring wallet")

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceEthAccount})
	if err != nil {
		return err
	}

	_, err = c.Get(ctx, pubkey)
	if err == nil {
		return nil
	}

	acc, err := c.db.GetDeleted(ctx, pubkey)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Wallets) error {
		err = dbtx.Restore(ctx, pubkey)
		if err != nil {
			return err
		}

		err = c.store.Restore(ctx, acc.KeyID)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support restoring, we only restore in DB
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("wallet restored successfully")
	return nil
}
