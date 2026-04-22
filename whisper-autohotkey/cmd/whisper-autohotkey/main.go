package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"unicode/utf16"
)

func writeTextToClipboard(text string) error {
	// Create a temporary file
	tmpFile, err := ioutil.TempFile("", "clip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after we're done

	utf16leBytes := utf16leEncode(text)
	// Omit the BOM as it may cause issues with some programs
	if len(utf16leBytes) >= 2 {
		utf16leBytes = utf16leBytes[2:]
	}
	if _, err = tmpFile.Write(utf16leBytes); err != nil {
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}

	// Use clip to copy the contents of the file to the clipboard
	clipCmd := exec.Command("cmd", "/c", fmt.Sprintf("type %s | clip", tmpFile.Name()))
	return clipCmd.Run()
}

// utf16leEncode encodes a string in UTF-16LE with a Byte Order Mark (BOM)
func utf16leEncode(s string) []byte {
	// Encode the string as UTF-16LE
	encoded := utf16.Encode([]rune(s))

	// Prepend the BOM
	bom := []uint16{0xFEFF}
	encoded = append(bom, encoded...)

	// Convert uint16 slice to byte slice
	b := make([]byte, 2*len(encoded))
	for i, runeValue := range encoded {
		b[2*i] = byte(runeValue)
		b[2*i+1] = byte(runeValue >> 8)
	}
	return b
}

func main() {

	// Open a file for logging
	logFile, e := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		panic(e)
	}
	defer logFile.Close()

	// Set the log output to the file
	log.SetOutput(logFile)
	log.Println("")
	log.Println("")
	log.Println("========================================")
	log.Println("Starting whisper-autohotkey")
	err := assertThatConfigFileExists()
	if err != nil {
		log.Fatal("Error when creating config file: ", err)
	}

	content, err := readConfigFile()
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during JSON parse: ", err)
	}

	// print config to log
	asrSettings := config.ResolveASRSettings()
	log.Println("Config:")
	log.Println("  ASR effective key configured: " + fmt.Sprintf("%t", asrSettings.APIKey != ""))
	log.Println("  ASR effective baseURL: " + asrSettings.BaseURL)
	log.Println("  ASR effective model: " + asrSettings.Model)
	log.Println("  AutoHotKeyExec: " + config.AutoHotKeyExec)

	// argLength := len(os.Args[1:])
	inputFileName := "rec.mp3"

	// if argLength > 1 {
	// 	inputFileName = os.Args[1:][1]
	// 	log.Println("Processing file " + inputFileName)
	// 	stats, err := os.Stat(inputFileName)
	// 	if errors.Is(err, os.ErrNotExist) {
	// 		log.Fatal("Input file does not exist")
	// 	} else {
	// 		log.Printf("File size %v", stats.Size())
	// 	}
	// }

	text, err := Transcribe(inputFileName, config)

	if err != nil {
		log.Fatal("Cannot transcribe text: ", err)
		return
	}
	log.Println("Transcription:\n" + text)

	log.Println("Ready to paste:\n" + text)
	err = writeTextToClipboard(text)
	if err != nil {
		log.Fatal("Failed to write text to clipboard:", err)
	}

	log.Println("Text copied to clipboard")

	ahkScript := `#Requires AutoHotkey v2.0
Send("^v")
ExitApp()
`

	// Assuming you have AutoHotKey installed and `paste.ahk` is in the same directory.
	_, err = RunCommand(config, ahkScript)
	if err != nil {
		log.Fatal("Cannot run AutoHotKey command", err)
	}
}

func readConfigFile() ([]byte, error) {
	content, err := os.ReadFile("./config.json")
	return content, err
}

func assertThatConfigFileExists() error {
	if !exists("./config.json") {
		template, err := os.ReadFile("./config.template.json")
		if err != nil {
			return fmt.Errorf("cannot read template config file: %w", err)
		}
		err = os.WriteFile("./config.json", template, 0666)
		if err != nil {
			return fmt.Errorf("cannot write new config file: %w", err)
		}
		return nil
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
