package formatters

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lugondev/signer-key-manager/src/stores/api/types"
	"github.com/lugondev/signer-key-manager/src/stores/entities"
	"github.com/lugondev/tx-builder/blockchain/bitcoin"
)

func FormatWalletResponse(wallet *entities.Wallet) *types.WalletResponse {
	resp := &types.WalletResponse{
		PublicKey:           hexutil.Encode(wallet.PublicKey),
		CompressedPublicKey: hexutil.Encode(wallet.CompressedPublicKey),
		Tags:                wallet.Tags,
		Disabled:            wallet.Metadata.Disabled,
		CreatedAt:           wallet.Metadata.CreatedAt,
		UpdatedAt:           wallet.Metadata.UpdatedAt,
		Addresses:           *FormatAddressesResponse(wallet),
	}

	if !wallet.Metadata.DeletedAt.IsZero() {
		resp.DeletedAt = &wallet.Metadata.DeletedAt
	}

	return resp
}

func FormatAddressesResponse(wallets *entities.Wallet) *types.AddressesResponse {
	pubkey, err := crypto.UnmarshalPubkey(wallets.PublicKey)
	if err != nil {
		return nil
	}

	parsePubKey, err := btcec.ParsePubKey(wallets.PublicKey)
	if err != nil {
		return nil
	}
	testnet3 := bitcoin.PubkeyToAddresses(parsePubKey, &chaincfg.TestNet3Params)
	mainnet := bitcoin.PubkeyToAddresses(parsePubKey, &chaincfg.MainNetParams)
	return &types.AddressesResponse{
		Evm: crypto.PubkeyToAddress(*pubkey),
		Bitcoin: map[types.BitcoinNet]bitcoin.KeyAddresses{
			types.BtcMainnet:  mainnet,
			types.BtcTestnet3: testnet3,
		},
	}
}
