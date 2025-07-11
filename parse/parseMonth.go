package parse

import (
	"time"
)

var m = map[string]string{
	"January":   "января",
	"February":  "февраля",
	"March":     "марта",
	"April":     "апреля",
	"May":       "мая",
	"June":      "июня",
	"July":      "июля",
	"August":    "августа",
	"September": "сентября",
	"October":   "октября",
	"November":  "ноября",
	"December":  "декабря",
}

func ParseMonth(month time.Month) string {
	return m[month.String()]
}
