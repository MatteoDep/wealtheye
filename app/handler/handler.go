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

func (h *Handler) ServeIndex(c *fiber.Ctx) error {
	app_name := "WealthEye"
	assets, err := h.Svc.GetAssets()
    wallets, err := h.Svc.GetWallets()
	if err != nil {
		return err
	}

	return c.Render("index", fiber.Map{
		"Title":  app_name,
		"Assets": assets,
        "wallets": wallets,
        "DefaultWalletName": "Wallet 0",
	}, "layouts/main")
}

func (h *Handler) ServeNewWalletForm(c *fiber.Ctx) error {
    wallets, err := h.Svc.GetWallets()
    if err != nil {
        return err
    }

    numbers := []uint64{}
    re, err := regexp.Compile(`^Wallet ([0-9]+)`)
    if err != nil {
        return err
    }
    for _, wallet := range wallets {
        if re.MatchString(wallet.Name) {
            num, err := strconv.ParseUint(re.SubexpNames()[1], 10, 64)
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
    DefaultWalletName := fmt.Sprintf("Wallet %d", nextnum)

	return c.Render("index", fiber.Map{
        "DefaultWalletName": DefaultWalletName,
	}, "layouts/main")
}

func (h *Handler) ServeSubmitWallet(c *fiber.Ctx) error {
    wallet := new(app.Wallet)
    if err := c.BodyParser(wallet); err != nil {
        return err
    }
    if err := h.Svc.PostWallet(*wallet); err != nil {
        return err
    }
    return nil
}
