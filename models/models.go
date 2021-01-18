package models

import (
	"fmt"
	"log"
	"time"

	"gin-demo/pkg/setting"

	"github.com/jinzhu/gorm"
	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updateAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

var db *gorm.DB

func init() {
	var (
		err                                               error
		adminPassword                                     string
		dbType, dbName, user, password, host, tablePrefix string
	)

	adminPassword = setting.GConfig.APP.AdminSecret
	dbType = setting.GConfig.DataBase.Type
	dbName = setting.GConfig.DataBase.Name
	user = setting.GConfig.DataBase.User
	password = setting.GConfig.DataBase.Password
	host = setting.GConfig.DataBase.Host
	tablePrefix = setting.GConfig.DataBase.TablePrefix

	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))
	if err != nil {
		log.Fatalln(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true) // if needed, switch to true
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.AutoMigrate(
		SysUser{},
		SysAudit{},
		SysMenu{},
		SysRole{},
		SysAudit{},
	)

	admin, err := GetSysUser(1)
	if err != nil {
		log.Fatalln("获取管理员失败，请检查数据库")
	}
	if admin == nil {
		if len(adminPassword) == 0 {
			log.Fatalln("管理员用户不存在，必须设置管理员密码 admin-secret")
		}
		data := map[string]interface{}{
			"username":     "admin",
			"display_name": "管理员",
			"password":     adminPassword,
			"email":        "admin@admin.com",
			"phone":        "12345678901",
			"status":       1,
			"remark":       "管理员账户",
		}
		err = AddSysUser(data)
		if err != nil {
			log.Fatalln("增加管理员失败，请检查数据库")
		}
	}
	if len(adminPassword) != 0 {
		data := make(map[string]interface{})
		data["password"] = adminPassword
		err = UpdateSysUserPassword(1, data)
		if err != nil {
			log.Fatalln("更改管理员密码失败，请检查数据库")
		}
	}

}

//CloseDB close db
func CloseDB() {
	defer db.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now()
		if createTimeField, ok := scope.FieldByName("CreatedAt"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("UpdatedAt"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("UpdatedAt", time.Now())
	}
}

// deleteCallback will set `DeletedOn` where deleting
func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedAt")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
