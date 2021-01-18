package sys_audit_service

import (
	"gin-demo/models"
	"time"
)

type SysAudit struct {
	User    string
	UserId  uint
	Ip      string
	Method  string
	Path    string
	Status  int
	Body    string
	Resp    string
	Latency time.Duration

	Page     int
	PageSize int
}

func (a *SysAudit) Add() error {
	data := make(map[string]interface{})
	data["user"] = a.User
	data["user_id"] = a.UserId
	data["ip"] = a.Ip
	data["method"] = a.Method
	data["path"] = a.Path
	data["status"] = a.Status
	data["body"] = a.Body
	data["latency"] = a.Latency

	return models.AddSysAudit(data)
}

func (a *SysAudit) GetAll() ([]*models.SysAudit, error) {
	query := make(map[string]interface{})
	query["user"] = a.User

	pageSize := a.PageSize
	page := a.Page

	return models.GetSysAudits(query, page, pageSize)
}

func (a *SysAudit) Count() (int, error) {
	query := make(map[string]interface{})
	query["user"] = a.User

	return models.GetSysAuditCount(query)
}
