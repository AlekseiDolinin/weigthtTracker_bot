package main

import (
	"fmt"
	"os"

	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"weightTrack_bot/backup"
	"weightTrack_bot/engine"
	"weightTrack_bot/models"
	"weightTrack_bot/parse"
)

func main() {

	backup.WriteLog("Приложение запущено")
	defer backup.WriteLog("Приложение остановлено")

	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		msg := fmt.Sprintf("Ошибка загрузки файла .env:  %v", err)
		backup.WriteLog(msg)
	}

	// Получаем из переменной окружения токен для подключения к телеграм-боту
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		msg := fmt.Sprintf("Ошибка получения токена:  %v", err)
		backup.WriteLog(msg)
	}

	// Резервное сохранение данных
	var path, _ = parse.ParseInt(os.Getenv("TELEGRAM_BOT_PATH"))
	var filePath1 = os.Getenv("TELEGRAM_BOT_BACKUP_FILE1")
	var filePath2 = os.Getenv("TELEGRAM_BOT_BACKUP_FILE2")
	var filePath3 = os.Getenv("TELEGRAM_BOT_BACKUP_FILE3")
	var filePath4 = os.Getenv("TELEGRAM_BOT_BACKUP_FILE4")
	go backup.StartDailyBackup(bot, filePath1, path)
	go backup.StartDailyBackup(bot, filePath2, path)
	go backup.StartDailyBackup(bot, filePath3, path)
	go backup.StartDailyBackup(bot, filePath4, path)

	// Настраиваем получение обновлений, устаналиваем время ожидания новых сообщений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Получаем канал обновлений
	updates := bot.GetUpdatesChan(u)

	// Не завершаем main пока все горутины не завершатся
	var wg sync.WaitGroup
	// Создание списка горутин
	goroutines := make(map[int64]*models.Goroutine)

	// Получаем сообщения из канала updates в бесконечном цикле
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		// Если горутина для этого chatID ещё не создана, создаём её
		if goroutines[chatID] == nil {
			//  Увеличиваем счетчик горутин на 1
			wg.Add(1)
			goroutines[chatID] = engine.StartBotGoroutine(
				chatID,
				bot,
				&wg,
				func(update tgbotapi.Update, bot *tgbotapi.BotAPI, wg *sync.WaitGroup) {
					engine.Engine(update, bot, wg) // Oсновная логика
				})
		}
		goroutines[chatID].Input <- update
	}
	wg.Wait() // Ждём завершения всех горутин
}
