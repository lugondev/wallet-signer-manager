package wallets

import (
	"context"

	"github.com/lugondev/signer-key-manager/src/auth/entities"
)

func (c Connector) Encrypt(ctx context.Context, pubkey string, data []byte) ([]byte, error) {
	logger := c.logger.With("pubkey", pubkey)

	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionEncrypt, Resource: entities.ResourceEthAccount})
	if err != nil {
		return nil, err
	}

	acc, err := c.db.Get(ctx, pubkey)
	if err != nil {
		return nil, err
	}

	result, err := c.store.Encrypt(ctx, acc.KeyID, data)
	if err != nil {
		return nil, err
	}

	logger.Debug("data encrypted successfully")
	return result, nil
}