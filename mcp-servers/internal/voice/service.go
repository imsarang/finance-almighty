package voice

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

type Service struct {
	model whisper.Model
}

func VoiceService(model whisper.Model) *Service {
	return &Service{model: model}
}

// TranscribeVoiceInput records for a fixed duration, then transcribes.
func (s *Service) TranscribeVoiceInput(duration int) (string, error) {
	samples, err := CaptureMicrophonePCM(duration)
	if err != nil {
		return "", fmt.Errorf("audio capture helper failed: %w", err)
	}
	return s.textFromSamples(samples)
}

// TranscribeWhileKeyHeld records while Right Ctrl is held, then transcribes.
func (s *Service) TranscribeWhileKeyHeld(ctx context.Context) (string, error) {
	hold, cleanup, err := WaitForHold(ctx)
	if err != nil {
		return "", err
	}
	defer cleanup()

	select {
	case <-hold.PressedChannel:
	case <-ctx.Done():
		return "", ctx.Err()
	}

	samples, err := CaptureUntil(ctx, hold.ReleasedChannel)
	if err != nil {
		return "", fmt.Errorf("audio capture helper failed: %w", err)
	}
	return s.textFromSamples(samples)
}

// textFromSamples runs the loaded Whisper model on microphone samples.
func (s *Service) textFromSamples(samples []float32) (string, error) {
	whisperContext, err := s.model.NewContext()
	if err != nil {
		return "", fmt.Errorf("failed to create engine processing context: %w", err)
	}

	if err := whisperContext.SetLanguage("en"); err != nil {
		return "", fmt.Errorf("failed to set whisper language: %w", err)
	}

	if err := whisperContext.Process(samples, nil, nil, nil); err != nil {
		return "", fmt.Errorf("local Whisper computation failed: %w", err)
	}

	var transcript strings.Builder
	for {
		segment, err := whisperContext.NextSegment()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("read whisper segment: %w", err)
		}
		transcript.WriteString(segment.Text)
		transcript.WriteByte(' ')
	}

	return strings.TrimSpace(transcript.String()), nil
}
