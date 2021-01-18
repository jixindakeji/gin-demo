package system

import (
	"fmt"
	"gin-demo/models"
	"gin-demo/pkg/app"
	"gin-demo/pkg/e"
	"gin-demo/pkg/util"
	"gin-demo/service/sys_user_service"
	"github.com/astaxie/beego/validation"
	"github.com/unknwon/com"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authForm struct {
	Username string `form:"username" valid:"Required;MaxSize(50)"`
	Password string `form:"password" valid:"Required;MaxSize(50)"`
}

type userPermission struct {
	Uid          uint              `json:"uid"`
	MenuTreeList []*models.SysMenu `json:"menuTreeList"`
	ButtonList   []string          `json:"buttonList"`
}

type userInfo struct {
	UserName       string          `json:"username"`
	DisplayName    string          `json:"displayName"`
	AccessToken    string          `json:"accessToken"`
	RefreshToken   string          `json:"refreshToken"`
	UserPermission *userPermission `json:"userInfo"`
}

// @Summary Get Auth
// @Produce json
// @Param username string true "username"
// @Param password string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /login [post]
func Login(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     authForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)
	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}

	redirectUrl := c.Query("redirectUrl")
	if len(redirectUrl) > 0 {
		fmt.Println(redirectUrl)
	}

	sysUserService := sys_user_service.SysUser{
		Username: form.Username,
		Password: form.Password,
	}

	sysUser, err := sysUserService.GetSysUserByUsername()
	if err != nil {
		errCode = e.UserGetFailed
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}
	if sysUser.Status == 3 {
		errCode = e.UserDisabled
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}
	if sysUser != nil {
		ok := util.VerifyRawPassword(form.Password, sysUser.Password, sysUser.Salt)
		if ok {
			token, err := util.GenerateToken(sysUser.Username, sysUser.ID)
			if err != nil {
				errCode = e.TokenError
				httpCode = http.StatusInternalServerError
				appG.Response(httpCode, errCode, nil)
				return
			}

			sysUserService.ID = sysUser.ID
			menu, button, err := sysUserService.GetSysMenuByUser()
			if err != nil {
				errCode = e.UserMenuGetFailed
				httpCode = http.StatusInternalServerError
				appG.Response(httpCode, errCode, nil)
				return
			}

			userPermission := &userPermission{
				Uid:          sysUser.ID,
				MenuTreeList: menu,
				ButtonList:   button,
			}
			userInfo := userInfo{
				UserName:       sysUser.Username,
				DisplayName:    sysUser.DisplayName,
				AccessToken:    token,
				RefreshToken:   token,
				UserPermission: userPermission,
			}

			resp := make(map[string]interface{})
			resp["user"] = &userInfo
			appG.Response(httpCode, errCode, resp)
			return
		}
	}
	errCode = e.TokenFailed
	appG.Response(httpCode, errCode, nil)
}

type addUserForm struct {
	Username    string `form:"username" valid:"Required;MaxSize(100)"`
	DisplayName string `form:"displayName" valid:"MaxSize(255)"`
	Phone       string `form:"phone" valid:"Numeric"`
	Email       string `form:"email" valid:"Email;MaxSize(100)"`
	Status      int    `form:"status" valid:"Range(2, 3)"`
	Remark      string `form:"remark" valid:"MaxSize(255)"`
}

// @Summary Add User
// @Produce json
// @Param username string true "username"
// @Param password string true "password"
// @Param isSuperUser int false "isSuperUser"
// @Param display_name string false "display_name"
// @Param email string false "email"
// @Param phone string false "phone"
// @Param status int false "status"
// @Param remark string false "remark"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/user [post]
func AddUser(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     addUserForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}

	sysUserService := sys_user_service.SysUser{
		Username:    form.Username,
		DisplayName: form.DisplayName,
		Password:    form.Username, // init password by username
		Phone:       form.Phone,
		Email:       form.Email,
		Status:      form.Status,
		Remark:      form.Remark,
	}

	err := sysUserService.Add()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.UserAddFailed
		e.Error(err.Error())
	}

	appG.Response(httpCode, errCode, nil)
}

// @Summary Get user
// @Produce json
// @Param id int true "id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/user/{id} [post]
func GetUser(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.InvalidParams, nil)
		return
	}

	sysUserService := sys_user_service.SysUser{
		ID: uint(id),
	}
	user, err := sysUserService.Get()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.UserGetFailed
		e.Error(err.Error())
	}
	data := make(map[string]interface{})
	data["user"] = user
	appG.Response(httpCode, errCode, data)
}

type queryUsersForm struct {
	UserName string `form:"username" valid:"MaxSize(255)"`
	Page     int    `form:"current" valid:"Required;Min(1)"`
	PageSize int    `form:"size" valid:"Required;Range(10, 50)"`
}

// @Summary Get users
// @Produce json
// @Param username string false "username"
// @Param status int false "status"
// @Param page int true "page"
// @Param pageSize int true "pageSize"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/user [get]
func GetUsers(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryUsersForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)
	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		fmt.Println(httpCode)
		appG.Response(httpCode, errCode, nil)
		return
	}
	sysUserService := sys_user_service.SysUser{
		Username: form.UserName,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	data := make(map[string]interface{})
	users, err := sysUserService.GetAll()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.UserGetFailed
		e.Error(err.Error())
	} else {
		count, err := sysUserService.Count()
		if err != nil {
			httpCode = http.StatusInternalServerError
			errCode = e.UserGetCountFailed
		} else {
			data["records"] = users
			data["total"] = count
		}
	}
	appG.Response(httpCode, errCode, data)
}

type updateUserForm struct {
	ID          int    `form:"id" valid:"Required;Min(1)"`
	DisplayName string `form:"displayName" valid:"MaxSize(255)"`
	Password    string `form:"password"`
	Email       string `form:"email" valid:"Email;MaxSize(100)"`
	Phone       string `form:"phone" valid:"Numeric"`
	Status      int    `form:"status" valid:"Range(1, 3)"`
	Remark      string `form:"remark" valid:"MaxSize(255)"`
}

// @Summary Update User
// @Produce json
// @Param username string true "username"
// @Param password string true "password"
// @Param isSuperUser int false "isSuperUser"
// @Param display_name string false "display_name"
// @Param email string false "email"
// @Param phone string false "phone"
// @Param status int false "status"
// @Param remark string false "remark"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/user [put]
func UpdateUser(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     updateUserForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}

	sysUserService := sys_user_service.SysUser{
		ID:          uint(form.ID),
		DisplayName: form.DisplayName,
		Password:    form.Password,
		Email:       form.Email,
		Phone:       form.Phone,
		Status:      form.Status,
		Remark:      form.Remark,
	}

	err := sysUserService.Update()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.UserAddFailed
		e.Error(err.Error())
	}
	appG.Response(httpCode, errCode, nil)
}

// @Summary Delete user
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/user/{id} [delete]
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.InvalidParams, nil)
		return
	}

	sysUserService := sys_user_service.SysUser{
		ID: uint(id),
	}

	err := sysUserService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.UserDeleteFailed, nil)
		e.Error(err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary get role by user
// @Produce json
// @Param id uint true "id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/user/{id}/role [get]
func GetUserRole(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	httpCode := http.StatusOK
	errCode := e.SUCCESS

	sysRoleService := sys_user_service.SysUser{
		ID: uint(id),
	}
	role, err := sysRoleService.GetRole()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.UserRoleGetFailed
		e.Error(err.Error())
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = role
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type RoleIds struct {
	Roles []uint `json:"ids" binding:"required"`
}

// @Summary Update role by user
// @Produce json
// @Param userId uint true "userId"
// @Param ids []uint true "ids"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/user/{id}/role [put]
func UpdateUserRole(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	var role RoleIds
	err := c.BindJSON(&role)
	if err != nil {
		fmt.Println(err)
		httpCode = http.StatusInternalServerError
		errCode = e.RoleMenuUpdateFailed
		appG.Response(httpCode, errCode, nil)
		return
	}
	sysUserService := sys_user_service.SysUser{
		ID:      uint(id),
		RoleIds: role.Roles,
	}
	if err := sysUserService.UpdateRole(); err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.RoleMenuUpdateFailed
		e.Error(err.Error())
	}

	appG.Response(httpCode, errCode, nil)
}
