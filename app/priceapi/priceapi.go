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
	"golang.org/x/exp/slices"
)

type PriceApi struct {
	Cfg *app.Config
}


type UndefinedTypeError struct {
    asset *app.Asset
}

func (ufe UndefinedTypeError) Error() string {
    msg := fmt.Sprintf(
        "Could not find price for asset %s (%s). Undefined type '%s'.",
        ufe.asset.Symbol,
        ufe.asset.Name,
        ufe.asset.Type,
    )
    return msg
}

func (p *PriceApi) GetDailyPrices(
    asset app.Asset,
    timestamps []time.Time,
) ([]app.Price, error) {
    prices := []app.Price{}
    if asset.Symbol == "USD" || len(timestamps) == 0 {
        return prices, nil
    }

    var priceLabel string
    timeSeriesLabel := "Time Series "
	reqUrl := "https://www.alphavantage.co/query"
    if asset.Type == "forex" {
        reqUrl += "?function=FX_DAILY&from_symbol=" + asset.Symbol
        reqUrl += "&to_symbol=USD"
        priceLabel = "4. close"
        timeSeriesLabel += "FX (Daily)"
    } else if asset.Type == "crypto" {
        reqUrl += "?function=DIGITAL_CURRENCY_DAILY&symbol=" + asset.Symbol
        reqUrl += "&market=USD"
        priceLabel = "4a. close (USD)"
        timeSeriesLabel += "(Digital Currency Daily)"
    } else {
        return nil, UndefinedTypeError{asset: &asset}
    }
	reqUrl += "&apikey=" + p.Cfg.PriceApi.Key

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
    timeSeries := result[timeSeriesLabel].(map[string]any)
    for date, priceMap := range timeSeries {
        timestamp, err := time.Parse(time.DateOnly, date)
        if err != nil {
            return nil, err
        }

        if !slices.Contains(timestamps, timestamp) {
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

    return prices, nil
}
