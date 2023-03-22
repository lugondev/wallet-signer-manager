package entities

type Wallet struct {
	Namespaces          string
	PublicKey           []byte
	CompressedPublicKey []byte
	Metadata            *Metadata
	Tags                map[string]string
}
