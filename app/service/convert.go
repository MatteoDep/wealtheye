package service

import (
	"fmt"
	"time"
)


func (s *Service) Convert(ammount float64, fromAssetSymbol string, toAssetSymbol string, timestamp time.Time) (float64, error) {
    fromAsset, err := s.GetAsset(fromAssetSymbol)
    if err != nil {
        return 0, err
    }
    toAsset, err := s.GetAsset(toAssetSymbol)
    if err != nil {
        return 0, err
    }

    fromAssetPrice, err := s.GetPrice(fromAsset.Symbol, timestamp)
    if err != nil {
        return 0, err
    }
    toAssetPrice, err := s.GetPrice(toAsset.Symbol, timestamp)
    if err != nil {
        return 0, err
    }

    if toAssetPrice.ValueUsd == 0 {
        return 0, fmt.Errorf("Zero division: %s price is zero so could not perform the conversion.", toAssetPrice)
    }
    convertedAmmount := ammount * fromAssetPrice.ValueUsd / toAssetPrice.ValueUsd

    return convertedAmmount, nil
}
