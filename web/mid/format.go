package mid

import (
	"bytes"
	"fmt"

	"github.com/guestin/kboot/web/internal"
	"github.com/guestin/kboot/web/kerrors"
	"github.com/guestin/mob/merrors"
	"github.com/labstack/echo/v4"

	"net/http"
	"runtime"
)

type (
	// FormatConfig defines the config for Format middleware.
	FormatConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper Skipper
	}
)

var (
	// DefaultFormatConfig is the default Format middleware config.
	DefaultFormatConfig = FormatConfig{
		Skipper: DefaultSkipper,
	}
)

func Format() echo.MiddlewareFunc {
	return FormatWithConfig(DefaultFormatConfig)
}
func FormatWithConfig(config FormatConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultFormatConfig.Skipper
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			defer func() {
				err := recover()
				if err != nil {
					internal.Log.Errorf("panic recover : \n%s\n", PanicTrace(PanicTraceSizeKb))
					ctx.Set(internal.CtxStatusKey, http.StatusInternalServerError)
					ctx.Set(internal.CtxErrorKey, kerrors.InternalErr("Server is busy"))
				}
				checkErrAndFlush(ctx, config)
			}()
			ctxErr := ctx.Get(internal.CtxErrorKey)
			if ctxErr != nil {
				return nil
			}
			err := next(ctx)
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					ctx.Set(internal.CtxStatusKey, he.Code)
				} else if _, ok := err.(merrors.Error); ok {
					ctx.Set(internal.CtxStatusKey, http.StatusOK)
				} else {
					ctx.Set(internal.CtxStatusKey, http.StatusInternalServerError)
				}
				ctx.Set(internal.CtxErrorKey, err)
			}
			return nil
		}
	}
}
func checkErrAndFlush(ctx echo.Context, config FormatConfig) {
	//requestId := ctx.Request().Header.Get(echo.HeaderXRequestID)
	statusCode := 200
	statusI := ctx.Get(internal.CtxStatusKey)
	if statusI != nil {
		if status, ok := statusI.(int); ok && status > 0 {
			statusCode = status
		}
	}
	var resp interface{} = nil
	ctxErr := ctx.Get(internal.CtxErrorKey)
	if ctxErr != nil {
		if config.Skipper(ctx) {
			resp = ctxErr
			goto flush
		}
		switch ctxErr.(type) {
		case merrors.Error:
			resp = ctxErr
		case *echo.HTTPError:
			he := ctxErr.(*echo.HTTPError)
			resp = merrors.Errorf0(he.Code, "%s", fmt.Sprint(he.Message))
		default:
			resp = kerrors.ErrInternalf("Server is busy : %s", ctxErr)
		}
		goto flush
	}
	resp = ctx.Get(internal.CtxRespKey)
	if config.Skipper(ctx) {
		goto flush
	}
	if resp != nil {
		switch resp.(type) {
		case merrors.Error:
			goto flush
		default:
			resp = kerrors.OkResult(resp)
		}
	} else {
		resp = kerrors.OkResult(nil)
	}

flush:
	// has rsp & no error need write response,otherwise err handler will handle
	if !ctx.Response().Committed && resp != nil {
		_ = ctx.JSON(statusCode, resp)
		//_ = ctx.JSONPretty(statusCode, resp,jsonIndent)
	}
}

const PanicTraceSizeKb = 2

func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}
