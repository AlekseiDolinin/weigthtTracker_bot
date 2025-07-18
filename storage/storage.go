package storage

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"weightTrack_bot/backup"
	"weightTrack_bot/models"
	"weightTrack_bot/parse"
)

const fileName = "data/dataBase.txt"

// проверка существования файла
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		msg := fmt.Sprintf("Файл %s не существует: %v", filename, err)
		backup.WriteLog(msg)
	}
	return !os.IsNotExist(err)
}

// добавляет в файл f запись r
func AddRecordToDB(r models.Record) (err error) {
	var file *os.File

	//проверяем существует ли файл
	if fileExists(fileName) { //если существует - открываем
		file, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			msg := fmt.Sprintf("Ошибка открытия файла: %v", err)
			backup.WriteLog(msg)
			return fmt.Errorf("ошибка открытия файла: %v", err)
		}
	} else { //если не существует - создаем
		file, err = os.Create(fileName)
		if err != nil {
			msg := fmt.Sprintf("Ошибка создания файла: %v", err)
			backup.WriteLog(msg)
			return fmt.Errorf("ошибка создания файла: %v", err)
		}
	}
	defer file.Close()
	//преобразует запись в строку
	record := fmt.Sprintf("%014d %06.2f %s %d\n", r.GetId(), r.GetWeight(), r.GetTime().Format(time.RFC3339), r.GetStatus())

	//записывает строку в файл
	_, err = file.WriteString(record)
	if err != nil {
		msg := fmt.Sprintf("Ошибка записи в файл: %v", err)
		backup.WriteLog(msg)
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}
	return nil
}

// возвращает слайс записей из файла f
func ReadRecords(chatID int) (records []models.Record, err error) {

	// Проверяем существует ли файл
	if !fileExists(fileName) {
		msg := fmt.Sprintf("Файл не существует: %s", fileName)
		backup.WriteLog(msg)
		return nil, fmt.Errorf("файл не существует: %s", fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		msg := fmt.Sprintf("Ошибка открытия файла: %v", err)
		backup.WriteLog(msg)
		return records, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if struc, err := parse.ParseRecord(scanner.Text()); struc.GetId() == chatID {
			if err != nil {
				continue
			}
			records = append(records, struc)
		}
	}

	if err := scanner.Err(); err != nil {
		msg := fmt.Sprintf("Ошибка чтения файла: %v", err)
		backup.WriteLog(msg)
		return records, fmt.Errorf("ошибка чтения файла: %v", err)
	} else if len(records) == 0 {
		return nil, fmt.Errorf("отсутствуют записи")
	}
	return
}

// возвращает слайс записей из файла f
func ReadAllRecords() (records []models.Record, err error) {

	// Проверяем существует ли файл
	if !fileExists(fileName) {
		msg := fmt.Sprintf("Файл не существует: %s", fileName)
		backup.WriteLog(msg)
		return nil, fmt.Errorf("файл не существует: %s", fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		msg := fmt.Sprintf("Ошибка открытия файла: %v", err)
		backup.WriteLog(msg)
		return records, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if struc, err := parse.ParseRecord(scanner.Text()); err == nil {
			records = append(records, struc)
		}
	}

	if err := scanner.Err(); err != nil {
		msg := fmt.Sprintf("Ошибка чтения файла: %v", err)
		backup.WriteLog(msg)
		return records, fmt.Errorf("ошибка чтения файла: %v", err)
	}
	return
}

// отправляет в чат предыдущую запись о весе
func ShowPreviousEntry(chatID int64) (result string, err error) {
	records, err := ReadRecords(int(chatID))
	if err != nil {
		msg := fmt.Sprintf("Ошибка чтения записей: %v", err)
		backup.WriteLog(msg)
		return result, err
	}
	//поиск последней неудаленной записи
	record, position := FindLastEntry(records, 0)

	if position != -1 {
		result := fmt.Sprintf("Ваш вес: %.2f кг\nЗапись от %d %s %dг. %02d:%02d ",
			record.GetWeight(),
			record.GetTime().Day(),
			parse.ParseMonth(record.GetTime().Month()),
			record.GetTime().Year(),
			record.GetTime().Hour(),
			record.GetTime().Minute(),
		)
		return result, nil
	}
	return
}

// удаляет(deleted = 0)/восстанавливает(deleted = 1) последнюю запись
func DeleteRestorePreviousEntry(chatID int64, delete int) error {

	allRecords, err := ReadAllRecords()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	record, positionOfAll := FindLastPosition(chatID, allRecords, delete)
	if positionOfAll == -1 {
		return fmt.Errorf("отсутствуют записи")
	}

	switch delete {
	case 0:
		record.SetStatus(1)
	case 1:
		record.SetStatus(0)
	}

	recordStr := fmt.Sprintf("%014d %06.2f %s %d\n", record.GetId(), record.GetWeight(), record.GetTime().Format(time.RFC3339), record.GetStatus())

	_, err = file.WriteAt([]byte(recordStr), int64(positionOfAll)*int64(len(recordStr))) // смещение: произведение длинны строки на количество строк
	if err != nil {
		msg := fmt.Sprintf("Ошибка перезаписи статуса записи в файле: %v", err)
		backup.WriteLog(msg)
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

// ищет позицию последней записи: deleted = 1 ищет последнюю удаленную, deleted = 0 - последнюю не удаленную
func FindLastPosition(chatID int64, records []models.Record, deleted int) (record models.Record, position int) {

	for i := len(records) - 1; i >= 0; i-- {
		if records[i].GetStatus() != deleted || records[i].GetId() != int(chatID) {
			continue
		} else {
			record = records[i]
			position = i
			return record, position
		}
	}
	return record, -1 //если не найдено
}

func DiffWeight(chatID int64) (weight float64, err error) {
	records, err := ReadRecords(int(chatID))
	if err != nil {
		msg := fmt.Sprintf("Ошибка чтения записей: %v", err)
		backup.WriteLog(msg)
		return 0.0, err
	}

	//поиск последней неудаленной записи
	record, position := FindLastEntry(records, 0)
	if position == -1 {
		return 0.0, fmt.Errorf("отсутствуют записи")
	}
	return record.GetWeight(), nil
}

// возвращает слайс средних данных показателей веса по дням за период
func FindPeriod(chatID int64, period int) (result []models.AvgRecordsPeriod, err error) {
	records, err := ReadRecords(int(chatID))
	if err != nil {
		msg := fmt.Sprintf("Ошибка чтения записей: %v", err)
		backup.WriteLog(msg)
		return nil, err
	} else if len(records) == 0 {
		return nil, fmt.Errorf("отсутствуют записи")
	}

	var dayAVG float64
	var countDays int
	lastEntry, _ := FindLastEntry(records, 0)
	currentDate := lastEntry.GetTime()

	for i := len(records) - 1; i >= 0 && (period > 0 || period <= -1); i-- {
		if records[i].GetStatus() != 0 {
			continue
		}
		if currentDate.Year() == records[i].GetTime().Year() &&
			currentDate.Month() == records[i].GetTime().Month() &&
			currentDate.Day() == records[i].GetTime().Day() {

			dayAVG += records[i].GetWeight()
			countDays++
		} else {
			dayAVG /= float64(countDays)
			result = append(result, models.NewAvgRecordsPeriod(dayAVG, currentDate))
			currentDate = records[i].GetTime()
			dayAVG = records[i].GetWeight()
			countDays = 1
			period--

		}
	}
	dayAVG /= float64(countDays)
	result = append(result, models.NewAvgRecordsPeriod(dayAVG, currentDate))
	return result, err
}

// формирует форматироваанную строку из слайса средних данных по дням за период []models.AvgRecordsPeriod
func ShowPeriod(result []models.AvgRecordsPeriod, period int) (s string) {
	s = "Среднее значение за день: \n"
	for i, rec := range result {
		if i >= period {
			continue
		}

		var diff float64
		if i < len(result)-1 {
			diff = rec.GetWeight() - result[i+1].GetWeight()
		}

		var diffStr string
		if diff != 0 {
			diffStr = fmt.Sprintf("%+06.2f", diff)
		}

		s += fmt.Sprintf("%02d %s %dг.: %.2fкг. %s \n",
			rec.GetTime().Day(),
			parse.ParseMonth(rec.GetTime().Month()),
			rec.GetTime().Year(),
			rec.GetWeight(),
			diffStr,
		)
	}
	return
}
