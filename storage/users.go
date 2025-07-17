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
	record := fmt.Sprintf("%014d %03d %06.2f\n", r.GetId(), r.GetAge(), r.GetHeight())

	//записывает строку в файл
	i, err := file.WriteString(record)

	if err != nil || i != 26 {
		return fmt.Errorf("ошибка записи нового пользователя: %v", err)
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
	if err != nil {
		return err
	}
	if position == -1 {
		return fmt.Errorf("отсутствуют записи")
	}

	userStr := fmt.Sprintf("%014d %03d %06.2f\n", user.GetId(), age, height)

	i, err := file.WriteAt([]byte(userStr), int64(position)*int64(len(userStr))-int64(len(userStr))) // смещение: произведение длинны строки на количество строк
	if err != nil || i != 26 {
		return fmt.Errorf("ошибка обновления записи пользователя: %v", err)
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

// Рассчитывает индекс массы тела и дает оценку индекса массы тела в зависимсти от возраста
func FindBMI(u models.User, r models.Record) (bmi float64, assessment string) {
	bmi = r.GetWeight() / (u.GetHeight() / 100.0 * u.GetHeight() / 100.0)
	age := u.GetAge()
	switch {
	case age > 0 && age < 18:
		assessment = "Оценка ИМТ для лиц младше 18 лет отсутствует"
	case age >= 18 && age < 26:
		switch {
		case bmi < 18.5:
			assessment = "Недостаточность питания"
		case bmi >= 18.5 && bmi <= 19.4:
			assessment = "Вес в норме\nПониженное питание"
		case bmi >= 19.5 && bmi <= 22.9:
			assessment = "Вес в норме\nНормальное соотношение роста и массы тела"
		case bmi >= 23 && bmi <= 27.4:
			assessment = "Вес в норме\nПовышенное питание"
		case bmi >= 27.5 && bmi <= 29.9:
			assessment = "Ожирение I степени"
		case bmi >= 30 && bmi <= 34.9:
			assessment = "Ожирение II степени"
		case bmi >= 35 && bmi <= 39.9:
			assessment = "Ожирение III степени"
		case bmi >= 40:
			assessment = "Ожирение IV степени"
		}
	case age >= 26:
		switch {
		case bmi < 19:
			assessment = "Недостаточность питания"
		case bmi >= 19 && bmi <= 19.9:
			assessment = "Пониженное питание"
		case bmi >= 20 && bmi <= 25.9:
			assessment = "Нормальное соотношение роста и массы тела"
		case bmi >= 26 && bmi <= 27.9:
			assessment = "Повышенное питание"
		case bmi >= 28 && bmi <= 30.9:
			assessment = "Ожирение I степени"
		case bmi >= 31 && bmi <= 35.9:
			assessment = "Ожирение II степени"
		case bmi >= 36 && bmi <= 40.9:
			assessment = "Ожирение III степени"
		case bmi >= 41:
			assessment = "Ожирение IV степени"
		}
	default:
		assessment = "Возможно, вы неверно указали свой возраст"
	}
	return
}
