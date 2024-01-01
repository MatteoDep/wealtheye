package app

import "time"

type Repository interface {
	GetAssets() ([]Asset, error)
	GetAsset(symbol string) (*Asset, error)

	GetLastPriceTimestamp(symbol string) (time.Time, error)
	GetPrices(symbol string, fromTimestampUtc time.Time, toTimestampUtc time.Time) ([]Price, error)
	PostPrices(prices []Price) error
	UpdatePricesValue(prices []Price) error

	GetWallets() ([]Wallet, error)
	GetWallet(id int) (*Wallet, error)
	PostWallet(name string) error
	UpdateWalletName(id int, name string) error
	UpdateWalletValue(id int, valueUsd float64) error

    GetTransfers() ([]Transfer, error)
    GetWalletTransfers(walletId int) ([]Transfer, error)
    PostTransfer(transfer *Transfer) error
    UpdateTransfer(transfer *Transfer) error
}

type Service interface {
    GetAssets() ([]Asset, error)
    GetAsset(symbol string) (*Asset, error)
    GetPrice(symbol string, timestampUtc time.Time) (*Price, error)
    GetPrices(symbol string, fromTimestampUtc time.Time, toTimestampUtc time.Time) ([]Price, error)
    GetWallets() ([]Wallet, error)
    GetWallet(id int) (*Wallet, error)
    GetTransfers() ([]Transfer, error)
    GetWalletTransfers(walletId int) ([]Transfer, error)
    UpdateWalletsValue() (error)
    UpdateWalletValue(id int) (error)
	UpdateWalletName(id int, name string) error
	PostWallet(name string) error
    PostTransfer(transfer *Transfer) error
    UpdateTransfer(transfer *Transfer) error

    TransferToWalletTransferDTO(transfer *Transfer, walletId int) (*WalletTransferDTO, error)
    WalletTransferDTOToTransfer(walletTransferDTO *WalletTransferDTO, walletId int) (*Transfer, error)
    GetExternalWalletName(walletTransferType WalletTransferType) string
}

type PriceApi interface {
	GetDailyPricesUsd(asset *Asset, fromTimestamp time.Time, toTimestamp time.Time) ([]Price, error)
}
