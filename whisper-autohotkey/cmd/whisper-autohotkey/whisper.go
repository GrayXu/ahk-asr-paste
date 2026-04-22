package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sashabaranov/go-openai"
)

func Transcribe(inputFileName string, config Config) (string, error) {

	prompt := config.ResolveTranscriptionPrompt()
	languageSymbol := "zh" // default language symbol

	argLength := len(os.Args[1:])
	if argLength >= 1 {
		languageSymbol = os.Args[1]
		log.Printf("Language Symbol: %s", languageSymbol)
	} else {
		log.Println("Insufficient arguments. Using default values.")
	}

	if argLength >= 2 {
		promptFileName, err := resolvePromptFileName(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		content, err := os.ReadFile(promptFileName)
		if err != nil {
			log.Fatal("Failed reading file: ", err)
		}
		prompt = string(content)
	}

	apiSettings := config.ResolveASRSettings()
	if err := apiSettings.Validate("ASR"); err != nil {
		return "", err
	}

	c := newOpenAIClient(apiSettings)
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    apiSettings.Model,
		Prompt:   prompt,
		Language: languageSymbol,
		FilePath: inputFileName,
	}
	response, err := c.CreateTranscription(ctx, req)
	if err != nil {
		return "", err
	}

	return response.Text, nil
}

func resolvePromptFileName(requestedName string) (string, error) {
	if _, err := os.Stat(requestedName); err == nil {
		return requestedName, nil
	}

	switch filepath.Base(requestedName) {
	case "transcriptionPrompt.txt":
		legacyName := filepath.Join(filepath.Dir(requestedName), "transtriptionPrompt.txt")
		if _, err := os.Stat(legacyName); err == nil {
			return legacyName, nil
		}
	case "transtriptionPrompt.txt":
		currentName := filepath.Join(filepath.Dir(requestedName), "transcriptionPrompt.txt")
		if _, err := os.Stat(currentName); err == nil {
			return currentName, nil
		}
	}

	return "", fmt.Errorf("prompt file %q does not exist", requestedName)
}
