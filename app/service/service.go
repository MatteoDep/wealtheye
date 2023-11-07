package service

import (
	"database/sql"
	"log"
	"time"

	"github.com/MatteoDep/wealtheye/app"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slices"
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
    if len(prices) < 1 {
        prices, err = s.GetPrices(asset, dayStart.AddDate(0, 0, -1), dayStart)
    }
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

    if len(missingTimestamps) > 0 {
        app.SortTimestamp(missingTimestamps)
        newPrices, err := s.PA.GetDailyPricesUsd(
            asset,
            missingTimestamps[0],
            missingTimestamps[len(missingTimestamps)-1],
        )
        if err != nil {
            return prices, err
        }

        missingPrices := []app.Price{}
        for _, price := range newPrices {
            if slices.Contains(missingTimestamps, price.TimestampUtc) {
                prices = append(prices, price)
                missingPrices = append(missingPrices, price)
            }
        }

        err = s.Rep.PostPrices(missingPrices)
        if err != nil {
            log.Println("Error during prices insert.", err)
        }
    }

    app.SortPrices(prices)

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

        if transfer.FromWalletId.Valid && int(transfer.FromWalletId.Int64) == id {
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

func (s *Service) TransferToWalletTransferDTO(transfer app.Transfer, walletId int) (app.WalletTransferDTO, error) {
    walletTransferDTO := app.WalletTransferDTO{}
    walletTransferDTO.Timestamp = transfer.TimestampUtc.Local()
    walletTransferDTO.Ammount = transfer.Ammount
    walletTransferDTO.AssetSymbol = transfer.AssetSymbol

    var otherWalletId int
    if int(transfer.ToWalletId.Int64) == walletId {
        walletTransferDTO.Type = app.Deposit
        if transfer.FromWalletId.Valid {
            otherWalletId = int(transfer.FromWalletId.Int64)
        } else {
            otherWalletId = -1
        }
    } else {
        walletTransferDTO.Type = app.Withdrawal
        if transfer.ToWalletId.Valid {
            otherWalletId = int(transfer.ToWalletId.Int64)
        } else {
            otherWalletId = -1
        }
    }
    walletTransferDTO.OtherWalletId = otherWalletId

    if otherWalletId == -1 {
        walletTransferDTO.OtherWalletName = s.GetExternalWalletName(walletTransferDTO.Type)
    } else {
        otherWallet, err := s.GetWallet(otherWalletId)
        if err != nil {
            return walletTransferDTO, err
        }
        walletTransferDTO.OtherWalletName = otherWallet.Name
    }


    return walletTransferDTO, nil
}

func (s *Service) WalletTransferDTOToTransfer(walletTransferDTO app.WalletTransferDTO, walletId int) (app.Transfer, error) {
    transfer := app.Transfer{}
    transfer.TimestampUtc = walletTransferDTO.Timestamp.UTC()
    transfer.Ammount = walletTransferDTO.Ammount
    // transfer.AssetSymbol = walletTransferDTO.AssetSymbol
    transfer.AssetSymbol = "USD"

    var otherWalletId sql.NullInt64
    if walletTransferDTO.OtherWalletId == -1 {
        otherWalletId = sql.NullInt64{
            Valid: false,
        }
    } else {
        otherWalletId = sql.NullInt64{
            Int64: int64(walletTransferDTO.OtherWalletId),
            Valid: true,
        }
    }
    if walletTransferDTO.Type == app.Deposit {
        transfer.FromWalletId = otherWalletId
        transfer.ToWalletId = sql.NullInt64{
            Int64: int64(walletId),
            Valid: true,
        }
    } else {
        transfer.FromWalletId = sql.NullInt64{
            Int64: int64(walletId),
            Valid: true,
        }
        transfer.ToWalletId = otherWalletId
    }

    return transfer, nil
}

func (s *Service) GetExternalWalletName(walletTransferType app.WalletTransferType) string {
    var externalWalletName string
    if (walletTransferType == app.Deposit) {
        externalWalletName = "Income/Gift"
    } else {
        externalWalletName = "Expense"
    }
    return externalWalletName
}
