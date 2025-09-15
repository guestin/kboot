package apis

import (
	"github.com/guestin/kboot"
	"github.com/guestin/kboot/web"
	"github.com/guestin/kboot/web/mid"
	"github.com/guestin/log"
	"github.com/labstack/echo/v4"
)

var logger log.ClassicLog
var zapLogger log.ZapLog

func init() {
	web.DependsOn("api")
	web.Router(routeBuilder)
	kboot.RegisterUnit("api", _init, kboot.DependsOn("db"))
}
func _init(unit kboot.Unit) (kboot.ExecFunc, error) {
	logger = unit.GetClassicLogger()
	zapLogger = unit.GetZapLogger()

	return func(unit kboot.Unit) kboot.ExitResult {
		<-unit.Done()
		return kboot.ExitResult{
			Code:  0,
			Error: nil,
		}
	}, nil
}

func routeBuilder(eCtx *echo.Echo) error {
	eCtx.POST("/echo", mid.Wrap(Echo))
	return nil
}
