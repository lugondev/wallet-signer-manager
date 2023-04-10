package roles

import (
	"context"

	"github.com/lugondev/wallet-signer-manager/src/auth/entities"
)

func (i *Roles) Get(ctx context.Context, name string, _ *entities.UserInfo) (*entities.Role, error) {
	logger := i.logger.With("name", name)

	// TODO: Implement {Resource/Role}BAC for roles

	vault, err := i.getRole(ctx, name)
	if err != nil {
		return nil, err
	}

	logger.Debug("role found successfully")
	return vault, nil
}
