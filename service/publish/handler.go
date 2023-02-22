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
	"toktik/kitex_gen/douyin/comment"
	"toktik/kitex_gen/douyin/comment/commentservice"
	"toktik/kitex_gen/douyin/favorite"
	favoriteService "toktik/kitex_gen/douyin/favorite/favoriteservice"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/user"
	"toktik/kitex_gen/douyin/user/userservice"
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

var UserClient userservice.Client
var CommentClient commentservice.Client
var FavoriteClient favoriteService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	UserClient, err = userservice.NewClient(config.UserServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
	CommentClient, err = commentservice.NewClient(config.CommentServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
	FavoriteClient, err = favoriteService.NewClient(config.FavoriteServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
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
		"user_id":  req.UserId,
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
	fileName := fmt.Sprintf("%d/%s.%s", req.UserId, generatedUUID.String(), "mp4")
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
	coverName := fmt.Sprintf("%d/%s.%s", req.UserId, generatedUUID.String(), "jpg")
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
		UserId:    req.UserId,
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

	u := gen.Q.User
	_, err = u.WithContext(ctx).Where(u.ID.Eq(req.UserId)).Update(u.WorkCount, u.WorkCount.Add(1))
	if err != nil {
		logger.Debug("failed to update the number of user works")
		return &publish.CreateVideoResponse{
			StatusCode: biz.Unable2CreateDBEntry,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, err
	}

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

	find, err := gen.Q.Video.WithContext(ctx).
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

	rVideo := make([]*feed.Video, 0, len(find))
	for _, m := range find {

		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  m.UserId,
			ActorId: req.ActorId,
		})
		if err != nil || userResponse.StatusCode != biz.OkStatusCode {
			log.Println(fmt.Errorf("failed to get user info: %w", err))
			continue
		}

		playUrl, err := storage.GetLink(m.FileName)
		if err != nil {
			log.Println(fmt.Errorf("failed to fetch play url: %w", err))
			continue
		}

		coverUrl, err := storage.GetLink(m.CoverName)
		if err != nil {
			log.Println(fmt.Errorf("failed to fetch cover url: %w", err))
			continue
		}

		commentCount, err := CommentClient.CountComment(ctx, &comment.CountCommentRequest{
			ActorId: req.ActorId,
			VideoId: m.ID,
		})
		if err != nil {
			log.Println(fmt.Errorf("failed to fetch comment count: %w", err))
			continue
		}

		isFavorite, err := FavoriteClient.IsFavorite(ctx, &favorite.IsFavoriteRequest{
			UserId:  req.UserId,
			VideoId: m.ID,
		})
		if err != nil {
			log.Println(fmt.Errorf("unable to determine if the user liked the video : %w", err))
			continue
		}

		rVideo = append(rVideo, &feed.Video{
			Id:            m.ID,
			Author:        userResponse.User,
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: m.FavoriteCount,
			CommentCount:  commentCount.CommentCount,
			IsFavorite:    isFavorite.Result,
			Title:         m.Title,
		})
	}

	logger.WithFields(logrus.Fields{
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return &publish.ListVideoResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		VideoList:  rVideo,
	}, nil
}
