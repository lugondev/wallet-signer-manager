package src

import (
	"context"
	"crypto/x509"
	"github.com/lugondev/wallet-signer-manager/pkg/app"
	authapp "github.com/lugondev/wallet-signer-manager/src/auth/app"
	authtypes "github.com/lugondev/wallet-signer-manager/src/auth/entities"
	"github.com/lugondev/wallet-signer-manager/src/infra/api-key/csv"
	"github.com/lugondev/wallet-signer-manager/src/infra/jwt"
	"github.com/lugondev/wallet-signer-manager/src/infra/jwt/jose"
	"github.com/lugondev/wallet-signer-manager/src/infra/log"
	"github.com/lugondev/wallet-signer-manager/src/infra/postgres/client"
	tls "github.com/lugondev/wallet-signer-manager/src/infra/tls/filesystem"
	storesapp "github.com/lugondev/wallet-signer-manager/src/stores/app"
	vaultsapp "github.com/lugondev/wallet-signer-manager/src/vaults/app"
)

func New(ctx context.Context, cfg *Config, logger log.Logger) (*app.App, error) {
	// Infra layer
	pgClient, err := client.New(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	var jwtValidator jwt.Validator
	var apikeyClaims map[string]*authtypes.UserClaims
	var rootCAs *x509.CertPool
	if cfg.OIDC != nil {
		jwtValidator, err = getJWTValidator(cfg.OIDC, logger)
		if err != nil {
			return nil, err
		}
	}

	if cfg.APIKey != nil {
		apikeyClaims, err = getAPIKeys(ctx, cfg.APIKey, logger)
		if err != nil {
			return nil, err
		}
	}

	if cfg.TLS != nil {
		rootCAs, err = getRootCAs(ctx, cfg.TLS, logger)
		if err != nil {
			return nil, err
		}
	}

	// Register Services
	a := app.New(&app.Config{HTTP: cfg.HTTP}, logger.WithComponent("app"))
	router := a.Router()

	authService, err := authapp.RegisterService(a, logger.WithComponent("auth"), jwtValidator, apikeyClaims, rootCAs)
	if err != nil {
		return nil, err
	}

	vaultsService := vaultsapp.RegisterService(logger.WithComponent("vaults"), authService)
	storesService := storesapp.RegisterService(router, logger.WithComponent("stores"), pgClient, authService, vaultsService)

	err = initialize(ctx, cfg.Manifest, authService, vaultsService, storesService)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func getAPIKeys(ctx context.Context, cfg *csv.Config, logger log.Logger) (map[string]*authtypes.UserClaims, error) {
	apiKeyReader, err := csv.New(cfg)
	if err != nil {
		return nil, err
	}

	apikeyClaims, err := apiKeyReader.Load(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info("API key authentication enabled")

	return apikeyClaims, nil
}

func getJWTValidator(cfg *jose.Config, logger log.Logger) (*jose.Validator, error) {
	jwtValidator, err := jose.New(cfg)
	if err != nil {
		return nil, err
	}

	logger.Info("JWT authentication enabled")

	return jwtValidator, nil
}

func getRootCAs(ctx context.Context, cfg *tls.Config, logger log.Logger) (*x509.CertPool, error) {
	tlsReader, err := tls.New(cfg)
	if err != nil {
		return nil, err
	}

	rootCAs, err := tlsReader.Load(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info("TLS authentication enabled")

	return rootCAs, nil
}
