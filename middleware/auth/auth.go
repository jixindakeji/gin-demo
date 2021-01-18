package auth

import (
	"fmt"
	"gin-demo/pkg/e"
	"gin-demo/service/sys_user_service"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		method := c.Request.Method
		var authItems []string
		found := false

		if url[len(url)-1] == '/' {
			url = url[:len(url)-1]
		}

		authUrls := strings.Split(url[1:], "/")
		for _, item := range authUrls {
			ok, _ := regexp.MatchString("\\d+", item)
			if !ok {
				authItems = append(authItems, item)
			}
		}
		auth := strings.Join(authItems, ":") + ":" + strings.ToLower(method)
		userId, ok := c.Get("userId")
		if !ok {
			fmt.Println("get user failed")
			c.AbortWithStatus(401)
			return
		}
		sysUserService := sys_user_service.SysUser{
			ID: userId.(uint),
		}
		user, err := sysUserService.Get()
		if err != nil {
			c.AbortWithStatus(500)
			e.Error(err.Error())
			return
		}
		if user == nil || user.Status == 0 {
			c.AbortWithStatus(401)
			return
		}

		_, apiAuth, err := sysUserService.GetSysMenuByUser()
		if err != nil {
			c.AbortWithStatus(500)
			e.Error(err.Error())
			return
		}
		for _, item := range apiAuth {
			if item == auth {
				found = true
				break
			}
		}
		if !found {
			c.AbortWithStatus(401)
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
