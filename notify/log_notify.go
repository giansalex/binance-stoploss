package notify

import (
	"log"
	"regexp"
)

const regex = `<.*?>`

type LogNotify struct {
}

// Send notify to log
func (lnotify *LogNotify) Send(message string) error {
	log.Println(lnotify.clearHtml(message))

	return nil
}

// This method uses a regular expresion to remove HTML tags.
func (lnotify *LogNotify) clearHtml(s string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}
