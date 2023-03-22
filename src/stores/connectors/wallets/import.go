package wallets

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/lugondev/signer-key-manager/src/stores/database/models"

	"github.com/lugondev/signer-key-manager/pkg/errors"

	authentities "github.com/lugondev/signer-key-manager/src/auth/entities"

	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

func (c Connector) Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.Wallet, error) {
	logger := c.logger.With("id", id)
	logger.Debug("importing wallet")

	if len(privKey) == 0 {
		errMessage := "private key must be provided"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	logger.Debug("checking permissions")
	err := c.authorizator.CheckPermission(&authentities.Operation{Action: authentities.ActionWrite, Resource: authentities.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	logger.Debug("store importing wallet")
	key, err := c.store.Import(ctx, id, privKey, attr)
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

	logger.With("pubkey", hexutil.Encode(acc.CompressedPublicKey)).Info("wallet imported successfully")
	return acc, nil
}
