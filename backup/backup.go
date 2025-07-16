package backup

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// startDailyBackup запускает фоновый таймер для отправки файла в 00:00
func StartDailyBackup(bot *tgbotapi.BotAPI, filePath string, chatID int64) {
	for {
		now := time.Now()
		nextMidnight := time.Date(
			now.Year(), now.Month(), now.Day()+1,
			0, 0, 0, 0, now.Location(),
		)
		timeUntilMidnight := nextMidnight.Sub(now)

		log.Printf("Следующая отправка файла в: %v", nextMidnight)
		time.Sleep(timeUntilMidnight)

		if err := sendFile(bot, filePath, chatID); err != nil {
			log.Println("Ошибка отправки файла:", err)
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

	log.Println("Файл успешно отправлен в", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
