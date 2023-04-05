package wallets

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/lugondev/signer-key-manager/src/auth/entities"

	"github.com/lugondev/signer-key-manager/pkg/errors"
	"github.com/lugondev/signer-key-manager/src/stores/database"
)

func (c Connector) Destroy(ctx context.Context, pubkey string) error {
	logger := c.logger.With("pubkey", pubkey)
	logger.Debug("destroying wallet")

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionDestroy, Resource: entities.ResourceWallets})
	if err != nil {
		return err
	}

	acc, err := c.db.GetDeleted(ctx, pubkey)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Wallets) error {
		err = dbtx.Purge(ctx, pubkey)
		if err != nil {
			return err
		}

		err = c.store.Destroy(ctx, hexutil.Encode(acc.CompressedPublicKey))
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("pubkey was permanently deleted")
	return nil
}
