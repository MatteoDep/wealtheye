package priceapi

import (
	"fmt"
	"strings"
	"time"

	"github.com/MatteoDep/wealtheye/app"
)


type UndefinedAssetTypeError struct {
    asset *app.Asset
}

func (ufe UndefinedAssetTypeError) Error() string {
    msg := fmt.Sprintf(
        "Could not find price for asset %s (%s). Undefined type '%s'.",
        ufe.asset.Symbol,
        ufe.asset.Name,
        ufe.asset.Type,
    )
    return msg
}

type UnsupportedAssetTypeError struct {
    asset *app.Asset
    provider string
}

func (ufe UnsupportedAssetTypeError) Error() string {
    msg := fmt.Sprintf(
        "Could not find price for asset %s (%s). Type '%s' is not supported by %s provider.",
        ufe.asset.Symbol,
        ufe.asset.Name,
        ufe.asset.Type,
        ufe.provider,
    )
    return msg
}

type NotImplementedAssetTypeError struct {
    asset *app.Asset
    provider string
}

func (ufe NotImplementedAssetTypeError) Error() string {
    msg := fmt.Sprintf(
        "Could not find price for asset %s (%s). Type '%s' is not yet implemented for %s provider.",
        ufe.asset.Symbol,
        ufe.asset.Name,
        ufe.asset.Type,
        ufe.provider,
    )
    return msg
}

func GetPriceApi(cfg *app.PriceApiConfig) app.PriceApi {
    switch strings.ToLower(cfg.Provider) {
    case "alphavantage":
        return &AlphaVantageApi{
            Cfg: cfg,
        }
    case "twelvedata":
        return &TwelveDataApi{
            Cfg: cfg,
        }
    default:
        return &TwelveDataApi{
            Cfg: cfg,
        }
    }
}

func GetPriceUsd(timestamp time.Time) app.Price {
    return app.Price{
        AssetSymbol: "USD",
        TimestampUtc: timestamp,
        ValueUsd: 1,
    }
}

func ffillPrices(prices []app.Price, toTimestamp time.Time) []app.Price {
    app.SortPrices(prices)
    prevPrice := prices[0]
    nextTimestamp := prices[0].TimestampUtc.AddDate(0, 0, 1)
    filledPrices := []app.Price{}
    for _, price := range prices[1:] {
        for price.TimestampUtc.After(nextTimestamp) {
            filledPrices = append(filledPrices, app.Price{
                AssetSymbol: prevPrice.AssetSymbol,
                TimestampUtc: nextTimestamp,
                ValueUsd: prevPrice.ValueUsd,
            })
            nextTimestamp = nextTimestamp.AddDate(0, 0, 1)
        }

        filledPrices = append(filledPrices, price)
        nextTimestamp = nextTimestamp.AddDate(0, 0, 1)
        prevPrice = price
    }

    for prevPrice.TimestampUtc.Before(toTimestamp) {
        filledPrices = append(filledPrices, app.Price{
            AssetSymbol: prevPrice.AssetSymbol,
            TimestampUtc: nextTimestamp,
            ValueUsd: prevPrice.ValueUsd,
        })
        nextTimestamp = nextTimestamp.AddDate(0, 0, 1)
        prevPrice = filledPrices[len(filledPrices) - 1]
    }

    return filledPrices
}
