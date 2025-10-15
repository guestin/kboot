package kboot

import (
	"strings"
)

type BootOption Option[*_ctx]

// AutoFindConfig auto find config file (*.fileType) from special search paths
// will override by ConfigFromFile
func AutoFindConfig(fileType string, searchPaths ...string) BootOption {
	return optionFunc[*_ctx](func(ctx *_ctx) {
		if len(strings.Trim(fileType, " ")) > 0 {
			ctx.configFileType = strings.Trim(fileType, " ")
		}
		if len(searchPaths) > 0 {
			ctx.configSearchPaths = append(ctx.configSearchPaths, searchPaths...)
		}
		ctx.configFile = ""
	})
}

// ConfigFromFile load config from special file
// override AutoFindConfig
func ConfigFromFile(file string) BootOption {
	return optionFunc[*_ctx](func(ctx *_ctx) {
		ctx.configFile = file
		ctx.configName = ""
	})
}

// ConfigFromBytes load config from special in bytes
func ConfigFromBytes(in []byte) BootOption {
	return optionFunc[*_ctx](func(ctx *_ctx) {
		ctx.configData = in[:]
	})
}

// ConfigEnvOverride enable load config from system env to override ,
// default is true and default prefix is defined by DefaultConfigEnvPrefix
func ConfigEnvOverride(enable bool, prefix ...string) BootOption {
	return optionFunc[*_ctx](func(ctx *_ctx) {
		ctx.enableEnvOverride = enable
		if len(prefix) > 0 && len(strings.Trim(prefix[0], " ")) > 0 {
			ctx.envPrefix = strings.Trim(prefix[0], " ")
		}
	})
}
