package service

import (
	"time"
)


func (s *Service) Convert(ammount float64, fromAssetSymbol string, toAssetSymbol string) (float64, error) {
    fromAsset, err := s.GetAsset(fromAssetSymbol)
    if err != nil {
        return 0, err
    }
    toAsset, err := s.GetAsset(toAssetSymbol)
    if err != nil {
        return 0, err
    }

    fromAssetPrice, err := s.GetPrice(fromAsset, time.Now().UTC())
    if err != nil {
        return 0, err
    }
    toAssetPrice, err := s.GetPrice(toAsset, time.Now().UTC())
    if err != nil {
        return 0, err
    }

    convertedAmmount := ammount * toAssetPrice.ValueUsd / fromAssetPrice.ValueUsd

    return convertedAmmount, nil
}
