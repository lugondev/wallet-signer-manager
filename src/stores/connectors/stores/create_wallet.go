package stores

import (
	"context"
	"github.com/lugondev/signer-key-manager/pkg/errors"
	authtypes "github.com/lugondev/signer-key-manager/src/auth/entities"
	"github.com/lugondev/signer-key-manager/src/stores"
	"github.com/lugondev/signer-key-manager/src/stores/entities"
	"github.com/lugondev/signer-key-manager/src/stores/store/wallets/hashicorp"

	hashicorpinfra "github.com/lugondev/signer-key-manager/src/infra/hashicorp"
)

//func (c *Connector) CreateWallet(ctx context.Context, name, keyStore string, allowedTenants []string, userInfo *auth.UserInfo) error {
//	logger := c.logger.With("name", name, "key_store", keyStore)
//	logger.Debug("creating wallet store")
//
//	resolver := authorizator.New(userInfo.Permissions, userInfo.Tenant, c.logger)
//
//	store, err := c.getStore(ctx, keyStore, resolver)
//	if err != nil {
//		logger.Debug("error getting store", "err", err)
//		return err
//	}
//	logger.Debug("store found")
//
//	c.createStore(name, entities.WalletStoreType, store, allowedTenants)
//
//	logger.Info("wallet store created successfully")
//	return nil
//}

func (c *Connector) CreateWallet(ctx context.Context, name, walletStore string, allowedTenants []string, userInfo *authtypes.UserInfo) error {
	logger := c.logger.With("name", name, "walletStore", walletStore)
	logger.Debug("creating key store")

	//if name != "" && walletStore != "" {
	//	errMessage := "cannot specify vault and secret store simultaneously. Please choose one option"
	//	logger.Error(errMessage)
	//	return errors.InvalidParameterError(errMessage)
	//}

	//resolver := authorizator.New(userInfo.Permissions, userInfo.Tenant, c.logger)

	// If vault is specified, it is a remote key store, otherwise it's a local key store
	var store stores.WalletStore
	switch {
	case name != "":
		vault, err := c.vaults.Get(ctx, name, userInfo)
		if err != nil {
			return err
		}
		store = hashicorp.New(vault.Client.(hashicorpinfra.PluginClient), logger)
	//case walletStore != "":
	//	store, err := c.getStore(ctx, walletStore, resolver)
	//	if err != nil {
	//		return err
	//	}
	//
	//	store = localkeys.New(secretstore, c.db.Secrets(secretStore), c.logger)
	default:
		errMessage := "either vault or secret store must be specified. Please choose one option"
		logger.Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	c.createStore(name, entities.WalletStoreType, store, allowedTenants)

	logger.Info("key store created successfully")
	return nil
}
