package wallets

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/lugondev/wallet-signer-manager/src/auth/entities"

	"github.com/lugondev/wallet-signer-manager/pkg/errors"
	"github.com/lugondev/wallet-signer-manager/src/stores/database"
)

func (c Connector) Delete(ctx context.Context, pubkey string) error {
	logger := c.logger.With("pubkey", pubkey)
	logger.Debug("deleting ethereum account")

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionDelete, Resource: entities.ResourceWallets})
	if err != nil {
		return err
	}

	acc, err := c.db.Get(ctx, pubkey)
	if err != nil {
		return err
	}

	err = c.db.RunInTransaction(ctx, func(dbtx database.Wallets) error {
		err = dbtx.Delete(ctx, pubkey)
		if err != nil {
			return err
		}

		err = c.store.Delete(ctx, hexutil.Encode(acc.CompressedPublicKey))
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support deleting, we only delete in DB
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	logger.Info("ethereum account deleted successfully")
	return nil
}
