package wallets

import (
	"github.com/lugondev/wallet-signer-manager/src/auth"
	"github.com/lugondev/wallet-signer-manager/src/entities"
	"github.com/lugondev/wallet-signer-manager/src/infra/log"
	"github.com/lugondev/wallet-signer-manager/src/stores"
	"github.com/lugondev/wallet-signer-manager/src/stores/database"
)

type Connector struct {
	store        stores.WalletStore
	logger       log.Logger
	db           database.Wallets
	authorizator auth.Authorizator
}

var _ stores.WalletStore = Connector{}

var ethAlgo = &entities.Algorithm{
	Type:          entities.Ecdsa,
	EllipticCurve: entities.Secp256k1,
}

func NewConnector(store stores.WalletStore, db database.Wallets, authorizator auth.Authorizator, logger log.Logger) *Connector {
	return &Connector{
		store:        store,
		logger:       logger,
		db:           db,
		authorizator: authorizator,
	}
}
