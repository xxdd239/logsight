package dashboard

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/talkincode/logsight/webserver"
)

func InitRouter() {
	webserver.GET("/admin/sysstatus", func(c echo.Context) error {
		return c.Render(http.StatusOK, "sysstatus", map[string]string{})
	})
	webserver.GET("/admin/overview", func(c echo.Context) error {
		return c.Render(http.StatusOK, "overview", map[string]string{})
	})

	initSystemMetricsRouter()
}
