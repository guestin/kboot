package web

import (
	"github.com/guestin/kboot/web/mid"
	"github.com/labstack/echo/v4"
)

var _options = make([]Option, 0)
var _rspFmtCfg mid.FormatConfig = mid.DefaultFormatConfig

var _dependencies = make([]string, 0)

type Option interface {
	apply(e *echo.Echo) error
}

type optionFunc func(e *echo.Echo) error

func (f optionFunc) apply(e *echo.Echo) error {
	return f(e)
}

type RouteBuilder func(eCtx *echo.Echo) error

func RspFormat(config mid.FormatConfig) {
	if config.Skipper != nil {
		_rspFmtCfg.Skipper = config.Skipper
	}
}

func Router(fn RouteBuilder) {
	_options = append(_options, optionFunc(func(e *echo.Echo) error {
		return fn(e)
	}))
}

func DependsOn(m ...string) {
	_dependencies = append(_dependencies, m...)
}
