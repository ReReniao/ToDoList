package service

import (
	"time"
	"todo_list/model"
	"todo_list/serializer"
)

type CreateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"status" form:"status"` // 0 表示未做 1 表示已做
}

func (service *CreateTaskService) Create(id uint) serializer.Response {
	var user model.User
	code := 200
	model.DB.First(&user, id)
	task := model.Task{
		User:      user,
		Uid:       id,
		Title:     service.Title,
		Status:    0,
		Content:   service.Content,
		StartTime: time.Now().Unix(),
		EndTime:   0,
	}
	if model.DB.Create(&task).Error != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Msg:    "创建备忘录失败",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    "创建成功",
	}
}

type ShowTaskService struct {
}

type ListTaskService struct {
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
}

type UpdateTaskService struct {
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
	Status  int    `json:"status" form:"status"` // 0 表示未做 1 表示已做
}

type SearchTaskService struct {
	Info     string `json:"info" form:"info"`
	PageNum  int    `json:"page_num" form:"page_num"`
	PageSize int    `json:"page_size" form:"page_size"`
}

type DeleteTaskService struct {
}

func (service *ShowTaskService) Show(tid string) serializer.Response {
	var task model.Task
	code := 200
	if model.DB.First(&task, tid).Error != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Data:   nil,
			Msg:    "查询失败",
			Error:  "",
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildTask(task),
		Msg:    "",
		Error:  "",
	}
}

// 列表返回用户所有备忘录
func (service *ListTaskService) List(uid uint) serializer.Response {
	var tasks []model.Task
	var count int64 = 0
	code := 200
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	err := model.DB.Model(&model.Task{}).Preload("User").Where("uid = ?", uid).Count(&count).
		Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&tasks).Error
	if err != nil {
		code = 500
		return serializer.Response{
			Status: code,
			Data:   nil,
			Msg:    "查询失败",
			Error:  "",
		}
	}
	return serializer.BuildListResponse(serializer.BuildTasks(tasks), count)
}

// 更新备忘录
func (service *UpdateTaskService) Update(tid string) serializer.Response {
	var task model.Task
	code := 200
	if err := model.DB.First(&task, tid).Error; err != nil {
		code = 30001
		return serializer.Response{
			Status: code,
			Msg:    "数据库操作错误",
		}
	}
	task.Content = service.Content
	task.Title = service.Title
	task.Status = service.Status
	if err := model.DB.Save(&task).Error; err != nil {
		code = 30002
		return serializer.Response{
			Status: code,
			Msg:    "数据库操作错误",
		}
	}
	return serializer.Response{
		Status: code,
		Data:   serializer.BuildTask(task),
		Msg:    "更新完成",
	}
}

// 查询备忘录
func (service *SearchTaskService) Search(uid uint) serializer.Response {
	var tasks []model.Task
	var count int64 = 0
	code := 200
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	err := model.DB.Model(&model.Task{}).Preload("User").Where("uid = ?", uid).Where("title like ? or content like ?", "%"+service.Info+"%", "%"+service.Info+"%").
		Count(&count).Limit(service.PageSize).Offset((service.PageNum - 1) * service.PageSize).Find(&tasks).Error
	if err != nil {
		code = 30001
		return serializer.Response{
			Status: code,
			Msg:    "数据库操作错误",
		}
	}

	return serializer.BuildListResponse(serializer.BuildTasks(tasks), count)
}

// 删除备忘录
func (service *DeleteTaskService) Delete(tid string, uid uint) serializer.Response {
	var task model.Task
	err := model.DB.Where("uid = ?", uid).Delete(&task, tid).Error
	if err != nil {
		return serializer.Response{
			Status: 500,
			Data:   nil,
			Msg:    "删除失败",
			Error:  "",
		}
	}
	return serializer.Response{
		Status: 200,
		Data:   nil,
		Msg:    "删除成功",
		Error:  "",
	}
}
