package service

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"todo_list/model"
	"todo_list/serializer"
	"todo_list/utils"
)

// UserRegisterService 用户注册服务
type UserService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=16"`
}

func (service *UserService) Register() serializer.Response {
	var user model.User
	var count int64
	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).
		First(&user).Count(&count)
	if count == 1 {
		return serializer.Response{
			Status: http.StatusBadRequest,
			Msg:    "用户已存在，无需再注册",
		}
	}
	user.UserName = service.UserName
	// 密码的加密
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Status: http.StatusBadRequest,
			Msg:    err.Error(),
		}
	}
	// 创建用户
	if model.DB.Create(&user).Error != nil {
		return serializer.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库操作错误",
		}
	}
	return serializer.Response{
		Status: http.StatusOK,
		Msg:    "用户注册成功",
	}
}

func (service *UserService) Login() serializer.Response {
	var user model.User
	// 先去数据库找这个user

	// 用户不存在，先注册
	if err := model.DB.Where("user_name = ?", service.UserName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return serializer.Response{
				Status: 400,
				Msg:    "用户不存在，请先注册",
			}
		}

		// 其他不可抗错误
		return serializer.Response{
			Status: 500,
			Msg:    "数据库操作错误",
		}
	}

	// 用户存在，但密码错误
	if user.CheckPassword(service.Password) == false {
		return serializer.Response{
			Status: 400,
			Msg:    "密码错误",
		}
	}

	// 用户存在且密码正确
	// 发一个token，为了其他功能需要身份验证所给前端存储的。
	// 创建备忘录需要token，不然不知道是谁创建的
	token, err := utils.GenerateToken(user.ID, service.UserName, service.Password)
	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "token签发错误",
		}
	}
	return serializer.Response{
		Status: 200,
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user),
			Token: token,
		},
		Msg: "登录成功",
	}
}
