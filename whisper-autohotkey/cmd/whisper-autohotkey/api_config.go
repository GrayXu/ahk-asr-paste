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
	API                 APISettings `json:"api"`
	ASR                 APISettings `json:"asr"`
	Command             APISettings `json:"command"`
	TranscriptionPrompt string      `json:"transcriptionPrompt"`
	OpenapiKey          string      `json:"openapiKey"`
	AutoHotKeyExec      string      `json:"autoHotKeyExec"`
	Coding              bool        `json:"coding"`
}

type ResolvedAPISettings struct {
	APIKey  string
	BaseURL string
	Model   string
}

const (
	defaultASRModel            = openai.Whisper1
	defaultCommandModel        = openai.GPT4
	defaultTranscriptionPrompt = "这是中文转写，主要内容是中文口述、日常输入，以及编程和软件开发相关内容。"
)

func (config Config) ResolveASRSettings() ResolvedAPISettings {
	return config.resolveAPISettings(config.ASR, defaultASRModel)
}

func (config Config) ResolveCommandSettings() ResolvedAPISettings {
	return config.resolveAPISettings(config.Command, defaultCommandModel)
}

func (config Config) ResolveTranscriptionPrompt() string {
	return firstNonEmpty(config.TranscriptionPrompt, defaultTranscriptionPrompt)
}

func (config Config) resolveAPISettings(override APISettings, defaultModel string) ResolvedAPISettings {
	apiKey := firstNonEmpty(override.APIKey, config.API.APIKey, config.OpenapiKey)
	baseURL := normalizeBaseURL(firstNonEmpty(override.BaseURL, config.API.BaseURL))
	model := firstNonEmpty(override.Model, config.API.Model, defaultModel)

	return ResolvedAPISettings{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
	}
}

func (settings ResolvedAPISettings) Validate(name string) error {
	if strings.TrimSpace(settings.APIKey) == "" {
		return fmt.Errorf("%s API key is empty; configure %s.apiKey, api.apiKey, or openapiKey in config.json", name, strings.ToLower(name))
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
