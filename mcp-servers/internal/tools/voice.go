package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type VoiceService interface {
	TranscribeVoiceInput(duration int) (string, error)
	TranscribeWhileKeyHeld(ctx context.Context) (string, error)
}

type VoiceTools struct {
	svc VoiceService
}

type captureVoice struct {
	DurationSeconds int `jsonschema:"How long to keep the microphone active. Default is 5 seconds."`
}

type transcribe struct {
	Transcription string
}

func VoiceToolInit(svc VoiceService) *VoiceTools {
	return &VoiceTools{svc: svc}
}

func (t *VoiceTools) Register(s *mcp.Server) {
	registerVoiceTools(s, t)
}

func registerVoiceTools(s *mcp.Server, t *VoiceTools) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "capture_local_voice",
		Description: "Records local microphone input for a specified duration and converts it to text offline using your machine's CPU/GPU.",
	}, t.handleCaptureLocalVoice)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "capture_voice_hold",
		Description: "Wait until Right Ctrl is held, record while it is held, and transcribe when it is released.",
	}, t.handleHoldToTalk)
}

func (t *VoiceTools) handleCaptureLocalVoice(_ context.Context, _ *mcp.CallToolRequest, in captureVoice) (*mcp.CallToolResult, transcribe, error) {
	duration := in.DurationSeconds
	if duration == 0 {
		duration = 5
	}

	text, err := t.svc.TranscribeVoiceInput(duration)
	if err != nil {
		return nil, transcribe{}, err
	}

	return nil, transcribe{Transcription: text}, nil
}

func (t *VoiceTools) handleHoldToTalk(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, transcribe, error) {
	text, err := t.svc.TranscribeWhileKeyHeld(ctx)
	if err != nil {
		return nil, transcribe{}, err
	}

	return nil, transcribe{Transcription: text}, nil
}
