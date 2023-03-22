package formatters

import (
	"github.com/lugondev/signer-key-manager/src/stores/api/types"
	"github.com/lugondev/signer-key-manager/src/stores/entities"
)

func FormatWalletResponse(wallet *entities.Wallet) *types.WalletResponse {
	resp := &types.WalletResponse{
		PublicKey:           wallet.PublicKey,
		CompressedPublicKey: wallet.CompressedPublicKey,
		Tags:                wallet.Tags,
		Disabled:            wallet.Metadata.Disabled,
		CreatedAt:           wallet.Metadata.CreatedAt,
		UpdatedAt:           wallet.Metadata.UpdatedAt,
	}

	if !wallet.Metadata.DeletedAt.IsZero() {
		resp.DeletedAt = &wallet.Metadata.DeletedAt
	}

	return resp
}
