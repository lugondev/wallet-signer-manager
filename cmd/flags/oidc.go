package flags

import (
	"fmt"
	"strings"
	"time"

	"github.com/lugondev/signer-key-manager/src/infra/jwt/jose"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(authOIDCIssuerURLViperKey, authOIDCIssuerURLEnv)
	_ = viper.BindEnv(AuthOIDCAudienceViperKey, authOIDCAudienceEnv)
	_ = viper.BindEnv(authOIDCCustomClaimsViperKey, authOIDCCustomClaimsEnv)
}

const (
	authOIDCIssuerURLFlag     = "auth-oidc-issuer-url"
	authOIDCIssuerURLViperKey = "auth.oidc.issuer.url"
	authOIDCIssuerURLDefault  = ""
	authOIDCIssuerURLEnv      = "AUTH_OIDC_ISSUER_URL"
)

const (
	authOIDCAudienceFlag     = "auth-oidc-audience"
	AuthOIDCAudienceViperKey = "auth.oidc.audience"
	authOIDCAudienceEnv      = "AUTH_OIDC_AUDIENCE"
)

const (
	authOIDCCustomClaimsFlag     = "auth-oidc-custom-claims"
	authOIDCCustomClaimsViperKey = "auth.oidc.custom.claims"
	authOIDCCustomClaimsEnv      = "AUTH_OIDC_CUSTOM_CLAIMS"
)

func OIDCFlags(f *pflag.FlagSet) {
	authOIDCIssuerServer(f)
	authOIDCAudience(f)
	authOIDCCustomClaimsPath(f)
}

func authOIDCIssuerServer(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`OpenID Connect issuer server domain (ie. https://quorum-key-manager.eu.auth0.com).
Environment variable: %q`, authOIDCIssuerURLEnv)
	f.String(authOIDCIssuerURLFlag, authOIDCIssuerURLDefault, desc)
	_ = viper.BindPFlag(authOIDCIssuerURLViperKey, f.Lookup(authOIDCIssuerURLFlag))
}

func authOIDCAudience(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Expected audience ("aud" field) of JWT tokens.
Environment variable: %q`, authOIDCAudienceEnv)
	f.String(authOIDCAudienceFlag, "", desc)
	_ = viper.BindPFlag(AuthOIDCAudienceViperKey, f.Lookup(authOIDCAudienceFlag))
}

func authOIDCCustomClaimsPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path to for Quorum Key Manager custom claims in the JWT.
Environment variable: %q`, authOIDCCustomClaimsEnv)
	f.String(authOIDCCustomClaimsFlag, "", desc)
	_ = viper.BindPFlag(authOIDCCustomClaimsViperKey, f.Lookup(authOIDCCustomClaimsFlag))
}

func NewOIDCConfig(vipr *viper.Viper) *jose.Config {
	issuerURL := vipr.GetString(authOIDCIssuerURLViperKey)

	var aud []string
	if vipr.GetString(AuthOIDCAudienceViperKey) != "" {
		aud = strings.Split(vipr.GetString(AuthOIDCAudienceViperKey), ",")
	}

	if issuerURL != "" {
		return jose.NewConfig(
			issuerURL,
			aud,
			vipr.GetString(authOIDCCustomClaimsViperKey),
			5*time.Minute, // TODO: Make the cache ttl an ENV var if needed
		)
	}

	return nil
}
