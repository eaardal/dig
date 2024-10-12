package localstorage

import (
	"fmt"
	"github.com/eaardal/dig/config"
	"github.com/eaardal/dig/utils"
	"os"
	"path/filepath"
)

func CacheDirPath(cacheId string) (string, error) {
	defaultConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	filePath := filepath.Join(defaultConfigDir, config.AppName, "cache", cacheId)
	return filePath, nil
}

func CacheDirExists(cacheId string) bool {
	filePath, err := CacheDirPath(cacheId)
	if err != nil {
		return false
	}

	_, err = os.Stat(filePath)
	return !os.IsNotExist(err)
}

func CreateCacheDirIfNotExist(cacheId string) error {
	if CacheDirExists(cacheId) {
		return nil
	}

	cacheDirPath, err := CacheDirPath(cacheId)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cacheDirPath, 0700); err != nil {
		return fmt.Errorf("failed to create local storage directory at %s: %w", cacheDirPath, err)
	}

	return nil
}

// writeToCacheFile writes content to a file in the cache directory.
// If the file already exists, it will be overwritten.
// If the file does not exist, it will be created.
func writeToCacheFile(cacheId, fileName string, content []byte) error {
	cacheDirPath, err := CacheDirPath(cacheId)
	if err != nil {
		return err
	}

	filePath := filepath.Join(cacheDirPath, fileName)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write cache file at %s: %w", filePath, err)
	}

	return nil
}

// appendToCacheFile appends content to a file in the cache directory.
// If the file does not exist, it will be created.
// If the file exists, content will be appended to the end of the file.
func appendToCacheFile(cacheId, fileName string, content []byte) error {
	cacheDirPath, err := CacheDirPath(cacheId)
	if err != nil {
		return err
	}

	filePath := filepath.Join(cacheDirPath, fileName)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open cache file at %s: %w", filePath, err)
	}

	if _, err := file.Write(content); err != nil {
		utils.CloseOrPanic(file.Close)
		return fmt.Errorf("failed to append log to cache file: %w", err)
	}

	utils.CloseOrPanic(file.Close)
	return nil
}
