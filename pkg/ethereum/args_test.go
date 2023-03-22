package walletsereum

import (
	"encoding/json"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertSendTxMsgEquals(t *testing.T, expectedMsg, msg *SendTxMsg) {
	assert.Equal(t, expectedMsg.From, msg.From, "From should be correct")
	assert.Equal(t, expectedMsg.To, msg.To, "To should be correct")
	assert.Equal(t, expectedMsg.Gas, msg.Gas, "Gas should be correct")
	assert.Equal(t, expectedMsg.GasPrice, msg.GasPrice, "GasPrice should be correct")
	assert.Equal(t, expectedMsg.Value, msg.Value, "Value should be correct")
	assert.Equal(t, expectedMsg.Nonce, msg.Nonce, "Nonce should be correct")
	assert.Equal(t, expectedMsg.Data, msg.Data, "Data should be correct")
	assert.Equal(t, expectedMsg.PrivateFrom, msg.PrivateFrom, "PrivateFrom should be correct")
	assert.Equal(t, expectedMsg.PrivateFor, msg.PrivateFor, "PrivateFor should be correct")
	assert.Equal(t, expectedMsg.PrivacyFlag, msg.PrivacyFlag, "PrivacyFlag should be correct")
	assert.Equal(t, expectedMsg.PrivacyGroupID, msg.PrivacyGroupID, "PrivacyGroupID should be correct")
}

func TestSendTxMsg(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the response
		body []byte

		expectedSendTxMsg SendTxMsg
		expectedIsPrivate bool

		expectedErrMsg string
	}{
		{
			desc: "all public fields",
			body: []byte(`{"from":"0xc94770007dda54cf92009bff0de90c06f603a09f","to":"0xfe3b557e8fb62b89f4916b721be55ceb828dbd73","gas":"0x5208","gasPrice":"0x3e8","value":"0x1","nonce":"0xf","data":"0xabcdef"}`),
			expectedSendTxMsg: SendTxMsg{
				From:     ethcommon.HexToAddress("0xc94770007dda54cf92009bff0de90c06f603a09f"),
				To:       func(addr ethcommon.Address) *ethcommon.Address { return &addr }(ethcommon.HexToAddress("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73")),
				Gas:      func(i uint64) *uint64 { return &i }(21000),
				GasPrice: big.NewInt(1000),
				Value:    big.NewInt(1),
				Nonce:    func(i uint64) *uint64 { return &i }(15),
				Data:     func(b []byte) *[]byte { return &b }(ethcommon.FromHex("0xabcdef")),
			},
		},
		{
			desc: "partial public and private fields",
			body: []byte(`{"from":"0xc94770007dda54cf92009bff0de90c06f603a09f","to":"0xfe3b557e8fb62b89f4916b721be55ceb828dbd73","value":"0x1","privateFrom":"GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=","privateFor":["KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=","eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="],"privacyFlag":3,"privacyGroupId":"kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="}`),
			expectedSendTxMsg: SendTxMsg{
				From:  ethcommon.HexToAddress("0xc94770007dda54cf92009bff0de90c06f603a09f"),
				To:    func(addr ethcommon.Address) *ethcommon.Address { return &addr }(ethcommon.HexToAddress("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73")),
				Value: big.NewInt(1),
				PrivateArgs: PrivateArgs{
					PrivateFrom:    func(s string) *string { return &s }("GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY="),
					PrivateFor:     func(s []string) *[]string { return &s }([]string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}),
					PrivacyFlag:    func(i PrivacyFlag) *PrivacyFlag { return &i }(3),
					PrivacyGroupID: func(s string) *string { return &s }("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="),
				},
			},
			expectedIsPrivate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(SendTxMsg)
			err := json.Unmarshal(tt.body, msg)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "Unmarshal must not fail")
				assertSendTxMsgEquals(t, &tt.expectedSendTxMsg, msg)

				if tt.expectedIsPrivate {
					assert.True(t, msg.IsPrivate(), "IsPrivate")
				} else {
					assert.False(t, msg.IsPrivate(), "IsPrivate")
				}

				b, err := json.Marshal(msg)
				require.NoError(t, err, "Marshal must not fail")
				assert.Equal(t, tt.body, b, "Marshal body should match")
			}
		})
	}
}

func assertSendEEATxMsgEquals(t *testing.T, expectedMsg, msg *SendEEATxMsg) {
	assert.Equal(t, expectedMsg.From, msg.From, "From should be correct")
	assert.Equal(t, expectedMsg.To, msg.To, "To should be correct")
	assert.Equal(t, expectedMsg.Nonce, msg.Nonce, "Nonce should be correct")
	assert.Equal(t, expectedMsg.Data, msg.Data, "Data should be correct")
	assert.Equal(t, expectedMsg.PrivateFrom, msg.PrivateFrom, "PrivateFrom should be correct")
	assert.Equal(t, expectedMsg.PrivateFor, msg.PrivateFor, "PrivateFor should be correct")
	assert.Equal(t, expectedMsg.PrivacyGroupID, msg.PrivacyGroupID, "PrivacyGroupID should be correct")
}

func TestSendEEATxMsg(t *testing.T) {
	tests := []struct {
		desc string

		// JSON body of the response
		body []byte

		expectedSendEEATxMsg SendEEATxMsg

		expectedErrMsg string
	}{
		{
			desc: "partial public and private fields",
			body: []byte(`{"from":"0xc94770007dda54cf92009bff0de90c06f603a09f","to":"0xfe3b557e8fb62b89f4916b721be55ceb828dbd73","privateFrom":"GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=","privateFor":["KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=","eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="],"privacyGroupId":"kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="}`),
			expectedSendEEATxMsg: SendEEATxMsg{
				From: ethcommon.HexToAddress("0xc94770007dda54cf92009bff0de90c06f603a09f"),
				To:   func(addr ethcommon.Address) *ethcommon.Address { return &addr }(ethcommon.HexToAddress("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73")),
				PrivateArgs: PrivateArgs{
					PrivateFrom:    func(s string) *string { return &s }("GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY="),
					PrivateFor:     func(s []string) *[]string { return &s }([]string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}),
					PrivacyFlag:    func(i PrivacyFlag) *PrivacyFlag { return &i }(3),
					PrivacyGroupID: func(s string) *string { return &s }("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			msg := new(SendEEATxMsg)
			err := json.Unmarshal(tt.body, msg)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err, "Unmarshal must not fail")
				assertSendEEATxMsgEquals(t, &tt.expectedSendEEATxMsg, msg)

				b, err := json.Marshal(msg)
				require.NoError(t, err, "Marshal must not fail")
				assert.Equal(t, tt.body, b, "Marshal body should match")
			}
		})
	}
}
