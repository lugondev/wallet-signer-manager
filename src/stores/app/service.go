package app

import (
	"github.com/gorilla/mux"
	"github.com/lugondev/signer-key-manager/src/auth"
	"github.com/lugondev/signer-key-manager/src/infra/log"
	"github.com/lugondev/signer-key-manager/src/infra/postgres"
	"github.com/lugondev/signer-key-manager/src/stores/api/http"
	"github.com/lugondev/signer-key-manager/src/stores/connectors/stores"
	db "github.com/lugondev/signer-key-manager/src/stores/database/postgres"
	"github.com/lugondev/signer-key-manager/src/vaults"
)

func RegisterService(router *mux.Router, logger log.Logger, postgresClient postgres.Client, roles auth.Roles, vaultsService vaults.Vaults) *stores.Connector {
	// Data layer
	storesDB := db.New(logger, postgresClient)

	// Business layer
	storesService := stores.NewConnector(roles, storesDB, vaultsService, logger)

	// Service layer
	http.NewWalletsHandler(storesService).Register(router)

	return storesService
}
