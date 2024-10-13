package localstorage

import (
	"fmt"
	"github.com/eaardal/dig/config"
	"github.com/eaardal/dig/utils"
	"os"
	"path/filepath"
	"sync"
)

func getCacheDirPath(cacheId string) (string, error) {
	defaultConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	filePath := filepath.Join(defaultConfigDir, config.AppName, "cache", cacheId)
	return filePath, nil
}

func cacheDirExists(cacheId string) bool {
	filePath, err := getCacheDirPath(cacheId)
	if err != nil {
		return false
	}

	_, err = os.Stat(filePath)
	return !os.IsNotExist(err)
}

func createCacheDirIfNotExist(cacheId string) error {
	if cacheDirExists(cacheId) {
		return nil
	}

	cacheDirPath, err := getCacheDirPath(cacheId)
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
	cacheDirPath, err := getCacheDirPath(cacheId)
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
	cacheDirPath, err := getCacheDirPath(cacheId)
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

func readCacheFiles(cacheId string, sinkCh chan<- *CacheFile) error {
	cacheDirPath, err := getCacheDirPath(cacheId)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(cacheDirPath)
	if err != nil {
		return fmt.Errorf("failed to read cache directory at %s: %w", cacheDirPath, err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	for _, file := range files {
		wg.Add(1)

		go func(fileName string) {
			defer wg.Done()

			content, err := os.ReadFile(filepath.Join(cacheDirPath, fileName))
			if err != nil {
				errCh <- fmt.Errorf("failed to read cache file at %s: %w", fileName, err)
				return
			}

			sinkCh <- &CacheFile{
				FileName: fileName,
				Content:  content,
			}
		}(file.Name())
	}

	go func() {
		wg.Wait()
		close(doneCh)
		close(sinkCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-doneCh:
		return nil
	}
}
