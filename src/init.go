package src

import (
	"context"
	"fmt"

	"github.com/lugondev/wallet-signer-manager/src/auth"
	rolesapi "github.com/lugondev/wallet-signer-manager/src/auth/api/manifest"
	"github.com/lugondev/wallet-signer-manager/src/entities"
	manifestreader "github.com/lugondev/wallet-signer-manager/src/infra/manifests/yaml"
	"github.com/lugondev/wallet-signer-manager/src/stores"
	storesapi "github.com/lugondev/wallet-signer-manager/src/stores/api/manifest"
	"github.com/lugondev/wallet-signer-manager/src/vaults"
	vaultsapi "github.com/lugondev/wallet-signer-manager/src/vaults/api/manifest"
)

func initialize(
	ctx context.Context,
	cfg *manifestreader.Config,
	rolesService auth.Roles,
	vaultsService vaults.Vaults,
	storesService stores.Stores,
) error {
	manifestReader, err := manifestreader.New(cfg)
	if err != nil {
		return err
	}

	manifests, err := manifestReader.Load(ctx)
	if err != nil {
		return err
	}
	fmt.Println("manifests: ", manifests)

	// Note that order is important here as stores depend on the existing vaults, do not use a switch!

	err = rolesapi.NewRolesHandler(rolesService).Register(ctx, manifests[entities.RoleKind])
	if err != nil {
		return err
	}

	err = vaultsapi.NewVaultsHandler(vaultsService).Register(ctx, manifests[entities.VaultKind])
	if err != nil {
		return err
	}

	err = storesapi.NewStoresHandler(storesService).Register(ctx, manifests[entities.StoreKind])
	if err != nil {
		return err
	}

	return nil
}
