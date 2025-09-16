package mid

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/guestin/kboot/web/kerrors"
	"github.com/guestin/mob"
	"github.com/labstack/echo/v4"
)

type (
	ACLSessionInfo struct {
		IsAnonymous bool
		Uid         string
		SessionId   string
		ExpireAt    int64
		UserData    interface{}
	}
	SessionProviderFunc func(ctx context.Context, sessionId string) (*ACLSessionInfo, error)
	ACLConfig           struct {
		Enable          bool     `toml:"enable"`    //是否启用，启用后将解析session info
		EnableAcl       bool     `toml:"enableAcl"` //是否启用权限控制 ，启用后将校验系统权限
		Whitelist       []string `toml:"whitelist"`
		SessionIdKey    string   `toml:"sessionIdKey"`
		SessionProvider SessionProviderFunc
	}
)

var DefaultACLConfig = ACLConfig{
	Enable:       false,
	EnableAcl:    false,
	Whitelist:    []string{},
	SessionIdKey: "kt-session-id",
}

func anonymousSession() *ACLSessionInfo {
	now := time.Now()
	randomId := strings.ReplaceAll(fmt.Sprintf("ANONYMOUS_%s", now.Format("20060102150405.000000")), ".", "")
	return &ACLSessionInfo{
		IsAnonymous: true,
		Uid:         randomId,
		SessionId:   "",
		ExpireAt:    now.Add(time.Hour * 2).Unix(),
		UserData:    nil,
	}
}

func CurrentACLSession(ctx echo.Context) *ACLSessionInfo {
	session := ctx.Get(CtxCallerInfoKey)
	if session == nil {
		return anonymousSession()
	}
	return session.(*ACLSessionInfo)
}

func ACL() echo.MiddlewareFunc {
	return ACLWithConfig(DefaultACLConfig)
}

func ACLWithConfig(config ACLConfig) echo.MiddlewareFunc {
	if strings.Trim(config.SessionIdKey, " ") == "" {
		config.SessionIdKey = DefaultACLConfig.SessionIdKey
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
			var sessionInfo *ACLSessionInfo = nil
			var err error
			if len(token) > 0 {
				if config.SessionProvider != nil {
					sessionInfo, err = config.SessionProvider(UnwrapContext(ctx), token)
					if err != nil && !ignore {
						return err
					}
				}
			}
			ctx.RealIP()
			if sessionInfo == nil {
				sessionInfo = anonymousSession()
			}
			sessionInfo.SessionId = token
			ctx.Set(CtxCallerInfoKey, sessionInfo)
			return next(ctx)
		}
	}
}
