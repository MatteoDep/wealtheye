package sqliteservice

import (
	"database/sql"
	"log"
	"time"

	"github.com/MatteoDep/wealtheye/app"
	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	DB *sql.DB
    PA app.PriceApi
}

func (s *Service) GetAssets() ([]app.Asset, error) {
	query_str := `
        select *
        from asset
    `
	rows, err := s.DB.Query(query_str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []app.Asset
    var asset app.Asset
	for rows.Next() {
		err := rows.Scan(
			&asset.Symbol,
            &asset.Name,
            &asset.Type,
			&asset.ValueUsd,
			&asset.LastSynched,
		)
		if err != nil {
			return assets, err
		}
		assets = append(assets, asset)
	}

	if err := rows.Err(); err != nil {
		return assets, err
	}

	return assets, nil
}

func (s *Service) GetAsset(symbol string) (app.Asset, error) {
	query_str := `
        SELECT *
        FROM asset
        WHERE symbol = $1
    `
	row := s.DB.QueryRow(query_str, symbol)

	var asset app.Asset
	err := row.Scan(
		&asset.Symbol,
        &asset.Name,
        &asset.Type,
		&asset.ValueUsd,
		&asset.LastSynched,
	)
	if err != nil {
		return asset, err
	}

	if err := row.Err(); err != nil {
		return asset, err
	}

	return asset, nil
}

func (s *Service) GetPrices(
    asset app.Asset,
    fromTimestamp time.Time,
    toTimestamp time.Time,
) ([]app.Price, error) {
    prices := []app.Price{}
    if asset.Symbol == "USD" || toTimestamp.Sub(fromTimestamp) < 24 * time.Hour {
        return prices, nil
    }

	query_str := `
        SELECT *
        FROM price_daily
        WHERE asset_symbol = $1
        AND timestamp_utc >= $2
        AND timestamp_utc <= $3
    `
	rows, err := s.DB.Query(query_str, asset.Symbol, fromTimestamp, toTimestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

    var price app.Price
	for rows.Next() {
		err := rows.Scan(
			&price.Id,
			&price.AssetSymbol,
			&price.TimestampUtc,
			&price.ValueUsd,
		)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

    missingTimestamps := app.GetMissingTimestamps(
        prices,
        fromTimestamp,
        toTimestamp,
        asset.Type != "crypto",
    )

    pricesToAppend, err := s.PA.GetDailyPrices(asset, missingTimestamps)
    if err != nil {
        return nil, err
    }

    prices = append(prices, pricesToAppend...)

    app.SortPrices(prices)

    err = s.PushPrices(prices)
    if err != nil {
        log.Println("Error during prices insert.", err)
    }

	return prices, nil
}

func (s *Service) PushPrices(prices []app.Price) error {
    insertsql := `
    INSERT INTO price_daily (asset_symbol, timestamp_utc, value_usd)
    VALUES ($1, $2, $3);
    `
    updatesql := `
    UPDATE asset
    SET
        value_usd = $3,
        last_synched = $2
    WHERE symbol = $1;
    `
    for _, price := range prices {
        _, err := s.DB.Exec(
            insertsql,
            price.AssetSymbol,
            price.TimestampUtc,
            price.ValueUsd,
        )
        if err != nil {
            return err
        }

        now := time.Now().UTC()
        if now.Sub(price.TimestampUtc).Abs().Hours() < 24 {
            _, err := s.DB.Exec(
                updatesql,
                price.AssetSymbol,
                now,
                price.ValueUsd,
            )
            if err != nil {
                return err
            }
        }
    }
    return nil
}
