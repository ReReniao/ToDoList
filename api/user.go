package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo_list/service"
)

func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	// 绑定请求体中的数据
	if err := c.ShouldBind(&userRegister); err == nil {
		// 执行注册方法，并返回响应数据结构体
		res := userRegister.Register()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
	}
}

func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	// 绑定请求体中的数据
	if err := c.ShouldBind(&userLogin); err == nil {
		// 执行注册方法，并返回响应数据结构体
		res := userLogin.Login()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
	}
}
