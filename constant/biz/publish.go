package biz

import "github.com/cloudwego/hertz/pkg/protocol/consts"

const (
	InvalidContentType     = 400101
	Unable2GenerateUUID    = 500101
	Unable2CreateThumbnail = 500102
	Unable2UploadVideo     = 500103
	Unable2UploadCover     = 500104
	Unable2CreateDBEntry   = 500105
)

var (
	PublishActionSuccess = "Uploaded successfully!"
)

var (
	OpenFileFailedError = GWError{HTTPStatusCode: consts.StatusInternalServerError, StatusCode: 400005, StatusMsg: "Open file failed"}

	SizeNotMatchError = GWError{HTTPStatusCode: consts.StatusInternalServerError, StatusCode: 400006, StatusMsg: "Size not match"}
)
