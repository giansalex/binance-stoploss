package stoploss

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Notify notify stoploss
type Notify struct {
	tlgToken string
	chatID   int64
}

// NewNotify create Notify instance
func NewNotify(telegramToken string, channelID int64) *Notify {
	return &Notify{telegramToken, channelID}
}

// Send send message
func (notify *Notify) Send(message string) {
	log.Println(message)

	if notify.tlgToken == "" {
		return
	}

	bot, err := tgbotapi.NewBotAPI(notify.tlgToken)
	if err != nil {
		fmt.Println("Cannot connect to Telegram:", err.Error())

		return
	}

	msg := tgbotapi.NewMessage(notify.chatID, message)

	_, err = bot.Send(msg)

	if err != nil {
		fmt.Println("Cannot send message to telegram:", err.Error())
	}
}
