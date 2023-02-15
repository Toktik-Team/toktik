package biz

import "github.com/cloudwego/hertz/pkg/protocol/consts"

const (
	VideoCount = 30
)

var (
	InvalidLatestTime = GWError{HTTPStatusCode: consts.StatusBadRequest, StatusCode: 400001, StatusMsg: "Invalid latest_time"}
)
