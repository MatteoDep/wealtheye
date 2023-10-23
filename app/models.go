package app

import (
	"database/sql"
	"time"
)

type Asset struct {
	Symbol      string
	Name        string
	Type        string
}

type Price struct {
	Id           int
	AssetSymbol  string
	TimestampUtc time.Time
	ValueUsd     float64
}

type Wallet struct {
	Id       int
	Name     string
	ValueUsd float64
}

type Transfer struct {
	Id           int
	TimestampUtc time.Time
	Ammount      float64
	AssetSymbol  string
	FromWalletId sql.NullInt64
	ToWalletId   sql.NullInt64
}
