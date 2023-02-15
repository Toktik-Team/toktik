package biz

import httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"

var (
	UnauthorizedError   = GWError{HTTPStatusCode: httpStatus.StatusUnauthorized, StatusCode: 400003, StatusMsg: "Unauthorized"}
	BadRequestError     = GWError{HTTPStatusCode: httpStatus.StatusBadRequest, StatusCode: 400004, StatusMsg: "Bad request"}
	InternalServerError = GWError{HTTPStatusCode: httpStatus.StatusInternalServerError, StatusCode: 500006, StatusMsg: "Internal server error"}
)
