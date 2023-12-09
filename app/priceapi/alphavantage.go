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

func (p *AlphaVantageApi) GetDailyPricesUsd(
    asset app.Asset,
    fromTimestamp time.Time,
    toTimestamp time.Time,
) ([]app.Price, error) {
    prices := []app.Price{}
    if asset.Symbol == "USD" {
        // todo
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
    fromTimestamp = fromTimestamp.Round(24 * time.Hour).Add(-time.Hour)
    toTimestamp = toTimestamp.Round(24 * time.Hour).Add(time.Hour)
    json.Unmarshal([]byte(body), &result)
    timeSeries := result[timeSeriesLabel].(map[string]any)
    // -2 days to account for weekends
    earlierFromTimestamp := fromTimestamp.AddDate(0, 0, -2)
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