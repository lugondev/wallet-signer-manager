package wallets

import (
	"context"
	authtypes "github.com/lugondev/wallet-signer-manager/src/auth/entities"

	"github.com/ethereum/go-ethereum/crypto"
)

func (c Connector) Sign(ctx context.Context, pubkey string, data []byte) ([]byte, error) {
	logger := c.logger.With("pubkey", pubkey)

	if len(data) != 32 {
		data = crypto.Keccak256(data)
	}
	signature, err := c.sign(ctx, pubkey, data)
	if err != nil {
		return nil, err
	}

	logger.Debug("signed payload successfully")
	return signature, nil
}

func (c Connector) sign(ctx context.Context, pubkey string, data []byte) ([]byte, error) {
	logger := c.logger.With("pubkey", pubkey)

	err := c.authorizator.CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceWallets})
	if err != nil {
		return nil, err
	}

	if _, err := c.db.Get(ctx, pubkey); err != nil {
		logger.WithError(err).Error("failed to get account")
		return nil, err
	}

	signature, err := c.store.Sign(ctx, pubkey, data)
	if err != nil {
		logger.WithError(err).Error("failed to sign payload")
		return nil, err
	}

	return signature, nil
}
