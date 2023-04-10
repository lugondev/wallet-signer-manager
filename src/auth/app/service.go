package app

import (
	"crypto/x509"

	"github.com/justinas/alice"
	"github.com/lugondev/wallet-signer-manager/pkg/app"
	"github.com/lugondev/wallet-signer-manager/src/auth/api/http"
	"github.com/lugondev/wallet-signer-manager/src/auth/entities"
	"github.com/lugondev/wallet-signer-manager/src/auth/service/authenticator"
	"github.com/lugondev/wallet-signer-manager/src/auth/service/roles"
	"github.com/lugondev/wallet-signer-manager/src/infra/jwt"
	"github.com/lugondev/wallet-signer-manager/src/infra/log"
)

func RegisterService(
	a *app.App,
	logger log.Logger,
	jwtValidator jwt.Validator,
	apikeyClaims map[string]*entities.UserClaims,
	rootCAs *x509.CertPool,
) (*roles.Roles, error) {
	// Business layer
	// TODO: Create authorizator service here

	var authmid alice.Constructor
	if jwtValidator != nil || apikeyClaims != nil || rootCAs != nil {
		autheServ := authenticator.New(jwtValidator, apikeyClaims, rootCAs, logger)
		authmid = http.NewAuth(autheServ).Middleware
		logger.Info("authentication middleware is enabled")
	} else {
		authmid = http.NewNoAuth().Middleware
		logger.Warn("authentication is disabled")
	}

	rolesService := roles.New(logger)

	// Service layer
	httpMid := alice.New(
		http.NewAccessLog(logger.WithComponent("accesslog")).Middleware, // TODO: Move to correct domain when it exists
		authmid,
	)
	err := a.SetMiddleware(httpMid.Then)
	if err != nil {
		return nil, err
	}

	return rolesService, nil
}
