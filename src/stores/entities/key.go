package entities

import "github.com/lugondev/wallet-signer-manager/src/entities"

// Key public part of a key
type Key struct {
	ID        string
	PublicKey []byte
	Algo      *entities.Algorithm
	Metadata  *Metadata
	Tags      map[string]string
}

func (k *Key) IsETHAccount() bool {
	return k.Algo.EllipticCurve == entities.Secp256k1 && k.Algo.Type == entities.Ecdsa
}
