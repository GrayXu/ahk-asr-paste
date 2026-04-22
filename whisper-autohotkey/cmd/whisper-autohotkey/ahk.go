package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunCommand(config Config, script string) (string, error) {

	if err := os.WriteFile("script.ahk", []byte(script), 0666); err != nil {
		return "", err
	}

	autoHotKeyPath, err := resolveAutoHotKeyExec(config)
	if err != nil {
		return "", err
	}
	data, err := exec.Command(autoHotKeyPath, "script.ahk").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run AutoHotkey with %q: %w: %s", autoHotKeyPath, err, string(data))
	}

	output := string(data)

	return output, nil
}

func resolveAutoHotKeyExec(config Config) (string, error) {
	if path := strings.TrimSpace(config.AutoHotKeyExec); path != "" {
		if err := assertAutoHotKeyV2(path); err != nil {
			return "", fmt.Errorf("AutoHotKeyExec %q is not a usable AutoHotkey v2 executable: %w", path, err)
		}
		return path, nil
	}

	for _, name := range []string{"AutoHotkey64.exe", "AutoHotkey32.exe", "AutoHotkey.exe"} {
		path, err := exec.LookPath(name)
		if err == nil && assertAutoHotKeyV2(path) == nil {
			return path, nil
		}
	}

	for _, baseDir := range []string{os.Getenv("ProgramFiles"), os.Getenv("ProgramFiles(x86)")} {
		if strings.TrimSpace(baseDir) == "" {
			continue
		}

		for _, suffix := range []string{
			filepath.Join("AutoHotkey", "v2", "AutoHotkey64.exe"),
			filepath.Join("AutoHotkey", "v2", "AutoHotkey32.exe"),
			filepath.Join("AutoHotkey", "AutoHotkey64.exe"),
			filepath.Join("AutoHotkey", "AutoHotkey32.exe"),
			filepath.Join("AutoHotkey", "AutoHotkey.exe"),
		} {
			candidate := filepath.Join(baseDir, suffix)
			if _, err := os.Stat(candidate); err == nil && assertAutoHotKeyV2(candidate) == nil {
				return candidate, nil
			}
		}
	}

	return "", fmt.Errorf("AutoHotkey v2 executable not found; install AutoHotkey v2 or set AutoHotKeyExec in config.json")
}

func assertAutoHotKeyV2(autoHotKeyPath string) error {
	probeFile, err := os.CreateTemp("", "ahk-v2-probe-*.ahk")
	if err != nil {
		return err
	}
	probePath := probeFile.Name()
	defer os.Remove(probePath)

	probeScript := "#Requires AutoHotkey v2.0\nExitApp(0)\n"
	if _, err := probeFile.WriteString(probeScript); err != nil {
		probeFile.Close()
		return err
	}
	if err := probeFile.Close(); err != nil {
		return err
	}

	data, err := exec.Command(autoHotKeyPath, probePath).CombinedOutput()
	if err != nil {
		output := strings.TrimSpace(string(data))
		if output == "" {
			return err
		}
		return fmt.Errorf("%w: %s", err, output)
	}

	return nil
}
