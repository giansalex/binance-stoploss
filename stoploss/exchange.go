package stoploss

// Exchange wrapper to connect to exchange
type Exchange interface {
	GetBalance(coin string) (string, error)
	GetMarketPrice(market string) (string, error)
	Sell(market string, quantity string) (string, error)
	Buy(market string, quantity string) (string, error)
}
