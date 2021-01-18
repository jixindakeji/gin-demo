package app

import (
	"fmt"
	"gin-demo/pkg/e"
	"github.com/astaxie/beego/validation"

	"github.com/gin-gonic/gin"
)

// BindAndValid bind and valid
func BindAndValid(c *gin.Context, form interface{}) int {
	err := c.Bind(form)
	if err != nil {
		fmt.Println(err)
		return e.InvalidParams
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		fmt.Println(err)
		return e.ERROR
	}
	if !check {
		fmt.Println(err)
		return e.InvalidParams
	}
	return e.SUCCESS
}
