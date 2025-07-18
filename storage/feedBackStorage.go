package storage

import (
	"fmt"
	"os"
	"time"
	"weightTrack_bot/backup"
	"weightTrack_bot/models"
)

const fileFeedBack = "data/feedBack.txt"

// добавляет в файл f запись r
func AddFeedBack(f models.FeedBack) (err error) {
	var file *os.File

	//проверяем существует ли файл
	if fileExists(fileFeedBack) { //если существует - открываем
		file, err = os.OpenFile(fileFeedBack, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			msg := fmt.Sprintf("Ошибка открытия файла: %v", err)
			backup.WriteLog(msg)
			return fmt.Errorf("ошибка открытия файла: %v", err)
		}
	} else { //если не существует - создаем
		file, err = os.Create(fileFeedBack)
		if err != nil {
			msg := fmt.Sprintf("Ошибка создания файла: %v", err)
			backup.WriteLog(msg)
			return fmt.Errorf("ошибка создания файла: %v", err)
		}
	}
	defer file.Close()
	//преобразует запись в строку
	record := fmt.Sprintf("%s %014d %s\n", f.GetTime().Format(time.RFC3339), f.GetUseID(), f.GetMsg())

	//записывает строку в файл
	_, err = file.WriteString(record)
	if err != nil {
		msg := fmt.Sprintf("Ошибка записи в файл: %v", err)
		backup.WriteLog(msg)
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}
	return nil
}
