package backup

import (
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const logFile = "data/log.txt"

// startDailyBackup запускает фоновый таймер для отправки файла в 00:00
func StartDailyBackup(bot *tgbotapi.BotAPI, filePath string, chatID int64) {
	for {
		now := time.Now()
		nextMidnight := time.Date(
			now.Year(), now.Month(), now.Day()+1,
			0, 0, 0, 0, now.Location(),
		)
		timeUntilMidnight := nextMidnight.Sub(now)

		msg := fmt.Sprintf("Следующая отправка файла %s в %v", filePath, nextMidnight)
		WriteLog(msg)

		time.Sleep(timeUntilMidnight)

		if err := sendFile(bot, filePath, chatID); err != nil {
			msg := fmt.Sprintf("Ошибка отправки файла %s: %v", filePath, err)
			WriteLog(msg)
			time.Sleep(10 * time.Minute)
		} else {
			time.Sleep(10 * time.Minute)
		}
	}
}

func sendFile(bot *tgbotapi.BotAPI, filePath string, chatID int64) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		return err
	}

	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{
		Name:  fileInfo.Name(),
		Bytes: fileBytes,
	})
	if _, err := bot.Send(doc); err != nil {
		return err
	}
	msg := fmt.Sprintf("Файл %s успешно отправлен в %v", filePath, time.Now().Format("2006-01-02 15:04:05"))
	WriteLog(msg)
	return nil
}

// проверка существования файла
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// добавляет в файл f запись r
func WriteLog(s string) (err error) {
	var file *os.File

	//проверяем существует ли файл
	if fileExists(logFile) { //если существует - открываем
		file, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			msg := fmt.Sprintf("ошибка открытия файла: %v", err)
			WriteLog(msg)
			return
		}
	} else { //если не существует - создаем
		file, err = os.Create(logFile)
		if err != nil {
			msg := fmt.Sprintf("ошибка создания файла: %v", err)
			WriteLog(msg)
			return
		}
	}
	defer file.Close()

	//записывает строку в файл
	_, err = file.WriteString(time.Now().Format("2006-01-02 15:04:05") + " " + s + "\n")
	if err != nil {
		msg := fmt.Sprintf("ошибка записи лога: %v", err)
		WriteLog(msg)
		return
	}
	return nil
}
