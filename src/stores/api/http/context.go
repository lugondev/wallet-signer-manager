package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lugondev/tx-builder/pkg/utils"
	"github.com/lugondev/wallet-signer-manager/pkg/common"
	"net/http"
	"strings"
)

type ctxKeyType string

const storeNameCtxKey ctxKeyType = "storeName"

func WithStoreName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, storeNameCtxKey, name)
}

func StoreNameFromContext(ctx context.Context) string {
	name, ok := ctx.Value(storeNameCtxKey).(string)
	if ok {
		return name
	}

	return ""
}

func generateRandomKeyID() string {
	return fmt.Sprintf("%s%s", "signer", common.RandString(15))
}

func getPubkey(request *http.Request) string {
	pubkey, err := utils.ParseHexToCompressedPublicKey(mux.Vars(request)["pubkey"])
	if err != nil {
		return ""
	}
	return strings.ToLower(pubkey)
}
