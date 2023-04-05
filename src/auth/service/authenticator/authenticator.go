package authenticator

import (
	"context"
	"crypto/sha256"
	tls2 "crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/lugondev/signer-key-manager/pkg/errors"
	"github.com/lugondev/signer-key-manager/pkg/tls"
	"github.com/lugondev/signer-key-manager/src/auth"
	"github.com/lugondev/signer-key-manager/src/auth/entities"
	"github.com/lugondev/signer-key-manager/src/infra/jwt"
	"github.com/lugondev/signer-key-manager/src/infra/log"
)

const (
	APIKeyAuthMode = "apikey"
	JWTAuthMode    = "jwt"
	TLSAuthMode    = "tls"
)

type Authenticator struct {
	logger       log.Logger
	jwtValidator jwt.Validator
	apiKeyClaims map[string]*entities.UserClaims
	rootCAs      *x509.CertPool
}

var _ auth.Authenticator = &Authenticator{}

func New(jwtValidator jwt.Validator, apiKeyClaims map[string]*entities.UserClaims, rootCAs *x509.CertPool, logger log.Logger) *Authenticator {
	return &Authenticator{
		jwtValidator: jwtValidator,
		apiKeyClaims: apiKeyClaims,
		rootCAs:      rootCAs,
		logger:       logger,
	}
}

func (a *Authenticator) AuthenticateJWT(ctx context.Context, token string) (*entities.UserInfo, error) {
	if a.jwtValidator == nil {
		errMessage := "jwt authentication method is not enabled"
		a.logger.Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	a.logger.Debug("extracting user info from jwt token")

	tokenClaims, err := a.jwtValidator.ValidateToken(ctx, token)
	if err != nil {
		errMessage := "failed to validate jwt token"
		a.logger.WithError(err).Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	claims, err := a.jwtValidator.ParseClaims(tokenClaims)
	if err != nil {
		errMessage := "failed to parse jwt token claims"
		a.logger.WithError(err).Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	return a.userInfoFromClaims(JWTAuthMode, claims), nil
}

func (a *Authenticator) AuthenticateAPIKey(_ context.Context, apiKey []byte) (*entities.UserInfo, error) {
	if a.apiKeyClaims == nil {
		errMessage := "api key authentication method is not enabled"
		a.logger.Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	a.logger.Debug("extracting user info from api key")

	apiKeySha256 := fmt.Sprintf("%x", sha256.Sum256(apiKey))
	claims, ok := a.apiKeyClaims[apiKeySha256]
	if !ok {
		errMessage := "invalid api key"
		a.logger.Warn(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	return a.userInfoFromClaims(APIKeyAuthMode, claims), nil
}

// AuthenticateTLS checks rootCAs and retrieve user info
func (a *Authenticator) AuthenticateTLS(_ context.Context, connState *tls2.ConnectionState) (*entities.UserInfo, error) {
	if a.rootCAs == nil {
		errMessage := "tls authentication method is not enabled"
		a.logger.Error(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	if !connState.HandshakeComplete {
		errMessage := "request must complete valid handshake"
		a.logger.Warn(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	err := tls.VerifyCertificateAuthority(connState.PeerCertificates, connState.ServerName, a.rootCAs, true)
	if err != nil {
		errMessage := "invalid tls certificate"
		a.logger.WithError(err).Warn(errMessage)
		return nil, errors.UnauthorizedError(errMessage)
	}

	// first array element is the leaf
	clientCert := connState.PeerCertificates[0]
	claims := &entities.UserClaims{
		Tenant:      clientCert.Subject.CommonName,
		Permissions: clientCert.Subject.OrganizationalUnit,
		Roles:       clientCert.Subject.Organization,
	}
	return a.userInfoFromClaims(TLSAuthMode, claims), nil
}

func (a *Authenticator) userInfoFromClaims(authMode string, claims *entities.UserClaims) *entities.UserInfo {
	userInfo := &entities.UserInfo{AuthMode: authMode}

	// If more than one element in subject, then the username has been specified
	subject := strings.Split(claims.Tenant, "|")
	if len(subject) > 1 {
		userInfo.Username = subject[1]
	}
	userInfo.Tenant = subject[0]

	for _, permission := range claims.Permissions {
		if !strings.Contains(permission, ":") {
			// Ignore invalid permissions
			continue
		}

		if strings.Contains(permission, "*") {
			userInfo.Permissions = append(userInfo.Permissions, entities.ListWildcardPermission(permission)...)
		} else {
			userInfo.Permissions = append(userInfo.Permissions, entities.Permission(permission))
		}
	}

	userInfo.Roles = claims.Roles

	return userInfo
}
