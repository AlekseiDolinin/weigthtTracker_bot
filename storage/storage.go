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
	record := fmt.Sprintf("%d %.2f %s\n", r.GetId(), r.GetWeight(), r.GetTime().Format(time.RFC3339))

	//записывает строку в файл
	_, err = file.WriteString(record)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}
	return nil
}

// возвращает слайс записей из файла f
func ReadRecords(chatID int) (records []string, err error) {

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
		if struc := parse.ParseRecord(scanner.Text()).GetId(); struc == chatID {
			records = append(records, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return records, fmt.Errorf("ошибка чтения файла: %v", err)
	}
	return
}
