package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Asset struct {
    Value string
    Label string
}

var assets = []Asset {
    {Value: "eur", Label: "EUR"},
    {Value: "usd", Label: "USD"},
    {Value: "btc", Label: "BTC"},
}
var default_asset = assets[0]

func plot_handler(c *fiber.Ctx) error {
    asset := default_asset
    chosen_asset := c.Query("asset")
    if chosen_asset != "" {
        for i := range assets {
            if chosen_asset == assets[i].Value {
                asset = assets[i]
                break
            }
        }
    }
    log.Println("Chosen ", asset.Label)
    return c.Render("balance-plot", fiber.Map{
        "AssetLabel": asset.Label,
    })
}

func root_handler(c *fiber.Ctx) error {
	app_name := "WealthEye"

	return c.Render("index", fiber.Map{
		"Title": app_name,
        "Assets": assets,
        "AssetLabel": default_asset.Label,
	}, "layouts/main")
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

    app.Static("/static/", "./static")

	app.Get("/", root_handler)
    app.Get("/plot:asset?", plot_handler)

	log.Fatal(app.Listen(":3000"))
}
