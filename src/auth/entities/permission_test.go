package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListWildcardPermission(t *testing.T) {
	list := ListWildcardPermission("*:*")
	assert.Equal(t, list, ListPermissions())

	list = ListWildcardPermission("read:*")
	assert.Equal(t, list, []Permission{ReadWallet})

	list = ListWildcardPermission("*:wallets")
	assert.Equal(t, list, []Permission{ReadWallet, WriteWallet, DeleteWallet, DestroyWallet, SignWallet})
}
