package kboot

import (
	"context"
	"fmt"
	"strings"

	"github.com/guestin/log"
	"github.com/ooopSnake/assert.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// HideBanner hide the bootstrap banner
func HideBanner() {
	_gCtx.hideBanner = true
}

func Version() string {
	return __version
}

func GetContext() Context {
	return _gCtx
}

func GetViper() *viper.Viper {
	return GetContext().GetViper()
}

func GetActivatedProfile() string {
	return GetContext().GetActivatedProfile()
}

func GetRootLogger() *zap.Logger {
	return _gCtx.rootLogger
}

func GetTaggedZapLogger(tag string, opt ...log.Opt) log.ZapLog {
	return GetContext().GetTaggedZapLogger(tag, opt...)
}

func GetTaggedLogger(tag string, opt ...log.Opt) log.ClassicLog {
	return GetContext().GetTaggedLogger(tag, opt...)
}

func UnmarshalSubConfig(key string, i interface{}, options ...CfgOption) (err error) {
	return GetContext().UnmarshalSubConfig(key, i, options...)
}

func RegisterUnit(name string, fn InitFunc, options ...UnitOption) {
	assert.Must(len(strings.TrimSpace(name)) != 0, "name must not empty or blank").Panic()
	assert.Must(fn != nil, "init func must not be nil").Panic()
	for _, u := range _gCtx.units {
		if u.GetName() == name {
			assert.Must(false, fmt.Sprintf("name '%s' already exist", name)).Panic()
		}
	}
	unit := &unitImpl{
		rootCtx:  _gCtx,
		name:     name,
		initFunc: fn,
	}
	for _, opt := range options {
		opt.apply(unit)
	}
	_gCtx.units = append(_gCtx.units, unit)
}

func Bootstrap(ctx context.Context, app Application, options ...BootOption) {
	if !_gCtx.hideBanner {
		fmt.Print(_BANNER)
	}
	_gCtx.logger.Info("Bootstrap ... ", zap.String("app", app.GetAppName()), zap.String("tz", app.GetTimezone().String()))
	assert.Must(ctx != nil, "root ctx must not be nil").Panic()
	assert.Must(app != nil, "app must not be nil").Panic()
	_ctx, cancel := context.WithCancel(ctx)
	_gCtx.ctx = _ctx
	_gCtx.cancel = cancel
	_gCtx.Application = app
	for _, opt := range options {
		opt.apply(_gCtx)
	}
	_gCtx.bootStrap()
}
