package main

import (
	"context"
	"log"
	"os"

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
