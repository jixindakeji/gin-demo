package jwt

import (
	"gin-demo/pkg/e"
	"gin-demo/pkg/setting"
	"gin-demo/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//JWT jwt
func JWT() gin.HandlerFunc {
	tokenName := setting.GConfig.APP.JwtTokenName
	return func(c *gin.Context) {
		var code int
		var data interface{}
		var id uint
		code = e.SUCCESS
		token := c.Request.Header.Get(tokenName)
		if token == "" {
			code = e.TokenMissed
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.TokenInvalid
			} else if claims == nil {
				code = e.TokenError
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.TokenExpired
			} else {
				id = claims.ID
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Set("userId", id)

		c.Next()
	}
}
