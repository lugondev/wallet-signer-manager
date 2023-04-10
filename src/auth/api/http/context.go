package http

import (
	"context"

	"github.com/lugondev/signer-key-manager/src/auth/entities"
)

type contextKey struct{}

func UserInfoFromContext(ctx context.Context) *entities.UserInfo {
	if reqCtx, ok := ctx.Value(contextKey{}).(*entities.UserInfo); ok {
		return reqCtx
	}
	return nil
}

func UserInfoToMap(userInfo *entities.UserInfo) map[string]interface{} {
	return map[string]interface{}{
		"username":    userInfo.Username,
		"roles":       userInfo.Roles,
		"authMode":    userInfo.AuthMode,
		"tenant":      userInfo.Tenant,
		"permissions": userInfo.Permissions,
	}
}

func WithUserInfo(ctx context.Context, reqCtx *entities.UserInfo) context.Context {
	return context.WithValue(ctx, contextKey{}, reqCtx)
}
