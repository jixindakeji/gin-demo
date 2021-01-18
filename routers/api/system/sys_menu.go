package system

import (
	"gin-demo/pkg/app"
	"gin-demo/pkg/e"
	"gin-demo/service/sys_menu_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

type addSysMenuForm struct {
	ParentId uint   `form:"parentId"`
	Name     string `form:"name" valid:"Required;MaxSize(255)"`
	Url      string `form:"url" valid:"MaxSize(255)"`
	Type     int    `form:"type" valid:"Required;Range(1,3)"`
	Code     string `form:"code" valid:"MaxSize(255)"`
	Icon     string `form:"icon" valid:"MaxSize(255)"`
	Sort     int    `form:"sort" valid:"Range(0,10000)"`
	Status   int    `form:"status"`

	Remark string `form:"remark" valid:"MaxSize(255)"`
}

// @Summary Add Menu
// @Produce json
// @Param parentId string true "parentId"
// @Param name string true "name"
// @Param url string true "url"
// @Param code string true "code"
// @Param type string true "type"
// @Param icon string false "icon"
// @Param sort int true "sort"
// @Param remark string false "remark"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/menu [post]
func AddMenu(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     addSysMenuForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}

	sysMenuService := sys_menu_service.SysMenu{
		ParentId: form.ParentId,
		Name:     form.Name,
		Url:      form.Url,
		Type:     form.Type,
		Code:     form.Code,
		Icon:     form.Icon,
		Sort:     form.Sort,
		Status:   form.Status,

		Remark: form.Remark,
	}
	if err := sysMenuService.Add(); err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.MenuAddFailed
		e.Error(err.Error())
	}

	appG.Response(httpCode, errCode, nil)
}

// @Summary Get all menus
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/menu [get]
func GetMenus(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)
	sysMenuService := sys_menu_service.SysMenu{}
	menus, err := sysMenuService.GetAll()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.MenuGetFailed
		e.Error(err.Error())
	}
	data := make(map[string]interface{})
	data["records"] = menus
	appG.Response(httpCode, errCode, data)
}

type updateSysMenuForm struct {
	ID       int    `form:"id" valid:"Min(0)"`
	ParentId int    `form:"parentId"`
	Name     string `form:"name" valid:"Required;MaxSize(255)"`
	Url      string `form:"url" valid:"MaxSize(255)"`
	Type     int    `form:"type" valid:"Required;Range(1,3)"`
	Code     string `form:"code" valid:"MaxSize(255)"`
	Icon     string `form:"icon"`
	Sort     int    `form:"sort" valid:"Required;Range(1,1000)"`
	Remark   string `form:"remark" valid:"MaxSize(255)"`
}

// @Summary Update Menu
// @Produce json
// @Param parentId string true "parentId"
// @Param name string true "name"
// @Param url string true "url"
// @Param code string true "code"
// @Param type string true "type"
// @Param icon string false "icon"
// @Param sort int true "sort"
// @Param remark string false "remark"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/menu [put]
func UpdateMenu(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     updateSysMenuForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)

	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}

	sysMenuService := sys_menu_service.SysMenu{
		ID:       uint(form.ID),
		ParentId: uint(form.ParentId),
		Name:     form.Name,
		Url:      form.Url,
		Type:     form.Type,
		Code:     form.Code,
		Icon:     form.Icon,
		Sort:     form.Sort,
		Remark:   form.Remark,
	}
	if err := sysMenuService.Update(); err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.MenuUpdateFailed
		e.Error(err.Error())
	}

	appG.Response(httpCode, errCode, nil)
}

// @Summary Delete menu
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/menu/{id} [delete]
func DeleteMenu(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.InvalidParams, nil)
		return
	}

	sysMenuService := sys_menu_service.SysMenu{
		ID: uint(id),
	}

	err := sysMenuService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.MenuDeleteFailed, nil)
		e.Error(err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Get menu
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /sys/menu/{id} [get]
func GetMenu(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.InvalidParams, nil)
		return
	}

	sysMenuService := sys_menu_service.SysMenu{
		ID: uint(id),
	}

	menu, err := sysMenuService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.MenuGetFailed, nil)
		e.Error(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["menu"] = menu
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
