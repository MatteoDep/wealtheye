package app

import (
	"fmt"
	"sort"
	"time"
)

func SortTimestamp(timestamps []time.Time) {
    sort.Slice(timestamps, func(i, j int) bool {
        return timestamps[i].Before(timestamps[j])
    })
}

func SortPrices(prices []Price) {
    sort.Slice(prices, func(i, j int) bool {
        return prices[i].TimestampUtc.Before(prices[j].TimestampUtc)
    })
}

func GetPriceMostIdx(prices []Price, rel func(int, int) bool) (int, error) {
    if len(prices) == 0 {
        return 0, fmt.Errorf("Cannot get any value from empty price slice.")
    }
    most := 0
    for i := range prices {
        if rel(i, most) {
            most = i
        }
    }
    return most, nil
}

func GetOldestPriceIdx(prices []Price) (int, error) {
    return GetPriceMostIdx(prices, func(i, j int) bool {
        return prices[i].TimestampUtc.Before(prices[j].TimestampUtc)
    })
}

func GetNewestPriceIdx(prices []Price) (int, error) {
    return GetPriceMostIdx(prices, func(i, j int) bool {
        return prices[i].TimestampUtc.After(prices[j].TimestampUtc)
    })
}

func GetMissingTimestamps(
    prices []Price,
    fromTimestamp time.Time,
    toTimestamp time.Time,
) []time.Time {
    missingTimestamps := []time.Time{}
    fromTimestamp = fromTimestamp.Round(24 * time.Hour)
    toTimestamp = toTimestamp.Round(24 * time.Hour)

    // todo: use lookup table
    for ts := fromTimestamp; ts.Before(toTimestamp); ts = ts.AddDate(0, 0, 1) {
        if _, err := GetPriceAtTimestampIdx(prices, ts); err != nil {
            missingTimestamps = append(missingTimestamps, ts)
        }
    }
    return missingTimestamps
}

func GetPriceAtTimestampIdx(
    prices []Price,
    timestampUtc time.Time,
) (int, error) {
    for i, price := range prices {
        if price.TimestampUtc.Equal(timestampUtc) {
            return i, nil
        }
    }

    return 0, fmt.Errorf("Could not find time stamp %v in prices %v", timestampUtc, prices)
}
