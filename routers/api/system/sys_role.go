package system

import (
	"gin-demo/pkg/app"
	"gin-demo/pkg/e"
	"gin-demo/service/sys_role_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

type addSysRoleForm struct {
	Name   string `form:"name" valid:"Required;MaxSize(255)"`
	Remark string `form:"remark" valid:"MaxSize(255)"`
}

// @Summary Add Role
// @Produce json
// @Param name string true "name"
// @Param remark string false "remark"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/role [post]
func AddRole(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     addSysRoleForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}

	sysRoleService := sys_role_service.SysRole{
		Name:   form.Name,
		Remark: form.Remark,
	}
	if err := sysRoleService.Add(); err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.RoleAddFailed
		e.Error(err.Error())
	}

	appG.Response(httpCode, errCode, nil)
}

type queryRolesForm struct {
	Name     string `form:"name" valid:"MaxSize(255)"`
	Page     int    `form:"current" valid:"Required;Min(1)"`
	PageSize int    `form:"size" valid:"Required;Range(10, 50)"`
}

// @Summary Get roles
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/role [get]
func GetRoles(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryRolesForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)
	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}

	sysRoleService := sys_role_service.SysRole{
		Name:     form.Name,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	data := make(map[string]interface{})
	roles, err := sysRoleService.GetAll()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.MenuGetFailed
		e.Error(err.Error())
	} else {
		count, err := sysRoleService.Count()
		if err != nil {
			httpCode = http.StatusInternalServerError
			errCode = e.MenuGetFailed
			e.Error(err.Error())
		} else {
			data["records"] = roles
			data["total"] = count
		}
	}
	appG.Response(httpCode, errCode, data)
}

// @Summary Get role
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/role/{id} [get]
func GetRole(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	sysRoleService := sys_role_service.SysRole{
		ID: uint(id),
	}

	role, err := sysRoleService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.MenuGetFailed, nil)
		e.Error(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["role"] = role
	appG.Response(http.StatusOK, e.SUCCESS, data)

}

type updateSysRoleForm struct {
	ID     int    `form:"id" valid:"Required;Min(1)"`
	Name   string `form:"string" valid:"Required;MaxSize(255)"`
	Remark string `form:"string" valid:"MaxSize(255)"`
}

// @Summary Update Role
// @Produce json
// @Param id uint true "id"
// @Param name string true "name"
// @Param remark string false "remark"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/role [put]
func UpdateRole(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     updateSysRoleForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}

	sysRoleService := sys_role_service.SysRole{
		ID:     uint(form.ID),
		Name:   form.Name,
		Remark: form.Remark,
	}
	if err := sysRoleService.Update(); err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.RoleUpdateFailed
		e.Error(err.Error())
	}

	appG.Response(httpCode, errCode, nil)
}

// @Summary Delete Role
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/role/{id} [delete]
func DeleteRole(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.InvalidParams, nil)
		return
	}

	sysRoleService := sys_role_service.SysRole{
		ID: uint(id),
	}

	err := sysRoleService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.RoleDeleteFailed, nil)
		e.Error(err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type MenuIds struct {
	Menus []uint `json:"ids" binding:"required"`
}

// @Summary Update role by menu
// @Produce json
// @Param id uint true "id"
// @Param ids []uint true "ids"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/role/{id}/menu [put]
func UpdateRoleByMenu(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	var menu MenuIds
	err := c.BindJSON(&menu)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.RoleMenuUpdateFailed
		appG.Response(httpCode, errCode, nil)
		return
	}
	sysRoleService := sys_role_service.SysRole{
		ID:      uint(id),
		MenuIds: menu.Menus,
	}
	if err := sysRoleService.UpdateMenu(); err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.RoleMenuUpdateFailed
		e.Error(err.Error())
	}

	appG.Response(httpCode, errCode, nil)
}

// @Summary get role by role
// @Produce json
// @Param id uint true "id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /sys/role/{id}/menu/ [get]
func GetMenuByRole(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	httpCode := http.StatusOK
	errCode := e.SUCCESS

	sysRoleService := sys_role_service.SysRole{
		ID: uint(id),
	}
	menu, err := sysRoleService.GetMenu()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.RoleMenuGetFailed
		e.Error(err.Error())
		appG.Response(httpCode, errCode, nil)
		return
	}

	data := make(map[string]interface{})
	data["records"] = menu
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
