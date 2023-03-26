package types

import (
	"github.com/lugondev/tx-builder/blockchain/bitcoin"
)

type BitcoinNet string

const (
	BtcMainnet  BitcoinNet = "mainnet"
	BtcTestnet3            = "testnet3"
)

type BitcoinAddresses map[BitcoinNet]bitcoin.KeyAddresses
