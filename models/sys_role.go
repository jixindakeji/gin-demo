package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type SysRole struct {
	Model

	Name   string `json:"name" gorm:"column:name"`
	Remark string `json:"remark" gorm:"column:remark"`
	Status string `json:"status" gorm:"column:status"`

	Menu []SysMenu `gorm:"many2many:role_menu;"`
}

func AddSysRole(data map[string]interface{}) error {
	sysRole := SysRole{
		Name:   data["name"].(string),
		Remark: data["remark"].(string),
	}

	if err := db.Create(&sysRole).Error; err != nil {
		return err
	}

	return nil
}

func GetSysRole(id uint) (*SysRole, error) {
	var sysRole SysRole
	if err := db.Where("id = ?", id).Find(&sysRole).Error; err != nil {
		return nil, err
	}
	return &sysRole, nil
}

func UpdateSysRole(id uint, data map[string]interface{}) error {
	if err := db.Model(&SysRole{}).Where("id = ?", id).Update(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteSysRole(id uint) error {
	var users []SysUser
	err := db.Preload("Role", "id=?", id).Find(&users).Error
	if err != nil {
		return err
	}
	for _, user := range users {
		if len(user.Role) > 0 {
			/* exist user link */
			return errors.New("role has linked with user")
		}
	}

	role := SysRole{}
	role.Model.ID = id
	/* delete menu linked with this role */
	tx := db.Begin()
	err = db.Model(&role).Association("Menu").Clear().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	/*now, delete role */
	err = db.Where("id = ?", id).Delete(SysRole{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func GetSysRoles(query map[string]interface{}, page int, pageSize int) ([]*SysRole, error) {
	var sysRoles []*SysRole
	var err error
	pageNum := (page - 1) * pageSize
	name := query["name"].(string)
	if len(name) > 0 {
		name = "%" + name + "%"
	}
	if len(name) > 0 {
		err = db.Where("name like ? ", name).Offset(pageNum).Limit(pageSize).Find(&sysRoles).Error
	} else {
		err = db.Offset(pageNum).Limit(pageSize).Find(&sysRoles).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return sysRoles, nil
}

func GetSysRoleCount(query map[string]interface{}) (int, error) {
	var err error
	count := 0
	name := query["name"].(string)
	if len(name) > 0 {
		name = "%" + name + "%"
	}
	if len(name) > 0 {
		err = db.Model(&SysRole{}).Where("name like ? ", name).Count(&count).Error
	} else {
		err = db.Model(&SysRole{}).Count(&count).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}

func UpdateSysRoleMenu(data map[string]interface{}) error {
	role := SysRole{}
	role.Model.ID = data["role_id"].(uint)
	menus := data["menu_ids"].([]uint)
	var menu []*SysMenu
	for _, id := range menus {
		temp := SysMenu{}
		temp.Model.ID = id
		menu = append(menu, &temp)
	}

	err := db.Model(&role).Association("menu").Replace(menu).Error
	if err != nil {
		return err
	}
	return nil
}

func GetSysMenuByRole(id uint) ([]*SysMenu, error) {
	role := SysRole{}
	role.Model.ID = id
	menu := []*SysMenu{}
	err := db.Model(&role).Association("menu").Find(&menu).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return menu, nil
}
