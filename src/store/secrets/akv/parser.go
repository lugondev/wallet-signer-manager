package akv

import (
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

func parseSecretBundle(secretBundle keyvault.SecretBundle) *entities.Secret {
	secret := &entities.Secret{
		Value:    *secretBundle.Value,
		Tags:     tomapstr(secretBundle.Tags),
		Metadata: &entities.Metadata{},
	}

	if secretBundle.ID != nil {
		// path.Base to only retrieve the secretVersion instead of https://<vaultName>.vault.azure.net/secrets/<secretName>/<secretVersion>
		chunks := strings.Split(*secretBundle.ID, "/")
		secret.Metadata.Version = chunks[len(chunks)-1]
		secret.ID = chunks[len(chunks)-2]
	}
	if expires := secretBundle.Attributes.Expires; expires != nil {
		secret.Metadata.ExpireAt = time.Unix(0, expires.Duration().Nanoseconds()).In(time.UTC)
	}
	if created := secretBundle.Attributes.Created; created != nil {
		secret.Metadata.CreatedAt = time.Unix(0, created.Duration().Nanoseconds()).In(time.UTC)
	}
	if updated := secretBundle.Attributes.Updated; updated != nil {
		secret.Metadata.UpdatedAt = time.Unix(0, updated.Duration().Nanoseconds()).In(time.UTC)
	}
	if enabled := secretBundle.Attributes.Enabled; enabled != nil {
		secret.Metadata.Disabled = !*enabled
	}

	return secret
}

func parseErrorResponse(err error) error {
	aerr, _ := err.(autorest.DetailedError)

	switch aerr.StatusCode.(int) {
	case http.StatusNotFound:
		return errors.NotFoundError("%v", aerr.Original)
	case http.StatusBadRequest:
		return errors.InvalidFormatError("%v", aerr.Original)
	case http.StatusUnprocessableEntity:
		return errors.InvalidParameterError("%v", aerr.Original)
	case http.StatusConflict:
		return errors.AlreadyExistsError("%v", aerr.Original)
	default:
		return errors.AKVConnectionError("%v", aerr.Original)
	}
}

func tomapstrptr(m map[string]string) map[string]*string {
	nm := make(map[string]*string)
	for k, v := range m {
		nm[k] = &(&struct{ x string }{v}).x
	}
	return nm
}

func tomapstr(m map[string]*string) map[string]string {
	nm := make(map[string]string)
	for k, v := range m {
		nm[k] = *v
	}
	return nm
}