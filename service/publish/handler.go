package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
	"toktik/constant/biz"
	"toktik/logging"
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
	methodFields := logrus.Fields{
		"user_id":  req.UserId,
		"title":    req.Title,
		"time":     time.Now(),
		"function": "CreateVideo",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	detectedContentType := http.DetectContentType(req.Data)
	if detectedContentType != "video/mp4" {
		logger.WithFields(logrus.Fields{
			"time":         time.Now(),
			"content_type": detectedContentType,
		}).Debug("invalid content type")
		return &publish.CreateVideoResponse{
			StatusCode: biz.InvalidContentType,
			StatusMsg:  biz.BadRequestStatusMsg,
		}, nil
	}
	// byte[] -> reader
	reader := bytes.NewReader(req.Data)

	// V7 based on timestamp
	generatedUUID, err := uuid.NewV7()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
			"err":  err,
		}).Debug("error generating uuid")
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2GenerateUUID,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}
	logger = logger.WithFields(logrus.Fields{
		"uuid": generatedUUID,
	})
	logger.WithFields(logrus.Fields{
		"time": time.Now(),
	}).Debug("generated uuid")

	// Upload video file
	fileName := fmt.Sprintf("%d/%s.%s", req.UserId, generatedUUID.String(), "mp4")
	_, err = storage.Upload(fileName, reader)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time":      time.Now(),
			"file_name": fileName,
			"err":       err,
		}).Debug("failed to upload video")
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2UploadVideo,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}
	logger.WithFields(logrus.Fields{
		"time":      time.Now(),
		"file_name": fileName,
	}).Debug("uploaded video")

	// Generate thumbnail
	coverName := fmt.Sprintf("%d/%s.%s", req.UserId, generatedUUID.String(), "jpg")
	thumbData, err := getThumbnail(reader)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time":       time.Now(),
			"file_name":  fileName,
			"cover_name": coverName,
			"err":        err,
		}).Debug("failed to create thumbnail")
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2CreateThumbnail,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}
	logger.WithFields(logrus.Fields{
		"time":       time.Now(),
		"cover_name": coverName,
		"data_size":  len(thumbData),
	}).Debug("generated thumbnail")

	// Upload thumbnail
	_, err = storage.Upload(coverName, bytes.NewReader(thumbData))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time":       time.Now(),
			"file_name":  fileName,
			"cover_name": coverName,
			"err":        err,
		}).Debug("failed to upload cover")
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2UploadCover,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}
	logger.WithFields(logrus.Fields{
		"time":       time.Now(),
		"cover_name": coverName,
		"data_size":  len(thumbData),
	}).Debug("uploaded thumbnail")

	publishModel := model.Video{
		UserId:    req.UserId,
		FileName:  fileName,
		CoverName: coverName,
		Title:     req.Title,
	}

	err = gen.Q.Video.WithContext(ctx).Create(&publishModel)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time":       time.Now(),
			"file_name":  fileName,
			"cover_name": coverName,
			"err":        err,
		}).Debug("failed to create db entry")
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2CreateDBEntry,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}
	logger.WithFields(logrus.Fields{
		"time":  time.Now(),
		"entry": publishModel,
	}).Debug("saved db entry")

	resp = &publish.CreateVideoResponse{StatusCode: 0, StatusMsg: biz.PublishActionSuccess}
	logger.WithFields(logrus.Fields{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return resp, nil
}
