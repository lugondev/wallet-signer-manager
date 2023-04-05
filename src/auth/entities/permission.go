package entities

import (
	"fmt"
	"strings"
)

type Permission string

const ReadWallet Permission = "read:wallets"
const WriteWallet Permission = "write:wallets"
const DeleteWallet Permission = "delete:wallets"
const DestroyWallet Permission = "destroy:wallets"
const SignWallet Permission = "sign:wallets"

func ListPermissions() []Permission {
	return []Permission{
		ReadWallet,
		WriteWallet,
		DeleteWallet,
		DestroyWallet,
		SignWallet,
	}
}

func ListWildcardPermission(p string) []Permission {
	all := ListPermissions()
	parts := strings.Split(p, ":")
	action, resource := parts[0], parts[1]
	if action == "*" && resource == "*" {
		return all
	}

	var included []Permission
	for _, ip := range all {
		if action == "*" && strings.Contains(string(ip), fmt.Sprintf(":%s", resource)) {
			included = append(included, ip)
		}
		if resource == "*" && strings.Contains(string(ip), fmt.Sprintf("%s:", action)) {
			included = append(included, ip)
		}
	}

	return included
}
