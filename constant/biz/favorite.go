package biz

import (
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const (
	UnauthorizedAction            = 401001
	FailedToLikeVideo             = 500001
	FailedToAddVideoFavoriteCount = 500002
	FailedToCancelLike            = 500003
	FailedToSubVideoFavoriteCount = 500004
	FailedToGetVideoList          = 500005
)

var (
	AuthFailed = GWError{
		HTTPStatusCode: consts.StatusUnauthorized,
		StatusCode:     UnauthorizedAction,
		StatusMsg:      "Unauthorized Action",
	}
)
