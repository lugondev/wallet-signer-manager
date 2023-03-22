package wallets

import (
	"context"

	"github.com/lugondev/signer-key-manager/src/stores/database/models"

	"github.com/lugondev/signer-key-manager/pkg/errors"
	authentities "github.com/lugondev/signer-key-manager/src/auth/entities"

	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

func (c Connector) Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.Wallet, error) {
	logger := c.logger.With("id", id)
	logger.Debug("creating wallet")

	err := c.authorizator.CheckPermission(&authentities.Operation{Action: authentities.ActionWrite, Resource: authentities.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	key, err := c.store.Create(ctx, id, attr)
	if err != nil && errors.IsAlreadyExistsError(err) {
		key, err = c.store.Get(ctx, id)
	}
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Add(ctx, models.NewWalletFromKey(key, attr))
	if err != nil {
		return nil, err
	}

	logger.With("pubkey", acc.CompressedPublicKey, "key_id", acc.KeyID).Info("wallet created successfully")
	return acc, nil
}
