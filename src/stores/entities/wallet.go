package entities

type Wallet struct {
	KeyID               string
	PublicKey           []byte
	CompressedPublicKey []byte
	Metadata            *Metadata
	Tags                map[string]string
}
