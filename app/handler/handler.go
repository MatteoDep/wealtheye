package handler

import (
	"fmt"
	"log"
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
    c.Set("HX-Push-Url", "/landing-page")
	return c.Render("index", fiber.Map{
        "PageGet": "/landing-page",
		"Title":  "WealthEye",
	}, "layouts/main")
}

func (h *Handler) ServeHoldingsPage(c *fiber.Ctx) error {
    if c.Get("HX-Request") != "true" {
        return c.Render("index", fiber.Map{
            "PageGet": "/landing-page",
            "Title": "WealthEye",
        }, "layouts/main")
    }

	assets, err := h.Svc.GetAssets()
	if err != nil {
        log.Println(err)
        return fmt.Errorf("On GetAssets: %s.", err.Error())
	}
    err = h.Svc.UpdateWalletsValue()
	if err != nil {
        log.Println(err)
        return fmt.Errorf("On UpdateWalletsValue: %s.", err.Error())
	}
    wallets, err := h.Svc.GetWallets()
	if err != nil {
        log.Println(err)
        return fmt.Errorf("On GetWallets: %s.", err.Error())
	}

	return c.Render("landing-page", fiber.Map{
		"Assets": assets,
        "AssetSelectionGet": "/plot",
        "AssetSelectionTarget": "#balance-plot",
        "Wallets": wallets,
	})
}

func (h *Handler) ServeWalletPage(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
	if err != nil {
        log.Println(err)
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
        err = fmt.Errorf("On UpdateWalletValue: %s.", err.Error())
        log.Println(err)
        return err
    }

    wallet, err := h.Svc.GetWallet(walletId)
	if err != nil {
        err = fmt.Errorf("On GetWallet: %s.", err.Error())
        log.Println(err)
        return err
	}
    transfers, err := h.Svc.GetWalletTransfers(walletId)
	if err != nil {
        err = fmt.Errorf("On GetWalletTransfers: %s.", err.Error())
        log.Println(err)
        return err
	}

    walletTransfersDTO := []app.WalletTransferDTO{}
    for _, transfer := range transfers {
        walletTransferDTO, err := h.Svc.TransferToWalletTransferDTO(&transfer, walletId)
        if err != nil {
            log.Println(err)
            return fmt.Errorf("On TransferToWalletTransferDTO: %s.", err.Error())
        }
        walletTransfersDTO = append(walletTransfersDTO, *walletTransferDTO)
    }

	return c.Render("wallet-page", fiber.Map{
        "Wallet": wallet,
        "WalletTransfers": walletTransfersDTO,
	})
}

func (h *Handler) ServeBalancePlot(c *fiber.Ctx) error {
    assetSymbol := c.Query("AssetSymbol")
	asset, err := h.Svc.GetAsset(assetSymbol)
	if err != nil {
		return err
	}
    today := time.Now().UTC().Round(24 * time.Hour)
    prices, err := h.Svc.GetPrices(
        asset.Symbol,
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

    assets, err := h.Svc.GetAssets()
    if err != nil {
        return err
    }

	return c.Render("wallet-transfer-create", fiber.Map{
        "WalletId": walletId,
        "Ammount": 0,
        "Assets": assets,
        "Types": []app.WalletTransferType{
            app.Deposit,
            app.Withdrawal,
        },
        "Wallets": otherWallets,
	})
}

func (h *Handler) ServeWalletTransferEditForm(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
    if err != nil {
        return err
    }
    transferId, err := strconv.Atoi(c.Query("TransferId", "-1"))
    if err != nil {
        return err
    }
    if transferId == -1 {
        return fmt.Errorf("No transferId was provided.")
    }

    transfer, err := h.Svc.GetTransfer(transferId)
    if err != nil {
        return err
    }
    walletTransfer, err := h.Svc.TransferToWalletTransferDTO(transfer, walletId)
    if err != nil {
        return err
    }

    wallets, err := h.Svc.GetWallets()
    if err != nil {
        return err
    }

    otherWallets := []app.Wallet{}
    var wallet_tmp app.Wallet
    for i, wallet := range wallets {
        if wallet.Id != walletId {
            otherWallets = append(otherWallets, wallet)
        }
        if wallet.Id == walletTransfer.OtherWalletId {
            wallet_tmp = wallets[0]
            wallets[0] = wallet
            wallets[i] = wallet_tmp
        }
    }

    assets, err := h.Svc.GetAssets()
    if err != nil {
        return err
    }
    var asset_tmp app.Asset
    for i, asset := range assets {
        if asset.Symbol == walletTransfer.AssetSymbol {
            asset_tmp = assets[0]
            assets[0] = asset
            assets[i] = asset_tmp
        }
    }

    var types []app.WalletTransferType = make([]app.WalletTransferType, 2)
    if walletTransfer.Type == app.Deposit {
        types[0] = app.Deposit
        types[1] = app.Withdrawal
    } else {
        types[0] = app.Withdrawal
        types[1] = app.Deposit
    }

	return c.Render("wallet-transfer-edit", fiber.Map{
        "TransferId": transferId,
        "WalletId": walletId,
        "Ammount": walletTransfer.Ammount,
        "Assets": assets,
        "Types": types,
        "Wallets": otherWallets,
	})
}

func (h *Handler) ServeTransferDeleteForm(c *fiber.Ctx) error {
    transferId, err := strconv.Atoi(c.Params("transferId"))
    if err != nil {
        return err
    }

    message := "Do you really want to delete transfer?"
    endPoint := fmt.Sprintf("/transfer/%d", transferId)

	return c.Render("confirmation-form", fiber.Map{
        "Message": message,
        "EndPoint": endPoint,
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

    c.Set("HX-Trigger-After-Swap", "walletReload")
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
    c.Set("HX-Trigger-After-Swap", "walletReload")
    return nil
}

func (h *Handler) ServeDeleteWallet(c *fiber.Ctx) error {
    // walletId, err := strconv.Atoi(c.Params("walletId"))
    // if err != nil {
    //     return err
    // }
    //
    // wallet := new(app.Wallet)
    // if err := c.BodyParser(wallet); err != nil {
    //     return err
    // }
    // wallet.Id = walletId
    // if err := h.Svc.UpdateWalletName(wallet.Id, wallet.Name); err != nil {
    //     return err
    // }
    // c.Set("HX-Trigger-After-Swap", "walletReload")
    return nil
}

func (h *Handler) ServePostWalletTransfer(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
    if err != nil {
        return err
    }
    walletTrasferDTO := new(app.WalletTransferDTO)
    err = c.BodyParser(walletTrasferDTO)
    if err != nil {
        return err
    }
    walletTrasferDTO.WalletId = walletId

    transfer, err := h.Svc.WalletTransferDTOToTransfer(walletTrasferDTO)
    if err := h.Svc.PostTransfer(transfer); err != nil {
        return err
    }

    c.Set("HX-Trigger-After-Swap", "walletReload")
    return nil
}

func (h *Handler) ServePutWalletTransfer(c *fiber.Ctx) error {
    walletId, err := strconv.Atoi(c.Params("walletId"))
    if err != nil {
        return err
    }
    transferId, err := strconv.Atoi(c.Query("TransferId", "-1"))
    if err != nil {
        return err
    }
    if transferId == -1 {
        return fmt.Errorf("No transferId was provided.")
    }

    walletTrasferDTO := new(app.WalletTransferDTO)
    err = c.BodyParser(walletTrasferDTO)
    if err != nil {
        return err
    }
    walletTrasferDTO.WalletId = walletId

    transfer, err := h.Svc.WalletTransferDTOToTransfer(walletTrasferDTO)
    transfer.Id = transferId
    if err := h.Svc.UpdateTransfer(transfer); err != nil {
        return err
    }

    c.Set("HX-Trigger-After-Swap", "walletReload")
    return nil
}

func (h *Handler) ServeDeleteTransfer(c *fiber.Ctx) error {
    transferId, err := strconv.Atoi(c.Params("transferId"))
    if err != nil {
        return err
    }

    err = h.Svc.DeleteTransfer(transferId)
    if err != nil {
        return err
    }

    c.Set("HX-Trigger-After-Swap", "walletReload")
    return nil
}

func (h *Handler) GetExternalWalletName(c *fiber.Ctx) error {
    walletTransferType := app.WalletTransferType(c.Query("Type"))
    externalWalletName := h.Svc.GetExternalWalletName(walletTransferType)
	return c.Render("external-wallet-option", fiber.Map{
        "Name": externalWalletName,
	})
}
