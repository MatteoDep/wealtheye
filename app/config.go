package app

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

type PriceApiConfig struct {
	Provider string
	Key      string
}

type SyncConfig struct {
	StartTimestamp time.Time
	WaitingTime    time.Duration
}

type Config struct {
	PriceApi PriceApiConfig
	Sync     SyncConfig
}

type ConfigKey string

const (
	PriceApiProvider   ConfigKey = "PRICE_API_PROVIDER"
	PriceApiKey        ConfigKey = "PRICE_API_KEY"
	SyncStartTimestamp ConfigKey = "SYNC_START_TIMESTAMP"
	SyncWaitingTime    ConfigKey = "SYNC_WAITING_TIME"
)

var configKeys = []ConfigKey{
    PriceApiProvider,
    PriceApiKey,
    SyncStartTimestamp,
    SyncWaitingTime,
}

var defaults = map[ConfigKey]string{
    PriceApiProvider: "alphavantage",
    SyncStartTimestamp: "2023-12-01",
    SyncWaitingTime: "5s",
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	envmap := map[ConfigKey]string{}
	readEnviron(envmap)
    fmt.Println(envmap)
	readDotenv(envmap)
    fmt.Println(envmap)
	for k, v := range envmap {
		if v == "" {
            if val, ok := defaults[k]; ok {
                v = val
            } else {
                return nil, fmt.Errorf("Could not get configuration for %s.", k)
            }
		}
		switch k {
		case PriceApiProvider:
			cfg.PriceApi.Provider = v
		case PriceApiKey:
			cfg.PriceApi.Key = v
		case SyncStartTimestamp:
			ts, err := time.Parse(time.DateOnly, v)
			if err != nil {
				return nil, err
			}
			cfg.Sync.StartTimestamp = ts
		case SyncWaitingTime:
			duration, err := time.ParseDuration(v)
			if err != nil {
				return nil, err
			}
			cfg.Sync.WaitingTime = duration
		}
	}
    fmt.Println(cfg)
	return &cfg, nil
}

func readEnviron(envMap map[ConfigKey]string) {
	for _, key := range configKeys {
		envMap[key] = os.Getenv(string(key))
	}
}

func readDotenv(envMap map[ConfigKey]string) {
	if _, err := os.Stat(".env"); err == nil {
		dotenvMap, err := godotenv.Read()
		if err == nil {
			for k := range dotenvMap {
				if slices.Contains(configKeys, ConfigKey(k)) {
					envMap[ConfigKey(k)] = dotenvMap[k]
				}
			}
		}
	}
}
