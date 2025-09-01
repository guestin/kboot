package apis

import (
	"github.com/guestin/kboot/web"
	"github.com/guestin/kboot/web/mid"
	"github.com/labstack/echo/v4"
)

func init() {
	web.Router(routeBuilder)
	web.DependsOn("db")
}

func routeBuilder(eCtx *echo.Echo) error {
	eCtx.POST("/echo", mid.Wrap(Echo))
	return nil
}
