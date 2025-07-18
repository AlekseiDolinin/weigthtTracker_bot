package parse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"weightTrack_bot/backup"
	"weightTrack_bot/models"
)

// переводит строку в структуру Record
func ParseRecord(record string) (models.Record, error) {

	recordSplit := strings.Split(record, " ")
	if len(recordSplit) < 4 {
		backup.WriteLog("Ошибка чтения записи: отсутствуют необходимые данные")
		return models.NewRecord(0, 0, time.Now(), 0), fmt.Errorf("ошибка чтения записи: отсутствуют необходимые данные")
	}

	id, err := strconv.Atoi(recordSplit[0])
	if err != nil {
		msg := fmt.Sprintf("Ошибка преобразования строки в целое число %v", err)
		backup.WriteLog(msg)
	}

	weight, err := strconv.ParseFloat(recordSplit[1], 64)
	if err != nil {
		msg := fmt.Sprintf("Ошибка преобразования строки в число с плавающей точкой %v", err)
		backup.WriteLog(msg)
	}

	RFC3339 := "2006-01-02T15:04:05Z07:00"
	date, err := time.Parse(RFC3339, recordSplit[2])
	if err != nil {
		msg := fmt.Sprintf("Ошибка преобразования строки в формат даты %v", err)
		backup.WriteLog(msg)
	}

	var deleted int
	if recordSplit[3] == "1" {
		deleted = 1
	}
	result := models.NewRecord(id, weight, date, deleted)
	return result, nil
}
