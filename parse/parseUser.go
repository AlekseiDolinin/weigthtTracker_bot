package parse

import (
	"fmt"
	"strconv"
	"strings"
	"weightTrack_bot/backup"
	"weightTrack_bot/models"
)

// переводит строку в структуру Record
func ParseUser(record string) (result models.User, err error) {

	recordSplit := strings.Split(record, " ")
	if len(recordSplit) < 3 {
		backup.WriteLog("Ошибка чтения записи: отсутствуют необходимые данные")
		return models.NewUser(0, 0, 0.0), fmt.Errorf("ошибка чтения пользователя: отсутствуют необходимые данные")
	}
	//recordSplit := strings.Split(record, " ")
	id, err1 := strconv.ParseInt(recordSplit[0], 10, 64)
	age, err2 := strconv.ParseInt(recordSplit[1], 10, 64)
	height, err3 := strconv.ParseFloat(recordSplit[2], 64)

	result = models.NewUser(id, int(age), height)
	if err1 != nil && err2 != nil && err3 != nil {
		msg := fmt.Sprintf("Ошибка чтения записи о пользователе: err1:%v,\nerr2:%v,\nerr3:%v", err1, err2, err3)
		backup.WriteLog(msg)
		return result, fmt.Errorf("ошибка чтения записи о пользователе: err1:%v,\nerr2:%v,\nerr3:%v", err1, err2, err3)
	} else {
		return result, nil
	}
}
