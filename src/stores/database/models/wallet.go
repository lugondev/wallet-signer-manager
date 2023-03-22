package models

import (
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

type Wallet struct {
	tableName struct{} `pg:"wallets"` // nolint:unused,structcheck // reason

	StoreID             string `pg:",pk"`
	KeyID               string
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
	return &entities.Wallet{
		Tags:                attr.Tags,
		PublicKey:           key.PublicKey,
		CompressedPublicKey: crypto.CompressPubkey(pubKey),
		Metadata: &entities.Metadata{
			Disabled:  key.Metadata.Disabled,
			CreatedAt: key.Metadata.CreatedAt,
			UpdatedAt: key.Metadata.UpdatedAt,
		},
	}
}

func (eth *Wallet) ToEntity() *entities.Wallet {
	return &entities.Wallet{
		PublicKey:           eth.PublicKey,
		CompressedPublicKey: eth.CompressedPublicKey,
		Metadata: &entities.Metadata{
			Disabled:  eth.Disabled,
			CreatedAt: eth.CreatedAt,
			UpdatedAt: eth.UpdatedAt,
			DeletedAt: eth.DeletedAt,
		},
		Tags: eth.Tags,
	}
}
