package hashicorp

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/lugondev/wallet-signer-manager/src/stores/entities"
	"time"

	"github.com/hashicorp/vault/api"
)

func parseAPISecretToWallet(hashicorpSecret *api.Secret) (*entities.Wallet, error) {
	pubKey := hashicorpSecret.Data[publicKeyLabel].(string)
	compressedPublicKey := hashicorpSecret.Data[compressedPublicKeyLabel].(string)
	namespace := hashicorpSecret.Data[namespaceLabel].(string)

	key := &entities.Wallet{
		Namespaces:          namespace,
		Pubkey:              compressedPublicKey,
		PublicKey:           common.FromHex(pubKey),
		CompressedPublicKey: common.FromHex(compressedPublicKey),
		Metadata: &entities.Metadata{
			Disabled: false,
		},
		Tags:  make(map[string]string),
		Extra: make(map[string]interface{}),
	}

	if hashicorpSecret.Data[tagsLabel] != nil {
		tags := hashicorpSecret.Data[tagsLabel].(map[string]interface{})
		for k, v := range tags {
			key.Tags[k] = v.(string)
		}
	}

	if hashicorpSecret.Data[extraLabel] != nil {
		auth := hashicorpSecret.Data[extraLabel].(map[string]interface{})
		for k, v := range auth {
			key.Extra[k] = v
		}
	}

	if hashicorpSecret.Data[createdAtLabel] != nil {
		key.Metadata.CreatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[createdAtLabel].(string))
	}

	if hashicorpSecret.Data[updatedAtLabel] != nil {
		key.Metadata.UpdatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[updatedAtLabel].(string))
	}

	return key, nil
}
