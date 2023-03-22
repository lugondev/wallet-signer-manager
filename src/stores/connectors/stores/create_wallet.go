package stores

import (
	"context"

	auth "github.com/lugondev/signer-key-manager/src/auth/entities"
	"github.com/lugondev/signer-key-manager/src/auth/service/authorizator"
	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

func (c *Connector) CreateWallet(ctx context.Context, name, keyStore string, allowedTenants []string, userInfo *auth.UserInfo) error {
	logger := c.logger.With("name", name, "key_store", keyStore)
	logger.Debug("creating ethereum store")

	// TODO: Uncomment when authManager no longer a runnable
	// permissions := c.authManager.UserPermissions(userInfo)
	resolver := authorizator.New(userInfo.Permissions, userInfo.Tenant, c.logger)

	store, err := c.getStore(ctx, keyStore, resolver)
	if err != nil {
		return err
	}

	c.createStore(name, entities.WalletStoreType, store, allowedTenants)

	logger.Info("ethereum store created successfully")
	return nil
}
