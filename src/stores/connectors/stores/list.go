package stores

import (
	"context"

	"github.com/lugondev/signer-key-manager/src/auth/service/authorizator"
	"github.com/lugondev/signer-key-manager/src/stores/entities"

	authtypes "github.com/lugondev/signer-key-manager/src/auth/entities"
)

func (c *Connector) List(ctx context.Context, storeType string, userInfo *authtypes.UserInfo) ([]string, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	var storeNames []string
	for k, storeInfo := range c.stores {
		if storeType != "" && storeInfo.StoreType != storeType {
			continue
		}

		permissions := c.roles.UserPermissions(ctx, userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

		if err := resolver.CheckAccess(storeInfo.AllowedTenants); err != nil {
			continue
		}

		storeNames = append(storeNames, k)
	}

	return storeNames, nil
}

func (c *Connector) ListAllWallets(ctx context.Context, userInfo *authtypes.UserInfo) ([]string, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	var accs []string
	stores, err := c.List(ctx, entities.WalletStoreType, userInfo)
	if err != nil {
		return nil, err
	}

	for _, storeName := range stores {
		store, err := c.Wallet(ctx, storeName, userInfo)
		if err != nil {
			return nil, err
		}

		storeAccs, err := store.List(ctx, 0, 0)
		if err != nil {
			return nil, err
		}
		accs = append(accs, storeAccs...)
	}

	return accs, nil
}
