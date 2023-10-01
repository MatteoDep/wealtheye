package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PriceApi struct {
		Key string
	}
}

func GetConfig(cfg *Config) {
    if _, err := os.Stat(".env"); err == nil {
        cfgmap, err := godotenv.Read()
        if err != nil {
            log.Fatal("Error loading .env file")
        }
        log.Println("cfgmap", cfgmap)

        cfg.PriceApi.Key = cfgmap["PRICE_API_KEY"]
    }

    if os.Getenv("PRICE_API_KEY") != "" {
        cfg.PriceApi.Key = os.Getenv("PRICE_API_KEY")
    }
}
