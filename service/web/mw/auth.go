package mw

import (
	"context"
	"toktik/constant/biz"
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
	AuthResultKey = "authentication_result"
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
			rc.Set(AuthResultKey, AUTH_RESULT_NO_TOKEN)
			rc.Set(UserIdKey, 0)
			rc.Next(ctx)
			return
		}

		authResp, err := authHandler.Client.Authenticate(ctx, &auth.AuthenticateRequest{Token: token})
		if err != nil {
			rc.Set(AuthResultKey, AUTH_RESULT_UNKNOWN)
			rc.Set(UserIdKey, 0)
			rc.Next(ctx)
			return
		}
		if authResp.StatusCode == 0 && authResp.StatusMsg == AUTH_RESULT_SUCCESS {
			rc.Set(AuthResultKey, AUTH_RESULT_SUCCESS)
			rc.Set(UserIdKey, authResp.UserId)
		} else {
			rc.Set(AuthResultKey, AUTH_RESULT_UNKNOWN)
			rc.Set(UserIdKey, 0)
		}
		rc.Next(ctx)
	}
}

type config struct {
	authRequired bool
}

// Option opts for opentelemetry tracer provider
type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

func newConfig(opts []Option) *config {
	cfg := &config{}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	return cfg
}

func WithAuthRequired() Option {
	return option(func(cfg *config) {
		cfg.authRequired = true
	})
}

func Auth(c *app.RequestContext, opts ...Option) (actorIdPtr *uint32, ok bool) {
	cfg := newConfig(opts)

	switch c.GetString(AuthResultKey) {
	case AUTH_RESULT_SUCCESS:
		*actorIdPtr = c.GetUint32(UserIdKey)
	case AUTH_RESULT_NO_TOKEN:
		// CHATGPT GENERATED CR COMMENT:
		// 这段代码中 fallthrough 和 break 的使用方法是正确的，不会产生 bug。
		//
		// 在 switch 语句中，使用 fallthrough 语句可以使程序继续执行下一个 case 语句，而不进行判断。在上述代码中，如果
		// AuthResultKey 的值为 AUTH_RESULT_NO_TOKEN，程序会先判断 !cfg.authRequired，如果为 true，则跳出
		// switch 语句，执行后续代码；如果为 false，则会继续执行下一个 case 语句，即 default，然后执行
		// biz.UnAuthorized.LaunchError(c)，返回 nil 和 false。
		//
		// 而在 switch 语句中，使用 break 语句可以跳出 switch 语句，不再继续执行下一个 case 语句。在上述代码中，如果
		// AuthResultKey 的值为 AUTH_RESULT_SUCCESS，程序会执行 actorId =
		// c.GetUint32(UserIdKey)，然后直接跳出 switch 语句，返回 actorIdPtr 和 true。
		//
		// 因此，使用 fallthrough 和 break 的方式是正确的，不会产生 bug。
		if !cfg.authRequired {
			break
		}
		fallthrough
	default:
		biz.UnAuthorized.LaunchError(c)
		return nil, false
	}
	return actorIdPtr, true
}
