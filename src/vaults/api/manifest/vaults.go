package manifest

import (
	"context"

	"github.com/lugondev/wallet-signer-manager/pkg/errors"
	"github.com/lugondev/wallet-signer-manager/pkg/json"
	auth "github.com/lugondev/wallet-signer-manager/src/auth/entities"
	"github.com/lugondev/wallet-signer-manager/src/entities"
	"github.com/lugondev/wallet-signer-manager/src/vaults"
)

type VaultsHandler struct {
	vaults   vaults.Vaults
	userInfo *auth.UserInfo
}

func NewVaultsHandler(vaultsService vaults.Vaults) *VaultsHandler {
	return &VaultsHandler{
		vaults:   vaultsService,
		userInfo: auth.NewWildcardUser(),
	}
}

func (h *VaultsHandler) Register(ctx context.Context, mnfs []entities.Manifest) error {
	for _, mnf := range mnfs {
		var err error
		switch mnf.ResourceType {
		case entities.HashicorpVaultType:
			err = h.CreateHashicorp(ctx, mnf.Name, mnf.AllowedTenants, mnf.Specs)
		default:
			return errors.InvalidFormatError("invalid vault type")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (h *VaultsHandler) CreateHashicorp(ctx context.Context, name string, allowedTenants []string, specs interface{}) error {
	config := &entities.HashicorpConfig{}
	err := json.UnmarshalYAML(specs, config)
	if err != nil {
		return errors.InvalidFormatError(err.Error())
	}

	err = h.vaults.CreateHashicorp(ctx, name, config, allowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
