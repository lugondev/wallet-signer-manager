package client

import (
	"encoding/base64"
	"path"

	"github.com/hashicorp/vault/api"
)

func (c *HashicorpVaultClient) GetWallet(id string) (*api.Secret, error) {
	secret, err := c.client.Logical().Read(c.pathWallets(id))
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

func (c *HashicorpVaultClient) UpdateWallet(id string, data map[string]interface{}) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(c.pathWallets(id), data)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) DestroyWallet(id string) error {
	_, err := c.client.Logical().Delete(path.Join(c.pathWallets(id), "destroy"))
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

func (c *HashicorpVaultClient) Sign(id string, data []byte) (*api.Secret, error) {
	secret, err := c.client.Logical().Write(path.Join(c.pathWallets(id), "sign"), map[string]interface{}{
		dataLabel: base64.URLEncoding.EncodeToString(data),
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return secret, nil
}

func (c *HashicorpVaultClient) pathWallets(suffix string) string {
	return path.Join(c.mountPoint, "wallets", suffix)
}
