package hashicorp

import (
	"encoding/base64"
	"time"

	"github.com/lugondev/signer-key-manager/src/stores/entities"

	"github.com/lugondev/signer-key-manager/pkg/errors"

	"github.com/hashicorp/vault/api"
)

func parseAPISecretToKey(hashicorpSecret *api.Secret) (*entities.Wallet, error) {
	pubKey, err := base64.URLEncoding.DecodeString(hashicorpSecret.Data[publicKeyLabel].(string))
	if err != nil {
		return nil, errors.HashicorpVaultError("failed to decode public key")
	}

	key := &entities.Wallet{
		PublicKey:           pubKey,
		CompressedPublicKey: pubKey,
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

	key.Metadata.CreatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[createdAtLabel].(string))
	key.Metadata.UpdatedAt, _ = time.Parse(time.RFC3339, hashicorpSecret.Data[updatedAtLabel].(string))

	return key, nil
}
