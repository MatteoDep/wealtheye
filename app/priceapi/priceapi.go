package priceapi

import (
	"strings"

	"github.com/MatteoDep/wealtheye/app"
)


func GetPriceApi(cfg *app.PriceApiConfig) app.PriceApi {
    switch strings.ToLower(cfg.Provider) {
    case "alphavantage":
        return &AlphaVantageApi{
            Cfg: cfg,
        }
    default:
        return &AlphaVantageApi{
            Cfg: cfg,
        }
    }
}
