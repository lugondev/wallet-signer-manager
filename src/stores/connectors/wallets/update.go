package wallets

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"

	authtypes "github.com/lugondev/signer-key-manager/src/auth/entities"

	"github.com/lugondev/signer-key-manager/pkg/errors"
	"github.com/lugondev/signer-key-manager/src/stores/database"
	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

func (c Connector) Update(ctx context.Context, pubkey string, attr *entities.Attributes) (*entities.Wallet, error) {
	logger := c.logger.With("pubkey", pubkey)
	logger.Debug("updating wallet")

	err := c.authorizator.CheckPermission(&authtypes.Operation{Action: authtypes.ActionWrite, Resource: authtypes.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Get(ctx, pubkey)
	if err != nil {
		return nil, err
	}
	acc.Tags = attr.Tags

	err = c.db.RunInTransaction(ctx, func(dbtx database.Wallets) error {
		acc, err = dbtx.Update(ctx, acc)
		if err != nil {
			return err
		}

		_, err = c.store.Update(ctx, hexutil.Encode(acc.CompressedPublicKey), attr)
		if err != nil && !errors.IsNotSupportedError(err) {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logger.Info("wallet updated successfully")
	return acc, nil
}
