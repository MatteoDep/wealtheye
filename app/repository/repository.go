package repository

import (
	"database/sql"
	"time"

	"github.com/MatteoDep/wealtheye/app"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) GetAssets() ([]app.Asset, error) {
	queryStr := `
        select *
        from asset
    `
	rows, err := r.DB.Query(queryStr)
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

func (r *Repository) GetAsset(symbol string) (app.Asset, error) {
	queryStr := `
        SELECT *
        FROM asset
        WHERE symbol = $1
    `
	row := r.DB.QueryRow(queryStr, symbol)

	var asset app.Asset
	err := row.Scan(
		&asset.Symbol,
        &asset.Name,
        &asset.Type,
	)
	if err != nil {
		return asset, err
	}

	if err := row.Err(); err != nil {
		return asset, err
	}

	return asset, nil
}

func (r *Repository) UpdateAssetValue(symbol string, valueUsd float64, lastSynched time.Time) (error) {
    updateStr := `
        UPDATE asset
        SET
            value_usd = $1,
            last_synched = $2
        WHERE
            symbol = $3;
    `
    _, err := r.DB.Exec(
        updateStr,
        valueUsd,
        lastSynched,
        symbol,
    )
    if err != nil {
        return err
    }
    return nil
}

func (r *Repository) GetPrice(
    asset app.Asset,
    timestampUtc time.Time,
) (app.Price, error) {
    dayStart := timestampUtc.Truncate(24 * time.Hour)
    dayEnd := dayStart.AddDate(0, 0, 1)
    prices, err := r.GetPrices(asset, dayStart, dayEnd)
    if err != nil || len(prices) < 1 {
        return app.Price{}, err
    }
    return prices[0], nil
}

func (r *Repository) GetPrices(
    asset app.Asset,
    fromTimestampUtc time.Time,
    toTimestampUtc time.Time,
) ([]app.Price, error) {
    prices := []app.Price{}
    if asset.Symbol == "USD" || toTimestampUtc.Sub(fromTimestampUtc) < 24 * time.Hour {
        return prices, nil
    }

	queryStr := `
        SELECT *
        FROM price_daily
        WHERE asset_symbol = $1
        AND timestamp_utc >= $2
        AND timestamp_utc <= $3
        ORDER BY timestamp_utc;
    `
	rows, err := r.DB.Query(queryStr, asset.Symbol, fromTimestampUtc, toTimestampUtc)
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
			return prices, err
		}
		prices = append(prices, price)
	}

	if err := rows.Err(); err != nil {
		return prices, err
	}

	return prices, nil
}

func (r *Repository) PostPrices(prices []app.Price) error {
    insertStr := `
    INSERT INTO price_daily (asset_symbol, timestamp_utc, value_usd)
    VALUES ($1, $2, $3);
    `
    updateStr := `
    UPDATE asset
    SET
        value_usd = $3,
        last_synched = $2
    WHERE symbol = $1;
    `
    for _, price := range prices {
        _, err := r.DB.Exec(
            insertStr,
            price.AssetSymbol,
            price.TimestampUtc,
            price.ValueUsd,
        )
        if err != nil {
            return err
        }

        now := time.Now().UTC()
        if now.Sub(price.TimestampUtc).Abs().Hours() < 24 {
            _, err := r.DB.Exec(
                updateStr,
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

func (r *Repository) GetWallets() ([]app.Wallet, error) {
	queryStr := `
        SELECT *
        FROM wallet
    `
	rows, err := r.DB.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []app.Wallet
    var wallet app.Wallet
	for rows.Next() {
		err := rows.Scan(
            &wallet.Id,
            &wallet.Name,
			&wallet.ValueUsd,
		)
		if err != nil {
			return wallets, err
		}
		wallets = append(wallets, wallet)
	}

	if err := rows.Err(); err != nil {
		return wallets, err
	}

	return wallets, nil
}

func (r *Repository) GetWallet(id int) (app.Wallet, error) {
	queryStr := `
        SELECT *
        FROM wallet
        WHERE id = $1
    `
	row := r.DB.QueryRow(queryStr, id)

	var wallet app.Wallet
	err := row.Scan(
        &wallet.Id,
        &wallet.Name,
		&wallet.ValueUsd,
	)
	if err != nil {
		return wallet, err
	}

	if err := row.Err(); err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (r *Repository) PostWallet(name string) error {
    insertStr := `
    INSERT INTO wallet (name)
    VALUES ($1);
    `
    _, err := r.DB.Exec(
        insertStr,
        name,
    )
    if err != nil {
        return err
    }
    return nil
}

func (r *Repository) UpdateWalletName(id int, name string) error {
    updateStr := `
    UPDATE wallet
    SET
        name = $1
    WHERE
        id = $2;
    `
    _, err := r.DB.Exec(
        updateStr,
        name,
        id,
    )
    if err != nil {
        return err
    }
    return nil
}

func (r *Repository) UpdateWalletValue(id int, valueUsd float64) error {
    updateStr := `
    UPDATE wallet
    SET
        value_usd = $1
    WHERE
        id = $2;
    `
    _, err := r.DB.Exec(
        updateStr,
        valueUsd,
        id,
    )
    if err != nil {
        return err
    }
    return nil
}

func (r *Repository) GetTransfers() ([]app.Transfer, error) {
	queryStr := `
        SELECT *
        FROM transfer;
    `
	rows, err := r.DB.Query(queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []app.Transfer
    var transfer app.Transfer
	for rows.Next() {
		err := rows.Scan(
            &transfer.Id,
            &transfer.TimestampUtc,
            &transfer.Ammount,
            &transfer.AssetSymbol,
            &transfer.FromWalletId,
            &transfer.ToWalletId,
		)
		if err != nil {
			return transfers, err
		}
		transfers = append(transfers, transfer)
	}

	if err := rows.Err(); err != nil {
		return transfers, err
	}

	return transfers, nil
}

func (r *Repository) GetWalletTransfers(walletId int) ([]app.Transfer, error) {
	queryStr := `
        SELECT *
        FROM transfer
        where from_wallet_id = $1
            or to_wallet_id = $1;
    `
	rows, err := r.DB.Query(queryStr, walletId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []app.Transfer
    var transfer app.Transfer
	for rows.Next() {
		err := rows.Scan(
            &transfer.Id,
            &transfer.TimestampUtc,
            &transfer.Ammount,
            &transfer.AssetSymbol,
            &transfer.FromWalletId,
            &transfer.ToWalletId,
		)
		if err != nil {
			return transfers, err
		}
		transfers = append(transfers, transfer)
	}

	if err := rows.Err(); err != nil {
		return transfers, err
	}

	return transfers, nil
}

func (r *Repository) PostTransfer(transfer app.Transfer) error {
    insertStr := `
    INSERT INTO transfer (ammount, from_wallet_id, to_wallet_id, asset_symbol)
    VALUES ($1, $2, $3, $4);
    `
    _, err := r.DB.Exec(
        insertStr,
        transfer.Ammount,
        transfer.FromWalletId,
        transfer.ToWalletId,
        transfer.AssetSymbol,
    )
    if err != nil {
        return err
    }
    return nil
}

func (r *Repository) UpdateTransfer(transfer app.Transfer) error {
    updateStr := `
    UPDATE transfer
    SET
        ammount = $1,
        from_wallet_id = $2,
        to_wallet_id = $3,
        asset_symbol = $4;
    `
    _, err := r.DB.Exec(
        updateStr,
        transfer.Ammount,
        transfer.FromWalletId,
        transfer.ToWalletId,
        transfer.AssetSymbol,
    )
    if err != nil {
        return err
    }
    return nil
}
