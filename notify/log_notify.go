package notify

import "log"

type LogNotify struct {
}

// Send notify to log
func (lnotify *LogNotify) Send(message string) error {
	log.Println(message)

	return nil
}
