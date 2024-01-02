package serializer

import "todo_list/model"

type User struct {
	ID       uint   `json:"id" form:"id" example:"1"`
	UserName string `json:"user_name" form:"user_name" example:"reniao"`
	Status   string `json:"status" form:"status"`
	CreateAt int64  `json:"create_at" form:"creat_at"`
}

// BuildUser 序列化用户
func BuildUser(user model.User) User {
	return User{
		ID:       user.ID,
		UserName: user.UserName,
		CreateAt: user.CreatedAt.Unix(),
	}
}
