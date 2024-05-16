package models

import (
	"encoding/json"
	"time"
)

func (d SysOprLog) MarshalJSON() ([]byte, error) {
	type Alias SysOprLog
	return json.Marshal(&struct {
		Alias
		OptTime string `json:"opt_time"`
	}{
		Alias:   (Alias)(d),
		OptTime: d.OptTime.Format(time.RFC3339),
	})
}
