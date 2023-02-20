package main

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var serverAddr = "https://toktik.xctra.cn"

// TODO: use safer way to store test user credentials
var testUserA = "douyinTestUserA"
var testUserB = "douyinTestUserB"

func newExpect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client:   http.DefaultClient,
		BaseURL:  serverAddr,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

func getTestUserToken(user string, e *httpexpect.Expect) (int, string) {
	registerResp := e.POST("/douyin/user/register/").
		WithQuery("username", user).WithQuery("password", user).
		WithFormField("username", user).WithFormField("password", user).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	userId := 0
	token := registerResp.Value("token").String().Raw()
	if len(token) == 0 {
		loginResp := e.POST("/douyin/user/login/").
			WithQuery("username", user).WithQuery("password", user).
			WithFormField("username", user).WithFormField("password", user).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		loginToken := loginResp.Value("token").String()
		loginToken.Length().Gt(0)
		token = loginToken.Raw()
		userId = int(loginResp.Value("user_id").Number().Raw())
	} else {
		userId = int(registerResp.Value("user_id").Number().Raw())
	}
	return userId, token
}

func ValidateUser(user *httpexpect.Object) {
	user.ContainsKey("id")
	user.Value("name").String().NotEmpty()
	user.ContainsKey("follow_count")
	user.ContainsKey("follower_count")
	user.ContainsKey("is_follow")
	user.Value("avatar").String().NotEmpty()
	user.Value("background_image").String().NotEmpty()
	user.Value("signature").String().NotEmpty()
	user.ContainsKey("total_favorited") // TODO: determine if this field should be string or int
	user.ContainsKey("work_count")
	user.ContainsKey("favorite_count")
}
