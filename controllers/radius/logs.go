package radius

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/talkincode/logsight/app"
	"github.com/talkincode/logsight/common"
	"github.com/talkincode/logsight/common/web"
	"github.com/talkincode/logsight/models"
	"github.com/talkincode/logsight/webserver"
)

func InitLogsRouter() {

	// 认证日志页面展示 assets/templates/radius_authlog.html
	webserver.GET("/admin/radius/authlog", func(c echo.Context) error {
		return c.Render(http.StatusOK, "radius_authlog", map[string]interface{}{})
	})

	// 记账日志页面展示 assets/templates/radius_accounting.html
	webserver.GET("/admin/radius/accounting", func(c echo.Context) error {
		return c.Render(http.StatusOK, "radius_accounting", nil)
	})

	// 记账日志查询
	webserver.GET("/admin/radius/accounting/query", func(c echo.Context) error {
		var count, start int
		web.NewParamReader(c).
			ReadInt(&start, "start", 0).
			ReadInt(&count, "count", 40)
		var data []models.TsRadiusAccounting
		prequery := web.NewPreQuery(c).
			DefaultOrderBy("acct_stop_time desc").
			DateRange2("starttime", "endtime", "acct_stop_time", time.Now().Add(-time.Hour*8), time.Now()).
			KeyFields("username", "framed_ipaddr", "mac_addr")

		var total int64
		common.Must(prequery.Query(app.GDB().Model(&models.TsRadiusAccounting{})).Count(&total).Error)

		query := prequery.Query(app.GDB().Debug().Model(&models.TsRadiusAccounting{})).Offset(start).Limit(count)
		if query.Find(&data).Error != nil {
			return c.JSON(http.StatusOK, common.EmptyList)
		}
		return c.JSON(http.StatusOK, &web.PageResult{TotalCount: total, Pos: int64(start), Data: data})
	})

	webserver.POST("/radius/accounting/add", func(c echo.Context) error {
		form := new(models.TsRadiusAccounting)
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = json.Unmarshal(body, form)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		acct := models.TsRadiusAccounting{}
		err = app.GDB().Debug().
			Where("username = ? and acct_session_id = ?",
				form.Username, form.AcctSessionId,
			).First(&acct).Error
		if err == nil {
			// update
			common.Must(app.GDB().
				Model(&models.TsRadiusAccounting{}).
				Where("username = ? and acct_session_id = ?",
					form.Username, form.AcctSessionId).Updates(map[string]interface{}{
				"acct_session_time":   form.AcctSessionTime,
				"acct_input_total":    form.AcctInputTotal,
				"acct_output_total":   form.AcctOutputTotal,
				"acct_input_packets":  form.AcctInputPackets,
				"acct_output_packets": form.AcctOutputPackets,
				"acct_stop_time":      form.AcctStopTime,
				"session_timeout":     form.SessionTimeout,
				"last_update":         time.Now(),
			}).Error)
		} else {
			form.ID = common.UUID()
			form.LastUpdate = time.Now()
			common.Must(app.GDB().Create(form).Error)
		}
		return c.JSON(http.StatusOK, web.RestSucc("success"))
	})

}
