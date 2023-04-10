package wallets

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/lugondev/wallet-signer-manager/src/stores/database/models"

	"github.com/lugondev/wallet-signer-manager/pkg/errors"
	authentities "github.com/lugondev/wallet-signer-manager/src/auth/entities"

	"github.com/lugondev/wallet-signer-manager/src/stores/entities"
)

func (c Connector) Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.Wallet, error) {
	logger := c.logger.With("id", id)
	logger.Debug("creating wallet")

	err := c.authorizator.CheckPermission(&authentities.Operation{Action: authentities.ActionWrite, Resource: authentities.ResourceWallets})
	if err != nil {
		return nil, err
	}

	wallet, err := c.store.Create(ctx, id, attr)
	if err != nil && errors.IsAlreadyExistsError(err) {
		wallet, err = c.store.Get(ctx, id)
	}
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Add(ctx, models.NewWalletFromKey(wallet, attr))
	if err != nil {
		return nil, err
	}

	logger.With("pubkey", hexutil.Encode(acc.CompressedPublicKey)).Info("wallet created successfully")
	return acc, nil
}
