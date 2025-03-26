package models

import "time"

type Data struct {
	Timestamp time.Time
	Error     []interface{} `json:"error" `
	Result    Result        `json:"result" `
}
type Result struct {
	Usdt_usd Usdt_usd `json:"USDTZUSD"`
}

type Usdt_usd struct {
	RawAsks [][]interface{} `json:"asks"`
	RawBids [][]interface{} `json:"bids"`
	Ask     *Ask
	Bid     *Bid
}

type Ask struct {
	Price     float64
	Quantity  float64
	Timestamp int64
}

type Bid struct {
	Price     float64
	Quantity  float64
	Timestamp int64
}
