package kboot

import (
	"bytes"
	"container/list"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/guestin/log"
	"github.com/guestin/mob/msync"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _gCtx *_ctx

type (
	Context interface {
		Application
		GetApplication() Application
		// GetViper get the viper instance
		GetViper() *viper.Viper
		GetActivatedProfile() string
		GetRootLogger() *zap.Logger
		GetTaggedZapLogger(tag string, opt ...log.Opt) log.ZapLog
		GetTaggedLogger(tag string, opt ...log.Opt) log.ClassicLog
		UnmarshalSubConfig(key string, i interface{}, options ...CfgOption) (err error)
		Shutdown(err error)
	}
)

type _ctx struct {
	Application
	ctx               context.Context
	cancel            context.CancelFunc
	viper             *viper.Viper
	logLevel          string
	hideBanner        bool
	rootLogger        *zap.Logger
	logger            log.ZapLog
	units             []*unitImpl
	configName        string
	configFileType    string
	configSearchPaths []string
	configFile        string
	configData        []byte
	enableEnvOverride bool
	envPrefix         string
}

func (this *_ctx) GetApplication() Application {
	return this.Application
}

func (this *_ctx) GetRootLogger() *zap.Logger {
	return this.rootLogger
}

func (this *_ctx) GetTaggedZapLogger(tag string, opt ...log.Opt) log.ZapLog {
	return log.NewTaggedZapLogger(this.GetRootLogger(), tag, opt...)
}

func (this *_ctx) GetTaggedLogger(tag string, opt ...log.Opt) log.ClassicLog {
	return log.NewTaggedClassicLogger(this.GetRootLogger(), tag, opt...)
}

func (this *_ctx) UnmarshalSubConfig(key string, any interface{}, options ...CfgOption) (err error) {
	defer func() {
		exitPanic := recover()
		if exitPanic != nil {
			err = errors.Wrapf(err, "umarshal [%s] config  panic :%v", key, exitPanic)
		}
	}()
	v := this.GetViper()
	subV := v.Sub(key)
	if subV == nil {
		return nil
	}
	for _, opt := range options {
		opt.apply(v)
	}
	if err := subV.Unmarshal(any); err != nil {
		return errors.Wrapf(err, "parser [%s] config failed ", key)
	}
	if err := MValidator().Validate(any); err != nil {
		return errors.Wrapf(err, "invlid [%s] config ", key)
	}
	return nil
}

func (this *_ctx) Shutdown(err error) {
	this.logger.Info("shutdown ", zap.Error(err))
	this.kill()
}

func (this *_ctx) GetViper() *viper.Viper {
	return this.viper
}

func (this *_ctx) GetActivatedProfile() string {
	return this.viper.GetString(CfgKeyProfilesActive)
}

func (this *_ctx) bootStrap() {
	defer func() {
		_ = this.rootLogger.Sync()
	}()
	err := this.autoConfig()
	if err != nil {
		this.logger.Fatal("auto config failed", zap.Error(err))
		return
	}
	err = this.reinitLoggerIfNeeded()
	if err != nil {
		this.logger.Fatal("reinit logger failed", zap.Error(err))
		return
	}
	this.execute()
}

func (this *_ctx) autoConfig() error {
	this.logger.Info("Load config ...")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}
	// load raw bytes first
	if len(this.configData) > 0 {
		err := this.viper.ReadConfig(bytes.NewReader(this.configData))
		if err != nil {
			if !errors.Is(err, &viper.ConfigFileNotFoundError{}) {
				return errors.Wrap(err, "load main config error")
			}
		}
	}
	this.prepareConfigPath()
	// load main config from file
	err = this.viper.MergeInConfig()
	if err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return err
		}
		// Config file not found; ignore error
	}
	activeProfile := this.GetActivatedProfile()
	this.logger.Info("active profile ", zap.Any("activeProfile", activeProfile))
	exts := make([]string, 0)
	if this.configFileType != "" {
		exts = append(exts, this.configFileType)
	} else {
		exts = append(exts, viper.SupportedExts...)
	}
	finder := newConfigFinder(this.logger)
	// load other config
	for _, dir := range this.configSearchPaths {
		configList, err := finder.FindConfigs(dir, exts...)
		if err != nil {
			return err
		}
		for cfgName := range configList {
			cfg := configList[cfgName]
			if cfg.Default != nil && cfgName != DefaultConfigName {
				this.logger.Info("apply default config ", zap.String("config", cfgName))
				if err := this.applyConfig(cfg.Default.FilePath); err != nil {
					return err
				}
			}
			for i := range cfg.Profiles {
				cfgFile := cfg.Profiles[i]
				// apply special profile
				if strings.EqualFold(cfgFile.Profile, activeProfile) {
					this.logger.Info("apply config ", zap.String("config", cfgName), zap.String("profile", activeProfile))
					if err := this.applyConfig(cfgFile.FilePath); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (this *_ctx) applyConfig(file string) error {
	this.viper.SetConfigFile(file)
	err := this.viper.MergeInConfig()
	if err != nil {
		return errors.Wrapf(err, "load config %s error", file)
	}
	return nil
}

func (this *_ctx) prepareConfigPath() {
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
}

func (this *_ctx) reinitLoggerIfNeeded() error {
	oldLv := this.logLevel
	this.logLevel = this.viper.GetString(CfgKeyAppLogLevel)
	lv, err := zapcore.ParseLevel(this.logLevel)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid app.log.level %s", this.logLevel))
	}
	if oldLv != this.logLevel {
		_ = this.rootLogger.Sync()
		this.rootLogger = this.rootLogger.WithOptions(zap.IncreaseLevel(lv))
		this.logger = log.NewTaggedZapLogger(this.rootLogger, LoggerTag)
		this.logger.Info("logger level changed", zap.String("level", this.logLevel))
	}
	return nil
}

func (this *_ctx) execute() {
	if len(this.units) == 0 {
		this.logger.Warn("no unit to execute ,exit...")
		return
	}
	this.handleKillSignal()
	group := msync.NewAsyncTaskGroup()
	defer group.Wait()
	// execute units
	taskStack := list.New()
	defer func() {
		for taskStack.Len() != 0 {
			item := taskStack.Front()
			taskStack.Remove(item)
			taskItem := item.Value.(*unitImpl)
			taskItem.Cancel()
			taskItem.Wait()
		}
	}()

	initer := func(unitItem *unitImpl) {
		defer func() {
			exitPanic := recover()
			if exitPanic != nil {
				this.logger.With(
					log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Red, true))).
					Panic("init panic", zap.Any("error", exitPanic))
			}
		}()
		this.logger.With(
			log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Yellow, true))).
			Info("start init...")
		err := unitItem.Init(this)
		if err != nil {
			this.logger.With(
				log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Red, true))).
				Panic("init failed  : ", zap.Error(err))
			return
		}
		this.logger.With(
			log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Green, true))).
			Info("init success!")
		taskStack.PushFront(unitItem)
	}
	runner := func(unitItem *unitImpl) {
		group.AddTask(func() {
			defer func() {
				exitPanic := recover()
				if exitPanic != nil {
					this.logger.With(
						log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Red, true))).
						Panic("exit unexpected", zap.Any("error", exitPanic))
				}
			}()
			this.logger.With(
				log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Cyan, true))).
				Info("running...")
			result := unitItem.Exec()
			exitTagColor := log.Cyan
			var logMeth = this.logger.With(
				log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), exitTagColor, true))).Info
			if result.Code != 0 {
				exitTagColor = log.Red
				logMeth = this.logger.With(
					log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), exitTagColor, true))).Warn
			}
			logMeth("exit", zap.Int("code", result.Code), zap.Error(result.Error))
		})
	}
	// sort by dependencies
	// a depends on b , then b should be before a
	sort.SliceStable(this.units, func(i, j int) bool {
		for _, dep := range this.units[j].depends {
			if dep == this.units[i].GetName() {
				return true
			}
			found := false
			for idx := 0; idx < j; idx++ {
				if dep == this.units[idx].GetName() {
					found = true
					break
				}
			}
			if !found {
				return true
			}
		}
		return false
	})

	for idx := range this.units {
		initer(this.units[idx])
	}
	for idx := range this.units {
		runner(this.units[idx])
	}
	<-this.ctx.Done()
}

func (this *_ctx) kill() {
	this.cancel()
}

func (this *_ctx) handleKillSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		sig := <-c
		this.logger.Info("Receive signal", zap.Any("signal", sig))
		this.Shutdown(errors.New(fmt.Sprintf("System signal : %v", sig)))
	}()
}
