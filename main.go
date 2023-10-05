package main

import (
	"database/sql"
	"log"

	"github.com/MatteoDep/wealtheye/app"
	"github.com/MatteoDep/wealtheye/app/alphavantageapi"
	"github.com/MatteoDep/wealtheye/app/handler"
	"github.com/MatteoDep/wealtheye/app/sqliteservice"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)


func main() {
    var cfg app.Config
    app.GetConfig(&cfg)

    const dbfile string = "db/dashboard.db"

    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    pa := alphavantageapi.PriceApi{
        Cfg: &cfg,
    }

    svc := sqliteservice.Service{
        DB: db,
        PA: &pa,
    }

    h := handler.Handler{
        Svc: &svc,
    }

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

    app.Static("/static/", "./static")

	app.Get("/", h.ServeIndex)
    app.Get("/plot:symbol?", h.ServeBalancePlot)
    app.Get("/new-wallet", h.ServeNewWalletForm)
    app.Post("/wallet", h.ServeSubmitWallet)

	log.Fatal(app.Listen(":4242"))
}
