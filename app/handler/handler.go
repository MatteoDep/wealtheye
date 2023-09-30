package handler

import (
	"github.com/MatteoDep/wealtheye/app"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
    Svc app.Service
}

func (h *Handler) ServeBalancePlot(c *fiber.Ctx) error {
    asset, err := h.Svc.GetAsset(c.Query("asset"))
    if err != nil {
        return err
    }
    return c.Render("balance-plot", fiber.Map{
        "AssetLabel": asset.Symbol,
    })
}

func (h *Handler) ServeIndex(c *fiber.Ctx) error {
	app_name := "WealthEye"
    assets, err := h.Svc.GetAssets()
    if err != nil {
        return err
    }

	return c.Render("index", fiber.Map{
		"Title": app_name,
        "Assets": assets,
	}, "layouts/main")
}
