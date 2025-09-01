package kboot

import (
	"strings"
)

type Option interface {
	apply(*ctxImpl)
}

type optionFunc func(*ctxImpl)

func (f optionFunc) apply(ctx *ctxImpl) {
	f(ctx)
}

// HideBanner hide the bootstrap banner
func HideBanner() Option {
	return optionFunc(func(ctx *ctxImpl) {
		ctx.hideBanner = true
	})
}

// AutoFindConfig auto find config file (name.fileType) from special search paths
// will override by ConfigFromFile
func AutoFindConfig(name string, fileType string, searchPaths ...string) Option {
	return optionFunc(func(ctx *ctxImpl) {
		ctx.configName = name
		if len(strings.Trim(fileType, " ")) > 0 {
			ctx.configFileType = strings.Trim(fileType, " ")
		}
		if len(searchPaths) > 0 {
			ctx.configSearchPaths = append(ctx.configSearchPaths, searchPaths...)
		}
	})
}

// ConfigFromFile load config from special file
// override AutoFindConfig
func ConfigFromFile(file string) Option {
	return optionFunc(func(ctx *ctxImpl) {
		ctx.configFile = file
	})
}

// ConfigFromBytes load config from special in bytes
func ConfigFromBytes(in []byte) Option {
	return optionFunc(func(ctx *ctxImpl) {
		ctx.configData = in[:]
	})
}

// ConfigEnvOverride enable load config from system env to override ,
// default is true and default prefix is defined by DefaultConfigEnvPrefix
func ConfigEnvOverride(enable bool, prefix ...string) Option {
	return optionFunc(func(ctx *ctxImpl) {
		ctx.enableEnvOverride = enable
		if len(prefix) > 0 && len(strings.Trim(prefix[0], " ")) > 0 {
			ctx.envPrefix = strings.Trim(prefix[0], " ")
		}
	})
}

type ConfigBinder func(unit Unit)

func WithConfigBinder(key string, binder ConfigBinder) Option {
	return optionFunc(func(ctx *ctxImpl) {

	})
}
