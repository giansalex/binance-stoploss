package stoploss

import (
	"log"

	bNotify "github.com/giansalex/binance-stoploss/notify"
)

// Notify notify stoploss
type Notify struct {
	inners []bNotify.SingleNotify
}

// NewNotify create Notify instance
func NewNotify(notifiers []bNotify.SingleNotify) *Notify {
	return &Notify{notifiers}
}

// Send send message
func (notify *Notify) Send(message string) error {
	for _, v := range notify.inners {
		err := v.Send(message)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
