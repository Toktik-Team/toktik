package storage

import (
	"io"
	"os"
	"path"
	"toktik/constant/config"
)

type FsStorage struct {
}

func (f FsStorage) Upload(fileName string, content io.Reader) (output *PutObjectOutput, err error) {
	all, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(path.Join(config.EnvConfig.LOCAL_FS_LOCATION, fileName), all, 0644)
	if err != nil {
		return nil, err
	}
	return &PutObjectOutput{}, nil
}

func (f FsStorage) GetLink(fileName string) (string, error) {
	return path.Join(config.EnvConfig.LOCAL_FS_BASEURL, fileName), nil
}
