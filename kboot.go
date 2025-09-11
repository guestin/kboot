package kboot

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/guestin/log"
	"github.com/guestin/mob/msync"
	"github.com/ooopSnake/assert.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var _gCtx *ctxImpl

type (
	Context interface {
		GetAppName() string
		// GetViper get the viper instance
		GetViper() *viper.Viper
		Shutdown(err error)
	}
)

type ctxImpl struct {
	ctx               context.Context
	cancel            context.CancelFunc
	viper             *viper.Viper
	appName           string
	timezone          *time.Location
	logLevel          string
	hideBanner        bool
	rootLogger        *zap.Logger
	logger            log.ClassicLog
	units             []*unitImpl
	unitsInitRes      *sync.Map
	configName        string
	configFileType    string
	configSearchPaths []string
	configFile        string
	configData        []byte
	enableEnvOverride bool
	envPrefix         string
}

func (this *ctxImpl) GetAppName() string {
	return this.appName
}

func (this *ctxImpl) Shutdown(err error) {
	this.logger.Infof("shutdown app cause of : %v", err)
	this.kill()
}

func (this *ctxImpl) GetViper() *viper.Viper {
	return this.viper
}

func RegisterUnit(name string, fn InitFunc) {
	assert.Must(len(strings.TrimSpace(name)) != 0, "name must not empty or blank").Panic()
	assert.Must(fn != nil, "init func must not be nil").Panic()
	for _, u := range _gCtx.units {
		if u.GetName() == name {
			assert.Must(true, fmt.Sprintf("name '%s' already exist", name)).Panic()
		}
	}
	_gCtx.units = append(_gCtx.units, &unitImpl{
		rootCtx:  _gCtx,
		name:     name,
		initFunc: fn,
	})
}

func Bootstrap(ctx context.Context, opts ...Option) {
	if !_gCtx.hideBanner {
		fmt.Print(_BANNER)
	}
	_gCtx.logger.Info("bootstrap ...")
	assert.Must(ctx != nil, "root ctx must not be nil").Panic()
	_ctx, cancel := context.WithCancel(ctx)
	_gCtx.ctx = _ctx
	_gCtx.cancel = cancel

	for _, opt := range opts {
		opt.apply(_gCtx)
	}
	_gCtx.bootStrap()
}

func (this *ctxImpl) bootStrap() {
	defer func() {
		_ = this.rootLogger.Sync()
	}()
	err := this.loadConfig()
	if err != nil {
		this.logger.Fatal("load config err", zap.Error(err))
	}
	this.execute()
}

func (this *ctxImpl) execute() {
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
	runner := func(unitItem *unitImpl) {
		this.logger.With(
			log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Yellow, true))).
			Info("start init...")
		err := unitItem.Init(this)
		//
		if err != nil {
			this.logger.With(
				log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Red, true))).
				Fatal("init failed ,err : ", err)
			this.unitsInitRes.Store(unitItem.GetName(), NewBadResult(err))
			return
		}
		this.unitsInitRes.Store(unitItem.GetName(), NewSuccessResult())
		this.logger.With(
			log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Green, true))).
			Info("init success!")
		taskStack.PushFront(unitItem)
		group.AddTask(func() {
			defer func() {
				exitPanic := recover()
				if exitPanic != nil {
					this.logger.With(
						log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Red, true))).
						Panicf("exit unexpected, panic:%v", exitPanic)
				}
			}()
			this.logger.With(
				log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), log.Cyan, true))).
				Info("running...")
			result := unitItem.Exec()
			exitTagColor := log.Cyan
			var logMeth = this.logger.With(
				log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), exitTagColor, true))).Infof
			if result.Code != 0 {
				exitTagColor = log.Red
				logMeth = this.logger.With(
					log.UseSubTag(log.NewFixStyleText(unitItem.GetName(), exitTagColor, true))).Warnf
			}
			logMeth("exit, code : %d ,err: %v", result.Code, result.Error)
		})
	}
	for idx := range this.units {
		go runner(this.units[idx])
	}
	<-this.ctx.Done()
}

func (this *ctxImpl) kill() {
	this.cancel()
}

func (this *ctxImpl) handleKillSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		sig := <-c
		this.logger.Infof("receive signal : %v", sig)
		this.Shutdown(errors.New(fmt.Sprintf("system signal : %v", sig)))
	}()
}
