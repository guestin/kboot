package web

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/guestin/kboot"
	"github.com/guestin/kboot/web/kerrors"
	"github.com/guestin/kboot/web/mid"
	"github.com/guestin/mob/merrors"
	"github.com/guestin/mob/mvalidate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type _EchoValidator struct {
	v mvalidate.Validator
}

func (this *_EchoValidator) Validate(i interface{}) error {
	k := reflect.TypeOf(i).Kind()
	if k == reflect.Struct || (k == reflect.Ptr && reflect.ValueOf(i).Elem().Kind() == reflect.Struct) {
		return this.v.Struct(i)
	}
	if k == reflect.Slice {
		return this.v.Var(i, "required,dive,required")
	}
	return this.v.Var(i, "required")
}

func _initEcho(unit kboot.Unit, cfg config) (kboot.ExecFunc, error) {
	ctx := unit.GetContext()
	eCtx := echo.New()
	eCtx.HideBanner = true
	eCtx.HidePort = false
	eCtx.DisableHTTP2 = true
	eCtx.HTTPErrorHandler = globalErrorHandle
	mValidator, err := mvalidate.NewValidator(mvalidate.DefaultTranslator)
	if err != nil {
		return nil, err
	}
	eCtx.Validator = &_EchoValidator{mValidator}

	eCtx.Use(mid.WithContext(ctx))
	//recovery
	eCtx.Use(middleware.Recover())
	// request id
	eCtx.Use(middleware.RequestID())
	//cors
	eCtx.Use(middleware.CORS())
	//gzip
	eCtx.Use(middleware.Gzip())
	//logger
	loggerOption := mid.LogNone
	if cfg.Debug {
		loggerOption = mid.LogReqHeader | mid.LogRespHeader | mid.LogJson | mid.LogForm
	}
	eCtx.Use(mid.Log(loggerOption))

	// start watcher
	exitChan := make(chan error)

	return func(unit kboot.Unit) kboot.ExitResult {
		for _, opt := range _options {
			err = opt.apply(eCtx)
			if err != nil {
				logger.Panic("use option failed ", err)
			}
		}
		go func() {
			exitChan <- eCtx.Start(cfg.ListenAddress)
		}()
		select {
		case err = <-exitChan:
			logger.Info("API server exit", zap.Error(err))
			return kboot.NewSuccessResult()
		case <-unit.Done():
			_ = eCtx.Shutdown(unit.GetContext())
		}
		return kboot.NewSuccessResult()
	}, nil
}

func globalErrorHandle(err error, ctx echo.Context) {
	if err == nil {
		return
	}
	errCategory := uint8(0) // means default
	switch err.(type) {
	case merrors.Error:
		errCategory = 1
		_ = ctx.JSON(http.StatusOK, err)
		// code = 0, means no error
		if err.(merrors.Error).GetCode() == 0 {
			return
		}
	case validator.ValidationErrors, mvalidate.ValidateError:
		errCategory = 2
		_ = ctx.JSON(http.StatusOK,
			merrors.ErrorWrap0(err, kerrors.CodeBadRequest, "bad request params"))
	case *echo.HTTPError:
		errCategory = 3
		he := err.(*echo.HTTPError)
		_ = ctx.JSON(http.StatusOK,
			merrors.Errorf0(kerrors.HttpStatus2Code(he.Code), "%s", fmt.Sprint(he.Message)))
	case error:
		errCategory = 4
		_ = ctx.JSON(http.StatusOK,
			merrors.ErrorWrap0(err, kerrors.CodeInternalServer, "unexpect error"))
	default:
		ctx.Echo().DefaultHTTPErrorHandler(err, ctx)
	}
	logger.Warn("api global error handler",
		zap.String("path", ctx.Path()),
		zap.Uint8("errCategory", errCategory),
		zap.Error(err))
}
