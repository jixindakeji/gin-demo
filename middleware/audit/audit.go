package audit

import (
	"bytes"
	"fmt"
	"gin-demo/models"
	"gin-demo/pkg/e"
	"gin-demo/service/sys_audit_service"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

// audit logger
func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		bodyBytes := []byte{}
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err)
		} else {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		user, ok := c.Get("user")
		if !ok {
			return
		}
		sysUser := user.(*models.SysUser)
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		httpCode := c.Writer.Status()
		cliIP := c.ClientIP()

		sysAuditService := sys_audit_service.SysAudit{
			User:    sysUser.Username,
			UserId:  sysUser.ID,
			Ip:      cliIP,
			Method:  reqMethod,
			Path:    reqURI,
			Status:  httpCode,
			Body:    string(bodyBytes),
			Latency: latencyTime,
		}
		err = sysAuditService.Add()
		if err != nil {
			e.Error(err.Error())
		}
	}
}
