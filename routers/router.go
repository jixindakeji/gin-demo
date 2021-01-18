package routers

import (
	"gin-demo/middleware/audit"
	"gin-demo/middleware/auth"
	"gin-demo/middleware/jwt"
	"gin-demo/pkg/setting"
	"gin-demo/routers/api/system"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

//InitRouter init router
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default()) //跨域
	gin.SetMode(setting.GConfig.RunMode)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(299, gin.H{
			"message": "pong",
		})
	})
	// Login api without jwt
	r.POST("/login", system.Login)
	// System api
	apiSys := r.Group("/system", jwt.JWT(), auth.Auth(), audit.Audit())
	{
		apiSys.POST("/menu", system.AddMenu)
		apiSys.GET("/menu", system.GetMenus)
		apiSys.GET("/menu/:id", system.GetMenu)
		apiSys.PUT("/menu", system.UpdateMenu)
		apiSys.DELETE("/menu/:id", system.DeleteMenu)

		apiSys.POST("/role", system.AddRole)
		apiSys.GET("/role/:id", system.GetRole)
		apiSys.GET("/role", system.GetRoles)
		apiSys.PUT("/role", system.UpdateRole)
		apiSys.DELETE("/role/:id", system.DeleteRole)

		apiSys.GET("/role/:id/menu", system.GetMenuByRole)
		apiSys.POST("/role/:id/menu", system.UpdateRoleByMenu)

		apiSys.POST("/user", system.AddUser)
		apiSys.PUT("/user", system.UpdateUser)
		apiSys.GET("/user/:id", system.GetUser)
		apiSys.GET("/user", system.GetUsers)
		apiSys.DELETE("/user/:id", system.DeleteUser)

		apiSys.GET("/user/:id/role", system.GetUserRole)
		apiSys.POST("/user/:id/role", system.UpdateUserRole)

		apiSys.GET("/audit", system.GetAudits)

	}

	return r
}
