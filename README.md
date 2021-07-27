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

## Notifications

- Telegram.
```sh
./binance -pair=BTC/USDT -percent=3 -interval=60 -telegram.chat=<user-id>
```
![telegram binance bot](https://user-images.githubusercontent.com/14926587/127190283-b7117dd2-dd03-421b-ae49-95997503ae67.png)

> For get user id, talk to the [userinfobot](https://t.me/userinfobot)


- Mailing.
```sh
./binance -pair=BTC/USDT -percent=3 \
      -mail.host="smtp.example.com" \
      -mail.port=587 \
      -mail.user="user@example.com" \
      -mail.pass="xxxx" \
      -mail.from="user@example.com" \
      -mail.to="bob@gmail.com"
```

> You can notify both: telegram, mail.

## Docker

You can run in docker container. 
```bash
docker pull giansalex/binance-stoploss
# create container
docker run -d --name binance_sell_BTC giansalex/binance-stoploss \
      -type=BUY \
      -pair=BTC/USDT \
      -percent=5 \
      -amount=0.01
```

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
