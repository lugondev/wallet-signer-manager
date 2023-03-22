package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type CreateWalletRequest struct {
	KeyID string            `json:"keyId,omitempty" example:"my-key-account"`
	Tags  map[string]string `json:"tags,omitempty"`
}

type ImportWalletRequest struct {
	KeyID      string            `json:"keyId,omitempty" example:"my-imported-key-account"`
	PrivateKey hexutil.Bytes     `json:"privateKey" validate:"required" example:"0x56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E" swaggertype:"string"`
	Tags       map[string]string `json:"tags,omitempty"`
}

type UpdateWalletRequest struct {
	Tags map[string]string `json:"tags,omitempty"`
}

type SignWalletRequest struct {
	Message hexutil.Bytes `json:"message" validate:"required" example:"0xfeade..." swaggertype:"string"`
}

type WalletResponse struct {
	PublicKey           hexutil.Bytes     `json:"publicKey" example:"0x1abae27a0cbfb02945720425d3b80c7e09728534" swaggertype:"string"`
	CompressedPublicKey hexutil.Bytes     `json:"compressedPublicKey" example:"0x6019a3c8..." swaggertype:"string"`
	CreatedAt           time.Time         `json:"createdAt" example:"2020-07-09T12:35:42.115395Z"`
	UpdatedAt           time.Time         `json:"updatedAt" example:"2020-07-09T12:35:42.115395Z"`
	DeletedAt           *time.Time        `json:"deletedAt,omitempty" example:"2020-07-09T12:35:42.115395Z"`
	KeyID               string            `json:"keyId" example:"my-key-id"`
	Tags                map[string]string `json:"tags,omitempty"`
	Disabled            bool              `json:"disabled" example:"false"`
}
