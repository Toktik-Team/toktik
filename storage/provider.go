package storage

import (
	"github.com/sirupsen/logrus"
	"io"
	"toktik/constant/config"
	"toktik/logging"
)

var Instance storageProvider = S3Storage{}

func init() {
	if config.EnvConfig.STORAGE_TYPE == "fs" {
		Instance = FSStorage{}
	}
	logging.Logger.WithFields(logrus.Fields{
		"storage_type": config.EnvConfig.STORAGE_TYPE,
	}).Info("storage init")
}

type PutObjectOutput struct {
}

type storageProvider interface {
	Upload(fileName string, content io.Reader) (*PutObjectOutput, error)
	GetLink(fileName string) (string, error)
}

// Upload to the s3 storage using given fileName
func Upload(fileName string, content io.Reader) (*PutObjectOutput, error) {
	return Instance.Upload(fileName, content)
}

func GetLink(fileName string) (string, error) {
	return Instance.GetLink(fileName)
}
