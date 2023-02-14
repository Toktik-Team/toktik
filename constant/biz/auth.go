package biz

import "github.com/cloudwego/hertz/pkg/protocol/consts"

var (
	TokenNotFoundMessage = "token not found"
)

const (
	UserNameExist       = 400001
	ServiceNotAvailable = 503001
	UserNotFound        = 400002
	PasswordIncorrect   = 401003
	RequestIsNil        = 500001
	TokenNotFound       = 401001
)

var (
	NoUserNameOrPassWord = GWError{HTTPStatusCode: consts.StatusBadRequest, StatusCode: 400007, StatusMsg: "no username or password"}
)
