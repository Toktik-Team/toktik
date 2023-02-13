package storage

import (
	"io"
	"net/url"
	"os"
	"path"
	"time"
	"toktik/constant/config"
	"toktik/logging"
)

type FSStorage struct {
}

func (f FSStorage) Upload(fileName string, content io.Reader) (output *PutObjectOutput, err error) {
	logger := logging.Logger.WithFields(map[string]interface{}{
		"time":      time.Now(),
		"function":  "FSStorage.Upload",
		"file_name": fileName,
	})
	logger.Debug("Process start")
	all, err := io.ReadAll(content)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"time": time.Now(),
			"err":  err,
		}).Debug("failed reading content")
		return nil, err
	}
	err = os.WriteFile(path.Join(config.EnvConfig.LOCAL_FS_LOCATION, fileName), all, 0644)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"time": time.Now(),
			"err":  err,
		}).Debug("failed writing content to file")
		return nil, err
	}
	return &PutObjectOutput{}, nil
}

func (f FSStorage) GetLink(fileName string) (string, error) {
	return url.JoinPath(config.EnvConfig.LOCAL_FS_BASEURL, fileName)
}
