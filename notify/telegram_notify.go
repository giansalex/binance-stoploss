package notify

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramNotify struct {
	token  string
	chatID int64
}

func NewTelegramNotify(token string, chatID int64) *TelegramNotify {
	return &TelegramNotify{token: token, chatID: chatID}
}

// Send notify to telegram
func (tgNotify *TelegramNotify) Send(message string) error {
	bot, err := tgbotapi.NewBotAPI(tgNotify.token)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(tgNotify.chatID, message)
	msg.ParseMode = "markdown"

	_, err = bot.Send(msg)

	return err
}
