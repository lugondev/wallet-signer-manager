package postgres

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"time"

	"github.com/lugondev/wallet-signer-manager/src/infra/postgres/client"
	"github.com/lugondev/wallet-signer-manager/src/stores/database/models"
	"github.com/lugondev/wallet-signer-manager/src/stores/entities"

	"github.com/lugondev/wallet-signer-manager/src/infra/log"
	"github.com/lugondev/wallet-signer-manager/src/infra/postgres"
	"github.com/lugondev/wallet-signer-manager/src/stores/database"

	"github.com/lugondev/wallet-signer-manager/pkg/errors"
)

type Wallets struct {
	storeID string
	logger  log.Logger
	client  postgres.Client
}

var _ database.Wallets = &Wallets{}

func NewWallets(storeID string, db postgres.Client, logger log.Logger) *Wallets {
	return &Wallets{
		storeID: storeID,
		logger:  logger,
		client:  db,
	}
}

func (ea *Wallets) RunInTransaction(ctx context.Context, persist func(dbTx database.Wallets) error) error {
	return ea.client.RunInTransaction(ctx, func(dbTx postgres.Client) error {
		ea.client = dbTx
		return persist(ea)
	})
}

func (ea *Wallets) Get(ctx context.Context, pubkey string) (*entities.Wallet, error) {
	wallet := &models.Wallet{Pubkey: pubkey, StoreID: ea.storeID}
	err := ea.client.SelectPK(ctx, wallet)
	if err != nil {

		errMessage := "failed to get account"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return wallet.ToEntity(), nil
}

func (ea *Wallets) GetDeleted(ctx context.Context, pubkey string) (*entities.Wallet, error) {
	wallet := &models.Wallet{Pubkey: pubkey, StoreID: ea.storeID}

	err := ea.client.SelectDeletedPK(ctx, wallet)
	if err != nil {
		errMessage := "failed to get deleted account"
		ea.logger.With("pubkey", pubkey).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return wallet.ToEntity(), nil
}

func (ea *Wallets) GetAll(ctx context.Context) ([]*entities.Wallet, error) {
	var wallets []*models.Wallet

	err := ea.client.SelectWhere(ctx, &wallets, "store_id = ?", []string{}, ea.storeID)
	if err != nil {
		errMessage := "failed to get all accounts"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var accounts []*entities.Wallet
	for _, acc := range wallets {
		accounts = append(accounts, acc.ToEntity())
	}

	return accounts, nil
}

func (ea *Wallets) GetAllDeleted(ctx context.Context) ([]*entities.Wallet, error) {
	var wallets []*models.Wallet

	err := ea.client.SelectDeletedWhere(ctx, &wallets, "store_id = ?", ea.storeID)
	if err != nil {
		errMessage := "failed to get all deleted accounts"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	var accounts []*entities.Wallet
	for _, acc := range wallets {
		accounts = append(accounts, acc.ToEntity())
	}

	return accounts, nil
}

func (ea *Wallets) SearchAddresses(ctx context.Context, isDeleted bool, limit, offset uint64) ([]string, error) {
	ids, err := client.QuerySearchIDs(ctx, ea.client, "wallets", "pubkey", "store_id = ?", []interface{}{ea.storeID}, isDeleted, limit, offset)
	if err != nil {
		errMessage := "failed to list of ethereum addresses"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return ids, nil
}

func (ea *Wallets) Add(ctx context.Context, account *entities.Wallet) (*entities.Wallet, error) {
	accModel := models.NewWallet(account)
	accModel.StoreID = ea.storeID
	accModel.CreatedAt = time.Now()
	accModel.UpdatedAt = time.Now()

	err := ea.client.Insert(ctx, accModel)
	if err != nil {
		errMessage := "failed to add account"
		ea.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return accModel.ToEntity(), nil
}

func (ea *Wallets) Update(ctx context.Context, account *entities.Wallet) (*entities.Wallet, error) {
	accModel := models.NewWallet(account)
	accModel.StoreID = ea.storeID
	accModel.UpdatedAt = time.Now()

	err := ea.client.UpdatePK(ctx, accModel)
	if err != nil {
		errMessage := "failed to update account"
		ea.logger.With("pubkey", account.CompressedPublicKey).WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return accModel.ToEntity(), nil
}

func (ea *Wallets) Delete(ctx context.Context, pubkey string) error {
	err := ea.client.DeletePK(ctx, &models.Wallet{Pubkey: pubkey, StoreID: ea.storeID})
	if err != nil {
		errMessage := "failed to delete account"
		ea.logger.With("pubkey", pubkey).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (ea *Wallets) Restore(ctx context.Context, pubkey string) error {
	accModel := &models.Wallet{
		CompressedPublicKey: common.FromHex(pubkey),
		StoreID:             ea.storeID,
	}
	err := ea.client.UndeletePK(ctx, accModel)
	if err != nil {
		errMessage := "failed to restore account"
		ea.logger.With("address", pubkey).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (ea *Wallets) Purge(ctx context.Context, pubkey string) error {
	err := ea.client.ForceDeletePK(ctx, &models.Wallet{Pubkey: pubkey, StoreID: ea.storeID})
	if err != nil {
		errMessage := "failed to permanently delete account"
		ea.logger.With("address", pubkey).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}
