package localstorage

import "fmt"

type CacheFile struct {
	FileName    string
	FileContent []byte
}

func (c CacheFile) Name() string {
	return c.FileName
}

func (c CacheFile) Content() []byte {
	return c.FileContent
}

func SaveFileToCache(cacheId string, sourceCh <-chan *CacheFile) error {
	if err := createCacheDirIfNotExist(cacheId); err != nil {
		return err
	}

	isFirstChunk := true

	for log := range sourceCh {
		if isFirstChunk {
			if err := writeToCacheFile(cacheId, log.FileName, log.FileContent); err != nil {
				return err
			}
			isFirstChunk = false
		} else {
			if err := appendToCacheFile(cacheId, log.FileName, log.FileContent); err != nil {
				return err
			}
		}
	}

	return nil
}

func GetCachedFiles(cacheId string) ([]*CacheFile, error) {
	if !cacheDirExists(cacheId) {
		return nil, fmt.Errorf("cache directory for cache ID %s does not exist", cacheId)
	}

	filesCh := make(chan *CacheFile)

	go func() {
		defer close(filesCh)

		if err := readCacheFiles(cacheId, filesCh); err != nil {
			return
		}
	}()

	files := make([]*CacheFile, 0)
	for file := range filesCh {
		files = append(files, file)
	}

	return files, nil
}

func StreamCacheFiles(cacheId string, sinkCh chan<- *CacheFile) error {
	if !cacheDirExists(cacheId) {
		return fmt.Errorf("cache directory for cache ID %s does not exist", cacheId)
	}

	return readCacheFiles(cacheId, sinkCh)
}
