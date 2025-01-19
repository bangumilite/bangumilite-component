package utils

import (
	"fmt"
	"time"
)

// Season -
type Season struct {
	ID     string `json:"id"`
	Season string `json:"season"`
	Month  string `json:"month"`
	Year   string `json:"year"`
}

const (
	Spring = "春"
	Summer = "夏"
	Autumn = "秋"
	Winter = "冬"
)

// GetSeason returns the current season information based on the current date
func GetSeason() Season {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	var season string
	var seasonStartMonth int

	switch {
	case month >= 1 && month <= 3:
		season = "冬"
		seasonStartMonth = 1
	case month >= 4 && month <= 6:
		season = "春"
		seasonStartMonth = 4
	case month >= 7 && month <= 9:
		season = "夏"
		seasonStartMonth = 7
	case month >= 10 && month <= 12:
		season = "秋"
		seasonStartMonth = 10
	}

	seasonMonth := fmt.Sprintf("%02d", seasonStartMonth)
	id := fmt.Sprintf("%d%s", year, seasonMonth)

	return Season{
		Season: season,
		Month:  seasonMonth,
		Year:   fmt.Sprintf("%d", year),
		ID:     id,
	}
}
