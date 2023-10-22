package service

import (
	"log"
	"time"

	"github.com/MatteoDep/wealtheye/app"
	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
    Rep  app.Repository
    PA app.PriceApi
}

func (s *Service) GetAssets() ([]app.Asset, error) {
    return s.Rep.GetAssets()
}

func (s *Service) GetAsset(symbol string) (app.Asset, error) {
    return s.Rep.GetAsset(symbol)
}

func (s *Service) GetPrice(
    asset app.Asset,
    timestampUtc time.Time,
) (app.Price, error) {
    dayStart := timestampUtc.Truncate(24 * time.Hour)
    prices, err := s.GetPrices(asset, dayStart, timestampUtc)
    if err != nil || len(prices) < 1 {
        return app.Price{}, err
    }
    return prices[0], nil
}

func (s *Service) GetPrices(
    asset app.Asset,
    fromTimestampUtc time.Time,
    toTimestampUtc time.Time,
) ([]app.Price, error) {
    if asset.Symbol == "USD" {
        prices := []app.Price{}
        timestamps := app.GetMissingTimestamps(
            prices,
            fromTimestampUtc,
            toTimestampUtc,
        )
        for _, timestamp := range timestamps {
            prices = append(prices, app.Price{
                AssetSymbol: asset.Symbol,
                TimestampUtc: timestamp,
                ValueUsd: 1,
            })
        }

        return prices, nil
    }

    prices, err := s.Rep.GetPrices(asset, fromTimestampUtc, toTimestampUtc)
    if err != nil {
        return prices, err
    }

    missingTimestamps := app.GetMissingTimestamps(
        prices,
        fromTimestampUtc,
        toTimestampUtc,
    )

    pricesToAppend, err := s.PA.GetDailyPrices(asset, missingTimestamps)
    if err != nil {
        return prices, err
    }
    for _, timestamp := range missingTimestamps {
        var dayShift int
        if timestamp.Weekday() == time.Saturday{
            dayShift = -1
        }
        if timestamp.Weekday() == time.Sunday {
            dayShift = -2
        }
        price, err := s.GetPrice(asset, timestamp.AddDate(0, 0, dayShift))
        if err != nil {
            return prices, err
        }
        pricesToAppend = append(pricesToAppend, price)
    }


    prices = append(prices, pricesToAppend...)
    app.SortPrices(prices)

    err = s.Rep.PostPrices(pricesToAppend)
    if err != nil {
        log.Println("Error during prices insert.", err)
    }

	return prices, nil
}

func (s *Service) GetWallets() ([]app.Wallet, error) {
    return s.Rep.GetWallets()
}

func (s *Service) GetWallet(id int) (app.Wallet, error) {
    return s.Rep.GetWallet(id)
}

func (s *Service) GetTransfers() ([]app.Transfer, error) {
    return s.GetTransfers()
}

func (s *Service) GetWalletTransfers(walletId int) ([]app.Transfer, error) {
	return s.Rep.GetWalletTransfers(walletId)
}

func (s *Service) UpdateWalletsValue() (error) {
    wallets, err := s.Rep.GetWallets()
	if err != nil {
		return err
	}
    for _, wallet := range wallets {
        err = s.UpdateWalletValue(wallet.Id)
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *Service) UpdateWalletValue(id int) (error) {
    transfers, err := s.Rep.GetWalletTransfers(id)
	if err != nil {
		return err
	}

    var valueUsd float64 = 0
    for _, transfer := range transfers {
        ammountUsd, err := s.Convert(transfer.Ammount, transfer.AssetSymbol, "USD")
        if err != nil {
            return err
        }

        if transfer.FromWalletId == id {
            ammountUsd *= -1
        }
        valueUsd += ammountUsd
    }

    return s.Rep.UpdateWalletValue(id, valueUsd)
}

func (s *Service) UpdateWalletName(id int, name string) (error) {
    return s.Rep.UpdateWalletName(id, name)
}

func (s *Service) PostWallet(name string) (error) {
    return s.Rep.PostWallet(name)
}

func (s *Service) PostTransfer(transfer app.Transfer) (error) {
    return s.Rep.PostTransfer(transfer)
}

func (s *Service) UpdateTransfer(transfer app.Transfer) (error) {
    return s.Rep.UpdateTransfer(transfer)
}