package main

import "testing"

func TestResolveASRSettingsFallsBackToLegacyOpenAPIKey(t *testing.T) {
	config := Config{
		OpenapiKey: "legacy-key",
	}

	settings := config.ResolveASRSettings()

	if settings.APIKey != "legacy-key" {
		t.Fatalf("expected legacy api key fallback, got %q", settings.APIKey)
	}
	if settings.Model != defaultASRModel {
		t.Fatalf("expected default ASR model %q, got %q", defaultASRModel, settings.Model)
	}
}

func TestResolveASRSettingsPrefersSpecificOverrides(t *testing.T) {
	config := Config{
		API: APISettings{
			APIKey:  "shared-key",
			BaseURL: "https://shared.example.com/v1/",
			Model:   "shared-model",
		},
		ASR: APISettings{
			APIKey:  "asr-key",
			BaseURL: "https://asr.example.com/v1/",
			Model:   "asr-model",
		},
	}

	settings := config.ResolveASRSettings()

	if settings.APIKey != "asr-key" {
		t.Fatalf("expected ASR-specific api key, got %q", settings.APIKey)
	}
	if settings.BaseURL != "https://asr.example.com/v1" {
		t.Fatalf("expected normalized ASR base URL, got %q", settings.BaseURL)
	}
	if settings.Model != "asr-model" {
		t.Fatalf("expected ASR-specific model, got %q", settings.Model)
	}
}

func TestResolveCommandSettingsFallsBackToSharedDefaults(t *testing.T) {
	config := Config{
		API: APISettings{
			APIKey:  "shared-key",
			BaseURL: "https://shared.example.com/v1/",
			Model:   "shared-model",
		},
	}

	settings := config.ResolveCommandSettings()

	if settings.APIKey != "shared-key" {
		t.Fatalf("expected shared api key, got %q", settings.APIKey)
	}
	if settings.BaseURL != "https://shared.example.com/v1" {
		t.Fatalf("expected normalized shared base URL, got %q", settings.BaseURL)
	}
	if settings.Model != "shared-model" {
		t.Fatalf("expected shared model fallback, got %q", settings.Model)
	}
}

func TestResolveCommandSettingsUsesDefaultModelWhenUnset(t *testing.T) {
	config := Config{
		API: APISettings{
			APIKey: "shared-key",
		},
	}

	settings := config.ResolveCommandSettings()

	if settings.Model != defaultCommandModel {
		t.Fatalf("expected default command model %q, got %q", defaultCommandModel, settings.Model)
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
