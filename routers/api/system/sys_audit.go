package system

import (
	"gin-demo/pkg/app"
	"gin-demo/pkg/e"
	"gin-demo/service/sys_audit_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type queryAuditForm struct {
	UserName string `form:"username" valid:"MaxSize(255)"`
	Page     int    `form:"current" valid:"Required;Min(1)"`
	PageSize int    `form:"size" valid:"Required;Range(10, 50)"`
}

// @Summary Get users
// @Produce json
// @Param username string false "username"
// @Param page int true "current"
// @Param pageSize int true "size"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /system/user/audit [get]
func GetAudits(c *gin.Context) {
	var (
		appG     = app.Gin{C: c}
		form     queryAuditForm
		httpCode = http.StatusOK
		errCode  = e.SUCCESS
	)
	errCode = app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		httpCode = http.StatusInternalServerError
		appG.Response(httpCode, errCode, nil)
		return
	}
	sysAuditService := sys_audit_service.SysAudit{
		User:     form.UserName,
		Page:     form.Page,
		PageSize: form.PageSize,
	}
	data := make(map[string]interface{})
	users, err := sysAuditService.GetAll()
	if err != nil {
		httpCode = http.StatusInternalServerError
		errCode = e.AuditsGetFailed
		e.Error(err.Error())
	} else {
		count, err := sysAuditService.Count()
		if err != nil {
			httpCode = http.StatusInternalServerError
			errCode = e.AuditsGetFailed
			e.Error(err.Error())
		} else {
			data["records"] = users
			data["total"] = count
		}
	}
	appG.Response(httpCode, errCode, data)
}
