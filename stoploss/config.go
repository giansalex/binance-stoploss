package stoploss

// Config stop-loss configuration
type Config struct {
	OrderType        string
	Market           string
	Price            float64
	Quantity         string
	StopFactor       float64
	NotifyStopChange bool
}
