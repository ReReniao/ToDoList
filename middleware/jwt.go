package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
	"todo_list/utils"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 200
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			claim, err := utils.ParseToken(token)
			if err != nil {
				code = 403 // 无权限，token是无权限的
			} else if time.Now().Unix() > claim.ExpiresAt {
				code = 401 // token过期了
			}
		}
		if code != 200 {
			c.JSON(200, gin.H{
				"Status": code,
				"Msg":    "Token解析错误",
			})
			c.Abort()
			return
		}
	}
}
