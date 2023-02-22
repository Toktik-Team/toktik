package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/relation"
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
func removeDups(elements []uint32) (nodups []uint32) {
	encountered := make(map[uint32]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}

func intersection(s1, s2 []uint32) (inter []uint32) {
	hash := make(map[uint32]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = removeDups(inter)
	return
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
	for _, m := range relationModels {
		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  m.TargetId,
			ActorId: req.UserId,
		})
		if err != nil || userResponse.StatusCode != biz.OkStatusCode {
			log.Println(fmt.Errorf("failed to get user info: %w", err))
			continue
		}
		userList = append(userList, userResponse.User)
	}

	logger.WithFields(map[string]interface{}{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")
	resp = &relation.FollowListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		UserList:   userList,
	}
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
			StatusCode: biz.UnableToQueryFollowerList,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
			UserList:   nil,
		}
		return
	}
	var userList []*user.User
	for _, m := range relationModels {
		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  m.UserId,
			ActorId: req.UserId,
		})
		if err != nil || userResponse.StatusCode != biz.OkStatusCode {
			logger.Error("failed to get user info: %w", err)
			continue
		}
		userList = append(userList, userResponse.User)
	}

	logger.WithFields(map[string]interface{}{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")
	resp = &relation.FollowerListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		UserList:   userList,
	}
	return
}

// GetFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFriendList(ctx context.Context, req *relation.FriendListRequest) (resp *relation.FriendListResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":  req.UserId,
		"time":     time.Now(),
		"function": "GetFriendList",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	r := repo.Q.Relation

	// Get all people followed current user
	followerModels, err := r.Where(r.TargetId.Eq(req.UserId)).Find()
	if err != nil {
		resp = &relation.FriendListResponse{
			StatusCode: biz.UnableToQueryFollowerList,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
			UserList:   nil,
		}
		return
	}

	// Get all people current user following to
	followModels, err := r.Where(r.UserId.Eq(req.UserId)).Find()
	if err != nil {
		resp = &relation.FriendListResponse{
			StatusCode: biz.UnableToQueryFollowList,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
			UserList:   nil,
		}
		return
	}

	var followerIds []uint32
	for _, followerModel := range followerModels {
		followerIds = append(followerIds, followerModel.UserId)
	}
	var followIds []uint32
	for _, followModel := range followModels {
		followIds = append(followIds, followModel.TargetId)
	}

	var ids = intersection(followerIds, followIds)

	var userList = []*user.User{biz.ChatGPTUser}
	for _, m := range ids {
		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  m,
			ActorId: req.UserId,
		})
		if err != nil || userResponse.StatusCode != biz.OkStatusCode {
			logger.Error("failed to get user info: %w", err)
			continue
		}
		userList = append(userList, userResponse.User)
	}

	resp = &relation.FriendListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		UserList:   userList,
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

	if req.ToUserId == req.UserId {
		resp = &relation.RelationActionResponse{
			StatusCode: biz.InvalidToUserId,
			StatusMsg:  biz.BadRequestStatusMsg,
		}
		return
	}

	relationModel := model.Relation{
		UserId:   req.UserId,
		TargetId: req.ToUserId,
	}

	// make sure target id exists
	userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
		UserId:  req.ToUserId,
		ActorId: req.UserId,
	})
	if err != nil || userResponse.StatusCode != biz.OkStatusCode || userResponse.User == nil {
		logger.Error("failed to get user info: %w", err)
		resp = &relation.RelationActionResponse{
			StatusCode: biz.UserNotFound,
			StatusMsg:  biz.BadRequestStatusMsg,
		}
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
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
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
	if req.ToUserId == req.UserId {
		resp = &relation.RelationActionResponse{
			StatusCode: biz.InvalidToUserId,
			StatusMsg:  biz.BadRequestStatusMsg,
		}
		return
	}

	r := repo.Q.Relation

	relationModel, err := r.WithContext(ctx).Where(r.UserId.Eq(req.UserId), r.TargetId.Eq(req.ToUserId)).First()
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"user_id":    req.UserId,
			"to_user_id": req.ToUserId,
			"time":       time.Now(),
		}).Debug("record not found")
		// Unfollow a user that was not followed before
		resp = &relation.RelationActionResponse{
			StatusCode: biz.RelationNotFound,
			StatusMsg:  biz.BadRequestStatusMsg,
		}
		return
	}

	if _, err = r.Where(r.UserId.Eq(req.UserId), r.TargetId.Eq(req.ToUserId)).Delete(); err != nil {
		logger.WithFields(map[string]interface{}{
			"time":  time.Now(),
			"entry": relationModel,
			"err":   err,
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

	logger.WithFields(map[string]interface{}{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")
	resp = &relation.RelationActionResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
	}
	return
}

// CountFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) CountFollowList(ctx context.Context, req *relation.CountFollowListRequest) (resp *relation.CountFollowListResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":  req.UserId,
		"time":     time.Now(),
		"function": "CountFollowList",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	r := repo.Q.Relation
	count, err := r.WithContext(ctx).Where(r.UserId.Eq(req.UserId)).Count()
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"user_id": req.UserId,
			"time":    time.Now(),
			"err":     err,
		}).Debug("failed to count follow list")

		resp = &relation.CountFollowListResponse{
			StatusCode: biz.UnableToQueryFollowList,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}

	return &relation.CountFollowListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		Count:      uint32(count),
	}, nil
}

// CountFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) CountFollowerList(ctx context.Context, req *relation.CountFollowerListRequest) (resp *relation.CountFollowerListResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":  req.UserId,
		"time":     time.Now(),
		"function": "CountFollowerList",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	r := repo.Q.Relation
	count, err := r.WithContext(ctx).Where(r.TargetId.Eq(req.UserId)).Count()
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"user_id": req.UserId,
			"time":    time.Now(),
			"err":     err,
		}).Debug("failed to count follower list")

		resp = &relation.CountFollowerListResponse{
			StatusCode: biz.UnableToQueryFollowerList,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}

	return &relation.CountFollowerListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		Count:      uint32(count),
	}, nil
}

// IsFollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) IsFollow(ctx context.Context, req *relation.IsFollowRequest) (resp *relation.IsFollowResponse, err error) {
	methodFields := logrus.Fields{
		"user_id":  req.UserId,
		"time":     time.Now(),
		"function": "IsFollow",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	r := repo.Q.Relation
	count, err := r.WithContext(ctx).Where(r.UserId.Eq(req.ActorId), r.TargetId.Eq(req.UserId)).Count()
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"user_id": req.UserId,
			"time":    time.Now(),
			"err":     err,
		}).Debug("failed to count follower list")

		resp = &relation.IsFollowResponse{
			Result: false,
		}
		return
	}

	return &relation.IsFollowResponse{
		Result: count > 0,
	}, nil

}
