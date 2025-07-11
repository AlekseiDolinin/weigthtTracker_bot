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
		if struc, _ := parse.ParseRecord(scanner.Text()); struc.GetId() == chatID {
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
	//поиск последней неудаленной записи
	record, position := FindLastEntry(records, 0)

	if position != -1 {
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

// удаляет(deleted = 0)/восстанавливает(deleted = 1) последнюю запись
func DeleteRestorePreviousEntry(chatID int64, delete int) error {
	file, err := os.OpenFile(FileName, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	records, err := ReadRecords(int(chatID))
	if err != nil {
		return err
	}

	//поиск последней удаленной/неудаленной записи
	record, position := FindLastEntry(records, delete)
	if position == -1 {
		return fmt.Errorf("отсутствуют записи")
	}

	switch delete {
	case 0:
		record.SetStatus(1)
	case 1:
		record.SetStatus(0)
	}

	recordStr := fmt.Sprintf("%d %.2f %s %d\n", record.GetId(), record.GetWeight(), record.GetTime().Format(time.RFC3339), record.GetStatus())

	_, err = file.WriteAt([]byte(recordStr), int64(position)*int64(len(recordStr))) // смещение длинна строки на количество строк
	if err != nil {
		panic(err)
	}
	return nil
}

// ищет последнюю запись: deleted = 1 ищет последнюю удаленную, deleted = 0 - последнюю не удаленную
func FindLastEntry(records []models.Record, deleted int) (record models.Record, position int) {
	for i := len(records) - 1; i >= 0; i-- {
		if records[i].GetStatus() != deleted {
			continue
		} else {
			record = records[i]
			position = i
			return record, position
		}
	}
	return record, -1 //если не найдено
}
