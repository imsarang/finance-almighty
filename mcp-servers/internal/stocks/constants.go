package stocks

import (
	"time"
	"errors"
)

const (
	TaskAnalyze = "stocks.analyze"
	QueueName = "stocks"
)

const (
	fundamentalsTimeout = 2 * time.Second
)

var (
	ErrTickerRequired = errors.New("Ticker is Required")
	ErrQuestionRequired = errors.New("Question is Required")
)