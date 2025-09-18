package kboot

import (
	"context"
	"sync/atomic"

	"github.com/guestin/log"
	"go.uber.org/zap"
)

type (
	Unit interface {
		GetGlobalContext() Context
		GetContext() context.Context
		GetName() string
		// GetRootLogger get the root zap logger
		GetRootLogger() *zap.Logger
		// GetClassicLogger get current unit classic logger
		GetClassicLogger() log.ClassicLog
		// GetZapLogger get current unit zap logger
		GetZapLogger() log.ZapLog

		UnMarshalConfig(key string, any interface{}) error
		Depends(dep ...string)
		// Done wait for  application exit
		Done() <-chan struct{}
	}

	ExitResult struct {
		Code  int
		Error error
	}
	ExecFunc func(unit Unit) ExitResult
	InitFunc func(unit Unit) (ExecFunc, error)
)

type unitImpl struct {
	rootCtx    *ctxImpl
	ctx        context.Context
	name       string
	initFunc   InitFunc
	cancelFunc context.CancelFunc
	exeFunc    ExecFunc
	done       chan struct{}
	closeOnce  uint32
	logger     log.ClassicLog
	zapLogger  log.ZapLog
	depends    []string
}

func (this *unitImpl) UnMarshalConfig(key string, any interface{}) error {
	v := this.GetGlobalContext().GetViper()
	err := v.UnmarshalKey(key, any)
	if err != nil {
		return err
	}
	return MValidator().Validate(v)
}

func (this *unitImpl) Depends(dep ...string) {
	if len(dep) == 0 {
		return
	}
	this.depends = append(this.depends, dep...)
}

func (this *unitImpl) GetGlobalContext() Context {
	return this.rootCtx
}

func (this *unitImpl) GetContext() context.Context {
	return this.ctx
}

func (this *unitImpl) GetName() string {
	return this.name
}

func (this *unitImpl) GetRootLogger() *zap.Logger {
	return this.rootCtx.rootLogger
}

func (this *unitImpl) GetClassicLogger() log.ClassicLog {
	return this.logger
}

func (this *unitImpl) GetZapLogger() log.ZapLog {
	return this.zapLogger
}

func (this *unitImpl) Done() <-chan struct{} {
	return this.ctx.Done()
}

func (this *unitImpl) Wait() {
	<-this.done
}

func (this *unitImpl) HasExecFunc() bool {
	return this.exeFunc != nil
}

func (this *unitImpl) Exec() ExitResult {
	defer func() {
		if this.done != nil && atomic.CompareAndSwapUint32(&this.closeOnce, 0, 1) {
			close(this.done)
		}
	}()
	if !this.HasExecFunc() {
		<-this.ctx.Done()
		return NewSuccessResult()
	}
	return this.exeFunc(this)
}

func (this *unitImpl) Cancel() {
	this.cancelFunc()
}

func (this *unitImpl) Init(rootCtx *ctxImpl) error {
	ctx, cancelFunc := context.WithCancel(rootCtx.ctx)
	this.rootCtx = rootCtx
	this.ctx = ctx
	this.cancelFunc = cancelFunc
	this.logger = log.NewTaggedClassicLogger(rootCtx.rootLogger, this.GetName())
	this.zapLogger = log.NewTaggedZapLogger(rootCtx.rootLogger, this.GetName())
	this.done = make(chan struct{})
	exeFunc, err := this.initFunc(this)
	if err != nil {
		this.Cancel()
		return err
	}
	this.exeFunc = exeFunc
	return nil
}
