package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/imsarang/finance-almighty/mcp-servers/internal/queue"
)

type JobsTools struct {
	results queue.ResultReader
}

type jobStatusIn struct {
	JobID string
}

func JobsToolsInit(r queue.ResultReader) *JobsTools {
	return &JobsTools{results: r}
}

// register the tools for the mcp client to use
func (j *JobsTools) Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "job_status",
		Description: "Get the status and result of a job enqueued by stocks_analyze",
	}, j.status)
}

func (j *JobsTools) status(ctx context.Context, _ *mcp.CallToolRequest, in jobStatusIn) (*mcp.CallToolResult, queue.Result, error) {
	res, err := j.results.Result(ctx, in.JobID)

	if err != nil {
		return nil, queue.Result{}, err
	}

	return nil, res, nil
}
