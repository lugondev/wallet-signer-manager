package stores

import (
	"context"
	"github.com/lugondev/signer-key-manager/src/stores/connectors/wallets"

	"github.com/lugondev/signer-key-manager/src/auth/service/authorizator"
	"github.com/lugondev/signer-key-manager/src/stores/entities"

	"github.com/lugondev/signer-key-manager/src/auth"

	"github.com/lugondev/signer-key-manager/pkg/errors"
	authtypes "github.com/lugondev/signer-key-manager/src/auth/entities"
	"github.com/lugondev/signer-key-manager/src/stores"
)

func (c *Connector) Wallet(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.WalletStore, error) {
	permissions := c.roles.UserPermissions(ctx, userInfo)
	resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

	store, err := c.getWalletStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("wallet store found successfully", "store_name", storeName)
	return wallets.NewConnector(store, c.db.Wallets(storeName), resolver, c.logger), nil
}

func (c *Connector) WalletByPubkey(ctx context.Context, pubkey string, userInfo *authtypes.UserInfo) (stores.WalletStore, error) {
	logger := c.logger.With("pubkey", pubkey)

	walletStores, err := c.List(ctx, entities.WalletStoreType, userInfo)
	if err != nil {
		return nil, err
	}

	for _, storeName := range walletStores {
		walletStore, err := c.Wallet(ctx, storeName, userInfo)
		if err != nil {
			return nil, err
		}

		// If the account is not found in this store, continue to next one
		if _, err = walletStore.Get(ctx, pubkey); err != nil && errors.IsNotFoundError(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		logger.Debug("wallet store found successfully", "store_name", storeName)
		return walletStore, nil
	}

	errMessage := "wallet store was not found for the given address"
	logger.Error(errMessage)
	return nil, errors.NotFoundError(errMessage)
}

func (c *Connector) getWalletStore(ctx context.Context, storeName string, resolver auth.Authorizator) (stores.WalletStore, error) {
	storeInfo, err := c.getStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	if storeInfo.StoreType != entities.WalletStoreType {
		errMessage := "not a wallet store"
		c.logger.Error(errMessage, "store_name", storeName)
		return nil, errors.NotFoundError(errMessage)
	}

	return storeInfo.Store.(stores.WalletStore), nil
}
