package biz

import "github.com/cloudwego/hertz/pkg/protocol/consts"

var (
	InvalidUserID = GWError{HTTPStatusCode: consts.StatusBadRequest, StatusCode: 400003, StatusMsg: "Invalid user_id"}
)
