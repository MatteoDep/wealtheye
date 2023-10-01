package sqlite_service

import (
	"database/sql"

	"github.com/MatteoDep/wealtheye/app"
	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	DB *sql.DB
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
	for rows.Next() {
		var asset app.Asset
		err := rows.Scan(
			&asset.Symbol,
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
        select *
        from asset
        where symbol = $1
    `
	row := s.DB.QueryRow(query_str, symbol)

	var asset app.Asset
	err := row.Scan(
		&asset.Symbol,
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

func (s *Service) GetPrices(assetSymbol string) ([]app.Price, error) {
	query_str := `
        select *
        from price_daily
        where asset_symbol = $1
    `
	rows, err := s.DB.Query(query_str, assetSymbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []app.Price
	for rows.Next() {
		var price app.Price
		err := rows.Scan(
			&price.Id,
			&price.AssetSymbol,
			&price.TimestampUtc,
			&price.ValueUsd,
		)
		if err != nil {
			return prices, err
		}
		prices = append(prices, price)
	}

	if err := rows.Err(); err != nil {
		return prices, err
	}

	return prices, nil
}
