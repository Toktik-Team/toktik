package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"toktik/constant/biz"
	"toktik/repo/model"

	"github.com/bakape/thumbnailer/v2"
	"github.com/gofrs/uuid"

	"image/jpeg"

	"toktik/kitex_gen/douyin/publish"
	gen "toktik/repo"
	"toktik/storage"
)

// getThumbnail Generate JPEG thumbnail from video
func getThumbnail(input io.ReadSeeker) ([]byte, error) {
	_, thumb, err := thumbnailer.Process(input, thumbnailer.Options{})
	if err != nil {
		return nil, errors.New("failed to create thumbnail")
	}
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, thumb, nil)
	if err != nil {
		return nil, errors.New("failed to create buffer")
	}
	return buf.Bytes(), nil
}

// PublishServiceImpl implements the last service interface defined in the IDL.
type PublishServiceImpl struct{}

// CreateVideo implements the PublishServiceImpl interface.
func (s *PublishServiceImpl) CreateVideo(ctx context.Context, req *publish.CreateVideoRequest) (resp *publish.CreateVideoResponse, err error) {
	if http.DetectContentType(req.Data) != "video/mp4" {
		_ = errors.New("invalid content type")
		return &publish.CreateVideoResponse{
			StatusCode: biz.InvalidContentType,
			StatusMsg:  biz.BadRequestStatusMsg,
		}, nil
	}
	// byte[] -> reader
	reader := bytes.NewReader(req.Data)

	// V7 based on timestamp
	uid, err := uuid.NewV7()
	if err != nil {
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2GenerateUUID,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}

	// Upload video file
	fileName := fmt.Sprintf("%d/%s.%s", req.UserId, uid.String(), "mp4")
	_, err = storage.Upload(fileName, reader)
	if err != nil {
		_ = fmt.Errorf("failed to upload video %s: %w", fileName, err)
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2UploadVideo,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}

	// Generate thumbnail
	coverName := fmt.Sprintf("%d/%s.%s", req.UserId, uid.String(), "jpg")
	thumbData, err := getThumbnail(reader)
	if err != nil {
		_ = fmt.Errorf("failed to create thumbnail %s: %w", fileName, err)
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2CreateThumbnail,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}

	// Upload thumbnail
	_, err = storage.Upload(coverName, bytes.NewReader(thumbData))
	if err != nil {
		_ = fmt.Errorf("failed to upload cover %s: %w", fileName, err)
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2UploadCover,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}

	publishModel := model.Video{
		UserId:    req.UserId,
		FileName:  fileName,
		CoverName: coverName,
		Title:     req.Title,
	}

	err = gen.Q.Video.WithContext(ctx).Create(&publishModel)
	if err != nil {
		_ = fmt.Errorf("failed to create db entry %s: %w", fileName, err)
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2CreateDBEntry,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}

	return &publish.CreateVideoResponse{StatusCode: 0}, nil
}
