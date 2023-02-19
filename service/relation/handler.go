package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"toktik/constant/biz"
	"toktik/constant/config"
	relation "toktik/kitex_gen/douyin/relation"
	"toktik/kitex_gen/douyin/user"
	"toktik/kitex_gen/douyin/user/userservice"
	"toktik/logging"
	"toktik/repo"
	"toktik/repo/model"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var UserClient userservice.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	UserClient, err = userservice.NewClient(config.UserServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// GetFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowList(ctx context.Context, req *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":  req.UserId,
		"time":     time.Now(),
		"function": "GetFollowList",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	r := repo.Q.Relation
	relationModels, err := r.Where(r.UserId.Eq(req.UserId)).Find()
	if err != nil {
		resp = &relation.FollowListResponse{
			StatusCode: biz.UnableToQueryFollowList,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
			UserList:   nil,
		}
		return
	}
	var userList []*user.User
	for i, m := range relationModels {
		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  m.TargetId,
			ActorId: req.UserId,
		})
		if err != nil || userResponse.StatusCode != biz.OkStatusCode {
			log.Println(fmt.Errorf("failed to get user info: %w", err))
			continue
		}
		userList[i] = userResponse.User
	}

	logger.WithFields(map[string]interface{}{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return
}

// GetFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowerList(ctx context.Context, req *relation.FollowerListRequest) (resp *relation.FollowerListResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":  req.UserId,
		"time":     time.Now(),
		"function": "GetFollowerList",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	r := repo.Q.Relation
	relationModels, err := r.Where(r.TargetId.Eq(req.UserId)).Find()
	if err != nil {
		resp = &relation.FollowerListResponse{
			StatusCode: biz.UnableToQueryFollowList,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
			UserList:   nil,
		}
		return
	}
	var userList []*user.User
	for i, m := range relationModels {
		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  m.TargetId,
			ActorId: req.UserId,
		})
		if err != nil || userResponse.StatusCode != biz.OkStatusCode {
			logger.Error("failed to get user info: %w", err)
			continue
		}
		userList[i] = userResponse.User
	}

	logger.WithFields(map[string]interface{}{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return
}

// Follow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Follow(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":    req.UserId,
		"to_user_id": req.ToUserId,
		"time":       time.Now(),
		"function":   "RelationAction",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	relationModel := model.Relation{
		UserId:   uint32(req.UserId),
		TargetId: uint32(req.ToUserId),
	}

	// make sure target id exists
	userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
		UserId:  req.ToUserId,
		ActorId: req.UserId,
	})
	if err != nil || userResponse.StatusCode != biz.OkStatusCode || userResponse.User == nil {
		logger.Error("failed to get user info: %w", err)
		return
	}
	r := repo.Q.Relation

	// Follow
	err = r.WithContext(ctx).Create(&relationModel)
	if err != nil {
		// Follow a user that is followed already
		resp = &relation.RelationActionResponse{
			StatusCode: biz.RelationAlreadyExists,
			StatusMsg:  biz.BadRequestStatusMsg,
		}
		return
	}
	logger.WithFields(map[string]interface{}{
		"time":  time.Now(),
		"entry": relationModel,
	}).Debug("create relation")

	resp = &relation.RelationActionResponse{
		StatusCode: 0,
		StatusMsg:  biz.RelationActionSuccess,
	}

	logger.WithFields(map[string]interface{}{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return
}

// Unfollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Unfollow(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":    req.UserId,
		"to_user_id": req.ToUserId,
		"time":       time.Now(),
		"function":   "RelationAction",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	// TODO make sure target id exists

	r := repo.Q.Relation

	// Unfollow
	relationModel, err := r.WithContext(ctx).Where(r.UserId.Eq(req.UserId)).First()
	if err != nil {
		// Unfollow a user that was not followed before
		resp = &relation.RelationActionResponse{
			StatusCode: biz.RelationNotFound,
			StatusMsg:  biz.BadRequestStatusMsg,
		}
		return
	}

	if _, err = r.Delete(relationModel); err != nil {
		logger.WithFields(map[string]interface{}{
			"time":  time.Now(),
			"entry": relationModel,
		}).Debug("failed to delete")

		resp = &relation.RelationActionResponse{
			StatusCode: biz.UnableToDeleteRelation,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}

	logger.WithFields(map[string]interface{}{
		"time":  time.Now(),
		"entry": relationModel,
	}).Debug("deleted db entry")
	resp = &relation.RelationActionResponse{
		StatusCode: 0,
		StatusMsg:  biz.RelationActionSuccess,
	}

	logger.WithFields(map[string]interface{}{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return
}
