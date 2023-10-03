package handler

import (
    "time"

	"github.com/MatteoDep/wealtheye/app"
	"github.com/gofiber/fiber/v2"
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
	if err != nil {
		return err
	}

	return c.Render("index", fiber.Map{
		"Title":  app_name,
		"Assets": assets,
        "DefaultWalletName": "Wallet 0",
	}, "layouts/main")
}

func (h *Handler) PostWallet(c *fiber.Ctx) error {
    name := c.Query("name")
    value := c.Query("value")
    println(name, value)
    return nil
}
