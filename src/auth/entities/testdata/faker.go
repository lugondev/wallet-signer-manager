package testdata

import (
	"github.com/lugondev/signer-key-manager/src/auth/entities"
)

func FakeUserClaims() *entities.UserClaims {
	return &entities.UserClaims{
		Tenant:      "TenantOne|Alice",
		Permissions: []string{"read:key", "write:key"},
		Roles:       []string{"guest", "admin"},
	}
}
