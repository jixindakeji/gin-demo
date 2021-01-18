package sys_role_service

import "gin-demo/models"

type SysRole struct {
	ID uint

	Name    string
	Remark  string
	MenuIds []uint

	Page     int
	PageSize int
}

func (r *SysRole) Add() error {
	data := make(map[string]interface{})
	data["name"] = r.Name
	data["remark"] = r.Remark

	return models.AddSysRole(data)
}

func (r *SysRole) Update() error {
	data := make(map[string]interface{})
	data["name"] = r.Name
	data["remark"] = r.Remark

	return models.UpdateSysRole(r.ID, data)
}

func (r *SysRole) Delete() error {
	return models.DeleteSysRole(r.ID)
}

func (r *SysRole) GetAll() ([]*models.SysRole, error) {
	query := make(map[string]interface{})
	query["name"] = r.Name

	return models.GetSysRoles(query, r.Page, r.PageSize)
}

func (r *SysRole) Get() (*models.SysRole, error) {
	return models.GetSysRole(r.ID)
}

func (r *SysRole) Count() (int, error) {
	query := make(map[string]interface{})
	query["name"] = r.Name
	query["remark"] = r.Remark

	return models.GetSysRoleCount(query)
}

func (r *SysRole) UpdateMenu() error {
	data := make(map[string]interface{})
	data["role_id"] = r.ID
	data["menu_ids"] = r.MenuIds

	return models.UpdateSysRoleMenu(data)
}

func (r *SysRole) GetMenu() ([]uint, error) {
	var data []uint
	menu, err := models.GetSysMenuByRole(r.ID)
	if err != nil {
		return data, err
	}
	for _, item := range menu {
		data = append(data, item.Model.ID)
	}
	return data, nil
}
