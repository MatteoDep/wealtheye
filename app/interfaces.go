package app

import "time"

type Service interface {
	GetAssets() ([]Asset, error)
	GetAsset(symbol string) (Asset, error)
	GetPrices(asset Asset, fromTimestamp time.Time, toTimestamp time.Time) ([]Price, error)
}

type PriceApi interface {
	GetDailyPrices(asset Asset, timestamps []time.Time) ([]Price, error)
}
