package biz

import "github.com/cloudwego/hertz/pkg/protocol/consts"

const (
	VideoCount = 30

	Unable2ParseLatestTimeStatusCode = 400001
	SQLQueryErrorStatusCode          = 500001
)

var (
	InvalidLatestTime = GWError{HTTPStatusCode: consts.StatusBadRequest, StatusCode: 400001, StatusMsg: "Invalid latest_time"}
)
