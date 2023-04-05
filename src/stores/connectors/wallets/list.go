package wallets

import (
	"context"

	"github.com/lugondev/signer-key-manager/src/auth/entities"
)

func (c Connector) List(ctx context.Context, limit, offset uint64) ([]string, error) {
	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceWallets})
	if err != nil {
		return nil, err
	}

	strAddr, err := c.db.SearchAddresses(ctx, false, limit, offset)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, addr := range strAddr {
		addrs = append(addrs, addr)
	}

	c.logger.Debug("ethereum accounts listed successfully")
	return addrs, nil
}

func (c Connector) ListDeleted(ctx context.Context, limit, offset uint64) ([]string, error) {
	err := c.authorizator.CheckPermission(&entities.Operation{Action: entities.ActionRead, Resource: entities.ResourceWallets})
	if err != nil {
		return nil, err
	}

	strAddr, err := c.db.SearchAddresses(ctx, true, limit, offset)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, addr := range strAddr {
		addrs = append(addrs, addr)
	}

	c.logger.Debug("deleted ethereum accounts listed successfully")
	return addrs, nil
}
