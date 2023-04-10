package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
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
	pubkeyRequest := strings.ToLower(mux.Vars(request)["pubkey"])
	if strings.HasPrefix(pubkeyRequest, "0x") {
		return pubkeyRequest[2:]
	}
	return pubkeyRequest
}
