package season

import (
	"fmt"
)

const (
	WinterName = "冬"
	SpringName = "春"
	SummerName = "夏"
	AutumnName = "秋"

	WinterMonth = "01"
	SpringMonth = "04"
	SummerMonth = "07"
	AutumnMonth = "10"
)

func winter(year int) Season {
	return Season{
		id:   fmt.Sprintf("%d%s", year, WinterMonth),
		name: WinterName,
		year: year,
	}
}

func spring(year int) Season {
	return Season{
		id:   fmt.Sprintf("%d%s", year, SpringMonth),
		name: SpringName,
		year: year,
	}
}

func summer(year int) Season {
	return Season{
		id:   fmt.Sprintf("%d%s", year, SummerMonth),
		name: SummerName,
		year: year,
	}
}

func autumn(year int) Season {
	return Season{
		id:   fmt.Sprintf("%d%s", year, AutumnMonth),
		name: AutumnName,
		year: year,
	}
}
