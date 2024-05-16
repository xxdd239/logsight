package app

import (
	"time"

	"github.com/talkincode/logsight/common"
	"github.com/talkincode/logsight/models"
)

func (a *Application) InitTest() {
	a.initTestSettings()
	a.initTestOpr()
}

func (a *Application) initTestSettings() {
	a.gormDB.Where("1 = 1").Delete(&models.SysConfig{})
	a.gormDB.Create(&models.SysConfig{ID: 0, Sort: 1, Type: "system", Name: "SystemTitle", Value: "ToughRADIUS management system", Remark: "System title"})
	a.gormDB.Create(&models.SysConfig{ID: 0, Sort: 1, Type: "system", Name: "SystemTheme", Value: "light", Remark: "System theme"})
	a.gormDB.Create(&models.SysConfig{ID: 0, Sort: 3, Type: "system", Name: "SystemLoginRemark", Value: "Recommended browser: Chrome/Edge", Remark: "Login page description"})
	a.gormDB.Create(&models.SysConfig{ID: 0, Sort: 3, Type: "system", Name: "SystemLoginSubtitle", Value: "ToughRADIUS community edition", Remark: "Login form title"})
}

func (a *Application) initTestOpr() {
	a.gormDB.Where("1 = 1").Delete(&models.SysOpr{})
	a.gormDB.Create(&models.SysOpr{
		ID:        common.UUIDint64(),
		Realname:  "管理员",
		Mobile:    "0000",
		Email:     "N/A",
		Username:  "admin",
		Password:  common.Sha256HashWithSalt("logsight", common.SecretSalt),
		Level:     "super",
		Status:    "enabled",
		Remark:    "super",
		LastLogin: time.Now(),
	})
}
