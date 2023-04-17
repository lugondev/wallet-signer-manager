package hashicorp

import (
	hashicorp "github.com/hashicorp/vault/api"
)

type Client interface {
	PluginClient
	SetToken(token string)
	UnwrapToken(token string) (*hashicorp.Secret, error)
	Mount(path string, mountInfo *hashicorp.MountInput) error
	HealthCheck() error
}

type PluginClient interface {
	GetWallet(id string) (*hashicorp.Secret, error)
	CreateWallet(data map[string]interface{}) (*hashicorp.Secret, error)
	ImportWallet(data map[string]interface{}) (*hashicorp.Secret, error)
	ListWallets() (*hashicorp.Secret, error)
	UpdateWallet(id string, data map[string]interface{}) (*hashicorp.Secret, error)
	DestroyWallet(id string) error
	Sign(id string, typeSign string, data []byte) (*hashicorp.Secret, error)
}
