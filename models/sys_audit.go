package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type SysAudit struct {
	Model

	User    string        `json:"user" gorm:"column:user;comment:'用户'"`
	UserId  uint          `json:"userId" gorm:"column:user_id;comment:'用户'"`
	Ip      string        `json:"ip" gorm:"column:ip;comment:'请求ip'"`
	Method  string        `json:"method" gorm:"column:method;comment:'请求方法'"`
	Path    string        `json:"path" gorm:"column:path;comment:'请求路径'"`
	Status  int           `json:"status" gorm:"column:status;comment:'请求状态'"`
	Latency time.Duration `json:"latency" gorm:"column:latency;comment:'延迟'"`
	Body    string        `json:"body" gorm:"type:longtext;column:body;comment:'请求Body'"`
}

func AddSysAudit(data map[string]interface{}) error {
	sysAudit := SysAudit{
		User:    data["user"].(string),
		UserId:  data["user_id"].(uint),
		Ip:      data["ip"].(string),
		Method:  data["method"].(string),
		Path:    data["path"].(string),
		Status:  data["status"].(int),
		Latency: data["latency"].(time.Duration),
		Body:    data["body"].(string),
	}

	err := db.Create(&sysAudit).Error
	if err != nil {
		return err
	}
	return nil
}

func GetSysAudits(query map[string]interface{}, page int, pageSize int) ([]*SysAudit, error) {
	sysAudit := []*SysAudit{}
	var err error
	pageNum := (page - 1) * pageSize
	username := query["user"].(string)
	if len(username) > 0 {
		sysUser, err := GetSysUserByUsername(username)
		if err != nil {
			return sysAudit, err
		}
		if sysUser != nil {
			err = db.Model(&SysAudit{}).Where("user_id = ?", sysUser.ID).Offset(pageNum).Limit(pageSize).Find(&sysAudit).Error
		}
	} else {
		err = db.Model(&SysAudit{}).Offset(pageNum).Limit(pageSize).Find(&sysAudit).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return sysAudit, nil
}

func GetSysAuditCount(query map[string]interface{}) (int, error) {
	var err error
	count := 0
	username := query["user"].(string)
	if len(username) > 0 {
		sysUser, err := GetSysUserByUsername(username)
		if err != nil {
			return 0, err
		}
		if sysUser != nil {
			err = db.Model(&SysAudit{}).Where("user_id = ?", sysUser.ID).Count(&count).Error
		}
	} else {
		err = db.Model(&SysAudit{}).Count(&count).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}
