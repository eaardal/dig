package digfile

import (
	"fmt"
	"github.com/eaardal/dig/config"
	"github.com/eaardal/dig/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func GetPath() (string, error) {
	homeDir := os.Getenv(config.HomeEnvVar)
	if homeDir != "" {
		return filepath.Join(homeDir, config.DigfileFileName), nil
	}

	defaultConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	filePath := filepath.Join(defaultConfigDir, config.AppName, config.DigfileFileName)
	return filePath, nil
}

func Exists() bool {
	filePath, err := GetPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(filePath)
	return !os.IsNotExist(err)
}

func Read() (*Digfile, error) {
	filePath, err := GetPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get digfile path: %v", err)
	}

	if !Exists() {
		if err := Write(*DefaultDigfile); err != nil {
			return nil, err
		}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read digfile at %s: %v", filePath, err)
	}

	var digfile Digfile
	if err = yaml.Unmarshal(content, &digfile); err != nil {
		return nil, fmt.Errorf("failed to unmarshal digfile: %v", err)
	}

	return &digfile, nil
}

func Write(digfile Digfile) error {
	filePath, err := GetPath()
	if err != nil {
		return fmt.Errorf("failed to get digfile path: %v", err)
	}

	var file *os.File

	if !Exists() {
		err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not create directory: %w", err)
		}

		existingFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}

		file = existingFile
		defer utils.CloseOrPanic(existingFile.Close)
	}

	if file == nil {
		newFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}

		file = newFile
		defer utils.CloseOrPanic(newFile.Close)
	}

	fileAsYaml, err := yaml.Marshal(digfile)
	if err != nil {
		return fmt.Errorf("could not marshal digfile to YAML: %w", err)
	}

	_, err = file.Write(fileAsYaml)
	if err != nil {
		return fmt.Errorf("could not write YAML to file: %w", err)
	}

	return nil
}
