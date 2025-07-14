package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"weightTrack_bot/messages"
	"weightTrack_bot/models"
	"weightTrack_bot/parse"
	"weightTrack_bot/plots"
	"weightTrack_bot/storage"
)

//const fileName = "dataBase.txt"

func main() {

	//получаем из переменной окружения токен для подключения к телеграм-боту
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	//настраиваем получение обновлений, устаналиваем время ожидания новых сообщений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//получаем канал обновлений
	updates := bot.GetUpdatesChan(u)

	//заглушка сохранения веса
	var weightInput float64 //заменить на сохранение в базе данных

	//олучаем сообщения из канала updates в бесконечном цикле
	for update := range updates {
		if update.Message == nil {
			continue
		}

		//извлекаем текст сообщенния и идентификатор чата для отправки ответа
		text := update.Message.Text
		chatID := update.Message.Chat.ID

		//парсим сообщение пользователя на наличие числа для записи веса
		if _, err := parse.ParseFloat(update.Message.Text); err == nil {
			weightInput, _ = parse.ParseFloat(update.Message.Text)
		}

		//выбираем ответ по запросу
		//вынести логику в отдельную функцию
		switch {
		case strings.EqualFold(text, "/start"):
			msg := tgbotapi.NewMessage(chatID, messages.WelcomeMsg)
			bot.Send(msg)
		case strings.EqualFold(text, "/show_week"):
			period, _ := storage.FindPeriod(chatID, 7)
			preMsg := storage.ShowPeriod(period, 7)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
		case strings.EqualFold(text, "/show_month"):
			period, _ := storage.FindPeriod(chatID, 31)
			preMsg := storage.ShowPeriod(period, 31)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
		case strings.EqualFold(text, "/show_progress"):
			period, _ := storage.FindPeriod(chatID, 31)
			// Создаем график в памяти
			imgBytes, err := plots.MakePlot(period)
			if err != nil {
				log.Panic(err)
			}
			// Создаем файл для отправки в Telegram
			file := tgbotapi.FileBytes{
				Name:  "plot.png",
				Bytes: imgBytes,
			}
			// Отправляем фото
			msg := tgbotapi.NewPhoto(chatID, file)
			msg.Caption = "График изменения веса"
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		case strings.EqualFold(text, "/show_weight"):
			preMsg, _ := storage.ShowPreviousEntry(chatID)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
		case strings.EqualFold(text, "/delete"):
			err := storage.DeleteRestorePreviousEntry(chatID, 0)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ошибка удаления: %v", err))
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "Последняя введенная запись удалена")
				bot.Send(msg)
			}
		case strings.EqualFold(text, "/restore"):
			err := storage.DeleteRestorePreviousEntry(chatID, 1)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ошибка восстановления: %v", err))
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "Последняя удаленная запись восстановлена")
				bot.Send(msg)
			}
		case strings.EqualFold(text, "/help"):
			msg := tgbotapi.NewMessage(chatID, messages.Help)
			bot.Send(msg)
		case err != nil:
			msg := tgbotapi.NewMessage(chatID, messages.ErrCommand)
			bot.Send(msg)
		case weightInput > 0:
			storedWeight, _ := storage.DiffWeight(chatID)
			storage.AddRecordToDB(models.NewRecord(int(chatID), weightInput, time.Now(), 0))

			diffWeight := weightInput - storedWeight
			var preMsg string
			if diffWeight == 0 {
				preMsg = fmt.Sprintf("Ваш вес %.2f кг записан\nРазница с прежним весом: %.2f кг",
					weightInput,
					diffWeight,
				)
			} else {
				preMsg = fmt.Sprintf("Ваш вес %.2f кг записан\nРазница с прежним весом: %+.2f кг",
					weightInput,
					diffWeight,
				)
			}
			msg := tgbotapi.NewMessage(chatID, preMsg)
			weightInput = 0
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(chatID, "Неизвестная команда")
			bot.Send(msg)
		}
	}
}
