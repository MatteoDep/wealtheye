package alphavantageapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MatteoDep/wealtheye/app"
)

type PriceApi struct {
	Cfg *app.Config
}

func (p *PriceApi) GetDailyPrices(assetSymbol string, numDays int) ([]app.Price, error) {
    prices := []app.Price{}
    if assetSymbol == "USD" {
        log.Println(assetSymbol, "is USD")
        // price := app.Price{
        //     TimestampUtc: timestamp,
        //     AssetSymbol: assetSymbol,
        //     ValueUsd: value,
        // }
        return prices, nil
    }

	reqUrl := "https://www.alphavantage.co/query"
	reqUrl += "?function=FX_DAILY&from_symbol=" + assetSymbol
	reqUrl += "&to_symbol=USD"
	reqUrl += "&apikey=" + p.Cfg.PriceApi.Key
    log.Println(assetSymbol)
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
	if err != nil {
        return prices, err
	}
    log.Println(body)

    var result map[string]any
    json.Unmarshal([]byte(body), &result)
    log.Println(result)
    timeSeries := result["Time Series FX (Daily)"].(map[string]any)
    daysCount := 0
    for date, priceMap := range timeSeries {
        timestamp, err := time.Parse(time.DateOnly, date)
        if err != nil {
            return prices, err
        }
        value, err := strconv.ParseFloat(
            priceMap.(map[string]any)["4. close"].(string),
            64,
        )
        if err != nil {
            return prices, err
        }
        price := app.Price{
            TimestampUtc: timestamp,
            AssetSymbol: assetSymbol,
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
