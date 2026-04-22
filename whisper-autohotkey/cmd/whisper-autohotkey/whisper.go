package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sashabaranov/go-openai"
)

func Transcribe(inputFileName string, config Config) (string, error) {

	// Define the prompt variable and the language symbol
	prompt := "This is a transcription in English, mainly about programming, coding and software development."
	promptFileName, err := resolvePromptFileName("transcriptionPrompt.txt")
	if err != nil {
		log.Fatal(err)
	}
	languageSymbol := "en" // default language symbol

	argLength := len(os.Args[1:])
	if argLength >= 2 { // Check if at least two arguments are provided
		languageSymbol = os.Args[1]
		promptFileName, err = resolvePromptFileName(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		// log.Println("Processing file " + promptFileName)
		stats, err := os.Stat(promptFileName)
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("Prompt file does not exist")
		} else {
			log.Printf("File size %v", stats.Size())
		}
		log.Printf("Language Symbol: %s", languageSymbol)
	} else {
		log.Println("Insufficient arguments. Using default values.")
	}

	content, err := os.ReadFile(promptFileName)
	if err != nil {
		log.Fatal("Failed reading file: ", err)
	}
	prompt = string(content)

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
