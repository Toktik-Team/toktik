package main

import (
	"context"
	"log"
	"sync"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/relation"
	"toktik/kitex_gen/douyin/user"
	"toktik/kitex_gen/douyin/user/userservice"
	"toktik/logging"
	"toktik/repo"
	"toktik/repo/model"

	"github.com/kitex-contrib/obs-opentelemetry/tracing"

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
	UserClient, err = userservice.NewClient(
		config.UserServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		log.Fatal(err)
	}
}
func removeDups(elements []*uint32) (nodups []*uint32) {
	encountered := make(map[uint32]bool)
	for _, element := range elements {
		if !encountered[*element] {
			nodups = append(nodups, element)
			encountered[*element] = true
		}
	}
	return
}

func intersection(s1, s2 []*uint32) (inter []*uint32) {
	hash := make(map[uint32]bool)
	for _, e := range s1 {
		hash[*e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[*e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = removeDups(inter)
	return
}

func queryUsers(
	ctx context.Context,
	logger *logrus.Entry,
	actorId uint32,
	userIds []*uint32,
) (respUserList []*user.User) {
	var wg sync.WaitGroup
	wg.Add(len(userIds))
	respUserList = make([]*user.User, len(userIds))
	for i, v := range userIds {
		go func(i int, v *uint32) {
			defer wg.Done()
			userResponse, localErr := UserClient.GetUser(ctx, &user.UserRequest{
				UserId:  *v,
				ActorId: actorId,
			})
			if localErr != nil || userResponse.StatusCode != biz.OkStatusCode {
				logger.WithFields(logrus.Fields{
					"actor_id": actorId,
					"user_id":  *v,
					"cause":    localErr,
				}).Warning("failed to get user info")
				return
			}
			respUserList[i] = userResponse.User
		}(i, v)
	}
	wg.Wait()
	return
}

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// GetFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetFollowList(ctx context.Context, req *relation.FollowListRequest) (resp *relation.FollowListResponse, err error) {
	logger := logging.Logger.WithFields(logrus.Fields{
		"actor_id": req.ActorId,
		"user_id":  req.UserId,
		"function": "GetFollowList",
	})
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

	userIds := make([]*uint32, len(relationModels))
	for i, m := range relationModels {
		userIds[i] = &m.TargetId
	}
	userList := queryUsers(ctx, logger, req.ActorId, userIds)

	logger.WithFields(logrus.Fields{
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
		"actor_id": req.ActorId,
		"user_id":  req.UserId,
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

	userIds := make([]*uint32, len(relationModels))
	for i, m := range relationModels {
		userIds[i] = &m.UserId
	}
	userList := queryUsers(ctx, logger, req.ActorId, userIds)

	logger.WithFields(logrus.Fields{
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
	logger := logging.Logger.WithFields(logrus.Fields{
		"actor_id": req.ActorId,
		"user_id":  req.UserId,
		"function": "GetFriendList",
	})
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

	var followerIds []*uint32
	for _, followerModel := range followerModels {
		followerIds = append(followerIds, &followerModel.UserId)
	}
	var followIds []*uint32
	for _, followModel := range followModels {
		followIds = append(followIds, &followModel.TargetId)
	}

	var ids = intersection(followerIds, followIds)

	userList := queryUsers(ctx, logger, req.ActorId, ids)
	userList = append(userList, biz.ChatGPTUser)

	resp = &relation.FriendListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		UserList:   userList,
	}

	logger.WithFields(logrus.Fields{
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return
}

// Follow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Follow(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	logger := logging.Logger.WithFields(logrus.Fields{
		"actor_id": req.ActorId,
		"user_id":  req.UserId,
		"function": "RelationAction",
	})
	logger.Debug("Process start")

	if req.UserId == req.ActorId {
		resp = &relation.RelationActionResponse{
			StatusCode: biz.InvalidToUserId,
			StatusMsg:  biz.BadRequestStatusMsg,
		}
		return
	}

	relationModel := model.Relation{
		UserId:   req.ActorId,
		TargetId: req.UserId,
	}

	// make sure target id exists
	userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
		UserId:  req.UserId,
		ActorId: req.ActorId,
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
	logger.WithFields(logrus.Fields{
		"entry": relationModel,
	}).Debug("create relation")

	resp = &relation.RelationActionResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
	}

	logger.WithFields(logrus.Fields{
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return
}

// Unfollow implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Unfollow(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	logger := logging.Logger.WithFields(logrus.Fields{
		"actor_id": req.ActorId,
		"user_id":  req.UserId,
		"function": "RelationAction",
	})
	logger.Debug("Process start")
	if req.ActorId == req.UserId {
		resp = &relation.RelationActionResponse{
			StatusCode: biz.InvalidToUserId,
			StatusMsg:  biz.BadRequestStatusMsg,
		}
		return
	}

	r := repo.Q.Relation

	relationModel, err := r.WithContext(ctx).Where(r.UserId.Eq(req.ActorId), r.TargetId.Eq(req.UserId)).First()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"actor_id": req.ActorId,
			"user_id":  req.UserId,
		}).Debug("record not found")
		// Unfollow a user that was not followed before
		resp = &relation.RelationActionResponse{
			StatusCode: biz.RelationNotFound,
			StatusMsg:  biz.BadRequestStatusMsg,
		}
		return
	}

	if _, err = r.Where(r.UserId.Eq(req.ActorId), r.TargetId.Eq(req.UserId)).Delete(); err != nil {
		logger.WithFields(logrus.Fields{
			"entry": relationModel,
			"err":   err,
		}).Debug("failed to delete")

		resp = &relation.RelationActionResponse{
			StatusCode: biz.UnableToDeleteRelation,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}

	logger.WithFields(logrus.Fields{
		"entry": relationModel,
	}).Debug("deleted db entry")

	logger.WithFields(logrus.Fields{
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
		"function": "CountFollowList",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	r := repo.Q.Relation
	count, err := r.WithContext(ctx).Where(r.UserId.Eq(req.UserId)).Count()
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"user_id": req.UserId,
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
		"function": "CountFollowerList",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	r := repo.Q.Relation
	count, err := r.WithContext(ctx).Where(r.TargetId.Eq(req.UserId)).Count()
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"user_id": req.UserId,
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
		"function": "IsFollow",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	r := repo.Q.Relation
	count, err := r.WithContext(ctx).Where(r.UserId.Eq(req.ActorId), r.TargetId.Eq(req.UserId)).Count()
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"user_id": req.UserId,
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
