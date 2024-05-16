package models

import (
	"time"
)

type TsRadiusAccounting struct {
	ID                string    `json:"id" gorm:"primaryKey"` // 主键 ID
	Username          string    `json:"username" gorm:"primaryKey"`
	AcctSessionId     string    `json:"acct_session_id" gorm:"primaryKey"`
	AcctStartTime     time.Time `json:"acct_start_time" gorm:"primaryKey"`
	NasId             string    `json:"nas_id"`
	NasAddr           string    `json:"nas_addr"`
	NasPaddr          string    `json:"nas_paddr"`
	SessionTimeout    int       `json:"session_timeout"`
	FramedIpaddr      string    `json:"framed_ipaddr"`
	FramedNetmask     string    `json:"framed_netmask"`
	MacAddr           string    `json:"mac_addr"`
	NasPort           int64     `json:"nas_port,string"`
	NasClass          string    `json:"nas_class"`
	NasPortId         string    `json:"nas_port_id"`
	NasPortType       int       `json:"nas_port_type"`
	ServiceType       int       `json:"service_type"`
	AcctSessionTime   int       `json:"acct_session_time"`
	AcctInputTotal    int64     `json:"acct_input_total,string"`
	AcctOutputTotal   int64     `json:"acct_output_total,string"`
	AcctInputPackets  int       `json:"acct_input_packets"`
	AcctOutputPackets int       `json:"acct_output_packets"`
	AcctStopTime      time.Time `json:"acct_stop_time"`
	LastUpdate        time.Time `json:"last_update"`
}

type TsSyslog struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	Timestamp       time.Time `json:"timestamp" gorm:"primaryKey"`
	Logtype         string    `json:"logtype"`
	MsgID           string    `json:"msg_id,omitempty"`
	ProcID          string    `json:"proc_id,omitempty"`
	Appname         string    `json:"appname,omitempty"`
	Hostname        string    `json:"hostname,omitempty"`
	Priority        int64     `json:"priority,omitempty"`
	Facility        int64     `json:"facility,omitempty"`
	FacilityMessage string    `json:"facility_message,omitempty"`
	Severity        int64     `json:"severity,omitempty"`
	SeverityMessage string    `json:"severity_message,omitempty"`
	Version         int64     `json:"version,omitempty"`
	Message         string    `json:"message"`
	Tags            string    `json:"tags,omitempty"`
}
