package voice

import (
	"context"
	"fmt"
	"time"

	"github.com/gordonklaus/portaudio"
)

const (
	SampleRate = 16000
	BufferSize = 4000
)

// CaptureUntil records until stopChannel is closed or the call is cancelled.
func CaptureUntil(ctx context.Context, stopChannel <-chan struct{}) ([]float32, error) {
	recordedSamples := make([]float32, 0, SampleRate)
	inputBuffer := make([]float32, BufferSize)

	stream, err := portaudio.OpenDefaultStream(1, 0, float64(SampleRate), len(inputBuffer), inputBuffer)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return nil, err
	}
	defer stream.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-stopChannel:
			return recordedSamples, nil
		default:
		}
		if err := stream.Read(); err != nil {
			return nil, err
		}
		frame := make([]float32, len(inputBuffer))
		copy(frame, inputBuffer)
		recordedSamples = append(recordedSamples, frame...)
	}
}

// CaptureMicrophonePCM records the microphone for a fixed number of seconds.
func CaptureMicrophonePCM(durationSeconds int) ([]float32, error) {
	if durationSeconds <= 0 {
		return nil, fmt.Errorf("duration must be positive")
	}

	recordedSamples := make([]float32, 0, SampleRate*durationSeconds)
	inputBuffer := make([]float32, BufferSize)

	stream, err := portaudio.OpenDefaultStream(1, 0, float64(SampleRate), len(inputBuffer), inputBuffer)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return nil, err
	}
	defer stream.Stop()

	endTime := time.Now().Add(time.Duration(durationSeconds) * time.Second)
	for time.Now().Before(endTime) {
		if err := stream.Read(); err != nil {
			return nil, err
		}
		frame := make([]float32, len(inputBuffer))
		copy(frame, inputBuffer)
		recordedSamples = append(recordedSamples, frame...)
	}

	return recordedSamples, nil
}
