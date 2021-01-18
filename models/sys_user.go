package models

import (
	"errors"
	"gin-demo/pkg/util"
	"github.com/jinzhu/gorm"
)

type SysUser struct {
	Model

	Username    string `json:"username" gorm:"column:username;not null;unique;comment:'用户名'"`
	DisplayName string `json:"displayName" gorm:"column:display_name;comment:'昵称'"`
	Password    string `json:"-" gorm:"column:password;not null;comment:'密码'"`
	Salt        string `json:"-" gorm:"column:salt;not null;comment:'密码盐'"`
	Email       string `json:"email" gorm:"column:email;comment:'邮箱'"`
	Phone       string `json:"phone" gorm:"column:phone;comment:'手机'"`
	Status      int    `json:"status" gorm:"column:status;default:'0';comment:'状态: 1-admin,2-normal,3-locked'"`
	Remark      string `json:"remark" gorm:"column:remark;comment:'备注'"`

	Role []SysRole `gorm:"many2many:user_role;"`
}

func GetSysUserByUsername(username string) (*SysUser, error) {
	var user SysUser
	err := db.Where(SysUser{Username: username}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if user.ID > 0 {
		return &user, nil
	}

	return nil, nil
}

func AddSysUser(data map[string]interface{}) error {
	salt, password := util.GetSaltAndEncodedPassword(data["password"].(string))
	sysUser := SysUser{
		Username:    data["username"].(string),
		DisplayName: data["display_name"].(string),
		Salt:        salt,
		Password:    password,
		Email:       data["email"].(string),
		Phone:       data["phone"].(string),
		Status:      data["status"].(int),
		Remark:      data["remark"].(string),
	}

	if err := db.Create(&sysUser).Error; err != nil {
		return err
	}

	return nil
}

func GetSysUser(id uint) (*SysUser, error) {
	var sysUser SysUser
	err := db.Where("id = ?", id).First(&sysUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if sysUser.ID > 0 {
		return &sysUser, nil
	}

	return nil, nil
}

func GetSysUsers(query map[string]interface{}, page int, pageSize int) ([]*SysUser, error) {
	var sysUsers []*SysUser
	var err error
	pageNum := (page - 1) * pageSize
	username := query["username"].(string)
	if len(username) > 0 {
		username = "%" + username + "%"
	}
	if len(username) > 0 {
		err = db.Where("username like ?", username).Offset(pageNum).Limit(pageSize).Find(&sysUsers).Error
	} else {
		err = db.Offset(pageNum).Limit(pageSize).Find(&sysUsers).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return sysUsers, nil
}

func GetSysUserCount(query map[string]interface{}) (int, error) {
	var err error
	count := 0
	username := query["username"].(string)
	if len(username) > 0 {
		username = "%" + username + "%"
	}
	if len(username) > 0 {
		err = db.Model(&SysUser{}).Where("username like ?", username).Count(&count).Error
	} else {
		err = db.Model(&SysUser{}).Count(&count).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return count, err
}

func UpdateSysUser(id uint, data map[string]interface{}) error {
	_, ok := data["password"]
	if ok {
		salt, password := util.GetSaltAndEncodedPassword(data["password"].(string))
		data["salt"] = salt
		data["password"] = password
	}
	if err := db.Model(&SysUser{}).Where("id = ?", id).Update(data).Error; err != nil {
		return err
	}
	return nil
}

func UpdateSysUserPassword(id uint, data map[string]interface{}) error {
	salt, password := util.GetSaltAndEncodedPassword(data["password"].(string))
	data["salt"] = salt
	data["password"] = password
	if err := db.Model(&SysUser{}).Where("id = ?", id).Update(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteUser(id uint) error {
	user, err := GetSysUser(id)
	if err != nil {
		return err
	}
	if user.Status == 1 {
		/* admin user can not be deleted */
		return errors.New("admin user can not be deleted")
	}

	/* delete roles linked with this user */
	tx := db.Begin()
	err = db.Model(&user).Association("role").Clear().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = db.Where("id = ?", user.Model.ID).Delete(SysUser{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func GetSysRoleByUser(id uint) ([]*SysRole, error) {
	role := []*SysRole{}
	user := SysUser{}
	user.Model.ID = id
	err := db.Model(&user).Association("role").Find(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return role, nil
}

func UpdateSysUserRole(data map[string]interface{}) error {
	user := SysUser{}
	user.Model.ID = data["user_id"].(uint)
	roles := data["role_ids"].([]uint)
	var role []*SysRole
	for _, id := range roles {
		temp := SysRole{}
		temp.Model.ID = id
		role = append(role, &temp)
	}

	err := db.Model(&user).Association("role").Replace(role).Error
	if err != nil {
		return err
	}
	return nil
}
