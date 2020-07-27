package stoploss

import (
	"context"
	"strconv"
	"strings"

	"github.com/adshao/go-binance"
)

// Binance exchange
type Binance struct {
	api *binance.Client
	ctx context.Context
}

// NewExchange create Exchange instance
func NewExchange(ctx context.Context, api *binance.Client) *Binance {
	return &Binance{api, ctx}
}

// GetBalance get balance for coin
func (exchange *Binance) GetBalance(coin string) (string, error) {
	coin = strings.ToUpper(coin)
	account, err := exchange.api.NewGetAccountService().Do(exchange.ctx)
	if err != nil {
		return "0", err
	}

	for _, balance := range account.Balances {
		if strings.ToUpper(balance.Asset) == coin {
			return balance.Free, nil
		}
	}

	return "0", nil
}

// GetMarketPrice get last price for market pair
func (exchange *Binance) GetMarketPrice(market string) (float64, error) {
	prices, err := exchange.api.NewListPricesService().Symbol(market).Do(exchange.ctx)
	if err != nil {
		return 0, err
	}

	for _, p := range prices {
		if p.Symbol == market {
			return strconv.ParseFloat(p.Price, 64)
		}
	}

	return 0, nil
}

// Sell create a sell order to market price
func (exchange *Binance) Sell(market string, quantity string) (string, error) {
	order, err := exchange.api.NewCreateOrderService().Symbol(market).
		Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).
		Quantity(quantity).Do(exchange.ctx)

	if err != nil {
		return "", err
	}

	return strconv.FormatInt(order.OrderID, 10), nil
}

// Buy create a buy order to market price
func (exchange *Binance) Buy(market string, quantity string) (string, error) {
	order, err := exchange.api.NewCreateOrderService().Symbol(market).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).
		Quantity(quantity).Do(exchange.ctx)

	if err != nil {
		return "", err
	}

	return strconv.FormatInt(order.OrderID, 10), nil
}
