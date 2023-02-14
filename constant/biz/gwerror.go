package biz

import (
	"fmt"
	"toktik/logging"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/sirupsen/logrus"
)

// GWError is the error struct for gateway.
type GWError struct {
	HTTPStatusCode int
	StatusCode     uint32
	StatusMsg      string
	Cause          *error
	Fields         *logrus.Fields
}

func (G GWError) Error() string {
	return fmt.Sprintf("http status code: %d, status code: %d, status msg: %s, cause: %v", G.HTTPStatusCode, G.StatusCode, G.StatusMsg, G.Cause)
}

// extractFieldsFromRequest extracts important fields from request context for debugging purpose.
func extractFieldsFromRequest(c *app.RequestContext) logrus.Fields {
	return logrus.Fields{"request_context_info": map[string]interface{}{
		"request_id":   c.Request.Header.Get("X-Request-Id"),
		"method":       c.Method(),
		"host":         c.Host(),
		"uri":          c.URI(),
		"ip":           c.ClientIP(),
		"ua":           c.UserAgent(),
		"query_args":   c.QueryArgs(),
		"post_args":    c.PostArgs(),
		"content_type": c.ContentType(),
		"body":         c.GetRawData(),
		"handler":      c.HandlerName(),
	}}
}

// LaunchError logs the error and returns the error.
func (G GWError) LaunchError(c *app.RequestContext) {
	logger := logging.Logger.WithFields(extractFieldsFromRequest(c))
	if G.Fields != nil {
		logger = logger.WithFields(*G.Fields)
	}
	if G.Cause != nil {
		logger = logger.WithField("cause", *G.Cause)
	}
	logger.Debugf("launch error: %v", G)
	c.JSON(G.HTTPStatusCode, map[string]interface{}{
		"status_code": G.StatusCode,
		"status_msg":  G.StatusMsg,
	})
}

// WithCause adds the cause to the error.
func (G GWError) WithCause(err error) GWError {
	G.Cause = &err
	return G
}

// WithFields adds the fields to the error.
func (G GWError) WithFields(fields *logrus.Fields) GWError {
	G.Fields = fields
	return G
}

var (
	RPCCallError = GWError{HTTPStatusCode: consts.StatusInternalServerError, StatusCode: 500001, StatusMsg: "RPC call error"}
)
