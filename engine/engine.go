package engine

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"weightTrack_bot/donate"
	"weightTrack_bot/messages"
	"weightTrack_bot/models"
	"weightTrack_bot/parse"
	"weightTrack_bot/plots"
	"weightTrack_bot/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UserState хранит состояние ввода для конкретного пользователя
type UserState struct {
	IsAgeInput    bool
	IsHeightInput bool
	IsWeightInput bool
	WeightInput   float64
	HeightInput   float64
	AgeInput      int64
}

// Reset сбрасывает все флаги ввода (IsAgeInput, IsHeightInput, IsWeightInput) в false.
func (us *UserState) Reset() {
	*us = UserState{} // сбрасывает все поля в их нулевые значения
}

// Хранилище состояний пользователей
var userStates = make(map[int64]*UserState)

// GetUserState возвращает состояние пользователя (создаёт, если не существует)
func GetUserState(chatID int64) *UserState {
	if userStates[chatID] == nil {
		userStates[chatID] = &UserState{}
	}
	return userStates[chatID]
}

func StartBotGoroutine(
	id int64,
	bot *tgbotapi.BotAPI,
	wg *sync.WaitGroup,
	handler func(update tgbotapi.Update, bot *tgbotapi.BotAPI, wg *sync.WaitGroup)) *models.Goroutine {
	g := &models.Goroutine{
		ID:    id,
		Input: make(chan any),
		Stop:  make(chan struct{}),
	}

	go func() {
		defer wg.Done() // Уменьшаем счётчик горутин WaitGroup
		for {
			select {
			case data := <-g.Input:
				if update, ok := data.(tgbotapi.Update); ok {
					handler(update, bot, wg) // Вызываем переданный обработчик
				}
			case <-g.Stop:
				return
			}
		}
	}()
	return g
}

func Engine(update tgbotapi.Update, bot *tgbotapi.BotAPI, wg *sync.WaitGroup) {

	// Получаем состояние для текущего пользователя
	state := GetUserState(update.Message.Chat.ID)

	//извлекаем текст сообщенния и идентификатор чата для отправки ответа
	text := update.Message.Text
	chatID := update.Message.Chat.ID
	//парсим сообщение пользователя на наличие числа для записи веса
	if _, err := parse.ParseFloat(update.Message.Text); err == nil && state.IsWeightInput {
		state.WeightInput, _ = parse.ParseFloat(update.Message.Text)
	}
	//парсим рост пользователя
	if state.IsHeightInput {
		state.HeightInput, _ = parse.ParseFloat(update.Message.Text)
	}
	//парсим возраст пользователя
	if state.IsAgeInput {
		state.AgeInput, _ = parse.ParseInt(update.Message.Text)
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
		state.Reset()
	case strings.EqualFold(text, "/show_bmi"):
		var save_weight string
		var edit_height string
		var edit_age string

		user, positionU, errU := storage.FindUserPosition(chatID)
		records, errR := storage.ReadRecords(int(chatID))
		record, positionR := storage.FindLastEntry(records, 0)

		if user.GetHeight() == 0 {
			edit_height = "/edit_height - редактировать рост\n"
		}
		if user.GetAge() == 0 {
			edit_age = "/edit_age - редактировать возраст\n"
		}
		if positionR == -1 {
			save_weight = "/save_weight - записать вес\n"
		}

		if errU != nil || errR != nil || positionU == -1 || positionR == -1 ||
			user.GetAge() == 0 || user.GetHeight() == 0 {
			preMsg := fmt.Sprintf("Недостаточно данных:\n%s%s%s",
				save_weight,
				edit_height,
				edit_age,
			)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
		}

		if positionR != -1 && positionU != -1 && user.GetHeight() != 0 && user.GetAge() != 0 {
			bmi, assessment := storage.FindBMI(user, record)
			preMsg := fmt.Sprintf("Ваш ИМТ равен: %.2f\n%s", bmi, assessment)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
		}
		state.Reset()
	case state.IsHeightInput && state.HeightInput > 0:
		if state.HeightInput > 999.0 {
			preMsg := fmt.Sprintf("Вы ввели %.2f\nРост не может быть больше 999 см", state.HeightInput)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
			state.Reset()
			return
		}
		user, position, err := storage.FindUserPosition(chatID)
		if position != -1 && err == nil {
			err := storage.UpdateUser(chatID, user, user.GetAge(), state.HeightInput)
			if err != nil {
				fmt.Printf("ошибка %v\n", err)
			} else {
				preMsg := fmt.Sprintf("Ваш рост %.2f см записан\n", state.HeightInput)
				msg := tgbotapi.NewMessage(chatID, preMsg)
				bot.Send(msg)
			}
		} else {
			err := storage.AddUserToDB(models.NewUser(chatID, int(state.AgeInput), state.HeightInput))
			if err != nil {
				fmt.Println("ошибка добавления роста пользователя")
			} else {
				preMsg := fmt.Sprintf("Ваш рост %.2f см записан\n", state.HeightInput)
				msg := tgbotapi.NewMessage(chatID, preMsg)
				bot.Send(msg)
			}
		}
		state.Reset()
	case strings.EqualFold(text, "/edit_height"):
		// При множественном вводе команд оставляем только последнюю
		state.IsAgeInput = false
		state.IsWeightInput = false

		state.IsHeightInput = true
		msg := tgbotapi.NewMessage(chatID, "Введите рост в сантиметрах")
		bot.Send(msg)
	case state.IsAgeInput && state.AgeInput > 0:
		if state.AgeInput > 999 {
			preMsg := fmt.Sprintf("Вы ввели %d\nВозраст не может быть больше 999 лет", state.AgeInput)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
			state.Reset()
			return
		}
		user, position, err := storage.FindUserPosition(chatID)
		if position != -1 && err == nil {
			err := storage.UpdateUser(chatID, user, int(state.AgeInput), user.GetHeight())
			if err != nil {
				fmt.Printf("ошибка %v\n", err)
			} else {
				preMsg := fmt.Sprintf("Ваш возраст %2d лет(года) записан\n", state.AgeInput)
				msg := tgbotapi.NewMessage(chatID, preMsg)
				bot.Send(msg)
			}
		} else {
			err := storage.AddUserToDB(models.NewUser(chatID, int(state.AgeInput), state.HeightInput))
			if err != nil {
				fmt.Println("ошибка добавления возраста пользователя")
			} else {
				preMsg := fmt.Sprintf("Ваш возраст %2d лет(года) записан\n", state.AgeInput)
				msg := tgbotapi.NewMessage(chatID, preMsg)
				bot.Send(msg)
			}
		}
		state.Reset()
	case strings.EqualFold(text, "/edit_age"):
		// При множественном вводе команд оставляем только последнюю
		state.IsHeightInput = false
		state.IsWeightInput = false

		state.IsAgeInput = true
		msg := tgbotapi.NewMessage(chatID, "Введите возраст (полных лет)")
		bot.Send(msg)
	case strings.EqualFold(text, "/start"):
		msg := tgbotapi.NewMessage(chatID, messages.WelcomeMsg)
		bot.Send(msg)
		state.Reset()
	case strings.EqualFold(text, "/show_week"):
		period, err := storage.FindPeriod(chatID, 7)
		if err != nil {
			preMsg := fmt.Sprintf("Не удалось прочитать данные: %v\n", err)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
			state.Reset()
			return
		}
		preMsg := storage.ShowPeriod(period, 7)
		msg := tgbotapi.NewMessage(chatID, preMsg)
		bot.Send(msg)
		state.Reset()
	case strings.EqualFold(text, "/show_month"):
		period, err := storage.FindPeriod(chatID, 31)
		if err != nil {
			preMsg := fmt.Sprintf("Не удалось прочитать данные: %v\n", err)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
			state.Reset()
			return
		}
		preMsg := storage.ShowPeriod(period, 31)
		msg := tgbotapi.NewMessage(chatID, preMsg)
		bot.Send(msg)
		state.Reset()
	case strings.EqualFold(text, "/show_progress"):
		period, err := storage.FindPeriod(chatID, 31)
		if err != nil {
			preMsg := fmt.Sprintf("Не удалось прочитать данные: %v\n", err)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
			state.Reset()
			return
		}
		// Создаем график в памяти
		imgBytes, err := plots.MakePlot(period)
		if err != nil {
			preMsg := fmt.Sprintf("Не удалось создать график: %v\n", err)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
			state.Reset()
			return
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
			preMsg := fmt.Sprintf("Не удалось отправить график: \n%v\n", err)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
		}
		state.Reset()
	case strings.EqualFold(text, "/show_weight"):
		preMsg, err := storage.ShowPreviousEntry(chatID)
		if err != nil {
			preMsg = fmt.Sprintf("Ошибка: %v", err)
		}
		msg := tgbotapi.NewMessage(chatID, preMsg)
		bot.Send(msg)
		state.Reset()
	case strings.EqualFold(text, "/delete"):
		err := storage.DeleteRestorePreviousEntry(chatID, 0)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ошибка удаления: %v", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Последняя введенная запись удалена")
			bot.Send(msg)
		}
		state.Reset()
	case strings.EqualFold(text, "/restore"):
		err := storage.DeleteRestorePreviousEntry(chatID, 1)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ошибка восстановления: %v", err))
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Последняя удаленная запись восстановлена")
			bot.Send(msg)
		}
		state.Reset()
	case strings.EqualFold(text, "/help"):
		msg := tgbotapi.NewMessage(chatID, messages.Help)
		bot.Send(msg)
		state.Reset()
	case state.IsWeightInput && state.WeightInput > 0 && (!state.IsAgeInput || !state.IsHeightInput):
		if state.WeightInput > 999.00 {
			preMsg := fmt.Sprintf("Вы ввели %.2f\nВес не может быть больше 999 кг", state.WeightInput)
			msg := tgbotapi.NewMessage(chatID, preMsg)
			bot.Send(msg)
			state.Reset()
			return
		}
		storedWeight, err := storage.DiffWeight(chatID)
		if err != nil {
			storedWeight = state.WeightInput
		}

		storage.AddRecordToDB(models.NewRecord(int(chatID), state.WeightInput, time.Now(), 0))

		diffWeight := state.WeightInput - storedWeight
		var preMsg string
		if diffWeight == 0 {
			preMsg = fmt.Sprintf("Ваш вес %.2f кг записан\n",
				state.WeightInput,
			)
		} else {
			preMsg = fmt.Sprintf("Ваш вес %.2f кг записан.\nРазница с прежним весом: %+.2f кг",
				state.WeightInput,
				diffWeight,
			)
		}
		msg := tgbotapi.NewMessage(chatID, preMsg)
		bot.Send(msg)
		state.Reset()
	case strings.EqualFold(text, "/save_weight"):
		// При множественном вводе команд оставляем только последнюю
		state.IsAgeInput = false
		state.IsHeightInput = false

		state.IsWeightInput = true
		msg := tgbotapi.NewMessage(chatID, "Введите вес в килограммах")
		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(chatID, "Неизвестная команда")
		bot.Send(msg)
		state.Reset()
	}
}
