package biz

import (
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const (
	UnauthorizedAction   = 401001
	FailedToGetVideoList = 500001
)

var (
	AuthFailed = GWError{
		HTTPStatusCode: consts.StatusUnauthorized,
		StatusCode:     UnauthorizedAction,
		StatusMsg:      "Unauthorized Action",
	}
)
