package storage

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"weightTrack_bot/models"
	"weightTrack_bot/parse"
)

const FileName = "dataBase.txt"

// проверка существования файла
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// добавляет в файл f запись r
func AddRecordToDB(r models.Record) (err error) {
	var file *os.File

	//проверяем существует ли файл
	if fileExists(FileName) { //если существует - открываем
		file, err = os.OpenFile(FileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("ошибка открытия файла: %v", err)
		}
	} else { //если не существует - создаем
		file, err = os.Create(FileName)
		if err != nil {
			return fmt.Errorf("ошибка создания файла: %v", err)
		}
	}
	defer file.Close()
	//преобразует запись в строку
	//record := string(rune(r.getId())) + ", " /*+ r.getNickname()*/ + string(rune(r.getWeight())) + ", " + r.getTime().String()
	record := fmt.Sprintf("%d %.2f %s %d\n", r.GetId(), r.GetWeight(), r.GetTime().Format(time.RFC3339), r.GetStatus())

	//записывает строку в файл
	_, err = file.WriteString(record)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}
	return nil
}

// возвращает слайс записей из файла f
func ReadRecords(chatID int) (records []models.Record, err error) {

	// Проверяем существует ли файл
	if !fileExists(FileName) {
		return nil, fmt.Errorf("файл не существует: %s", FileName)
	}

	file, err := os.Open(FileName)
	if err != nil {
		return records, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if struc, err := parse.ParseRecord(scanner.Text()); struc.GetId() == chatID && err == nil {
			records = append(records, struc)
		}
	}

	if err := scanner.Err(); err != nil {
		return records, fmt.Errorf("ошибка чтения файла: %v", err)
	}
	return
}

// отправляет в чат предыдущую запись о весе
func ShowPreviousEntry(chatID int64) (result string, err error) {
	records, err := ReadRecords(int(chatID))
	if err != nil {
		result := "ошибка при чтении данных"
		return result, err
	}
	if records != nil {
		record := records[len(records)-1]
		result := fmt.Sprintf("Предыдущая запись создана %d %s %d в %02d:%02d \nВаш вес: %.2f кг",
			record.GetTime().Day(),
			record.GetTime().Month(),
			record.GetTime().Year(),
			record.GetTime().Hour(),
			record.GetTime().Minute(),
			record.GetWeight(),
		)
		return result, nil

	} else {
		result = "Вы еще не записывали свой вес"
		return result, err
	}
}

func DeletePreviousEntry(chatID int64) error {

	return nil
}

func FindLastEntry(chatID int64) {

}
