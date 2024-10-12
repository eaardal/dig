package localstorage

type CacheFile struct {
	FileName string
	Content  []byte
}

func SaveFileToCache(cacheId string, sourceCh <-chan *CacheFile) error {
	if err := CreateCacheDirIfNotExist(cacheId); err != nil {
		return err
	}

	isFirstChunk := true

	for log := range sourceCh {
		if isFirstChunk {
			if err := writeToCacheFile(cacheId, log.FileName, log.Content); err != nil {
				return err
			}
			isFirstChunk = false
		} else {
			if err := appendToCacheFile(cacheId, log.FileName, log.Content); err != nil {
				return err
			}
		}
	}

	return nil
}
