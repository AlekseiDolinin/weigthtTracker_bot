package parse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"weightTrack_bot/models"
)

// переводит строку в структуру Record
func ParseRecord(record string) (models.Record, error) {

	recordSplit := strings.Split(record, " ")
	id, _ := strconv.Atoi(recordSplit[0])
	weight, _ := strconv.ParseFloat(recordSplit[1], 64)
	RFC3339 := "2006-01-02T15:04:05Z07:00"
	date, _ := time.Parse(RFC3339, recordSplit[2])

	var deleted int
	if recordSplit[3] == "1" {
		deleted = 1
	}
	result := models.NewRecord(id, weight, date, deleted)

	if deleted == 0 {
		return result, nil
	} else {
		return result, fmt.Errorf("запись удалена")
	}
}
