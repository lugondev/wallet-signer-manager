package http

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"net/http"

	"github.com/lugondev/wallet-signer-manager/src/stores/api/formatters"

	auth "github.com/lugondev/wallet-signer-manager/src/auth/api/http"

	infrahttp "github.com/lugondev/wallet-signer-manager/src/infra/http"

	"github.com/gorilla/mux"
	"github.com/lugondev/wallet-signer-manager/pkg/errors"
	jsonutils "github.com/lugondev/wallet-signer-manager/pkg/json"
	"github.com/lugondev/wallet-signer-manager/src/stores"
	"github.com/lugondev/wallet-signer-manager/src/stores/api/types"
	"github.com/lugondev/wallet-signer-manager/src/stores/entities"
)

type WalletsHandler struct {
	stores stores.Stores
}

func NewWalletsHandler(storesConnector stores.Stores) *WalletsHandler {
	return &WalletsHandler{
		stores: storesConnector,
	}
}

func (h *WalletsHandler) Register(r *mux.Router) {
	r.Methods(http.MethodPost).Path("").HandlerFunc(h.create)
	r.Methods(http.MethodGet).Path("").HandlerFunc(h.list)
	r.Methods(http.MethodPost).Path("/import").HandlerFunc(h.importAccount)
	r.Methods(http.MethodPost, http.MethodPut).Path("/{pubkey}/sign").HandlerFunc(h.sign)
	r.Methods(http.MethodPut).Path("/{pubkey}/restore").HandlerFunc(h.restore)
	r.Methods(http.MethodPatch, http.MethodPut).Path("/{pubkey}").HandlerFunc(h.update)
	r.Methods(http.MethodGet).Path("/{pubkey}").HandlerFunc(h.getOne)
	r.Methods(http.MethodDelete, http.MethodPut).Path("/{pubkey}").HandlerFunc(h.delete)
	r.Methods(http.MethodDelete).Path("/{pubkey}/destroy").HandlerFunc(h.destroy)
}

// @Summary      Create a wallet
// @Description  Create a new ECDSA Secp256k1 key representing a wallet
// @Tags         Ethereum
// @Accept       json
// @Produce      json
// @Param        storeName  path      string                         true  "Store ID"
// @Param        request    body      types.CreateEthAccountRequest  true  "Create Ethereum Account request"
// @Success      200        {object}  types.EthAccountResponse       "Created Ethereum Account"
// @Failure      400        {object}  infrahttp.ErrorResponse        "Invalid request format"
// @Failure      401        {object}  infrahttp.ErrorResponse        "Unauthorized"
// @Failure      403        {object}  infrahttp.ErrorResponse        "Forbidden"
// @Failure      404        {object}  infrahttp.ErrorResponse        "Store not found"
// @Failure      500        {object}  infrahttp.ErrorResponse        "Internal server error"
// @Router       /stores/{storeName}/wallets [post]
func (h *WalletsHandler) create(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	createReq := &types.CreateWalletRequest{}
	err := jsonutils.UnmarshalBody(request.Body, createReq)
	if err != nil && err.Error() != "EOF" {
		infrahttp.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userAuth := auth.UserInfoFromContext(ctx)
	walletStore, err := h.stores.Wallet(ctx, StoreNameFromContext(request.Context()), userAuth)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	var keyID string
	if createReq.KeyID != "" {
		keyID = createReq.KeyID
	} else {
		keyID = generateRandomKeyID()
	}

	wallet, err := walletStore.Create(ctx, keyID, &entities.Attributes{Tags: createReq.Tags, Auth: auth.UserInfoToMap(userAuth)})
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WriteJSON(rw, formatters.FormatWalletResponse(wallet))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary      Import a wallet
// @Description  Import an ECDSA Secp256k1 key representing a Wallet
// @Accept       json
// @Produce      json
// @Tags         Ethereum
// @Param        storeName  path      string                         true  "Store ID"
// @Param        request    body      types.ImportEthAccountRequest  true  "Create Ethereum Account request"
// @Success      200        {object}  types.EthAccountResponse       "Created Ethereum Account"
// @Failure      400        {object}  infrahttp.ErrorResponse        "Invalid request format"
// @Failure      401        {object}  infrahttp.ErrorResponse        "Unauthorized"
// @Failure      403        {object}  infrahttp.ErrorResponse        "Forbidden"
// @Failure      404        {object}  infrahttp.ErrorResponse        "Store not found"
// @Failure      500        {object}  infrahttp.ErrorResponse        "Internal server error"
// @Router       /stores/{storeName}/wallets/import [post]
func (h *WalletsHandler) importAccount(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	importReq := &types.ImportWalletRequest{}
	err := jsonutils.UnmarshalBody(request.Body, importReq)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userAuth := auth.UserInfoFromContext(ctx)
	walletStore, err := h.stores.Wallet(ctx, StoreNameFromContext(request.Context()), userAuth)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	var keyID string
	if importReq.KeyID != "" {
		keyID = importReq.KeyID
	} else {
		keyID = generateRandomKeyID()
	}

	wallet, err := walletStore.Import(ctx, keyID, importReq.PrivateKey, &entities.Attributes{Tags: importReq.Tags, Auth: auth.UserInfoToMap(userAuth)})
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WriteJSON(rw, formatters.FormatWalletResponse(wallet))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary      Update a wallet
// @Description  Update a wallet metadata
// @Accept       json
// @Produce      json
// @Tags         Ethereum
// @Param        storeName  path      string                         true  "Store ID"
// @Param        address    path      string                         true  "Ethereum address"
// @Param        request    body      types.UpdateEthAccountRequest  true  "Update Ethereum Account metadata request"
// @Failure      400        {object}  infrahttp.ErrorResponse        "Invalid request format"
// @Failure      401        {object}  infrahttp.ErrorResponse        "Unauthorized"
// @Failure      403        {object}  infrahttp.ErrorResponse        "Forbidden"
// @Failure      404        {object}  infrahttp.ErrorResponse        "Store/Account not found"
// @Failure      500        {object}  infrahttp.ErrorResponse        "Internal server error"
// @Success      200        {object}  types.EthAccountResponse       "Update Ethereum Account"
// @Router       /stores/{storeName}/wallets/{pubkey} [patch]
func (h *WalletsHandler) update(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	updateReq := &types.UpdateWalletRequest{}
	err := jsonutils.UnmarshalBody(request.Body, updateReq)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	userAuth := auth.UserInfoFromContext(ctx)
	walletStore, err := h.stores.Wallet(ctx, StoreNameFromContext(ctx), userAuth)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	wallet, err := walletStore.Update(ctx, getPubkey(request), &entities.Attributes{Tags: updateReq.Tags, Auth: auth.UserInfoToMap(userAuth)})
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WriteJSON(rw, formatters.FormatWalletResponse(wallet))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary      Sign a message (EIP-191)
// @Description  Sign a message, following EIP-191, using an existing Ethereum Account
// @Tags         Ethereum
// @Accept       json
// @Produce      plain
// @Param        storeName  path      string                    true  "Store ID"
// @Param        address    path      string                    true  "Ethereum address"
// @Param        request    body      types.SignMessageRequest  true  "Sign message request"
// @Success      200        {string}  string                    "Signed payload signature"
// @Failure      400        {object}  infrahttp.ErrorResponse   "Invalid request format"
// @Failure      401        {object}  infrahttp.ErrorResponse   "Unauthorized"
// @Failure      403        {object}  infrahttp.ErrorResponse   "Forbidden"
// @Failure      404        {object}  infrahttp.ErrorResponse   "Store/Account not found"
// @Failure      500        {object}  infrahttp.ErrorResponse   "Internal server error"
// @Router       /stores/{storeName}/wallets/{pubkey}/sign-message [post]
func (h *WalletsHandler) sign(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	signPayloadReq := &types.SignWalletRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signPayloadReq)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, errors.InvalidFormatError(err.Error()))
		return
	}

	walletStore, err := h.stores.Wallet(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	signature, err := walletStore.Sign(ctx, getPubkey(request), signPayloadReq.TypeSign, signPayloadReq.Data)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	//err = infrahttp.WriteJSON(rw, formatters.FormatSignatureResponse(signature, signPayloadReq, getPubkey(request)))
	_, err = rw.Write([]byte(hexutil.Encode(signature)))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary      Get a wallet
// @Description  Fetch a wallet data by its address
// @Tags         Ethereum
// @Accept       json
// @Produce      json
// @Param        storeName  path      string                    true   "Store ID"
// @Param        address    path      string                    true   "Ethereum address"
// @Param        deleted    query     bool                      false  "filter by only deleted accounts"
// @Failure      404        {object}  infrahttp.ErrorResponse   "Store/Account not found"
// @Failure      401        {object}  infrahttp.ErrorResponse   "Unauthorized"
// @Failure      403        {object}  infrahttp.ErrorResponse   "Forbidden"
// @Failure      500        {object}  infrahttp.ErrorResponse   "Internal server error"
// @Success      200        {object}  types.EthAccountResponse  "Ethereum Account data"
// @Router       /stores/{storeName}/wallets/{pubkey} [get]
func (h *WalletsHandler) getOne(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	walletStore, err := h.stores.Wallet(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
	getDeleted := request.URL.Query().Get("deleted")
	var wallet *entities.Wallet
	if getDeleted == "" {
		wallet, err = walletStore.Get(ctx, getPubkey(request))
	} else {
		wallet, err = walletStore.GetDeleted(ctx, getPubkey(request))
	}
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WriteJSON(rw, formatters.FormatWalletResponse(wallet))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary      List Ethereum accounts
// @Description  List Ethereum accounts addresses allocated in the targeted Store
// @Tags         Ethereum
// @Accept       json
// @Produce      json
// @Param        storeName   path      string                   true   "Store ID"
// @Param        deleted     query     bool                     false  "filter by only deleted accounts"
// @Param        chain_uuid  query     string                   false  "Chain UUID"
// @Param        limit       query     int                      false  "page size"
// @Param        page        query     int                      false  "page number"
// @Success      200         {array}   infrahttp.PageResponse   "Ethereum Account list"
// @Failure      401         {object}  infrahttp.ErrorResponse  "Unauthorized"
// @Failure      403         {object}  infrahttp.ErrorResponse  "Forbidden"
// @Failure      500         {object}  infrahttp.ErrorResponse  "Internal server error"
// @Router       /stores/{storeName}/wallets [get]
func (h *WalletsHandler) list(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	walletStore, err := h.stores.Wallet(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	limit, offset, err := getLimitOffset(request)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	getDeleted := request.URL.Query().Get("deleted")
	var addresses []string
	if getDeleted == "" {
		addresses, err = walletStore.List(ctx, limit, offset)
	} else {
		addresses, err = walletStore.ListDeleted(ctx, limit, offset)
	}
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = infrahttp.WritePagingResponse(rw, request, addresses)
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}
}

// @Summary      Delete Ethereum Account
// @Description  Soft delete a wallet, can be recovered
// @Tags         Ethereum
// @Accept       json
// @Param        storeName  path  string  true  "Store ID"
// @Param        address    path  string  true  "Ethereum address"
// @Success      204        "Deleted successfully"
// @Failure      401        {object}  infrahttp.ErrorResponse  "Unauthorized"
// @Failure      403        {object}  infrahttp.ErrorResponse  "Forbidden"
// @Failure      404        {object}  infrahttp.ErrorResponse  "Store/Account not found"
// @Failure      500        {object}  infrahttp.ErrorResponse  "Internal server error"
// @Router       /stores/{storeName}/wallets/{pubkey} [delete]
func (h *WalletsHandler) delete(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	walletStore, err := h.stores.Wallet(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = walletStore.Delete(ctx, getPubkey(request))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary      Destroy Wallet
// @Description  Hard delete a wallet, cannot be recovered
// @Tags         Ethereum
// @Accept       json
// @Param        storeName  path  string  true  "Store ID"
// @Param        address    path  string  true  "Ethereum address"
// @Success      204        "Destroyed successfully"
// @Failure      401        {object}  infrahttp.ErrorResponse  "Unauthorized"
// @Failure      403        {object}  infrahttp.ErrorResponse  "Forbidden"
// @Failure      404        {object}  infrahttp.ErrorResponse  "Store/Account not found"
// @Failure      500        {object}  infrahttp.ErrorResponse  "Internal server error"
// @Router       /stores/{storeName}/wallets/{pubkey}/destroy [delete]
func (h *WalletsHandler) destroy(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	walletStore, err := h.stores.Wallet(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = walletStore.Destroy(ctx, getPubkey(request))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary      Restore Wallet
// @Description  Recover a soft-deleted Ethereum Account
// @Tags         Ethereum
// @Accept       json
// @Param        storeName  path  string  true  "Store ID"
// @Param        address    path  string  true  "Ethereum address"
// @Success      204        "Restored successfully"
// @Failure      401        {object}  infrahttp.ErrorResponse  "Unauthorized"
// @Failure      403        {object}  infrahttp.ErrorResponse  "Forbidden"
// @Failure      404        {object}  infrahttp.ErrorResponse  "Store/Account not found"
// @Failure      500        {object}  infrahttp.ErrorResponse  "Internal server error"
// @Router       /stores/{storeName}/wallets/{pubkey}/restore [put]
func (h *WalletsHandler) restore(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	walletStore, err := h.stores.Wallet(ctx, StoreNameFromContext(ctx), auth.UserInfoFromContext(ctx))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	err = walletStore.Restore(ctx, getPubkey(request))
	if err != nil {
		infrahttp.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
