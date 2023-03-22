package stores

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/lugondev/signer-key-manager/src/auth/service/authorizator"

	arrays "github.com/lugondev/signer-key-manager/pkg/common"
	authtypes "github.com/lugondev/signer-key-manager/src/auth/entities"
	"github.com/lugondev/signer-key-manager/src/stores/database/models"
	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

func (c *Connector) ImportWallets(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) error {
	logger := c.logger.With("store_name", storeName)
	logger.Info("importing ethereum accounts...")

	// TODO: Uncomment when authManager no longer a runnable
	// permissions := c.authManager.UserPermissions(userInfo)
	resolver := authorizator.New(userInfo.Permissions, userInfo.Tenant, c.logger)

	store, err := c.getWalletStore(ctx, storeName, resolver)
	if err != nil {
		return err
	}

	storeIDs, err := store.List(ctx, 0, 0)
	if err != nil {
		return err
	}

	db := c.db.Wallets(storeName)
	dbAddresses, err := db.SearchAddresses(ctx, false, 0, 0)
	if err != nil {
		return err
	}
	addressMap := arrays.ToMap(dbAddresses)

	var nSuccesses uint
	var nFailures uint
	for _, id := range storeIDs {
		key, err := store.Get(ctx, id)
		if err != nil {
			nFailures++
			continue
		}

		acc := models.NewWalletFromKey(key, &entities.Attributes{})
		if _, found := addressMap[hexutil.Encode(acc.CompressedPublicKey)]; !found {
			_, err = db.Add(ctx, acc)
			if err != nil {
				nFailures++
				continue
			}

			nSuccesses++
		}
	}

	logger.Info("wallet import completed", "n_successes", nSuccesses, "n_failures", nFailures)
	return nil
}
