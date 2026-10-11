package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/imsarang/finance-almighty/mcp-servers/internal/fundamentals"
	"github.com/imsarang/finance-almighty/mcp-servers/internal/stocks"
)

type StockService interface {
	RequestAnalysis(ctx context.Context, req stocks.AnalysisRequest)(string, error)
	Fundamentals(ctx context.Context, ticker string)(fundamentals.Fundamentals, error)
}

type StockTools struct {
	svc StockService
}

type analyzeIn struct {
	Ticker   string
	Question string
}

type analyzeOut struct {
	JobID string
}

type fundamentalsIn struct {
	Ticker string
}


func StockToolsConstructor(svc StockService) *StockTools{
	return &StockTools{svc: svc}
}

// register the tools for the mcp client to use
func(t *StockTools) Register(s *mcp.Server){

	// rag driven tool
	mcp.AddTool(s, &mcp.Tool{
		Name: "stocks_analyze",
		Description: "Deep analysis using RAG (stream output as processed (WIP))",
	}, t.analyze)

	// sync driven tool
	mcp.AddTool(s, &mcp.Tool{
		Name: "stocks_fundamentals",
		Description: "Fundamental analysis of a stock",
	}, t.getFundamentals)
}

func(t *StockTools) analyze(ctx context.Context, _ *mcp.CallToolRequest, in analyzeIn)(*mcp.CallToolResult, analyzeOut, error){
	id, err := t.svc.RequestAnalysis(ctx, stocks.AnalysisRequest{Ticker: in.Ticker, Question: in.Question})

	if err != nil {
		return nil, analyzeOut{}, err
	}

	return nil, analyzeOut{JobID: id}, nil
}

func(t *StockTools) getFundamentals(ctx context.Context, _ *mcp.CallToolRequest, in fundamentalsIn)(*mcp.CallToolResult, fundamentals.Fundamentals, error){
	f, err := t.svc.Fundamentals(ctx, in.Ticker)

	if err!=nil{
		return nil, fundamentals.Fundamentals{}, err
	}

	return nil, f, nil
}