package main

import "testing"

func TestResolveASRSettingsUsesDirectSettings(t *testing.T) {
	config := Config{
		ASR: APISettings{
			APIKey:  "asr-key",
			BaseURL: "https://asr.example.com/v1/",
		},
	}

	settings := config.ResolveASRSettings()

	if settings.APIKey != "asr-key" {
		t.Fatalf("expected ASR api key, got %q", settings.APIKey)
	}
	if settings.BaseURL != "https://asr.example.com/v1" {
		t.Fatalf("expected normalized ASR base URL, got %q", settings.BaseURL)
	}
	if settings.Model != defaultASRModel {
		t.Fatalf("expected default ASR model %q, got %q", defaultASRModel, settings.Model)
	}
}

func TestResolveASRSettingsUsesConfiguredModel(t *testing.T) {
	config := Config{
		ASR: APISettings{
			APIKey:  "asr-key",
			BaseURL: "https://asr.example.com/v1/",
			Model:   "asr-model",
		},
	}

	settings := config.ResolveASRSettings()

	if settings.APIKey != "asr-key" {
		t.Fatalf("expected ASR api key, got %q", settings.APIKey)
	}
	if settings.BaseURL != "https://asr.example.com/v1" {
		t.Fatalf("expected normalized ASR base URL, got %q", settings.BaseURL)
	}
	if settings.Model != "asr-model" {
		t.Fatalf("expected ASR model, got %q", settings.Model)
	}
}

func TestResolvedAPISettingsValidateRequiresAPIKey(t *testing.T) {
	settings := ResolvedAPISettings{
		Model: defaultASRModel,
	}

	if err := settings.Validate("ASR"); err == nil {
		t.Fatal("expected validation error when api key is empty")
	}
}

func TestResolveTranscriptionPromptUsesConfiguredValue(t *testing.T) {
	config := Config{
		TranscriptionPrompt: "custom prompt",
	}

	prompt := config.ResolveTranscriptionPrompt()

	if prompt != "custom prompt" {
		t.Fatalf("expected configured prompt, got %q", prompt)
	}
}

func TestResolveTranscriptionPromptFallsBackToDefault(t *testing.T) {
	config := Config{}

	prompt := config.ResolveTranscriptionPrompt()

	if prompt != defaultTranscriptionPrompt {
		t.Fatalf("expected default prompt %q, got %q", defaultTranscriptionPrompt, prompt)
	}
}
