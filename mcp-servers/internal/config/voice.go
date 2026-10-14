package config

type VoiceConfig struct {
	Config
	WhisperModelPath string
}

func LoadVoice() (VoiceConfig, error) {
	modelPath, err := requireEnv("WHISPER_MODEL_PATH")
	if err != nil {
		return VoiceConfig{}, err
	}
	return VoiceConfig{
		Config:           Load("voice-mcp"),
		WhisperModelPath: modelPath,
	}, nil
}
