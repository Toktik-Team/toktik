package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"io"
	"log"
	"net/http"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/feed"
	feedService "toktik/kitex_gen/douyin/feed/feedservice"
	"toktik/logging"
	"toktik/repo/model"

	"github.com/sirupsen/logrus"

	"github.com/bakape/thumbnailer/v2"
	"github.com/gofrs/uuid"

	"image/jpeg"

	"toktik/kitex_gen/douyin/publish"
	gen "toktik/repo"
	"toktik/storage"
)

var FeedClient feedService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	FeedClient, err = feedService.NewClient(config.FeedServiceName, client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

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
		"user_id":  req.ActorId,
		"title":    req.Title,
		"function": "CreateVideo",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	detectedContentType := http.DetectContentType(req.Data)
	if detectedContentType != "video/mp4" {
		logger.WithFields(logrus.Fields{
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
			"err": err,
		}).Debug("error generating uuid")
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2GenerateUUID,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}
	logger = logger.WithFields(logrus.Fields{
		"uuid": generatedUUID,
	})
	logger.Debug("generated uuid")

	// Upload video file
	fileName := fmt.Sprintf("%d/%s.%s", req.ActorId, generatedUUID.String(), "mp4")
	_, err = storage.Upload(fileName, reader)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"file_name": fileName,
			"err":       err,
		}).Debug("failed to upload video")
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2UploadVideo,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}
	logger.WithFields(logrus.Fields{
		"file_name": fileName,
	}).Debug("uploaded video")

	// Generate thumbnail
	coverName := fmt.Sprintf("%d/%s.%s", req.ActorId, generatedUUID.String(), "jpg")
	thumbData, err := getThumbnail(reader)
	if err != nil {
		logger.WithFields(logrus.Fields{
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
		"cover_name": coverName,
		"data_size":  len(thumbData),
	}).Debug("generated thumbnail")

	// Upload thumbnail
	_, err = storage.Upload(coverName, bytes.NewReader(thumbData))
	if err != nil {
		logger.WithFields(logrus.Fields{
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
		"cover_name": coverName,
		"data_size":  len(thumbData),
	}).Debug("uploaded thumbnail")

	publishModel := model.Video{
		UserId:    req.ActorId,
		FileName:  fileName,
		CoverName: coverName,
		Title:     req.Title,
	}

	err = gen.Q.Video.WithContext(ctx).Create(&publishModel)
	if err != nil {
		logger.WithFields(logrus.Fields{
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
		"entry": publishModel,
	}).Debug("saved db entry")

	resp = &publish.CreateVideoResponse{StatusCode: 0, StatusMsg: biz.PublishActionSuccess}
	logger.WithFields(logrus.Fields{
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return resp, nil
}

// ListVideo implements the PublishServiceImpl interface.
func (s *PublishServiceImpl) ListVideo(ctx context.Context, req *publish.ListVideoRequest) (resp *publish.ListVideoResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":  req.UserId,
		"actor_id": req.ActorId,
		"function": "ListVideo",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	videos, err := gen.Q.Video.WithContext(ctx).
		Where(gen.Q.Video.UserId.Eq(req.UserId)).
		Order(gen.Q.Video.CreatedAt.Desc()).
		Find()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Debug("failed to query video")
		return &publish.ListVideoResponse{
			StatusCode: biz.UnableToQueryVideo,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}

	videoIds := make([]uint32, 0, len(videos))
	for _, video := range videos {
		videoIds = append(videoIds, video.ID)
	}

	queryVideoResp, err := FeedClient.QueryVideos(ctx, &feed.QueryVideosRequest{
		ActorId:  req.ActorId,
		VideoIds: videoIds,
	})

	logger.WithFields(logrus.Fields{
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return &publish.ListVideoResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		VideoList:  queryVideoResp.VideoList,
	}, nil
}

// CountVideo implements the PublishServiceImpl interface.
func (s *PublishServiceImpl) CountVideo(ctx context.Context, req *publish.CountVideoRequest) (resp *publish.CountVideoResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":  req.UserId,
		"function": "CountVideo",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	count, err := gen.Q.Video.WithContext(ctx).Where(gen.Q.Video.UserId.Eq(req.UserId)).Count()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Debug("failed to query video")
		return &publish.CountVideoResponse{
			StatusCode: biz.UnableToQueryVideo,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}

	return &publish.CountVideoResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		Count:      uint32(count),
	}, nil
}
