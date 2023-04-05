package hashicorp

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/lugondev/signer-key-manager/src/stores/entities"

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
		Tags: make(map[string]string),
	}

	if hashicorpSecret.Data[tagsLabel] != nil {
		tags := hashicorpSecret.Data[tagsLabel].(map[string]interface{})
		for k, v := range tags {
			key.Tags[k] = v.(string)
		}
	}

	//key.Metadata.CreatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[createdAtLabel].(string))
	//key.Metadata.UpdatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[updatedAtLabel].(string))

	return key, nil
}
