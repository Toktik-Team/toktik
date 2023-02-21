package biz

import "github.com/cloudwego/hertz/pkg/protocol/consts"

var (
	TokenNotFoundMessage = "token not found"
)

var (
	NoUserNameOrPassWord = GWError{HTTPStatusCode: consts.StatusBadRequest, StatusCode: 400007, StatusMsg: "no username or password"}
)
