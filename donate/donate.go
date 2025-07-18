package donate

import (
	"fmt"
	"os"
	"weightTrack_bot/backup"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/skip2/go-qrcode"
)

// Генерация QR-кода в виде []byte
func generateQRCode(content string) ([]byte, error) {
	qr, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		msg := fmt.Sprintf("Ошибка генерации QR-кода %v", err)
		backup.WriteLog(msg)
		return nil, err
	}
	return qr, nil
}

// Логика отправки доната
func DoDonate(sum float64, chatID int64) tgbotapi.PhotoConfig {

	donateLink := os.Getenv("TELEGRAM_BOT_DONATE")
	// Создаем QR-код
	qrCode, err := generateQRCode(donateLink)
	if err != nil {
		msg := fmt.Sprintf("Ошибка генерации QR-кода %v", err)
		backup.WriteLog(msg)
	}

	// Отправляем QR-код в чат
	file := tgbotapi.FileBytes{Name: "qrcode.png", Bytes: qrCode}
	photo := tgbotapi.NewPhoto(chatID, file)
	photo.Caption = fmt.Sprintf("🔹 Переведите **%.2f ₽**\n🔹 На поддержку бота\n\nОтсканируйте QR-код или перейдите по ссылке: %s", sum, donateLink)
	photo.ParseMode = "Markdown"

	return photo
}
