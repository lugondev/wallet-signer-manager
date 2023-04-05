package entities

type OpAction string
type OpResource string

var ActionRead OpAction = "read"
var ActionWrite OpAction = "write"
var ActionSign OpAction = "sign"
var ActionDelete OpAction = "delete"
var ActionDestroy OpAction = "destroy"

var ResourceWallets OpResource = "wallets"

type Operation struct {
	Action   OpAction
	Resource OpResource
}
