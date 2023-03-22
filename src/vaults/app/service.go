package app

import (
	"github.com/lugondev/signer-key-manager/src/auth"
	"github.com/lugondev/signer-key-manager/src/infra/log"
	"github.com/lugondev/signer-key-manager/src/vaults/service/vaults"
)

func RegisterService(logger log.Logger, roles auth.Roles) *vaults.Vaults {
	// Business layer
	return vaults.New(roles, logger)
}
