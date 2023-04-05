package manifest

import (
	"context"

	"github.com/lugondev/signer-key-manager/pkg/errors"
	authtypes "github.com/lugondev/signer-key-manager/src/auth/entities"
	entities2 "github.com/lugondev/signer-key-manager/src/entities"
	"github.com/lugondev/signer-key-manager/src/stores"
	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

type StoresHandler struct {
	stores   stores.Stores
	userInfo *authtypes.UserInfo
}

func NewStoresHandler(storesService stores.Stores) *StoresHandler {
	return &StoresHandler{
		stores:   storesService,
		userInfo: authtypes.NewWildcardUser(), // This handler always use the wildcard user because it's a manifest handler
	}
}

func (h *StoresHandler) Register(ctx context.Context, mnfs []entities2.Manifest) error {
	for _, mnf := range mnfs {
		var err error
		switch mnf.ResourceType {
		case entities.WalletStoreType:
			err = h.CreateWallet(ctx, mnf.Name, mnf.AllowedTenants)
		default:
			err = errors.InvalidFormatError("invalid store type")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (h *StoresHandler) CreateWallet(ctx context.Context, name string, allowedTenants []string) error {

	err := h.stores.CreateWallet(ctx, name, allowedTenants, h.userInfo)
	if err != nil {
		return err
	}

	return nil
}
