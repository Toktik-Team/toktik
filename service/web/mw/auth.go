package mw

import (
	"context"
	"toktik/kitex_gen/douyin/auth"
	authHandler "toktik/service/web/auth"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// AuthResult Authentication result enum
const (
	// AUTH_RESULT_SUCCESS Authentication success
	AUTH_RESULT_SUCCESS string = "success"
	// AUTH_RESULT_NO_TOKEN Authentication failed due to no token
	AUTH_RESULT_NO_TOKEN string = "no_token"
	// AUTH_RESULT_UNKNOWN Authentication failed due to unknown reason
	AUTH_RESULT_UNKNOWN string = "unknown"
)

const (
	authResultKey = "authentication_result"
	UserIdKey     = "user_id"
)

func init() {
	hlog.Info("using auth")
}

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, rc *app.RequestContext) {
		var token string
		formToken := string(rc.FormValue("token"))

		if formToken != "" {
			token = formToken
		}

		if token == "" {
			rc.Set(authResultKey, AUTH_RESULT_NO_TOKEN)
			rc.Set(UserIdKey, 0)
			rc.Next(ctx)
			return
		}

		authResp, err := authHandler.Client.Authenticate(ctx, &auth.AuthenticateRequest{Token: token})
		if err != nil {
			rc.Set(authResultKey, AUTH_RESULT_UNKNOWN)
			rc.Set(UserIdKey, 0)
			rc.Next(ctx)
			return
		}
		if authResp.StatusCode == 0 && authResp.StatusMsg == string(AUTH_RESULT_SUCCESS) {
			rc.Set(authResultKey, AUTH_RESULT_SUCCESS)
			rc.Set(UserIdKey, authResp.UserId)
		} else {
			rc.Set(authResultKey, AUTH_RESULT_UNKNOWN)
			rc.Set(UserIdKey, 0)
		}
		rc.Next(ctx)
	}
}

func GetAuthResult(c *app.RequestContext) string {
	authResult := c.GetString(authResultKey)
	if authResult == "" {
		return AUTH_RESULT_UNKNOWN
	}
	return authResult
}

func GetAuthActorId(c *app.RequestContext) uint32 {
	return c.GetUint32(UserIdKey)
}
