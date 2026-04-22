package main

import (
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type APISettings struct {
	APIKey  string `json:"apiKey"`
	BaseURL string `json:"baseURL"`
	Model   string `json:"model"`
}

type Config struct {
	ASR                 APISettings `json:"asr"`
	TranscriptionPrompt string      `json:"transcriptionPrompt"`
	AutoHotKeyExec      string      `json:"autoHotKeyExec"`
}

type ResolvedAPISettings struct {
	APIKey  string
	BaseURL string
	Model   string
}

const (
	defaultASRModel            = openai.Whisper1
	defaultTranscriptionPrompt = "这是中文转写，主要内容是中文口述、日常输入，以及编程和软件开发相关内容。"
)

func (config Config) ResolveASRSettings() ResolvedAPISettings {
	return config.resolveAPISettings(config.ASR, defaultASRModel)
}

func (config Config) ResolveTranscriptionPrompt() string {
	return firstNonEmpty(config.TranscriptionPrompt, defaultTranscriptionPrompt)
}

func (config Config) resolveAPISettings(settings APISettings, defaultModel string) ResolvedAPISettings {
	return ResolvedAPISettings{
		APIKey:  firstNonEmpty(settings.APIKey),
		BaseURL: normalizeBaseURL(firstNonEmpty(settings.BaseURL)),
		Model:   firstNonEmpty(settings.Model, defaultModel),
	}
}

func (settings ResolvedAPISettings) Validate(name string) error {
	if strings.TrimSpace(settings.APIKey) == "" {
		return fmt.Errorf("%s API key is empty; configure %s.apiKey in config.json", name, strings.ToLower(name))
	}

	return nil
}

func newOpenAIClient(settings ResolvedAPISettings) *openai.Client {
	config := openai.DefaultConfig(settings.APIKey)
	if settings.BaseURL != "" {
		config.BaseURL = settings.BaseURL
	}

	return openai.NewClientWithConfig(config)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}

	return ""
}

func normalizeBaseURL(baseURL string) string {
	trimmed := strings.TrimSpace(baseURL)
	return strings.TrimRight(trimmed, "/")
}
