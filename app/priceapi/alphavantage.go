package priceapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MatteoDep/wealtheye/app"
)

type AlphaVantageApi struct {
	Cfg *app.PriceApiConfig
}

// TODO align with twelvedata implementation
func (p *AlphaVantageApi) GetDailyPricesUsd(
    asset *app.Asset,
    fromTimestamp time.Time,
    toTimestamp time.Time,
) ([]app.Price, error) {
    if asset.Symbol == "USD" {
        return nil, nil
    }

    var priceLabel string
    timeSeriesLabel := "Time Series "
	reqUrl := "https://www.alphavantage.co/query"
    switch asset.Type {
    case "forex":
        reqUrl += "?function=FX_DAILY&from_symbol=" + asset.Symbol
        reqUrl += "&to_symbol=USD"
        priceLabel = "4. close"
        timeSeriesLabel += "FX (Daily)"
    case "crypto":
        reqUrl += "?function=DIGITAL_CURRENCY_DAILY&symbol=" + asset.Symbol
        reqUrl += "&market=USD"
        priceLabel = "4a. close (USD)"
        timeSeriesLabel += "(Digital Currency Daily)"
    case "stock":
        return nil, NotImplementedAssetTypeError{asset: asset, provider: p.Cfg.Provider}
    case "commodity":
        return nil, NotImplementedAssetTypeError{asset: asset, provider: p.Cfg.Provider}
    case "bond":
        return nil, NotImplementedAssetTypeError{asset: asset, provider: p.Cfg.Provider}
    case "etf":
        return nil, UnsupportedAssetTypeError{asset: asset, provider: p.Cfg.Provider}
    default:
        return nil, UndefinedAssetTypeError{asset: asset}
    }
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
    fromTimestamp = fromTimestamp.Round(24 * time.Hour)
    toTimestamp = toTimestamp.Round(24 * time.Hour)
    json.Unmarshal([]byte(body), &result)
    if result[timeSeriesLabel] == nil {
        return nil, fmt.Errorf("Error occured while retrieving data. %s returned %v.", p.Cfg.Provider, result)
    }
    timeSeries := result[timeSeriesLabel].(map[string]any)
    // -2 days to account for weekends
    earlierFromTimestamp := fromTimestamp.AddDate(0, 0, -2)
    prices := []app.Price{}
    for date, priceMap := range timeSeries {
        timestamp, err := time.Parse(time.DateOnly, date)
        if err != nil {
            return nil, err
        }

        if timestamp.Before(earlierFromTimestamp) || timestamp.After(toTimestamp) {
            continue
        }

        value, err := strconv.ParseFloat(
            priceMap.(map[string]any)[priceLabel].(string),
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

    prices = ffillPrices(prices, toTimestamp)

    startIndex := 0
    for _, price := range prices {
        if price.TimestampUtc.Before(fromTimestamp) {
            startIndex++
        } else {
            break
        }
    }

    return prices[startIndex:], nil
}
