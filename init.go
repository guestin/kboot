package kboot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/guestin/log"
	"github.com/spf13/pflag"
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
		_gCtx = &ctxImpl{
			ctx:               context.Background(),
			viper:             viper.New(),
			appName:           DefaultAppName,
			timezone:          DefaultAppTz,
			logLevel:          DefaultLogLevel,
			hideBanner:        false,
			rootLogger:        rootLogger,
			logger:            log.NewTaggedClassicLogger(rootLogger, LoggerTag),
			units:             make([]*unitImpl, 0),
			unitsInitRes:      new(sync.Map),
			configName:        DefaultConfigName,
			configFileType:    DefaultConfigFileType,
			configSearchPaths: []string{DefaultConfigFilePath},
			configFile:        "",
			configData:        nil,
			enableEnvOverride: true,
			envPrefix:         DefaultConfigEnvPrefix,
		}
		_gCtx.viper.SetDefault(CfgKeyAppName, DefaultAppName)
		_gCtx.viper.SetDefault(CfgKeyAppTz, DefaultAppTz.String())
		_gCtx.viper.SetDefault(CfgKeyAppLogLevel, DefaultLogLevel)
	})
}

func (this *ctxImpl) loadConfig() error {
	this.logger.Info("Load config ...")
	if this.configFile != "" {
		this.viper.SetConfigFile(this.configFile)
	}
	if this.configName != "" {
		this.viper.SetConfigName(this.configName)
		if this.configFileType != "" {
			this.viper.SetConfigType(this.configFileType)
		}
		for pi := range this.configSearchPaths {
			this.viper.AddConfigPath(this.configSearchPaths[pi])
		}
	}
	if this.enableEnvOverride {
		this.viper.AutomaticEnv()
		if this.envPrefix != "" {
			this.viper.SetEnvPrefix(this.envPrefix)
		}
		this.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}
	err := this.viper.ReadInConfig()
	if err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			this.logger.Fatalf("load config from file '%s' error: %v", this.configFile, err.Error())
			return err
		}
		// Config file not found; ignore error
	}

	if len(this.configData) > 0 {
		err = this.viper.MergeConfig(bytes.NewReader(this.configData))
		if err != nil {
			if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
				return err
			}
		}
	}
	pflag.Parse()
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}
	this.appName = this.viper.GetString(CfgKeyAppName)
	tz := this.viper.GetString(CfgKeyAppTz)
	this.timezone = time.FixedZone(tz, 0)
	oldLv := this.logLevel
	this.logLevel = this.viper.GetString(CfgKeyAppLogLevel)
	lv, err := zapcore.ParseLevel(this.logLevel)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid app.log.level %s", this.logLevel))
	}
	if oldLv != this.logLevel {
		_ = this.rootLogger.Sync()
		this.rootLogger = this.rootLogger.WithOptions(zap.IncreaseLevel(lv))
		this.logger = log.NewTaggedClassicLogger(this.rootLogger, LoggerTag)
		this.logger.Infof("logger level changed to %s", this.logLevel)
	}
	return nil
}
