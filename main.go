package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Asset struct {
    Value string
    Label string
}

var assets = map[string]Asset {
    "btc": {Value: "btc", Label: "BTC"},
    "usd": {Value: "usd", Label: "USD"},
    "eur": {Value: "eur", Label: "EUR"},
}

func plot_handler(c *fiber.Ctx) error {
    asset := assets[c.Query("asset")].Label
    fmt.Println(asset)
    return c.Render("asset-preview", fiber.Map{
        "Asset": asset,
    })
}

func root_handler(c *fiber.Ctx) error {
	app_name := "WealthEye"

	return c.Render("index", fiber.Map{
		"Title": app_name,
        "Assets": assets,
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
