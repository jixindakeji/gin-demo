package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type SysMenu struct {
	Model

	ParentId uint   `json:"parentId" gorm:"column:parent_id; default:'0'"`
	Name     string `json:"name" gorm:"column:name; default:''"`
	Url      string `json:"url" gorm:"column:url;default:''"`
	Type     int    `json:"type" gorm:"column:type;default:'0'"`
	Code     string `json:"code" gorm:"column:code;default:''"`
	Icon     string `json:"icon" gorm:"icon"`
	Sort     int    `json:"sort" gorm:"column:sort;default:'1'"`
	Remark   string `json:"remark" gorm:"column:remark;default:''"`
	Status   int    `json:"status" gorm:"column:status;default:0"`
	// Role     []SysRole `json:"-" gorm:"many2many:role_menu;"`

	Children []*SysMenu `json:"children"`
}

func AddSysMenu(data map[string]interface{}) error {
	sysMenu := SysMenu{
		ParentId: data["parent_id"].(uint),
		Name:     data["name"].(string),
		Url:      data["url"].(string),
		Type:     data["type"].(int),
		Code:     data["code"].(string),
		Icon:     data["icon"].(string),
		Sort:     data["sort"].(int),
		Remark:   data["remark"].(string),
	}

	err := db.Create(&sysMenu).Error
	if err != nil {
		return err
	}
	return nil
}

func GetSysMenu(id uint) (*SysMenu, error) {
	var sysMenu SysMenu
	if err := db.Where("id = ? ", id).Find(&sysMenu).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &sysMenu, nil
}

func GetSysMenuChildren(id uint) ([]*SysMenu, error) {
	var sysMenus []*SysMenu
	err := db.Where("parent_id = ?", id).Find(&sysMenus).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	for _, v := range sysMenus {
		items, err := GetSysMenuChildren(v.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		for _, v1 := range items {
			v.Children = append(v.Children, v1)
		}
	}
	return sysMenus, nil
}

func UpdateSysMenu(id uint, data map[string]interface{}) error {
	if err := db.Model(&SysMenu{}).Where("id = ?", id).Update(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteSysMenu(id uint) error {
	var count int
	if err := db.Model(&SysMenu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if count > 0 {
		/* exist child menu link */
		return errors.New("record exists inner link")
	}

	var roles []SysRole
	err := db.Preload("Menu", "id=?", id).Find(&roles).Error
	if err != nil {
		return err
	}
	for _, role := range roles {
		if len(role.Menu) > 0 {
			/* exist role link */
			return errors.New("role has linked with user")
		}
	}

	/* now, delete menu */
	err = db.Where("id = ?", id).Delete(SysMenu{}).Error
	if err != nil {
		return err
	}
	return nil
}
