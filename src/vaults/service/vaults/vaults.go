package vaults

import (
	"context"
	"sync"

	"github.com/lugondev/wallet-signer-manager/pkg/errors"
	"github.com/lugondev/wallet-signer-manager/src/auth"
	"github.com/lugondev/wallet-signer-manager/src/entities"
	"github.com/lugondev/wallet-signer-manager/src/infra/log"
	"github.com/lugondev/wallet-signer-manager/src/vaults"
)

type Vaults struct {
	logger log.Logger
	mux    sync.RWMutex
	vaults map[string]*entities.Vault
	roles  auth.Roles
}

var _ vaults.Vaults = &Vaults{}

func New(roles auth.Roles, logger log.Logger) *Vaults {
	return &Vaults{
		logger: logger,
		mux:    sync.RWMutex{},
		vaults: make(map[string]*entities.Vault),
		roles:  roles,
	}
}

// TODO: Move to in-memory data layer
func (c *Vaults) createVault(name, vaultType string, allowedTenants []string, cli interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.vaults[name] = &entities.Vault{
		Name:           name,
		Client:         cli,
		VaultType:      vaultType,
		AllowedTenants: allowedTenants,
	}
}

// TODO: Move to data layer
func (c *Vaults) getVault(_ context.Context, name string, resolver auth.Authorizator) (*entities.Vault, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if vault, ok := c.vaults[name]; ok {
		if err := resolver.CheckAccess(vault.AllowedTenants); err != nil {
			return nil, err
		}

		return vault, nil
	}

	errMessage := "vault was not found"
	c.logger.Error(errMessage, "name", name)
	return nil, errors.NotFoundError(errMessage)
}
