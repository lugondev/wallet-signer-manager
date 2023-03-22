package roles

import (
	"context"
	"fmt"

	"github.com/lugondev/signer-key-manager/pkg/errors"
	"github.com/lugondev/signer-key-manager/src/auth/entities"
)

func (i *Roles) Create(ctx context.Context, name string, permissions []entities.Permission, _ *entities.UserInfo) error {
	logger := i.logger.With("name", name, "permissions", permissions)
	logger.Debug("creating role")

	// TODO: Implement {Resource/Role}BAC for roles

	if _, ok := i.roles[name]; ok {
		errMessage := fmt.Sprintf("role %s already exist", name)
		logger.Error(errMessage)
		return errors.AlreadyExistsError(errMessage)
	}

	i.createRole(ctx, name, permissions)

	logger.Info("role created successfully")
	return nil
}
