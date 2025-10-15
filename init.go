package kboot

import (
	"context"
	"sync"

	"github.com/guestin/log"
	"github.com/spf13/viper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _initOnce = &sync.Once{}

func init() {
	// init logger
	lv, err := zapcore.ParseLevel(DefaultLogLevel)
	if err != nil {
		panic(err)
	}
	rootLogger, _ := log.EasyInitConsoleLogger(lv, zap.ErrorLevel)
	_initOnce.Do(func() {
		_gCtx = &_ctx{
			ctx:               context.Background(),
			viper:             viper.New(),
			logLevel:          DefaultLogLevel,
			hideBanner:        false,
			rootLogger:        rootLogger,
			logger:            log.NewTaggedZapLogger(rootLogger, LoggerTag),
			units:             make([]*unitImpl, 0),
			configName:        DefaultConfigName,
			configFileType:    "",
			configSearchPaths: []string{DefaultConfigFilePath},
			configFile:        "",
			configData:        nil,
			enableEnvOverride: true,
			envPrefix:         DefaultConfigEnvPrefix,
		}
		_gCtx.viper.SetDefault(CfgKeyAppTz, DefaultAppTz.String())
		_gCtx.viper.SetDefault(CfgKeyAppLogLevel, DefaultLogLevel)
	})
}
