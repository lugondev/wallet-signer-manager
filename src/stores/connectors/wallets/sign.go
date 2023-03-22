package wallets

import (
	"bytes"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	authtypes "github.com/lugondev/signer-key-manager/src/auth/entities"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lugondev/signer-key-manager/pkg/errors"
)

var (
	secp256k1N, _     = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN, _ = new(big.Int).SetString("7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0", 16)
)

func (c Connector) Sign(ctx context.Context, pubkey string, data []byte) ([]byte, error) {
	logger := c.logger.With("pubkey", pubkey)

	signature, err := c.sign(ctx, pubkey, crypto.Keccak256(data))
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

	acc, err := c.db.Get(ctx, pubkey)
	if err != nil {
		return nil, err
	}

	signature, err := c.store.Sign(ctx, hexutil.Encode(acc.CompressedPublicKey), data)
	if err != nil {
		return nil, err
	}
	signature = malleabilityECDSASignature(signature)

	// Recover the recID, please read: http://coders-errand.com/ecrecover-signature-verification-ethereum/
	for _, recID := range []byte{0, 1} {
		appendedSignature := append(signature, recID)
		recoveredPubKey, err := crypto.SigToPub(data, appendedSignature)
		if err != nil {
			errMessage := "failed to recover public key candidate with appended recID"
			logger.WithError(err).Error(errMessage, "recID", recID, "signature", hexutil.Encode(signature))
			return nil, errors.CryptoOperationError(errMessage)
		}

		if bytes.Equal(crypto.FromECDSAPub(recoveredPubKey), acc.PublicKey) {
			return appendedSignature, nil
		}
	}

	errMessage := "failed to recover public key candidate"
	logger.Error(errMessage)
	return nil, errors.DependencyFailureError(errMessage)
}

// Azure generates ECDSA signature that does not prevent malleability
// A malleable signature can be transformed into a new and valid one for a different message or key.
// https://docs.microsoft.com/en-us/azure/key-vault/keys/about-keys-details
// More info about the issue: http://coders-errand.com/malleability-ecdsa-signatures/
// More info about the fix: https://en.bitcoin.it/wiki/BIP_0062
func malleabilityECDSASignature(signature []byte) []byte {
	S := new(big.Int).SetBytes(signature[32:])
	if S.Cmp(secp256k1halfN) <= 0 {
		return signature
	}

	S2 := new(big.Int).Sub(secp256k1N, S)
	return append(signature[:32], common.LeftPadBytes(S2.Bytes(), 32)...)
}
