package notify

type SingleNotify interface {
	Send(message string) error
}
