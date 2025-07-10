package parse

import (
	"strconv"
	"strings"
	"time"
	"weightTrack_bot/models"
)

func ParseRecord(record string) models.Record {
	recordSplit := strings.Split(record, " ")
	id, _ := strconv.Atoi(recordSplit[0])
	weight, _ := strconv.ParseFloat(recordSplit[1], 64)
	RFC3339 := "2006-01-02T15:04:05Z07:00"
	date, _ := time.Parse(RFC3339, recordSplit[2])
	return models.NewRecord(id, weight, date)
}
