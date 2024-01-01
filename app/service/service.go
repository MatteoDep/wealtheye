package service

import (
	"database/sql"
	"fmt"
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

func (s *Service) GetAsset(symbol string) (*app.Asset, error) {
    return s.Rep.GetAsset(symbol)
}

func (s *Service) GetPrice(
    symbol string,
    timestampUtc time.Time,
) (*app.Price, error) {
    if symbol == "USD" {
        return &app.Price{
            AssetSymbol: symbol,
            ValueUsd: 1,
            TimestampUtc: timestampUtc,
        }, nil
    }
    dayStart := timestampUtc.Truncate(24 * time.Hour)
    prices, err := s.GetPrices(symbol, dayStart, timestampUtc)
    if len(prices) < 1 {
        prices, err = s.GetPrices(symbol, dayStart.AddDate(0, 0, -1), dayStart)
    }
    if err != nil {
        return nil, err
    }
    if len(prices) < 1 {
        return nil, fmt.Errorf("No prices found for %s", symbol)
    }
    return &prices[0], nil
}

func (s *Service) GetPrices(
    symbol string,
    fromTimestampUtc time.Time,
    toTimestampUtc time.Time,
) ([]app.Price, error) {
    if symbol == "USD" {
        return []app.Price{
            {
                AssetSymbol: symbol,
                ValueUsd: 1,
                TimestampUtc: fromTimestampUtc,
            },
            {
                AssetSymbol: symbol,
                ValueUsd: 1,
                TimestampUtc: toTimestampUtc,
            },
        }, nil
    }
    prices, err := s.Rep.GetPrices(symbol, fromTimestampUtc, toTimestampUtc)
    if err != nil {
        return prices, err
    }

    app.SortPrices(prices)

	return prices, nil
}

func (s *Service) GetWallets() ([]app.Wallet, error) {
    return s.Rep.GetWallets()
}

func (s *Service) GetWallet(id int) (*app.Wallet, error) {
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
        ammountUsd, err := s.Convert(transfer.Ammount, transfer.AssetSymbol, "USD", transfer.TimestampUtc)
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

func (s *Service) PostTransfer(transfer *app.Transfer) (error) {
    return s.Rep.PostTransfer(transfer)
}

func (s *Service) UpdateTransfer(transfer *app.Transfer) (error) {
    return s.Rep.UpdateTransfer(transfer)
}

func (s *Service) TransferToWalletTransferDTO(transfer *app.Transfer, walletId int) (*app.WalletTransferDTO, error) {
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
            return nil, err
        }
        walletTransferDTO.OtherWalletName = otherWallet.Name
    }


    return &walletTransferDTO, nil
}

func (s *Service) WalletTransferDTOToTransfer(walletTransferDTO *app.WalletTransferDTO, walletId int) (*app.Transfer, error) {
    transfer := app.Transfer{}
    transfer.TimestampUtc = walletTransferDTO.Timestamp.UTC()
    transfer.Ammount = walletTransferDTO.Ammount
    transfer.AssetSymbol = walletTransferDTO.AssetSymbol

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

    return &transfer, nil
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
