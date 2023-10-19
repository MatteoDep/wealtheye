package app

import "time"

type Service interface {
	GetAssets() ([]Asset, error)
	GetAsset(symbol string) (Asset, error)

	GetPrices(asset Asset, fromTimestamp time.Time, toTimestamp time.Time) ([]Price, error)
	PostPrices(prices []Price) error

	GetWallets() ([]Wallet, error)
	GetWallet(id int) (Wallet, error)
	PostWallet(wallet Wallet) error
	PutWallet(wallet Wallet) error

    GetTransfers() ([]Transfer, error)
    GetWalletTransfers(walletId int) ([]Transfer, error)
    PostTransfer(transfer Transfer) error
    PutTransfer(transfer Transfer) error
}

type PriceApi interface {
	GetDailyPrices(asset Asset, timestamps []time.Time) ([]Price, error)
}
