package client

import (
	"context"

	"github.com/lugondev/signer-key-manager/pkg/jsonrpc"
	storestypes "github.com/lugondev/signer-key-manager/src/stores/api/types"
)

type WalletClient interface {
	CreateWallet(ctx context.Context, storeName string, request *storestypes.CreateWalletRequest) (*storestypes.WalletResponse, error)
	ImportWallet(ctx context.Context, storeName string, request *storestypes.ImportWalletRequest) (*storestypes.WalletResponse, error)
	UpdateWallet(ctx context.Context, storeName, pubkey string, request *storestypes.UpdateWalletRequest) (*storestypes.WalletResponse, error)
	Sign(ctx context.Context, storeName, account string, request *storestypes.SignWalletRequest) (string, error)
	GetWallet(ctx context.Context, storeName, pubkey string) (*storestypes.WalletResponse, error)
	ListWallets(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	ListDeletedWallets(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	DeleteWallet(ctx context.Context, storeName, pubkey string) error
	DestroyWallet(ctx context.Context, storeName, pubkey string) error
	RestoreWallet(ctx context.Context, storeName, pubkey string) error
}

type JSONRPC interface {
	Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error)
}

type KeyManagerClient interface {
	WalletClient
	JSONRPC
}
