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
	if len(recordSplit) < 4 {
		return models.NewRecord(0, 0, time.Now(), 0), fmt.Errorf("ошибка чтения записи: отсутствуют необходимые данные")
	}
	id, _ := strconv.Atoi(recordSplit[0])
	weight, _ := strconv.ParseFloat(recordSplit[1], 64)
	RFC3339 := "2006-01-02T15:04:05Z07:00"
	date, _ := time.Parse(RFC3339, recordSplit[2])

	var deleted int
	if recordSplit[3] == "1" {
		deleted = 1
	}
	result := models.NewRecord(id, weight, date, deleted)

	//зачем?
	/*
		if deleted == 0 {
			return result, nil
		} else {
			return result, fmt.Errorf("запись удалена")
		}*/
	return result, nil

}
