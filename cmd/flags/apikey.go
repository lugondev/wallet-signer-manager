package flags

import (
	"fmt"

	"github.com/lugondev/signer-key-manager/src/infra/api-key/csv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(authAPIKeyFileViperKey, authAPIKeyFileEnv)
}

const (
	authAPIKeyFileFlag        = "auth-api-key-file"
	authAPIKeyFileViperKey    = "auth.api.key.file"
	authAPIKeyDefaultFileFlag = ""
	authAPIKeyFileEnv         = "AUTH_API_KEY_FILE"
)

func APIKeyFlags(f *pflag.FlagSet) {
	authAPIKeyFile(f)
}

func authAPIKeyFile(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`API key CSV file location.
Environment variable: %q`, authAPIKeyFileEnv)
	f.String(authAPIKeyFileFlag, authAPIKeyDefaultFileFlag, desc)
	_ = viper.BindPFlag(authAPIKeyFileViperKey, f.Lookup(authAPIKeyFileFlag))
}

func NewAPIKeyConfig(vipr *viper.Viper) *csv.Config {
	path := vipr.GetString(authAPIKeyFileViperKey)

	if path != "" {
		return csv.NewConfig(path)
	}

	return nil
}
