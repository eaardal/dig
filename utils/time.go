package utils

import "time"

func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}
