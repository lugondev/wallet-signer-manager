package client

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"path"

	"github.com/hashicorp/vault/api"
)

func (c *HashicorpVaultClient) GetWallet(pubkey string) (*api.Secret, error) {
	secret, err := c.client.Logical().Read(c.pathWallets(pubkey))
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) CreateWallet(data map[string]interface{}) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(c.pathWallets(""), data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) ImportWallet(data map[string]interface{}) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(c.pathWallets("import"), data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) UpdateWallet(pubkey string, data map[string]interface{}) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(c.pathWallets(pubkey), data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) DestroyWallet(pubkey string) error {
	_, err := c.client.Logical().Delete(path.Join(c.pathWallets(pubkey), "destroy"))
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *HashicorpVaultClient) ListWallets() (*api.Secret, error) {
	secret, err := c.client.Logical().List(c.pathWallets(""))
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) Sign(pubkey string, data []byte) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(path.Join(c.pathWallets(pubkey), "sign"), map[string]interface{}{
		dataLabel: hexutil.Encode(data),
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) pathWallets(suffix string) string {
	return path.Join(c.mountPoint, "wallets", suffix)
}
