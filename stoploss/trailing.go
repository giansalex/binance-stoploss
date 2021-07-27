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
	config    *Config
	market    string
	baseCoin  string
	countCoin string
	lastStop  float64
}

// NewTrailing new trailing instance
func NewTrailing(exchange Exchange, notify notify.SingleNotify, config *Config) *Trailing {
	pair := strings.Split(strings.ToUpper(config.Market), "/")

	tlg := &Trailing{
		exchange:  exchange,
		notify:    notify,
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
		msgFmt := "📉 **SELL**\n __Market:__ `%s`\n__Amount:__ %s %s\n__Price:__ %.6f\n__Order:__ %s"
		tlg.notify.Send(fmt.Sprintf(msgFmt, tlg.config.Market, quantity, tlg.baseCoin, marketPrice, order))
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
		msgFmt := "📈 **BUY**\n __Market:__ `%s`\n__Amount:__ %s %s\n__Price:__ %.6f\n__Order:__ %s"
		tlg.notify.Send(fmt.Sprintf(msgFmt, tlg.config.Market, quantity, tlg.baseCoin, marketPrice, order))
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
	if !tlg.config.NotifyStopChange {
		return
	}

	result := big.NewFloat(prev).Cmp(big.NewFloat(next))

	if result == 0 {
		return
	}

	tlg.notify.Send(fmt.Sprintf("Stop-loss %s (%s): %.6f - Market Price: %.6f", tlg.market, tlg.config.OrderType, next, price))
}
