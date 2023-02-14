package biz

import httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"

const (
	InvalidCommentActionType = 400001
	VideoNotFound            = 400002
	ActorIDNotMatch          = 403001
	UnableToCreateComment    = 500001
	UnableToDeleteComment    = 500002
	UnableToQueryVideo       = 500003
	UnableToQueryComment     = 500004
	UnableToQueryUser        = 500005
)

var (
	UnauthorizedError = GWError{HTTPStatusCode: httpStatus.StatusUnauthorized, StatusCode: 400003, StatusMsg: "Unauthorized"}

	BadRequestError = GWError{HTTPStatusCode: httpStatus.StatusBadRequest, StatusCode: 400004, StatusMsg: "Bad request"}

	InternalServerError = GWError{HTTPStatusCode: httpStatus.StatusInternalServerError, StatusCode: 500006, StatusMsg: "Internal server error"}
)
