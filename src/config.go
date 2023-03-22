package src

import (
	"encoding/json"
	"github.com/lugondev/signer-key-manager/pkg/http/server"
	"github.com/lugondev/signer-key-manager/src/infra/api-key/csv"
	"github.com/lugondev/signer-key-manager/src/infra/jwt/jose"
	"github.com/lugondev/signer-key-manager/src/infra/log/zap"
	manifestreader "github.com/lugondev/signer-key-manager/src/infra/manifests/yaml"
	"github.com/lugondev/signer-key-manager/src/infra/postgres/client"
	tls "github.com/lugondev/signer-key-manager/src/infra/tls/filesystem"
)

type Config struct {
	HTTP     *server.Config
	Logger   *zap.Config
	Postgres *client.Config
	OIDC     *jose.Config
	APIKey   *csv.Config
	TLS      *tls.Config
	Manifest *manifestreader.Config
}

func (c Config) ToJson() string {
	// convert struct to json
	jsonByte, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return string(jsonByte)
}
