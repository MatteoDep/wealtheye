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
	app_name := "WealthEye"
	return c.Render("index", fiber.Map{
		"Title":  app_name,
	}, "layouts/main")
}

func (h *Handler) ServeHoldingsPage(c *fiber.Ctx) error {
	assets, err := h.Svc.GetAssets()
    wallets, err := h.Svc.GetWallets()
	if err != nil {
		return err
	}

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

func (h *Handler) ServeWalletPage(c *fiber.Ctx) error {
    id := c.Params("id")
    wallet, err := h.Svc.GetWallet(id)
	if err != nil {
		return err
	}

	return c.Render("wallet-page", fiber.Map{
        "Wallet": wallet,
	})
}

func (h *Handler) ServeWalletInfoCard(c *fiber.Ctx) error {
    id := c.Params("id")
    wallet, err := h.Svc.GetWallet(id)
	if err != nil {
		return err
	}

	return c.Render("wallet-info-card", wallet)
}

func (h *Handler) ServeNewWalletForm(c *fiber.Ctx) error {
    wallets, err := h.Svc.GetWallets()
    if err != nil {
        return err
    }

    numbers := []uint64{}
    re, err := regexp.Compile(`^Wallet (?P<num>[0-9]+)$`)
    if err != nil {
        return err
    }
    for _, wallet := range wallets {
        if re.MatchString(wallet.Name) {
            num, err := strconv.ParseUint(re.ReplaceAllString(wallet.Name, "${num}"), 10, 64)
            if err != nil {
                return err
            }

            numbers = append(numbers, num)
        }
    }
    var nextnum uint64 = 0
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

func (h *Handler) ServeEditwWalletForm(c *fiber.Ctx) error {
    id := c.Params("id")
    wallet, err := h.Svc.GetWallet(id)
	if err != nil {
		return err
	}

	return c.Render("wallet-form-edit", fiber.Map{
        "Wallet": wallet,
	})
}

func (h *Handler) ServePostWallet(c *fiber.Ctx) error {
    wallet := new(app.Wallet)
    if err := c.BodyParser(wallet); err != nil {
        return err
    }
    if err := h.Svc.PostWallet(*wallet); err != nil {
        return err
    }

    c.Set("HX-Trigger-After-Swap", "walletCreated")
    return nil
}

func (h *Handler) ServePutWallet(c *fiber.Ctx) error {
    wallet := new(app.Wallet)
    if err := c.BodyParser(wallet); err != nil {
        return err
    }
    if err := h.Svc.PutWallet(*wallet); err != nil {
        return err
    }
    c.Set("HX-Trigger-After-Swap", "walletEdited")
    return nil
}
