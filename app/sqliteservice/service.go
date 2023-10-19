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
	queryStr := `
        select *
        from asset
    `
	rows, err := s.DB.Query(queryStr)
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
	queryStr := `
        SELECT *
        FROM asset
        WHERE symbol = $1
    `
	row := s.DB.QueryRow(queryStr, symbol)

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

	queryStr := `
        SELECT *
        FROM price_daily
        WHERE asset_symbol = $1
        AND timestamp_utc >= $2
        AND timestamp_utc <= $3
    `
	rows, err := s.DB.Query(queryStr, asset.Symbol, fromTimestamp, toTimestamp)
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

    missingTimestamps := app.GetMissingTimestamps(
        prices,
        fromTimestamp,
        toTimestamp,
        asset.Type != "crypto",
    )

    pricesToAppend, err := s.PA.GetDailyPrices(asset, missingTimestamps)
    if err != nil {
        return prices, err
    }

    prices = append(prices, pricesToAppend...)
    app.SortPrices(prices)

    err = s.PostPrices(pricesToAppend)
    if err != nil {
        log.Println("Error during prices insert.", err)
    }

	return prices, nil
}

func (s *Service) PostPrices(prices []app.Price) error {
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
        _, err := s.DB.Exec(
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
            _, err := s.DB.Exec(
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

func (s *Service) GetWallets() ([]app.Wallet, error) {
	queryStr := `
        SELECT *
        FROM wallet
    `
	rows, err := s.DB.Query(queryStr)
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

func (s *Service) GetWallet(id int) (app.Wallet, error) {
	queryStr := `
        SELECT *
        FROM wallet
        WHERE id = $1
    `
	row := s.DB.QueryRow(queryStr, id)

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

func (s *Service) PostWallet(wallet app.Wallet) error {
    insertStr := `
    INSERT INTO wallet (name, value_usd)
    VALUES ($1, $2);
    `
    _, err := s.DB.Exec(
        insertStr,
        wallet.Name,
        wallet.ValueUsd,
    )
    if err != nil {
        return err
    }
    return nil
}

func (s *Service) PutWallet(wallet app.Wallet) error {
    updateStr := `
    UPDATE wallet
    SET
        name = $1,
        value_usd = $2;
    `
    _, err := s.DB.Exec(
        updateStr,
        wallet.Name,
        wallet.ValueUsd,
    )
    if err != nil {
        return err
    }
    return nil
}

func (s *Service) GetTransfers() ([]app.Transfer, error) {
	queryStr := `
        SELECT *
        FROM transfer;
    `
	rows, err := s.DB.Query(queryStr)
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

func (s *Service) GetWalletTransfers(walletId int) ([]app.Transfer, error) {
	queryStr := `
        SELECT *
        FROM transfer
        where from_wallet_id = $1
            or to_wallet_id = $1;
    `
	rows, err := s.DB.Query(queryStr, walletId)
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

func (s *Service) PostTransfer(transfer app.Transfer) error {
    insertStr := `
    INSERT INTO transfer (ammount, from_wallet_id, to_wallet_id, asset_symbol)
    VALUES ($1, $2, $3, $4);
    `
    _, err := s.DB.Exec(
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

func (s *Service) PutTransfer(transfer app.Transfer) error {
    updateStr := `
    UPDATE transfer
    SET
        ammount = $1,
        from_wallet_id = $2,
        to_wallet_id = $3,
        asset_symbol = $4;
    `
    _, err := s.DB.Exec(
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
