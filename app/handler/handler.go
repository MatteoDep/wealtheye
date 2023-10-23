package handler

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/MatteoDep/wealtheye/app"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
)

type Handler struct {
	Svc app.Service
}

func (h *Handler) ServeIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
        "PageGet": "/holdings-page",
		"Title":  "WealthEye",
	}, "layouts/main")
}

func (h *Handler) ServeHoldingsPage(c *fiber.Ctx) error {
    if c.Get("HX-Request") != "true" {
        return c.Render("index", fiber.Map{
            "PageGet": "/holdings-page",
            "Title": "WealthEye",
        }, "layouts/main")
    }

	assets, err := h.Svc.GetAssets()
	if err != nil {
        return fmt.Errorf("On GetAssets: %s.", err.Error())
	}
    err = h.Svc.UpdateWalletsValue()
	if err != nil {
        return fmt.Errorf("On UpdateWalletsValue: %s.", err.Error())
	}
    wallets, err := h.Svc.GetWallets()
	if err != nil {
        return fmt.Errorf("On GetWallets: %s.", err.Error())
	}

    c.Set("HX-Push-Url", "/holdings-page")
	return c.Render("holdings-page", fiber.Map{
		"Assets": assets,
        "Wallets": wallets,
	})
}

func (h *Handler) ServeBalancePlot(c *fiber.Ctx) error {
    assetSymbol := c.Query("symbol")
	asset, err := h.Svc.GetAsset(assetSymbol)
	if err != nil {
		return err
	}
    today := time.Now().UTC().Round(24 * time.Hour)
    prices, err := h.Svc.GetPrices(
        asset,
        today.AddDate(0, -1, 0),
        today,
    )
	if err != nil {
		return err
	}
	return c.Render("balance-plot", fiber.Map{
		"AssetSymbol": asset.Symbol,
        "Prices": prices,
	})
}

func (h *Handler) ServeWalletInfoCard(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
	if err != nil {
		return err
	}
    wallet, err := h.Svc.GetWallet(walletId)
	if err != nil {
		return err
	}

	return c.Render("wallet-info-card", wallet)
}

func (h *Handler) ServeWalletPage(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
	if err != nil {
		return err
	}

    if c.Get("HX-Request") != "true" {
        return c.Render("index", fiber.Map{
            "PageGet": fmt.Sprintf("/wallet-page/%d", walletId),
            "Title": "WealthEye",
        }, "layouts/main")
    }

    err = h.Svc.UpdateWalletValue(walletId)
    if err != nil {
        return err
    }

    wallet, err := h.Svc.GetWallet(walletId)
	if err != nil {
		return err
	}
    transfers, err := h.Svc.GetWalletTransfers(walletId)
	if err != nil {
		return err
	}

    walletTransfersDTO := []app.WalletTransferDTO{}
    for _, transfer := range transfers {
        walletTransferDTO, err := h.TransferToWalletTransferDTO(transfer, walletId)
        if err != nil {
            return err
        }
        walletTransfersDTO = append(walletTransfersDTO, walletTransferDTO)
    }

    c.Set("HX-Push-Url", fmt.Sprintf("/wallet-page/%d", walletId))
	return c.Render("wallet-page", fiber.Map{
        "Wallet": wallet,
        "WalletTransfers": walletTransfersDTO,
	})
}

func (h *Handler) ServeWalletCreateForm(c *fiber.Ctx) error {
    wallets, err := h.Svc.GetWallets()
    if err != nil {
        return err
    }

    numbers := []int{}
    re, err := regexp.Compile(`^Wallet (?P<num>[0-9]+)$`)
    if err != nil {
        return err
    }
    for _, wallet := range wallets {
        if re.MatchString(wallet.Name) {
            num, err := strconv.Atoi(re.ReplaceAllString(wallet.Name, "${num}"))
            if err != nil {
                return err
            }

            numbers = append(numbers, num)
        }
    }
    var nextnum int = 0
    for slices.Contains(numbers, nextnum) {
        nextnum++
    }

    defaultWallet := app.Wallet{
        Name: fmt.Sprintf("Wallet %d", nextnum),
        ValueUsd: 0,
    }

	return c.Render("wallet-form-create", fiber.Map{
        "Wallet": defaultWallet,
	})
}

func (h *Handler) ServeWalletEditForm(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
	if err != nil {
		return err
	}
    wallet, err := h.Svc.GetWallet(walletId)
	if err != nil {
		return err
	}

	return c.Render("wallet-form-edit", fiber.Map{
        "Wallet": wallet,
	})
}

func (h *Handler) ServeWalletTransferCreateForm(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
    if err != nil {
        return err
    }

    wallets, err := h.Svc.GetWallets()
    if err != nil {
        return err
    }

    otherWallets := []app.Wallet{}
    for _, wallet := range wallets {
        if wallet.Id != walletId {
            otherWallets = append(otherWallets, wallet)
        }
    }
	return c.Render("wallet-transfer-create", fiber.Map{
        "WalletId": walletId,
        "Ammount": 0,
        "AssetSymbol": "USD",
        "Types": []app.WalletTransferType{
            app.Deposit,
            app.Withdrawal,
        },
        "Wallets": otherWallets,
	})
}

func (h *Handler) ServePostWallet(c *fiber.Ctx) error {
    wallet := new(app.Wallet)
    if err := c.BodyParser(wallet); err != nil {
        return err
    }
    if err := h.Svc.PostWallet(wallet.Name); err != nil {
        return err
    }

    c.Set("HX-Trigger-After-Swap", "walletCreated")
    return nil
}

func (h *Handler) ServePutWallet(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
    if err != nil {
        return err
    }

    wallet := new(app.Wallet)
    if err := c.BodyParser(wallet); err != nil {
        return err
    }
    wallet.Id = walletId
    if err := h.Svc.UpdateWalletName(wallet.Id, wallet.Name); err != nil {
        return err
    }
    c.Set("HX-Trigger-After-Swap", "walletEdited")
    return nil
}

func (h *Handler) ServePostWalletTransfer(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
    if err != nil {
        return err
    }
    walletTrasferDTO := new(app.WalletTransferDTO)
    if err := c.BodyParser(walletTrasferDTO); err != nil {
        return err
    }

    transfer, err := h.WalletTransferDTOToTransfer(*walletTrasferDTO, walletId)
    if err := h.Svc.PostTransfer(transfer); err != nil {
        return err
    }

    c.Set("HX-Trigger-After-Swap", "walletTransferCreated")
    return nil
}

func (h *Handler) ServePutWalletTransfer(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
    if err != nil {
        return err
    }
    walletTrasferDTO := new(app.WalletTransferDTO)
    if err := c.BodyParser(walletTrasferDTO); err != nil {
        return err
    }

    transfer, err := h.WalletTransferDTOToTransfer(*walletTrasferDTO, walletId)
    if err := h.Svc.UpdateTransfer(transfer); err != nil {
        return err
    }

    c.Set("HX-Trigger-After-Swap", "walletTransferEdited")
    return nil
}

func (h *Handler) TransferToWalletTransferDTO(transfer app.Transfer, walletId int) (app.WalletTransferDTO, error) {
    walletTransferDTO := app.WalletTransferDTO{}
    walletTransferDTO.Timestamp = transfer.TimestampUtc.Local()
    walletTransferDTO.Ammount = transfer.Ammount
    walletTransferDTO.AssetSymbol = transfer.AssetSymbol

    var otherWalletId int
    if transfer.ToWalletId == walletId {
        walletTransferDTO.Type = app.Deposit
        otherWalletId = transfer.FromWalletId
    } else {
        walletTransferDTO.Type = app.Withdrawal
        otherWalletId = transfer.ToWalletId
    }

    otherWallet, err := h.Svc.GetWallet(otherWalletId)
    if err != nil {
        return walletTransferDTO, err
    }

    walletTransferDTO.OtherWalletId = otherWallet.Id
    walletTransferDTO.OtherWalletName = otherWallet.Name

    return walletTransferDTO, nil
}

func (h *Handler) WalletTransferDTOToTransfer(walletTransferDTO app.WalletTransferDTO, walletId int) (app.Transfer, error) {
    transfer := app.Transfer{}
    transfer.TimestampUtc = walletTransferDTO.Timestamp.UTC()
    transfer.Ammount = walletTransferDTO.Ammount
    // transfer.AssetSymbol = walletTransferDTO.AssetSymbol
    transfer.AssetSymbol = "USD"

    if walletTransferDTO.Type == app.Deposit {
        transfer.FromWalletId = walletTransferDTO.OtherWalletId
        transfer.ToWalletId = walletId
    } else {
        transfer.FromWalletId = walletId
        transfer.ToWalletId = walletTransferDTO.OtherWalletId
    }

    return transfer, nil
}

