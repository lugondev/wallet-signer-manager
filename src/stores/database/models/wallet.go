package models

import (
	"github.com/status-im/keycard-go/hexutils"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

type Wallet struct {
	tableName struct{} `pg:"wallets"` // nolint:unused,structcheck // reason

	Pubkey              string `pg:",pk"`
	StoreID             string `pg:",pk"`
	PublicKey           []byte
	CompressedPublicKey []byte
	Tags                map[string]string
	Disabled            bool
	CreatedAt           time.Time `pg:"default:now()"`
	UpdatedAt           time.Time `pg:"default:now()"`
	DeletedAt           time.Time `pg:",soft_delete"`
}

func NewWallet(account *entities.Wallet) *Wallet {
	return &Wallet{
		Pubkey:              strings.ToLower(hexutils.BytesToHex(account.CompressedPublicKey)),
		PublicKey:           account.PublicKey,
		CompressedPublicKey: account.CompressedPublicKey,
		Tags:                account.Tags,
		Disabled:            account.Metadata.Disabled,
		CreatedAt:           account.Metadata.CreatedAt,
		UpdatedAt:           account.Metadata.UpdatedAt,
		DeletedAt:           account.Metadata.DeletedAt,
	}
}

func NewWalletFromKey(key *entities.Wallet, attr *entities.Attributes) *entities.Wallet {
	pubKey, _ := crypto.UnmarshalPubkey(key.PublicKey)
	compressedPubkey := crypto.CompressPubkey(pubKey)

	return &entities.Wallet{
		Tags:                attr.Tags,
		Pubkey:              strings.ToLower(hexutils.BytesToHex(compressedPubkey)),
		PublicKey:           key.PublicKey,
		CompressedPublicKey: compressedPubkey,
		Metadata: &entities.Metadata{
			Disabled:  key.Metadata.Disabled,
			CreatedAt: key.Metadata.CreatedAt,
			UpdatedAt: key.Metadata.UpdatedAt,
		},
	}
}

func (w *Wallet) ToEntity() *entities.Wallet {
	return &entities.Wallet{
		Pubkey:              hexutils.BytesToHex(w.PublicKey),
		PublicKey:           w.PublicKey,
		CompressedPublicKey: w.CompressedPublicKey,
		Metadata: &entities.Metadata{
			Disabled:  w.Disabled,
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
			DeletedAt: w.DeletedAt,
		},
		Tags: w.Tags,
	}
}
