package main

import (
	"database/sql"
	"log"

	"github.com/MatteoDep/wealtheye/app"
	"github.com/MatteoDep/wealtheye/app/handler"
	"github.com/MatteoDep/wealtheye/app/priceapi"
	"github.com/MatteoDep/wealtheye/app/repository"
	"github.com/MatteoDep/wealtheye/app/service"
	"github.com/MatteoDep/wealtheye/app/sync"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)


func main() {
    cfg, err := app.GetConfig()
    if err != nil {
        log.Fatal(err)
    }

    const dbfile string = "db/dashboard.db"

    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    rep := repository.Repository{
        DB: db,
    }

    pa := priceapi.GetPriceApi(&cfg.PriceApi)

    sm := sync.SyncManager{
        Rep: &rep,
        PA: pa,
        Cfg: &cfg.Sync,
    }
    go sm.Start()

    svc := service.Service{
        Rep: &rep,
        PA: pa,
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
	app.Get("/landing-page", h.ServeHoldingsPage)
    app.Get("/plot:AssetSymbol?", h.ServeBalancePlot)
    app.Get("/wallet-page/:walletId", h.ServeWalletPage)
    app.Get("/wallet-info-card/:walletId", h.ServeWalletInfoCard)
    app.Get("/wallet-create-form", h.ServeWalletCreateForm)
    app.Get("/wallet-edit-form/:walletId", h.ServeWalletEditForm)
    app.Post("/wallet", h.ServePostWallet)
    app.Put("/wallet/:walletId", h.ServePutWallet)
    app.Get("/wallet-transfer-create-form/:walletId", h.ServeWalletTransferCreateForm)
    app.Post("/wallet-transfer/:walletId", h.ServePostWalletTransfer)
    app.Put("/wallet-transfer/:walletId", h.ServePutWalletTransfer)
    app.Get("/external-wallet-name:Type?", h.GetExternalWalletName)

	log.Fatal(app.Listen(":4242"))
}
