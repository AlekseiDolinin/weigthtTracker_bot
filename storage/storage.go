package storage

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Record struct {
	id int
	//nickname string
	weight float64
	t      time.Time
}

func (r Record) getId() int {
	return r.id
}

//func (r Record) getNickname() string {
//	return r.nickname
//}

func (r Record) getWeight() int {
	return int(r.weight)
}

func (r Record) getTime() time.Time {
	return r.t
}

// возвращает экземпляр записи
func NewRecord(id int /*nickname string,*/, weight float64) Record {
	return Record{id: id /*nickname: nickname,*/, weight: weight, t: time.Now()}
}

// проверка существования файла
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// добавляет в файл f запись r
func AddRecordToDB(r Record, filename string) (err error) {
	var file *os.File

	//проверяем существует ли файл
	if fileExists(filename) { //если существует - открываем
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("ошибка открытия файла: %v", err)
		}
	} else { //если не существует - создаем
		file, err = os.Create(filename)
		if err != nil {
			return fmt.Errorf("ошибка создания файла: %v", err)
		}
	}
	defer file.Close()
	//преобразует запись в строку
	//record := string(rune(r.getId())) + ", " /*+ r.getNickname()*/ + string(rune(r.getWeight())) + ", " + r.getTime().String()
	record := fmt.Sprintf("%d, %.2f, %s\n", r.id, r.weight, r.t.Format(time.RFC3339))

	//записывает строку в файл
	_, err = file.WriteString(record)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}
	return nil
}

// возвращает слайс записей из файла f
func ReadRecords(filename string) (records []string, err error) {

	// Проверяем существует ли файл
	if !fileExists(filename) {
		return nil, fmt.Errorf("файл не существует: %s", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return records, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		records = append(records, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return records, fmt.Errorf("ошибка чтения файла: %v", err)
	}
	return
}
