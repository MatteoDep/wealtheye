package alphavantageapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MatteoDep/wealtheye/app"
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

func (p *PriceApi) GetDailyPrices(asset app.Asset, numDays int) ([]app.Price, error) {
    prices := []app.Price{}
    if asset.Symbol == "USD" {
        return prices, nil
    }

    var priceLabel string
    timeSeriesLabel := "Time Series "
	reqUrl := "https://www.alphavantage.co/query"
    if asset.Type == "physical currency" {
        reqUrl += "?function=FX_DAILY&from_symbol=" + asset.Symbol
        reqUrl += "&to_symbol=USD"
        priceLabel = "4. close"
        timeSeriesLabel += "FX (Daily)"
    } else if asset.Type == "digital currency" {
        reqUrl += "?function=DIGITAL_CURRENCY_DAILY&symbol=" + asset.Symbol
        reqUrl += "&market=USD"
        priceLabel = "4a. close (USD)"
        timeSeriesLabel += "(Digital Currency Daily)"
    } else {
        log.Println("Undefined asset type", asset.Type)
        return nil, UndefinedTypeError{asset: &asset}
    }
	reqUrl += "&apikey=" + p.Cfg.PriceApi.Key
    log.Println(asset.Symbol)
    log.Println(reqUrl)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
        return prices, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
        return prices, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil || (strings.Fields(string(body))[0] == "Error") {
        return prices, err
	}

    var result map[string]any
    json.Unmarshal([]byte(body), &result)
    log.Println(result)
    timeSeries := result[timeSeriesLabel].(map[string]any)
    daysCount := 0
    for date, priceMap := range timeSeries {
        timestamp, err := time.Parse(time.DateOnly, date)
        if err != nil {
            return prices, err
        }
        value, err := strconv.ParseFloat(
            priceMap.(map[string]any)[priceLabel].(string),
            64,
        )
        if err != nil {
            return prices, err
        }
        price := app.Price{
            TimestampUtc: timestamp,
            AssetSymbol: asset.Symbol,
            ValueUsd: value,
        }
        prices = append(prices, price)

        daysCount++
        if daysCount == numDays {
            break
        }
    }

    return prices, nil
}
