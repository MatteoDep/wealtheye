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
	assets, err := sm.Rep.GetAssets()
	if err != nil {
		return err
	}

    for true {
        for _, asset := range assets {
            log.Println(asset.Name)
            sm.syncAssetPrices(asset)
        }

        time.Sleep(sm.Cfg.WaitingTime)
    }
	return nil
}

func (sm *SyncManager) syncAssetPrices(asset app.Asset) error {
	fromTimestampUtc := sm.Cfg.StartTimestamp
	toTimestampUtc := time.Now().UTC()
    log.Println(fromTimestampUtc, toTimestampUtc)

	prices, err := sm.Rep.GetPrices(asset, fromTimestampUtc, toTimestampUtc)
	if err != nil {
		return err
	}

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
		existentPrice := app.GetPriceAtTimestamp(prices, price.TimestampUtc)
		if existentPrice == nil {
			missingPrices = append(missingPrices, price)
		} else if price.ValueUsd == existentPrice.ValueUsd {
			changedPrices = append(changedPrices, price)
		}
	}

	log.Println("Adding prices:", missingPrices)
	err = sm.Rep.PostPrices(missingPrices)
	if err != nil {
		return err
	}

	log.Println("Updating prices:", changedPrices)
	err = sm.Rep.PostPrices(changedPrices)
	if err != nil {
		return err
	}

	return nil
}
