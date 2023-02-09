package mw

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	authRPC "toktik/kitex_gen/douyin/auth"
	"toktik/service/web/auth"
)

func init() {
	hlog.Info("using auth")
}

// AuthResult Authentication result enum
type AuthResult string

const (
	// AUTH_RESULT_SUCCESS Authentication success
	AUTH_RESULT_SUCCESS AuthResult = "success"
	// AUTH_RESULT_NO_TOKEN Authentication failed due to no token
	AUTH_RESULT_NO_TOKEN AuthResult = "no_token"
	// AUTH_RESULT_UNKNOWN Authentication failed due to unknown reason
	AUTH_RESULT_UNKNOWN = "unknown"
)

const (
	AUTH_RESULT_KEY = "authentication_result"
	USER_ID_KEY     = "user_id"
)

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, rc *app.RequestContext) {
		token := rc.Query("token")
		if token == "" {
			rc.Set(AUTH_RESULT_KEY, AUTH_RESULT_NO_TOKEN)
			rc.Set(USER_ID_KEY, 0)
			rc.Next(ctx)
		}

		authResp, err := auth.Client.Authenticate(ctx, &authRPC.AuthenticateRequest{Token: token})
		if err != nil {
			rc.Set(AUTH_RESULT_KEY, AUTH_RESULT_UNKNOWN)
			rc.Set(USER_ID_KEY, 0)
			rc.Next(ctx)
		}
		if authResp.StatusCode == 0 && authResp.StatusMsg == string(AUTH_RESULT_SUCCESS) {
			rc.Set(AUTH_RESULT_KEY, AUTH_RESULT_SUCCESS)
			rc.Set(USER_ID_KEY, authResp.UserId)
		} else {
			rc.Set(AUTH_RESULT_KEY, AUTH_RESULT_UNKNOWN)
			rc.Set(USER_ID_KEY, 0)
		}
		rc.Next(ctx)
	}
}
