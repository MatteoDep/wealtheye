package main

import (
	"database/sql"
	"log"

    "github.com/MatteoDep/wealtheye/app/handler"
    "github.com/MatteoDep/wealtheye/app/sqlite_service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
    _ "github.com/mattn/go-sqlite3"
)


func main() {
    const dbfile string = "db/dashboard.db"

    db, err := sql.Open("sqlite3", dbfile)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    svc := &sqlite_service.Service{
        DB: db,
    }

    h := handler.Handler{
        Svc: svc,
    }

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

    app.Static("/static/", "./static")

	app.Get("/", h.ServeIndex)
    app.Get("/plot:asset?", h.ServeBalancePlot)

	log.Fatal(app.Listen(":4242"))
}
