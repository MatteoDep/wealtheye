package handler

import (
	"github.com/MatteoDep/wealtheye/app"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Svc app.Service
    PA app.PriceApi
}

func (h *Handler) ServeBalancePlot(c *fiber.Ctx) error {
    assetSymbol := c.Query("symbol")
	asset, err := h.Svc.GetAsset(assetSymbol)
	if err != nil {
		return err
	}
    prices, err := h.PA.GetDailyPrices(asset, 0)
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
	}, "layouts/main")
}
