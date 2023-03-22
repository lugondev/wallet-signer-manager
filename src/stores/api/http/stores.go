package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lugondev/signer-key-manager/pkg/errors"
	http2 "github.com/lugondev/signer-key-manager/src/infra/http"
	"github.com/lugondev/signer-key-manager/src/stores"
)

type StoresHandler struct {
	wallets *WalletsHandler
}

// NewStoresHandler creates a http.Handler to be served on /stores
func NewStoresHandler(s stores.Stores) *StoresHandler {
	return &StoresHandler{
		wallets: NewWalletsHandler(s),
	}
}

func (h *StoresHandler) Register(router *mux.Router) {
	// Create sub router for /stores
	storesSubRouter := router.PathPrefix("/stores").Subrouter()

	// Create sub router for /stores/{storeName}
	storeSubRouter := storesSubRouter.PathPrefix("/{storeName}").Subrouter()
	storeSubRouter.Use(storeSelector)

	// Register secrets handler on /stores/{storeName}/wallets
	walletsSubRouter := storeSubRouter.PathPrefix("/wallets").Subrouter()
	h.wallets.Register(walletsSubRouter)
}

func storeSelector(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(WithStoreName(r.Context(), mux.Vars(r)["storeName"])))
	})
}

func getLimitOffset(request *http.Request) (rLimit, rOffset uint64, err error) {
	limit := request.URL.Query().Get("limit")
	page := request.URL.Query().Get("page")
	if limit == "" {
		limit = http2.DefaultPageSize
	}

	rLimit, err = strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return 0, 0, errors.InvalidFormatError("invalid limit value")
	}

	iPage := uint64(0)
	rOffset = 0
	if page != "" {
		iPage, err = strconv.ParseUint(page, 10, 64)
		if err != nil {
			return 0, 0, errors.InvalidFormatError("invalid page value")
		}

		rOffset = iPage * rLimit
	}

	return rLimit, rOffset, nil
}
