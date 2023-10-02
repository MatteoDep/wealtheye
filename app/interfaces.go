package app

type Service interface {
	GetAssets() ([]Asset, error)
	GetAsset(symbol string) (Asset, error)
	GetPrices(assetSymbol string) ([]Price, error)
}

type PriceApi interface {
	GetDailyPrices(asset Asset, numDays int) ([]Price, error)
}
