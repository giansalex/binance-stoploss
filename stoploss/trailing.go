package stoploss

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/giansalex/binance-stoploss/notify"
)

// Trailing stop-loss runner
type Trailing struct {
	exchange  Exchange
	notify    notify.SingleNotify
	sLog      notify.SingleNotify
	config    *Config
	market    string
	baseCoin  string
	countCoin string
	lastStop  float64
}

// NewTrailing new trailing instance
func NewTrailing(exchange Exchange, notify notify.SingleNotify, logNotify notify.SingleNotify, config *Config) *Trailing {
	pair := strings.Split(strings.ToUpper(config.Market), "/")

	tlg := &Trailing{
		exchange:  exchange,
		notify:    notify,
		config:    config,
		sLog:      logNotify,
		market:    pair[0] + pair[1],
		baseCoin:  pair[0],
		countCoin: pair[1],
	}

	if tlg.config.OrderType == "BUY" {
		tlg.lastStop = math.MaxFloat64
	}

	return tlg
}

// RunStop check stop loss apply
func (tlg *Trailing) RunStop() bool {
	if tlg.config.OrderType == "BUY" {
		return tlg.runBuy()
	}

	return tlg.runSell()
}

func (tlg *Trailing) runSell() bool {
	marketPrice, err := tlg.exchange.GetMarketPrice(tlg.market)
	if err != nil {
		tlg.notify.Send("Cannot get market price, error:" + err.Error())
		return true
	}

	stop := tlg.getSellStop(marketPrice)

	if marketPrice > stop {
		tlg.notifyStopLossOnChange(tlg.lastStop, stop, marketPrice)

		tlg.lastStop = stop
		return false
	}

	quantity := tlg.config.Quantity
	if quantity == "" {
		quantity, err = tlg.exchange.GetBalance(tlg.baseCoin)
		if err != nil {
			tlg.notify.Send("Cannot get balance, error:" + err.Error())
			return true
		}
	}

	order, err := tlg.exchange.Sell(tlg.market, quantity)
	if err != nil {
		tlg.notify.Send("Cannot create sell order, error:" + err.Error())
	} else {
		msgFmt := "📉 ## <b>SELL</b> ##\n<i>Market:</i> <code>%s</code>\n<i>Amount:</i> %s <code>%s</code> \n<i>Price:</i> %.6f <code>%s</code>\n<i>Order:</i> %s"
		tlg.notify.Send(fmt.Sprintf(msgFmt, tlg.config.Market, quantity, tlg.baseCoin, marketPrice, tlg.countCoin, order))
		tlg.sLog.Send(msgFmt)
	}

	return true
}

func (tlg *Trailing) runBuy() bool {
	marketPrice, err := tlg.exchange.GetMarketPrice(tlg.market)
	if err != nil {
		tlg.notify.Send("Cannot get market price, error:" + err.Error())
		return true
	}

	stop := tlg.getBuyStop(marketPrice)

	if stop > marketPrice {
		tlg.notifyStopLossOnChange(tlg.lastStop, stop, marketPrice)

		tlg.lastStop = stop
		return false
	}

	quantity := tlg.config.Quantity
	if quantity == "" {
		quantity, err = tlg.exchange.GetBalance(tlg.countCoin)
		if err != nil {
			tlg.notify.Send("Cannot get balance, error:" + err.Error())
			return true
		}
	}

	order, err := tlg.exchange.Buy(tlg.market, quantity)
	if err != nil {
		tlg.notify.Send("Cannot create buy order, error:" + err.Error())
	} else {
		msgFmt := "📈 ## <b>BUY</b> ##\n<i>Market:</i> <code>%s</code>\n<i>Amount:</i> %s <code>%s</code> \n<i>Price:</i> %.6f <code>%s</code>\n<i>Order:</i> %s"
		tlg.notify.Send(fmt.Sprintf(msgFmt, tlg.config.Market, quantity, tlg.baseCoin, marketPrice, tlg.countCoin, order))
		tlg.sLog.Send(msgFmt)
	}

	return true
}

func (tlg *Trailing) getBuyStop(price float64) float64 {
	if tlg.config.StopFactor > 0 {
		return math.Min(tlg.lastStop, price*(1+tlg.config.StopFactor))
	}

	return tlg.config.Price
}

func (tlg *Trailing) getSellStop(price float64) float64 {
	if tlg.config.StopFactor > 0 {
		return math.Max(tlg.lastStop, price*(1-tlg.config.StopFactor))
	}

	return tlg.config.Price
}

func (tlg *Trailing) notifyStopLossOnChange(prev float64, next float64, price float64) {
	result := big.NewFloat(prev).Cmp(big.NewFloat(next))

	if result == 0 {
		return
	}

	msg := fmt.Sprintf("New stoploss %s (%s): %.6f - Market Price: %.6f", tlg.market, tlg.config.OrderType, next, price)
	tlg.sLog.Send(msg)

	if !tlg.config.NotifyStopChange {
		return
	}
	tlg.notify.Send(msg)
}
