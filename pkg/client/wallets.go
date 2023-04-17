package client

import (
	"context"
	"fmt"

	"github.com/lugondev/wallet-signer-manager/src/stores/api/types"
)

const walletPath = "wallets"

func (c *HTTPClient) CreateWallet(ctx context.Context, storeName string, req *types.CreateWalletRequest) (*types.WalletResponse, error) {
	ethAcc := &types.WalletResponse{}
	reqURL := fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), walletPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, ethAcc)
	if err != nil {
		return nil, err
	}

	return ethAcc, nil
}

func (c *HTTPClient) ImportWallet(ctx context.Context, storeName string, req *types.ImportWalletRequest) (*types.WalletResponse, error) {
	ethAcc := &types.WalletResponse{}
	reqURL := fmt.Sprintf("%s/%s/import", withURLStore(c.config.URL, storeName), walletPath)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, ethAcc)
	if err != nil {
		return nil, err
	}

	return ethAcc, nil
}

func (c *HTTPClient) UpdateWallet(ctx context.Context, storeName, address string, req *types.UpdateWalletRequest) (*types.WalletResponse, error) {
	ethAcc := &types.WalletResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), walletPath, address)
	response, err := patchRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, ethAcc)
	if err != nil {
		return nil, err
	}

	return ethAcc, nil
}

func (c *HTTPClient) Sign(ctx context.Context, storeName, address string, req *types.SignWalletRequest) (string, error) {
	reqURL := fmt.Sprintf("%s/%s/%s/sign", withURLStore(c.config.URL, storeName), walletPath, address)
	response, err := postRequest(ctx, c.client, reqURL, req)
	if err != nil {
		return "", err
	}

	defer closeResponse(response)
	return parseStringResponse(response)
}

func (c *HTTPClient) GetWallet(ctx context.Context, storeName, address string) (*types.WalletResponse, error) {
	acc := &types.WalletResponse{}
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), walletPath, address)

	response, err := getRequest(ctx, c.client, reqURL)
	if err != nil {
		return nil, err
	}

	defer closeResponse(response)
	err = parseResponse(response, acc)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (c *HTTPClient) ListWallets(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return listRequest(ctx, c.client, fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), walletPath), false, limit, page)
}

func (c *HTTPClient) ListDeletedWallets(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return listRequest(ctx, c.client, fmt.Sprintf("%s/%s", withURLStore(c.config.URL, storeName), walletPath), true, limit, page)
}

func (c *HTTPClient) DeleteWallet(ctx context.Context, storeName, address string) error {
	reqURL := fmt.Sprintf("%s/%s/%s", withURLStore(c.config.URL, storeName), walletPath, address)
	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) DestroyWallet(ctx context.Context, storeName, address string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/destroy", withURLStore(c.config.URL, storeName), walletPath, address)
	response, err := deleteRequest(ctx, c.client, reqURL)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}

func (c *HTTPClient) RestoreWallet(ctx context.Context, storeName, address string) error {
	reqURL := fmt.Sprintf("%s/%s/%s/restore", withURLStore(c.config.URL, storeName), walletPath, address)
	response, err := putRequest(ctx, c.client, reqURL, nil)
	if err != nil {
		return err
	}

	defer closeResponse(response)
	return parseEmptyBodyResponse(response)
}
