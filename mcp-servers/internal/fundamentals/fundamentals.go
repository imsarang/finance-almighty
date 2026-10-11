package fundamentals

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("Symbol not found")
	ErrRateLimited = errors.New("provider rate limit exceeded")
)

type Profile struct {
	Ticker string
	Name string
	Exchange string
	Country string
	Currency string
	Industry string
	IPO string
	MarketCap float64
	SharesOutstanding float64
	WebURL string
	Logo string
}

type Metrics struct {
	High52Week float64
	Low52Week float64
	Low52WeekDate string
	Beta float64
	AvgVolume10Day float64
	All map[string]any `json:",omitempty"`
}

type Fundamentals struct {
	Profile Profile
	Metrics Metrics
	Source string
}

// interface to call the api
type FundamentalsProvider interface {
	Fundamentals(ctx context.Context, ticker string) (Fundamentals, error)
}
