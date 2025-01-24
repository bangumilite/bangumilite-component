package utils

import (
	"fmt"
	"time"
)

func GetYear() int {
	return time.Now().Year()
}

func GetMonth() int {
	return int(time.Now().Month())
}

// Season -
type Season struct {
	ID       string `json:"id"`        // 202501
	SeasonCN string `json:"season_cn"` // 春
	Month    string `json:"month"`     // 01
	Year     string `json:"year"`      // 2025
}

func (s Season) ToString() string {
	return fmt.Sprintf("ID:%s, SeasonCN:%s, Month:%s, Year:%s", s.ID, s.SeasonCN, s.Month, s.Year)
}

const (
	YYYYMMDDDateFormatter = "2006-01-02"

	WinterSeasonCN = "冬"
	SpringSeasonCN = "春"
	SummerSeasonCN = "夏"
	AutumnSeasonCN = "秋"

	WinterMonth = "01"
	SpringMonth = "04"
	SummerMonth = "07"
	AutumnMonth = "10"
)

func winter(year string) Season {
	return Season{
		ID:       year + WinterMonth,
		SeasonCN: WinterSeasonCN,
		Month:    WinterMonth,
		Year:     year,
	}
}

func spring(year string) Season {
	return Season{
		ID:       year + SpringMonth,
		SeasonCN: SpringSeasonCN,
		Month:    SpringMonth,
		Year:     year,
	}
}

func summer(year string) Season {
	return Season{
		ID:       year + SummerMonth,
		SeasonCN: SummerSeasonCN,
		Month:    SummerMonth,
		Year:     year,
	}
}

func autumn(year string) Season {
	return Season{
		ID:       year + AutumnMonth,
		SeasonCN: AutumnSeasonCN,
		Month:    AutumnMonth,
		Year:     year,
	}
}

func GetSeason() Season {
	now := time.Now()
	year := fmt.Sprintf("%d", now.Year())
	month := int(now.Month())

	switch {
	case month >= 1 && month <= 3:
		return winter(year)
	case month >= 4 && month <= 6:
		return spring(year)
	case month >= 7 && month <= 9:
		return summer(year)
	case month >= 10 && month <= 12:
		return autumn(year)
	}

	return Season{}
}

func (s Season) Next() Season {
	now := time.Now()
	year := now.Year()

	var yyyy string
	if s.Month == "12" {
		year += 1
		yyyy = fmt.Sprintf("%d", year)
	} else {
		yyyy = s.Year
	}

	switch s.SeasonCN {
	case WinterSeasonCN:
		return spring(yyyy)
	case SpringSeasonCN:
		return summer(yyyy)
	case SummerSeasonCN:
		return autumn(yyyy)
	case AutumnSeasonCN:
		return winter(yyyy)
	}

	return Season{}
}
