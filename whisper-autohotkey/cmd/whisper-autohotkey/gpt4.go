package main

import (
	"context"
	"math"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func BuildCommand(config Config, prompt string) (string, error) {

	if strings.TrimSpace(prompt) == "" {
		return `#Requires AutoHotkey v2.0
MsgBox("No input detected! Is your microphone working correctly?")`, nil
	}

	systemContext, err := os.ReadFile("./prompt.txt")
	if err != nil {
		return "", err
	}

	apiSettings := config.ResolveCommandSettings()
	if err := apiSettings.Validate("Command"); err != nil {
		return "", err
	}

	c := newOpenAIClient(apiSettings)

	response, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: apiSettings.Model,
			// https://github.com/sashabaranov/go-openai#why-dont-we-get-the-same-answer-when-specifying-a-temperature-field-of-0-and-asking-the-same-question
			Temperature: math.SmallestNonzeroFloat32,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: string(systemContext),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "ACTION: " + prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
