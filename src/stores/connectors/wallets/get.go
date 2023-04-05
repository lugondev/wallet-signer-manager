package wallets

import (
	"context"

	authentities "github.com/lugondev/signer-key-manager/src/auth/entities"

	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

func (c Connector) Get(ctx context.Context, pubkey string) (*entities.Wallet, error) {
	logger := c.logger.With("pubkey", pubkey)

	err := c.authorizator.CheckPermission(&authentities.Operation{Action: authentities.ActionRead, Resource: authentities.ResourceWallets})
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Get(ctx, pubkey)
	if err != nil {
		return nil, err
	}

	logger.Debug("wallet retrieved successfully")
	return acc, nil
}

func (c Connector) GetDeleted(ctx context.Context, pubkey string) (*entities.Wallet, error) {
	logger := c.logger.With("pubkey", pubkey)

	err := c.authorizator.CheckPermission(&authentities.Operation{Action: authentities.ActionRead, Resource: authentities.ResourceWallets})
	if err != nil {
		return nil, err
	}

	acc, err := c.db.GetDeleted(ctx, pubkey)
	if err != nil {
		return nil, err
	}

	logger.Debug("deleted ethereum account retrieved successfully")
	return acc, nil
}
