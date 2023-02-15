package biz

import "github.com/cloudwego/hertz/pkg/protocol/consts"

var (
	PublishActionSuccess = "Uploaded successfully!"
)

var (
	OpenFileFailedError = GWError{HTTPStatusCode: consts.StatusInternalServerError, StatusCode: 400005, StatusMsg: "Open file failed"}
	SizeNotMatchError   = GWError{HTTPStatusCode: consts.StatusInternalServerError, StatusCode: 400006, StatusMsg: "Size not match"}
)
