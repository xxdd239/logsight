package logs

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/talkincode/logsight/app"
	"github.com/talkincode/logsight/common"
	"github.com/talkincode/logsight/common/web"
	"github.com/talkincode/logsight/models"
	"github.com/talkincode/logsight/webserver"
)

func initSyslogRouter() {
	webserver.GET("/admin/syslog", func(c echo.Context) error {
		return c.Render(http.StatusOK, "syslog", map[string]interface{}{})
	})

	webserver.GET("/admin/syslog/query", func(c echo.Context) error {
		var count, start int
		web.NewParamReader(c).
			ReadInt(&start, "start", 0).
			ReadInt(&count, "count", 40)
		var data []models.TsSyslog
		prequery := web.NewPreQuery(c).
			DefaultOrderBy("timestamp desc").
			QueryField("hostname", "hostname").
			DateRange2("starttime", "endtime", "timestamp", time.Now().Add(-time.Hour*8), time.Now()).
			KeyFields("message")

		var total int64
		common.Must(prequery.Query(app.GDB().Model(&models.TsSyslog{})).Count(&total).Error)

		query := prequery.Query(app.GDB().Debug().Model(&models.TsSyslog{})).Offset(start).Limit(count)
		if query.Find(&data).Error != nil {
			return c.JSON(http.StatusOK, common.EmptyList)
		}
		return c.JSON(http.StatusOK, &web.PageResult{TotalCount: total, Pos: int64(start), Data: data})
	})
}
