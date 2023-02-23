//go:build e2e

// Progress:
//   All tests are synced with 2023-2-20 version of the API. Not sure if the API will change in the future.

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/segmentio/ksuid"
)

// TestFeed tests the procedure of unauthenticated users getting feed.
func TestFeed(t *testing.T) {
	e := newExpect(t)
	feedResp := e.GET("/douyin/feed/").Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().Equal(0)
	feedResp.Value("status_msg").String().NotEmpty()
	feedResp.Value("video_list").Array().Length().Gt(0)
	nextTime := feedResp.Value("next_time").String()
	nextTimeInt, err := strconv.ParseInt((*nextTime).Raw(), 10, 64)
	if err != nil {
		t.Error(err)
	}
	// check if nextTimeInt is a valid timestamp later than 2023-01-01 00:00:00
	if nextTimeInt < 1672502400000 {
		t.Error("next_time is not a valid timestamp")
	}

	for _, element := range feedResp.Value("video_list").Array().Iter() {
		video := element.Object()
		video.ContainsKey("id")
		author := video.Value("author").Object()

		ValidateUser(author)

		video.Value("play_url").String().NotEmpty()
		video.Value("cover_url").String().NotEmpty()
		video.ContainsKey("favorite_count")
		video.ContainsKey("comment_count")
		video.ContainsKey("is_favorite")
		video.Value("title").String().NotEmpty()
	}
}

// TestUserAction tests a whole user register & login & fetch user info procedure.
func TestUserAction(t *testing.T) {
	e := newExpect(t)

	registerValue := fmt.Sprintf("douyin_test_%s", ksuid.New().String())

	registerResp := e.POST("/douyin/user/register/").
		WithQuery("username", registerValue).WithQuery("password", registerValue).
		WithFormField("username", registerValue).WithFormField("password", registerValue).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	registerResp.Value("status_code").Number().Equal(0)
	registerResp.Value("status_msg").String().NotEmpty()
	registerResp.Value("user_id").Number().Gt(0)
	registerResp.Value("token").String().NotEmpty()

	loginResp := e.POST("/douyin/user/login/").
		WithQuery("username", registerValue).WithQuery("password", registerValue).
		WithFormField("username", registerValue).WithFormField("password", registerValue).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	loginResp.Value("status_code").Number().Equal(0)
	registerResp.Value("status_msg").String().NotEmpty()
	loginResp.Value("user_id").Number().Gt(0)
	loginResp.Value("token").String().NotEmpty()

	token := loginResp.Value("token").String().Raw()
	userId := loginResp.Value("user_id").Number().Raw()

	userResp := e.GET("/douyin/user/").
		WithQuery("user_id", userId).
		WithQuery("token", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	userResp.Value("status_code").Number().Equal(0)
	userResp.Value("status_msg").String().NotEmpty()
	userInfo := userResp.Value("user").Object()
	ValidateUser(userInfo)
}

// TestPublish tests the procedure of a publish & check the publish list of the user.
func TestPublish(t *testing.T) {
	e := newExpect(t)

	userId, token := getTestUserToken(testUserA, e)

	publishResp := e.POST("/douyin/publish/action/").
		WithMultipart().
		WithFile("data", "../../service/publish/resources/bear.mp4").
		WithFormField("token", token).
		WithFormField("title", "Bear").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	publishResp.Value("status_code").Number().Equal(0)
	publishResp.Value("status_msg").String().NotEmpty()

	publishListResp := e.GET("/douyin/publish/list/").
		WithQuery("user_id", userId).WithQuery("token", token).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	publishListResp.Value("status_code").Number().Equal(0)
	publishListResp.Value("video_list").Array().Length().Gt(0)

	for _, element := range publishListResp.Value("video_list").Array().Iter() {
		video := element.Object()
		video.ContainsKey("id")
		video.ContainsKey("author")
		video.Value("play_url").String().NotEmpty()
		video.Value("cover_url").String().NotEmpty()
	}
}
