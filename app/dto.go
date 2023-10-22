package app

import "time"

type WalletTransferType string

const (
	Deposit = "from"
	Withdrawal = "to"
)

type WalletTransferDTO struct {
	Timestamp       time.Time
	Ammount         float64
	AssetSymbol     string
	Type            WalletTransferType
	OtherWalletId   int
	OtherWalletName string
}
