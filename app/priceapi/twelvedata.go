package priceapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MatteoDep/wealtheye/app"
)

type TwelveDataApi struct {
	Cfg *app.PriceApiConfig
}

func (p *TwelveDataApi) GetDailyPricesUsd(
    asset app.Asset,
    fromTimestamp time.Time,
    toTimestamp time.Time,
) ([]app.Price, error) {
    if asset.Symbol == "USD" {
        return nil, nil
    }

	reqUrl := "https://api.twelvedata.com/time_series"
    var symbol string
    switch asset.Type {
    case "forex":
        symbol = asset.Symbol + "/USD"
    case "crypto":
        symbol = asset.Symbol + "/USD"
    case "stock":
        symbol = asset.Symbol
    case "commodity":
        return nil, NotImplementedAssetTypeError{asset: &asset, provider: p.Cfg.Provider}
    case "bond":
        return nil, NotImplementedAssetTypeError{asset: &asset, provider: p.Cfg.Provider}
    case "etf":
        return nil, UnsupportedAssetTypeError{asset: &asset, provider: p.Cfg.Provider}
    default:
        return nil, UndefinedAssetTypeError{asset: &asset}
    }
	reqUrl += "?symbol=" + symbol
	reqUrl += "&interval=1day"
	reqUrl += "&timezone=UTC"
	reqUrl += "&start_date=" + fromTimestamp.Format(time.DateOnly)
	reqUrl += "&end_date=" + toTimestamp.Format(time.DateOnly)
	reqUrl += "&apikey=" + p.Cfg.Key

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
        return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
        return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil || (strings.Fields(string(body))[0] == "Error") {
        return nil, err
	}

    var result map[string]any
    json.Unmarshal([]byte(body), &result)
    if result["values"] == nil {
        log.Printf("No Prices found for %s from %v to %v.", asset.Symbol, fromTimestamp, toTimestamp)
        return nil, nil
    }
    timeSeries := result["values"].([]any)
    prices := []app.Price{}
    var priceRecord map[string]any
    for _, priceRecord_ := range timeSeries {
        priceRecord = priceRecord_.(map[string]any)
        timestamp, err := time.Parse(time.DateOnly, priceRecord["datetime"].(string))
        if err != nil {
            return nil, err
        }

        value, err := strconv.ParseFloat(
            priceRecord["close"].(string),
            64,
        )
        if err != nil {
            return nil, err
        }
        price := app.Price{
            TimestampUtc: timestamp,
            AssetSymbol: asset.Symbol,
            ValueUsd: value,
        }
        prices = append(prices, price)
    }

    return prices, nil
}
