package app

import (
	"github.com/lugondev/wallet-signer-manager/src/auth"
	"github.com/lugondev/wallet-signer-manager/src/infra/log"
	"github.com/lugondev/wallet-signer-manager/src/vaults/service/vaults"
)

func RegisterService(logger log.Logger, roles auth.Roles) *vaults.Vaults {
	// Business layer
	return vaults.New(roles, logger)
}
