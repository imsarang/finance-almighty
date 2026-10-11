package stocks

import (
	"context"
	"strings"

	"github.com/imsarang/finance-almighty/mcp-servers/internal/fundamentals"
	"github.com/imsarang/finance-almighty/mcp-servers/internal/queue"
)

type AnalysisRequest struct {
	Ticker string
	Question string
}

type Service struct {
	queue queue.Enqueuer
	fundamentals fundamentals.FundamentalsProvider
}

func StockService(q queue.Enqueuer, fp fundamentals.FundamentalsProvider) *Service {
	return &Service{queue: q, fundamentals: fp}
}

// RAG trigger point -> to enqueue into queue for async rag (future)
func (s *Service) RequestAnalysis(ctx context.Context, req AnalysisRequest) (string, error) {
	ticker := normalizeTicker(req.Ticker)

	question := strings.TrimSpace(req.Question)

	if ticker == "" {
		return "", ErrTickerRequired
	}

	if question == "" {
		return "", ErrQuestionRequired
	}

	// push to queue for async rag answers
	return s.queue.Enqueue(ctx, queue.Task{
		Task_type: TaskAnalyze,
		Queue: QueueName,
		Kwargs: map[string] any{"ticker": ticker, "question": question},
	})
}

// small outputs for direct questions
func (s *Service) Fundamentals(ctx context.Context, ticker string)(fundamentals.Fundamentals, error){
	ticker = normalizeTicker(ticker)
	
	if ticker == "" {
		return fundamentals.Fundamentals{}, ErrTickerRequired
	}

	// set timeout for 3rd party call
	ctx, cancel := context.WithTimeout(ctx, fundamentalsTimeout)

	defer cancel()

	return s.fundamentals.Fundamentals(ctx, ticker)
}

func normalizeTicker(str string) string{
	return strings.ToUpper(strings.TrimSpace(str))
}

