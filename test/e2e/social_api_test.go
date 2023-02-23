//go:build e2e

// Progress:
//   All tests are synced with 2023-2-20 version of the API. Not sure if the API will change in the future.

package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelation(t *testing.T) {
	e := newExpect(t)

	userIdA, tokenA := getTestUserToken(testUserA, e)
	userIdB, tokenB := getTestUserToken(testUserB, e)

	// Sample:
	//{
	//    "status_code": 0,
	//    "status_msg": "string"
	//}
	relationResp := e.POST("/douyin/relation/action/").
		WithQuery("token", tokenA).WithQuery("to_user_id", userIdB).WithQuery("action_type", 1).
		WithFormField("token", tokenA).WithFormField("to_user_id", userIdB).WithFormField("action_type", 1).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	relationResp.Value("status_code").Number().Equal(0)
	relationResp.Value("status_msg").String().NotEmpty()

	// Sample:
	//{
	//    "status_code": "string",
	//    "status_msg": "string",
	//    "user_list": [
	//        {
	//            "id": 0,
	//            "name": "string",
	//            "follow_count": 0,
	//            "follower_count": 0,
	//            "is_follow": true,
	//            "avatar": "string",
	//            "background_image": "string",
	//            "signature": "string",
	//            "total_favorited": "string",
	//            "work_count": 0,
	//            "favorite_count": 0
	//        }
	//    ]
	//}
	followListResp := e.GET("/douyin/relation/follow/list/").
		WithQuery("token", tokenA).WithQuery("user_id", userIdA).
		WithFormField("token", tokenA).WithFormField("user_id", userIdA).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	followListResp.Value("status_code").Number().Equal(0)
	followListResp.Value("user_list").Array().NotEmpty()

	containTestUserB := false
	for _, element := range followListResp.Value("user_list").Array().Iter() {
		user := element.Object()
		user.ContainsKey("id")
		ValidateUser(user)
		if int(user.Value("id").Number().Raw()) == userIdB {
			containTestUserB = true
		}
	}
	assert.True(t, containTestUserB, "Follow test user failed")

	followerListResp := e.GET("/douyin/relation/follower/list/").
		WithQuery("token", tokenB).WithQuery("user_id", userIdB).
		WithFormField("token", tokenB).WithFormField("user_id", userIdB).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	followerListResp.Value("status_code").Number().Equal(0)
	followerListResp.Value("user_list").Array().NotEmpty()

	containTestUserA := false
	for _, element := range followerListResp.Value("user_list").Array().Iter() {
		user := element.Object()
		user.ContainsKey("id")
		ValidateUser(user)
		if int(user.Value("id").Number().Raw()) == userIdA {
			containTestUserA = true
		}
	}
	assert.True(t, containTestUserA, "Follower test user failed")
}

func TestChat(t *testing.T) {
	e := newExpect(t)

	userIdA, tokenA := getTestUserToken(testUserA, e)
	userIdB, tokenB := getTestUserToken(testUserB, e)

	// Sample:
	//{
	//    "status_code": 0,
	//    "status_msg": "string"
	//}
	messageResp := e.POST("/douyin/message/action/").
		WithQuery("token", tokenA).WithQuery("to_user_id", userIdB).WithQuery("action_type", 1).WithQuery("content", "Send to UserB").
		WithFormField("token", tokenA).WithFormField("to_user_id", userIdB).WithFormField("action_type", 1).WithQuery("content", "Send to UserB").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	messageResp.Value("status_code").Number().Equal(0)
	messageResp.Value("status_msg").String().NotEmpty()

	// Sample:
	//{
	//    "status_code": "string",
	//    "status_msg": "string",
	//    "message_list": [
	//        {
	//            "id": 0,
	//            "to_user_id": 0,
	//            "from_user_id": 0,
	//            "content": "string",
	//            "pre_msg_time": "string",
	//            "create_time": 0
	//        }
	//    ]
	//}
	chatResp := e.GET("/douyin/message/chat/").
		WithQuery("token", tokenA).WithQuery("to_user_id", userIdB).
		WithFormField("token", tokenA).WithFormField("to_user_id", userIdB).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	chatResp.Value("status_code").Number().Equal(0)
	chatResp.Value("status_msg").String().NotEmpty()
	chatResp.Value("message_list").Array().Length().Gt(0)
	chatResp.Value("message_list").Array().First().Object().Value("content").String().Equal("Send to UserB")
	chatResp.Value("message_list").Array().First().Object().Value("create_time").String().NotEmpty()

	chatResp = e.GET("/douyin/message/chat/").
		WithQuery("token", tokenB).WithQuery("to_user_id", userIdA).
		WithFormField("token", tokenB).WithFormField("to_user_id", userIdA).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	chatResp.Value("status_code").Number().Equal(0)
	chatResp.Value("status_msg").String().NotEmpty()
	chatResp.Value("message_list").Array().Length().Gt(0)
	chatResp.Value("message_list").Array().First().Object().Value("content").String().Equal("Send to UserB")
	chatResp.Value("message_list").Array().First().Object().Value("create_time").String().NotEmpty()
}
