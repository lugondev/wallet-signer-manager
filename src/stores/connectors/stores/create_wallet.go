package stores

import (
	"context"
	"github.com/lugondev/wallet-signer-manager/pkg/errors"
	authtypes "github.com/lugondev/wallet-signer-manager/src/auth/entities"
	"github.com/lugondev/wallet-signer-manager/src/stores"
	"github.com/lugondev/wallet-signer-manager/src/stores/entities"
	"github.com/lugondev/wallet-signer-manager/src/stores/store/wallets/hashicorp"

	hashicorpinfra "github.com/lugondev/wallet-signer-manager/src/infra/hashicorp"
)

func (c *Connector) CreateWallet(ctx context.Context, name string, allowedTenants []string, userInfo *authtypes.UserInfo) error {
	logger := c.logger.With("name", name)
	logger.Debug("creating wallet")

	// If vault is specified, it is a remote key store, otherwise it's a local key store
	var store stores.WalletStore
	switch {
	case name != "":
		vault, err := c.vaults.Get(ctx, name, userInfo)
		if err != nil {
			return err
		}
		store = hashicorp.New(vault.Client.(hashicorpinfra.PluginClient), logger)
	default:
		errMessage := "wallet name is required"
		logger.Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	c.createStore(name, entities.WalletStoreType, store, allowedTenants)

	logger.Info("key store created successfully")
	return nil
}
