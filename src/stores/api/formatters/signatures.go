package formatters

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lugondev/wallet-signer-manager/src/stores/api/types"
)

func FormatSignatureResponse(signature []byte, payload *types.SignWalletRequest, pubkey string) map[string]interface{} {
	return map[string]interface{}{
		"signature": hexutil.Encode(signature),
		"payload":   payload,
		"pubkey":    pubkey,
	}
}
