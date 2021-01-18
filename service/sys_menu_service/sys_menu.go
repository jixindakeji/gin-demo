package sys_menu_service

import (
	"gin-demo/models"
)

type SysMenu struct {
	ID uint

	ParentId uint
	Name     string
	Url      string
	Type     int
	Code     string
	Icon     string
	Sort     int
	Remark   string
	Status   int
}

func (m *SysMenu) Add() error {
	data := make(map[string]interface{})
	data["parent_id"] = m.ParentId
	data["name"] = m.Name
	data["url"] = m.Url
	data["type"] = m.Type
	data["code"] = m.Code
	data["icon"] = m.Icon
	data["sort"] = m.Sort
	data["remark"] = m.Remark
	data["status"] = m.Status

	return models.AddSysMenu(data)
}

type menuItem struct {
	ID       uint   `json:"id"`
	ParentId uint   `json:"parent_id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Type     int    `json:"type"`
	Code     string `json:"code"`
	Icon     string `json:"icon"`
	Sort     uint   `json:"sort"`
	Remark   string `json:"remark"`

	Children []*menuItem `json:"children"`
}

func (m *SysMenu) GetAll() ([]*models.SysMenu, error) {
	data := make(map[string]interface{})
	data["name"] = m.Name

	menus, err := models.GetSysMenuChildren(0)
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (m *SysMenu) Update() error {
	id := m.ID
	data := make(map[string]interface{})
	data["parent_id"] = m.ParentId
	data["name"] = m.Name
	data["url"] = m.Url
	data["type"] = m.Type
	data["code"] = m.Code
	data["icon"] = m.Icon
	data["sort"] = m.Sort
	data["remark"] = m.Remark

	return models.UpdateSysMenu(id, data)
}

func (m *SysMenu) Delete() error {
	return models.DeleteSysMenu(m.ID)
}

func (m *SysMenu) Get() (*models.SysMenu, error) {
	return models.GetSysMenu(m.ID)
}
