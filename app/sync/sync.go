package sync

import (
	"log"
	"time"

	"github.com/MatteoDep/wealtheye/app"
)

type SyncManager struct {
	Rep app.Repository
	PA  app.PriceApi
	Cfg *app.SyncConfig
}

func (sm *SyncManager) Start() error {
    for true {
        assets, err := sm.Rep.GetAssets()
        if err != nil {
            return err
        }

        for _, asset := range assets {
            if asset.Symbol == "USD" {
                continue
            }
            log.Println("Started synching", asset.Name)
            err := sm.syncAssetPrices(&asset)
            if err != nil {
                log.Println("Sync Error:", err)
            }
            log.Println("Done synching", asset.Name)
        }

        time.Sleep(sm.Cfg.WaitingTime)
    }
	return nil
}

func (sm *SyncManager) syncAssetPrices(asset *app.Asset) error {
	fromTimestampUtc := sm.Cfg.StartTimestamp
	toTimestampUtc := time.Now().UTC()
	if lastTimestamp, err := sm.Rep.GetLastPriceTimestamp(asset.Symbol); err == nil {
        fromTimestampUtc = lastTimestamp
	} else {
        log.Println("Error getting last price timestamp:", err)
    }

    log.Printf("Getting updated prices from %v to %v.\n", fromTimestampUtc, toTimestampUtc)
	newPrices, err := sm.PA.GetDailyPricesUsd(
		asset,
		fromTimestampUtc,
		toTimestampUtc,
	)
    log.Println("found new prices", newPrices)
	if err != nil {
		return err
	}

	missingPrices := []app.Price{}
	changedPrices := []app.Price{}
	for _, price := range newPrices {
		if price.TimestampUtc == fromTimestampUtc {
			changedPrices = append(changedPrices, price)
		} else {
			missingPrices = append(missingPrices, price)
        }
	}

	log.Println("Adding prices:", missingPrices)
	err = sm.Rep.PostPrices(missingPrices)
	if err != nil {
		return err
	}

	log.Println("Updating prices:", changedPrices)
	err = sm.Rep.UpdatePricesValue(changedPrices)
	if err != nil {
		return err
	}

	return nil
}
