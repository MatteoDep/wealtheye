package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Asset struct {
    Name string
    Symbol string
}

func get_assets() ([]Asset, error) {
    query_str := `
        select *
        from asset
    `
    rows, err := db.Query(query_str)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var assets []Asset
    for rows.Next() {
        var asset Asset
        if err := rows.Scan(
            & asset.Name,
            & asset.Symbol); err != nil {
            return assets, err
        }
        assets = append(assets, asset)
    }

    if err := rows.Err(); err != nil {
        return assets, err
    }

    return assets, nil
}

func get_asset(name string) (Asset, error) {
    query_str := `
        select *
        from asset
        where name = $1
    `
    row := db.QueryRow(query_str, name)

    var asset Asset
    if err := row.Scan(
        & asset.Name,
        & asset.Symbol); err != nil {
        return asset, err
    }

    if err := row.Err(); err != nil {
        return asset, err
    }

    return asset, nil
}

func plot_handler(c *fiber.Ctx) error {
    var asset Asset
    chosen_asset_name := c.Query("asset")
    asset, err := get_asset(chosen_asset_name)
    if err != nil {
        return err
    }
    return c.Render("balance-plot", fiber.Map{
        "AssetLabel": asset.Symbol,
    })
}

func root_handler(c *fiber.Ctx) error {
	app_name := "WealthEye"
    assets, err := get_assets()
    if err != nil {
        return err
    }

	return c.Render("index", fiber.Map{
		"Title": app_name,
        "Assets": assets,
	}, "layouts/main")
}

func main() {
    const dbfile string = "db/dashboard.db"
    var err error
    db, err = sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

    app.Static("/static/", "./static")

	app.Get("/", root_handler)
    app.Get("/plot:asset?", plot_handler)

	log.Fatal(app.Listen(":4242"))
}
