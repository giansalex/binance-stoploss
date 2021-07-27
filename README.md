# Binance Stop-Loss Bot [![Go Report Card](https://goreportcard.com/badge/github.com/giansalex/binance-stoploss)](https://goreportcard.com/report/github.com/giansalex/binance-stoploss)

[Binance](https://binance.com/) Trailing Stop-Loss Bot and optional Telegram notifications. 

> A trailing stop order sets the stop price at a fixed amount below the market price with an attached "trailing" amount. As the market price rises, the stop price rises by the trail amount, but if the stock price falls, the stop loss price doesn't change, and a market order is submitted when the stop price is hit.

## Run

Require [API Keys](https://www.binance.com/en/usercenter/settings/api-management).    
Set Environment variables:
- `BINANCE_APIKEY`
- `BINANCE_SECRET`
- `TELEGRAM_TOKEN` (optional to notify)

Simple command to run bot stoploss
> Sell all BTC balance to market price when down 3%.
```sh
./binance -pair=BTC/USDT -percent=3
```

For buy orders
> Buy 0.01 `BTC` when price up 0.5%
```sh
./binance -type=BUY -pair=BTC/USDT -percent=0.5 -amount=0.01
```

For sell orders with static stoploss
> SELL 0.1 BTC when `BTC` down to 9400 USDT
```sh
./binance -pair=BTC/USDT -price=9400 -amount=0.1
```

Use telegram for notifications.
```sh
./binance -pair=BTC/USDT -percent=3 -interval=60 -telegram.chat=<user-id>
```
> For get user id, talk o the [userinfobot](https://t.me/userinfobot)

List available parameters 
```sh
  -amount string
        (optional) amount to order (sell or buy) on stoploss
  -interval int
        interval in seconds to update price, example: 30 (30 sec.) (default 30)
  -mail.from string
        (optional) email sender
  -mail.host string
        (optional) SMTP Host
  -mail.pass string
        (optional) SMTP Password
  -mail.port int
        (optional) SMTP Port (default 587)
  -mail.to string
        (optional) email receptor
  -mail.user string
        (optional) SMTP User
  -pair string
        market pair, example: BNB/USDT
  -percent float
        percent (for trailing stop loss), example: 3.0 (3%)
  -price float
        price (for static stop loss), example: 9200.00 (BTC price)
  -stop-change
        Notify on stoploss change (default: false)
  -telegram.chat int
        (optional) telegram User ID for notify
  -type string
        order type: SELL or BUY (default "SELL")

```
