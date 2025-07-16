package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"weightTrack_bot/backup"
	"weightTrack_bot/donate"
	"weightTrack_bot/messages"
	"weightTrack_bot/models"
	"weightTrack_bot/parse"
	"weightTrack_bot/plots"
	"weightTrack_bot/storage"
)

var isAgeInput bool
var isHeightInput bool
var isWeightInput bool

func main() {

	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем из переменной окружения токен для подключения к телеграм-боту
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	// Резервное сохранение данных
	var path, _ = parse.ParseInt(os.Getenv("TELEGRAM_BOT_PATH"))
	var filePath1 = os.Getenv("TELEGRAM_BOT_BACKUP_FILE1")
	var filePath2 = os.Getenv("TELEGRAM_BOT_BACKUP_FILE2")
	go backup.StartDailyBackup(bot, filePath1, path)
	go backup.StartDailyBackup(bot, filePath2, path)

	// Настраиваем получение обновлений, устаналиваем время ожидания новых сообщений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Получаем канал обновлений
	updates := bot.GetUpdatesChan(u)

	// Временное хранение ввода
	var weightInput float64 //вес
	var heightInput float64 //рост
	var ageInput int64      //возраст

	// Получаем сообщения из канала updates в бесконечном цикле
	for update := range updates {
		if update.Message == nil {
			continue
		}

		//извлекаем текст сообщенния и идентификатор чата для отправки ответа
		text := update.Message.Text
		chatID := update.Message.Chat.ID

		//парсим сообщение пользователя на наличие числа для записи веса
		if _, err := parse.ParseFloat(update.Message.Text); err == nil && isWeightInput {
			weightInput, _ = parse.ParseFloat(update.Message.Text)
		}
		//парсим рост пользователя
		if isHeightInput {
			heightInput, _ = parse.ParseFloat(update.Message.Text)
		}
		//парсим возраст пользователя
		if isAgeInput {
			ageInput, _ = parse.ParseInt(update.Message.Text)
		}

		//выбираем ответ по запросу
		switch {
		case strings.HasPrefix(update.Message.Text, "/donate"):
			amount, err := parse.ParseFloat(update.Message.Text)
			if err != nil {
				photo := donate.DoDonate(100.00, chatID)
				if _, err := bot.Send(photo); err != nil {
					log.Println("Ошибка отправки QR:", err)
				}
			} else {
				photo := donate.DoDonate(amount, chatID)
				if _, err := bot.Send(photo); err != nil {
					log.Println("Ошибка отправки QR:", err)
				}
			}
		case strings.EqualFold(text, "/show_bmi"):
			user, position, err := storage.FindUserPosition(chatID)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Укажите свой возраст и рост с помощью команд:\n/edit_height -редактировать рост,\n/edit_age - редактировать возраст")
				bot.Send(msg)
			}

			records, err := storage.ReadRecords(int(chatID))
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Укажите свой возраст и рост с помощью команд:\n/edit_height -редактировать рост,\n/edit_age - редактировать возраст")
				bot.Send(msg)
			}

			record, _ := storage.FindLastEntry(records, 0)
			if position != -1 && err == nil && user.GetHeight() != 0 {
				bmi, assessment := storage.FindBMI(user, record)
				preMsg := fmt.Sprintf("Ваш ИМТ равен: %.2f\n%s", bmi, assessment)
				msg := tgbotapi.NewMessage(chatID, preMsg)
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "Укажите свой возраст и рост с помощью команд:\n/edit_height -редактировать рост,\n/edit_age - редактировать возраст")
				bot.Send(msg)
			}
		case isHeightInput && heightInput > 0:
			if heightInput > 999.0 {
				preMsg := fmt.Sprintf("Вы ввели %.2f\nРост не может быть больше 999 см", heightInput)
				msg := tgbotapi.NewMessage(chatID, preMsg)
				bot.Send(msg)
				continue
			}
			user, position, err := storage.FindUserPosition(chatID)
			if position != -1 && err == nil {
				err := storage.UpdateUser(chatID, user, user.GetAge(), heightInput)
				if err != nil {
					fmt.Printf("ошибка %v\n", err)
				} else {
					preMsg := fmt.Sprintf("Ваш рост %.2f см записан\n", heightInput)
					msg := tgbotapi.NewMessage(chatID, preMsg)
					bot.Send(msg)
				}
			} else {
				err := storage.AddUserToDB(models.NewUser(chatID, int(ageInput), heightInput))
				if err != nil {
					fmt.Println("ошибка добавления роста пользователя")
				} else {
					preMsg := fmt.Sprintf("Ваш рост %.2f см записан\n", heightInput)
					msg := tgbotapi.NewMessage(chatID, preMsg)
					bot.Send(msg)
				}
			}
			isHeightInput = false
			heightInput = 0
		case strings.EqualFold(text, "/edit_height"):
			isHeightInput = true
		case isAgeInput && ageInput > 0:
			if ageInput > 999 {
				preMsg := fmt.Sprintf("Вы ввели %d\nВозраст не может быть больше 999 лет", ageInput)
				msg := tgbotapi.NewMessage(chatID, preMsg)
				bot.Send(msg)
				continue
			}
			user, position, err := storage.FindUserPosition(chatID)
			if position != -1 && err == nil {
				err := storage.UpdateUser(chatID, user, int(ageInput), user.GetHeight())
				if err != nil {
					fmt.Printf("ошибка %v\n", err)
				} else {
					preMsg := fmt.Sprintf("Ваш возраст %2d лет(года) записан\n", ageInput)
					msg := tgbotapi.NewMessage(chatID, preMsg)
					bot.Send(msg)
				}
			} else {
				err := storage.AddUserToDB(models.NewUser(chatID, int(ageInput), heightInput))
				if err != nil {
					fmt.Println("ошибка добавления возраста пользователя")
				} else {
					preMsg := fmt.Sprintf("Ваш возраст %2d лет(года) записан\n", ageInput)
					msg := tgbotapi.NewMessage(chatID, preMsg)
					bot.Send(msg)
				}
			}
			isAgeInput = false
			ageInput = 0
		case strings.EqualFold(text, "/edit_age"):
			isAgeInput = true
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
		case isWeightInput && weightInput > 0 && (!isAgeInput || !isHeightInput):
			if weightInput > 999.00 {
				preMsg := fmt.Sprintf("Вы ввели %.2f\nВес не может быть больше 999 кг", weightInput)
				msg := tgbotapi.NewMessage(chatID, preMsg)
				bot.Send(msg)
				continue
			}
			storedWeight, err := storage.DiffWeight(chatID)
			if err != nil {
				storedWeight = weightInput
			}

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
			isWeightInput = false
		case strings.EqualFold(text, "/weight"):
			isWeightInput = true
		default:
			msg := tgbotapi.NewMessage(chatID, "Неизвестная команда")
			bot.Send(msg)
		}
	}
}
