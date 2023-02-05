package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/bakape/thumbnailer/v2"
	"github.com/gofrs/uuid"

	"image/jpeg"
	"os"

	publish "toktik/kitex_gen/douyin/publish"
	gen "toktik/repo"
	"toktik/service/publish/model"
	"toktik/service/publish/storage"
)

// PublishServiceImpl implements the last service interface defined in the IDL.
type PublishServiceImpl struct{}

// CreateVideo implements the PublishServiceImpl interface.
func (s *PublishServiceImpl) CreateVideo(ctx context.Context, req *publish.CreateVideoRequest) (resp *publish.CreateVideoResponse, err error) {
	// byte[] -> reader
	reader := bytes.NewReader(req.Data)
	// V7 based on timestamp
	uid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	// Handle video file
	fileName := fmt.Sprintf("%s.%s", uid.String(), "mp4")
	_, err = storage.Upload(fileName, reader)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to upload video %s\n", fileName)
		return nil, err
	}
	playURL, err := storage.GetLink(fileName)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get link for video %s\n", fileName)
		return nil, err
	}

	// Handle cover file
	_, thumb, err := thumbnailer.Process(reader, thumbnailer.Options{})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create thumbnail %s\n", fileName)
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, thumb, nil)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create buffer %s\n", fileName)
		return nil, err
	}
	reader = bytes.NewReader(buf.Bytes())
	coverName := fmt.Sprintf("%s.%s", uid.String(), "jpg")
	_, err = storage.Upload(coverName, reader)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to upload cover %s\n", fileName)
		return nil, err
	}
	coverURL, err := storage.GetLink(coverName)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to get link for cover %s\n", fileName)
		return nil, err
	}
	publishModel := model.Publish{
		UserId:   req.UserId,
		PlayUrl:  playURL,
		CoverUrl: coverURL,
		Title:    req.Title,
	}
	err = gen.Q.Publish.WithContext(ctx).Create(&publishModel)
	if err != nil {
		return nil, err
	}

	return &publish.CreateVideoResponse{Id: int64(publishModel.ID)}, nil
}
