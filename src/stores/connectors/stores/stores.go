package stores

import (
	"context"
	"sync"

	"github.com/lugondev/signer-key-manager/src/vaults"

	"github.com/lugondev/signer-key-manager/pkg/errors"
	"github.com/lugondev/signer-key-manager/src/stores/entities"

	"github.com/lugondev/signer-key-manager/src/auth"
	"github.com/lugondev/signer-key-manager/src/infra/log"
	"github.com/lugondev/signer-key-manager/src/stores"
	"github.com/lugondev/signer-key-manager/src/stores/database"
)

type Connector struct {
	logger log.Logger
	mux    sync.RWMutex
	roles  auth.Roles
	stores map[string]*entities.Store
	vaults vaults.Vaults
	db     database.Database
}

var _ stores.Stores = &Connector{}

func NewConnector(roles auth.Roles, db database.Database, vaultsService vaults.Vaults, logger log.Logger) *Connector {
	return &Connector{
		logger: logger,
		mux:    sync.RWMutex{},
		roles:  roles,
		stores: make(map[string]*entities.Store),
		vaults: vaultsService,
		db:     db,
	}
}

// TODO: Move to data layer
func (c *Connector) createStore(name, storeType string, store interface{}, allowedTenants []string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.stores[name] = &entities.Store{
		Name:           name,
		AllowedTenants: allowedTenants,
		Store:          store,
		StoreType:      storeType,
	}
}

// TODO: Move to data layer
func (c *Connector) getStore(_ context.Context, name string, resolver auth.Authorizator) (*entities.Store, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if store, ok := c.stores[name]; ok {
		if err := resolver.CheckAccess(store.AllowedTenants); err != nil {
			return nil, err
		}

		return store, nil
	}

	errMessage := "store was not found"
	c.logger.Error(errMessage, "name", name)
	return nil, errors.NotFoundError(errMessage)
}
