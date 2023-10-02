package app

import (
	"sort"
	"time"
)

func SortPrices(prices []Price) {
    sort.Slice(prices, func(i, j int) bool {
        return prices[i].TimestampUtc.Before(prices[j].TimestampUtc)
    })
}

func GetMissingTimestamps(
    prices []Price,
    fromTimestamp time.Time,
    toTimestamp time.Time,
    skipFestivities bool,
) []time.Time {
    missingTimestamps := []time.Time{}
    fromTimestamp = fromTimestamp.Round(24 * time.Hour)
    toTimestamp = toTimestamp.Round(24 * time.Hour)

    found := false
    for ts := fromTimestamp; ts.Before(toTimestamp); ts = ts.AddDate(0, 0, 1) {
        if skipFestivities &&
            ((ts.Weekday() == time.Saturday) || (ts.Weekday() == time.Sunday)) {
            continue
        }

        for _, price := range prices {
            if price.TimestampUtc.Equal(ts) {
                found = true
                break
            }
        }

        if !found {
            missingTimestamps = append(missingTimestamps, ts)
        }
    }
    return missingTimestamps
}
