package util

import (
	"gin-demo/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

//GetPage get page by url
func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.GConfig.APP.PageSize
	}
	return result
}
