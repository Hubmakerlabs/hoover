package config

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"gopkg.in/h2non/gentleman.v2"
)

type TokenPrice struct {
	ChainType string `json:"chainType"`
	Address   string
	Currency  string
	Price     float64
	Symbol    string
	Timestamp int64 `json:"timestamp"`
}

func GetTokenPriceByRedstone(tokenSymbol string, currency string,
	timestamp string) (float64, error) {
	cli := gentleman.New()
	cli.URL("https://api.redstone.finance")
	req := cli.Request()
	req.AddPath("/prices")
	req.AddQuery("symbols", fmt.Sprintf("%s,%s", strings.ToUpper(tokenSymbol),
		strings.ToUpper(currency)))
	req.AddQuery("provider", "redstone")
	if timestamp != "" {
		req.AddQuery("toTimestamp", timestamp)
	}

	resp, err := req.Send()
	if err != nil {
		return 0.0, err
	}

	if !resp.Ok {
		return 0.0, fmt.Errorf("get token: %s currency: %s prices from redstone failed",
			tokenSymbol, currency)
	}
	defer resp.Close()
	tokenJsonPath := fmt.Sprintf("%s.value", strings.ToUpper(tokenSymbol))
	currencyJsonPath := fmt.Sprintf("%s.value", strings.ToUpper(currency))
	prices := gjson.GetManyBytes(resp.Bytes(), tokenJsonPath, currencyJsonPath)
	if len(prices) != 2 {
		return 0.0, fmt.Errorf("get token: %s currency: %s prices from redstone failed, response price number incorrect",
			tokenSymbol, currency)
	}
	tokenPrice := prices[0].Float()
	currencyPrice := prices[1].Float()
	if currencyPrice <= 0.0 {
		return 0.0, fmt.Errorf("get currency: %s price from redstone less than 0.0; currencyPrice: %f",
			currency, currencyPrice)
	}
	price := tokenPrice / currencyPrice
	return price, nil
}
