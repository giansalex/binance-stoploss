package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/adshao/go-binance"
	"github.com/giansalex/binance-stoploss/stoploss"
	"github.com/hashicorp/go-retryablehttp"
)

var (
	typePtr     = flag.String("type", "SELL", "order type: SELL or BUY")
	pairPtr     = flag.String("pair", "", "market pair, example: BNB/USDT")
	percentPtr  = flag.Float64("percent", 0.00, "percent (for trailing stop loss), example: 3.0 (3%)")
	pricePtr    = flag.Float64("price", 0.00, "price (for static stop loss), example: 9200.00 (BTC price)")
	intervalPtr = flag.Int("interval", 30, "interval in seconds to update price, example: 30 (30 sec.)")
	amountPtr   = flag.String("amount", "", "(optional) amount to order (sell or buy) on stoploss")
	chatPtr     = flag.Int64("telegram.chat", 0, "(optional) telegram User ID for notify")
)

func main() {
	flag.Parse()
	apiKey := os.Getenv("BINANCE_APIKEY")
	secret := os.Getenv("BINANCE_SECRET")

	if apiKey == "" || secret == "" {
		log.Fatal("BINANCE_APIKEY, BINANCE_SECRET are required")
	}

	if pairPtr == nil || *pairPtr == "" {
		log.Fatal("pair market is required")
	}

	if (percentPtr == nil || *percentPtr <= 0) && (pricePtr == nil || *pricePtr <= 0) {
		log.Fatal("a price or percent parameter is required")
	}

	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RetryMax = 10
	api := binance.NewClient(apiKey, secret)
	api.HTTPClient = retryClient.StandardClient()
	notify := stoploss.NewNotify(os.Getenv("TELEGRAM_TOKEN"), *chatPtr)
	trailing := stoploss.NewTrailing(stoploss.NewExchange(context.Background(), api), notify, *typePtr, *pairPtr, *percentPtr/100, *amountPtr, *pricePtr)

	for {
		if trailing.RunStop() {
			break
		}

		time.Sleep(time.Duration(*intervalPtr) * time.Second)
	}
}
