package storage

import (
	"bufio"
	"fmt"
	"os"
	"weightTrack_bot/models"
	"weightTrack_bot/parse"
)

const fileNameUsers = "data/users.txt"

// добавляет в файл f запись r
func AddUserToDB(r models.User) (err error) {
	var file *os.File

	//проверяем существует ли файл
	if fileExists(fileNameUsers) { //если существует - открываем
		file, err = os.OpenFile(fileNameUsers, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("ошибка открытия файла: %v", err)
		}
	} else { //если не существует - создаем
		file, err = os.Create(fileNameUsers)
		if err != nil {
			return fmt.Errorf("ошибка создания файла: %v", err)
		}
	}
	defer file.Close()
	//преобразует запись в строку
	record := fmt.Sprintf("%d %03d %06.2f\n", r.GetId(), r.GetAge(), r.GetHeight())

	//записывает строку в файл
	_, err = file.WriteString(record)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}
	return nil
}

// возвращает слайс записей из файла f
func ReadUser(chatID int64) (user models.User, err error) {

	// Проверяем существует ли файл
	if !fileExists(fileNameUsers) {
		return user, fmt.Errorf("файл не существует: %s", fileNameUsers)
	}

	file, err := os.Open(fileNameUsers)
	if err != nil {
		return user, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if struc, err := parse.ParseUser(scanner.Text()); struc.GetId() == chatID {
			user = struc
			return user, err
		} else {
			return user, fmt.Errorf("пользователь не найден: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return user, fmt.Errorf("ошибка чтения файла: %v", err)
	}
	return user, err
}

// обновляет информацию о пользователе
func UpdateUser(chatID int64, user models.User, age int, height float64) error {
	file, err := os.OpenFile(fileNameUsers, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//поиск позиции пользователя
	user, position, err := FindUserPosition(chatID)
	if position == -1 {
		return fmt.Errorf("отсутствуют записи: %v", err)
	}

	userStr := fmt.Sprintf("%d %03d %06.2f\n", user.GetId(), age, height)

	_, err = file.WriteAt([]byte(userStr), int64(position)*int64(len(userStr))-int64(len(userStr))) // смещение: произведение длинны строки на количество строк
	if err != nil {
		panic(err)
	}
	return nil
}

// выполняет поиск позиции записи с данными о пользователе
func FindUserPosition(chatID int64) (user models.User, position int, err error) {
	// Проверяем существует ли файл
	if !fileExists(fileNameUsers) {
		return user, -1, fmt.Errorf("файл не существует: %s", fileNameUsers)
	}

	file, err := os.Open(fileNameUsers)
	if err != nil {
		return user, -1, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		position++
		if user, err := parse.ParseUser(scanner.Text()); user.GetId() == chatID {
			return user, position, err
		}
	}
	return user, -1, err //если не найдено
}
