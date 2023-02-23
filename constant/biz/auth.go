package biz

import "github.com/cloudwego/hertz/pkg/protocol/consts"

var (
	TokenNotFoundMessage = "会话已过期"
)

var (
	NoUserNameOrPassWord = GWError{HTTPStatusCode: consts.StatusBadRequest, StatusCode: 400007, StatusMsg: "no username or password"}
)
