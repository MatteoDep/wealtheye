package app

import (
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

func GetMissingTimestamps(
    prices []Price,
    fromTimestamp time.Time,
    toTimestamp time.Time,
) []time.Time {
    missingTimestamps := []time.Time{}
    fromTimestamp = fromTimestamp.Round(24 * time.Hour)
    toTimestamp = toTimestamp.Round(24 * time.Hour)

    // todo: use lookup table
    found := false
    for ts := fromTimestamp; ts.Before(toTimestamp); ts = ts.AddDate(0, 0, 1) {
        for _, price := range prices {
            if price.TimestampUtc.Equal(ts) {
                found = true
                break
            }
        }

        if found {
            found = false
        } else {
            missingTimestamps = append(missingTimestamps, ts)
        }
    }
    return missingTimestamps
}
