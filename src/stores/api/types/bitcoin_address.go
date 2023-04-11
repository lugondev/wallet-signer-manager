package types

import (
	"github.com/lugondev/tx-builder/pkg/blockchain/bitcoin/chain"
)

type BitcoinNet string

const (
	BtcMainnet  BitcoinNet = "mainnet"
	BtcTestnet3            = "testnet3"
)

type BitcoinAddresses map[BitcoinNet]chain.KeyAddresses
