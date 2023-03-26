package wallets

import (
	"context"
	"math/big"

	authtypes "github.com/lugondev/signer-key-manager/src/auth/entities"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	secp256k1N, _     = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN, _ = new(big.Int).SetString("7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0", 16)
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

	err := c.authorizator.CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount})
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
