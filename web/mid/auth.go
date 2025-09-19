package mid

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/guestin/kboot/web/kerrors"
	"github.com/guestin/mob"
	"github.com/labstack/echo/v4"
)

type (
	AuthSessionInfo struct {
		IsAnonymous bool
		Uid         string
		SessionId   string
		ExpireAt    int64
		ClientIp    string
		ClientUA    string
		UserData    interface{}
	}
	SessionProviderFunc func(ctx echo.Context, sessionId string) (*AuthSessionInfo, error)
	AuthConfig          struct {
		Enable               bool     `toml:"enable" mapstructure:"enable"` //是否启用，启用后将解析session info
		Whitelist            []string `toml:"whitelist" mapstructure:"whitelist"`
		SessionIdKey         string   `toml:"sessionIdKey" mapstructure:"sessionIdKey"`
		SessionExpireInHours int      `toml:"sessionExpireInHours" validate:"gte=0,lte=720" mapstructure:"sessionExpireInHours"`
	}
)

var AuthSessionProvider SessionProviderFunc = nil

var DefaultAuthConfig = AuthConfig{
	Enable:       false,
	Whitelist:    []string{},
	SessionIdKey: "kt-session-id",
}

func anonymousSession() *AuthSessionInfo {
	now := time.Now()
	randomId := strings.ReplaceAll(fmt.Sprintf("ANONYMOUS_%s", now.Format("20060102150405.000000")), ".", "")
	return &AuthSessionInfo{
		IsAnonymous: true,
		Uid:         randomId,
		SessionId:   "",
		ExpireAt:    now.Add(time.Hour * 2).Unix(),
		UserData:    nil,
	}
}

func CurrentSession(ctx echo.Context) *AuthSessionInfo {
	session := ctx.Get(CtxCallerInfoKey)
	if session == nil {
		return anonymousSession()
	}
	return session.(*AuthSessionInfo)
}

func Auth() echo.MiddlewareFunc {
	return AuthWithConfig(DefaultAuthConfig)
}

func AuthWithConfig(config AuthConfig) echo.MiddlewareFunc {
	if strings.Trim(config.SessionIdKey, " ") == "" {
		config.SessionIdKey = DefaultAuthConfig.SessionIdKey
	}
	excludePathSet := mob.NewConcurrentSet()
	excludeRegList := make([]*regexp.Regexp, 0)
	for i := range config.Whitelist {
		p := config.Whitelist[i]
		if excludePathSet.Add(p) {
			reg, err := regexp.Compile(p)
			if err != nil {
				panic(fmt.Sprintf("witelist path %s is not a valid reg path", p))
			}
			excludeRegList = append(excludeRegList, reg)
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			reqPath := ctx.Request().URL.String()
			ignore := false
			for i := range excludeRegList {
				if excludeRegList[i].MatchString(reqPath) {
					ignore = true
					break
				}
			}
			token := ctx.Request().Header.Get(config.SessionIdKey)
			if len(token) == 0 {
				token = ctx.QueryParam(config.SessionIdKey)
			}
			if len(token) == 0 && !ignore {
				return kerrors.ErrUnauthorized()
			}
			var sessionInfo *AuthSessionInfo = nil
			var err error
			if len(token) > 0 {
				if AuthSessionProvider != nil {
					sessionInfo, err = AuthSessionProvider(ctx, token)
					if err != nil && !ignore {
						return err
					}
				}
			}
			if sessionInfo == nil {
				sessionInfo = anonymousSession()
			}
			sessionInfo.SessionId = token
			sessionInfo.ClientIp = ctx.RealIP()
			sessionInfo.ClientUA = ctx.Request().UserAgent()
			ctx.Set(CtxCallerInfoKey, sessionInfo)
			return next(ctx)
		}
	}
}
