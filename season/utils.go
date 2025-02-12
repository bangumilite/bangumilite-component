package season

import (
	"fmt"
	"time"
)

type Season struct {
	id   string
	name string
	year int

	t time.Time // internally used to track the system local time
}

func New(t time.Time) Season {
	y := t.Year()
	var s Season
	switch t.Month() {
	case 1, 2, 3:
		s = winter(y)
	case 4, 5, 6:
		s = spring(y)
	case 7, 8, 9:
		s = summer(y)
	case 10, 11, 12:
		s = autumn(y)
	}
	s.t = t
	return s
}

func (s Season) Next() Season {
	switch s.t.Month() {
	case 1, 2, 3:
		return spring(s.year)
	case 4, 5, 6:
		return summer(s.year)
	case 7, 8, 9:
		return autumn(s.year)
	case 10, 11, 12:
		return winter(s.year + 1)
	}
	return Season{}
}

func (s Season) ID() string {
	return s.id
}

func (s Season) Name() string {
	return s.name
}

func (s Season) Year() int {
	return s.year
}

func (s Season) ToString() string {
	return fmt.Sprintf("id:%s,name:%s,year:%d", s.id, s.name, s.year)
}
