package app

import "time"

type WalletTransferType struct {
	Action      string
	Preposition string
}

var Deposit = WalletTransferType{
	Action:      "Deposit",
	Preposition: "from",
}

var Withdrawal = WalletTransferType{
	Action:      "Withdrawal",
	Preposition: "to",
}

type WalletTransferDTO struct {
	Timestamp       time.Time
	Ammount         float64
	AssetSymbol     string
	TypeAction      string
	OtherWalletId   int
	OtherWalletName string
}
