package entities

const (
	WalletStoreType = "wallet"
)

type Store struct {
	Name           string
	AllowedTenants []string
	Store          interface{}
	StoreType      string
}
