package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"weightTrack_bot/messages"
	"weightTrack_bot/parse"
	"weightTrack_bot/storage"
)

const fileName = "dataBase.txt"

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

	// отправляем приветственное сообщение со списком команд
	//в разработке

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
		case strings.EqualFold(text, "/weight"):
			records, err := storage.ReadRecords(fileName)
			if err != nil {
				fmt.Printf("ошибка чтения файла: %e", err)
			}
			//for record := range records {
			//	fmt.Println(record)
			//}

			if err != nil {
				log.Printf("Ошибка чтения файла: %v", err)
				msg := tgbotapi.NewMessage(chatID, "Ошибка при чтении данных")
				bot.Send(msg)
				continue
			}

			for _, record := range records {
				//fmt.Println(record)
				// Или отправляем сообщение с каждой записью
				msg := tgbotapi.NewMessage(chatID, record)
				bot.Send(msg)
			}

			/*
				if weight > 0 {
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ваш вес: %.2f кг", weight))
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(chatID, "Вы еще не записывали свой вес")
					bot.Send(msg)
				}
			*/
		case strings.EqualFold(text, "/help"):
			msg := tgbotapi.NewMessage(chatID, messages.Help)
			bot.Send(msg)
		case err != nil:
			msg := tgbotapi.NewMessage(chatID, messages.ErrCommand)
			bot.Send(msg)
		case weightInput > 0:
			storage.AddRecordToDB(storage.NewRecord(int(chatID), weightInput), fileName)
			//weight = weightInput
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Ваш вес %.2f кг записан", weightInput))
			weightInput = 0
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(chatID, "Неизвестная команда")
			bot.Send(msg)
		}
	}
}
