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
        if err := rows.Scan(
            & asset.Name,
            & asset.Symbol); err != nil {
            return assets, err
        }
        assets = append(assets, asset)
    }

    if err := rows.Err(); err != nil {
        return assets, err
    }

    return assets, nil
}

func (s *Service) GetAsset(name string) (app.Asset, error) {
    query_str := `
        select *
        from asset
        where name = $1
    `
    row := s.DB.QueryRow(query_str, name)

    var asset app.Asset
    if err := row.Scan(
        & asset.Name,
        & asset.Symbol); err != nil {
        return asset, err
    }

    if err := row.Err(); err != nil {
        return asset, err
    }

    return asset, nil
}
