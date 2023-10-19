package app

import (
	"database/sql"
	"time"
)

type Asset struct {
	Symbol      string
	Name        string
	Type        string
	ValueUsd    sql.NullFloat64
	LastSynched sql.NullTime
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
	FromWalletId int
	ToWalletId   int
}
