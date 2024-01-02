package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"todo_list/serializer"
)

func ErrorResponse(err error) serializer.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := fmt.Sprintf("Field.%s", e.Field)
			tag := fmt.Sprintf("Tag.Valid.%s", e.Tag)
			return serializer.Response{
				Status: 40001,
				Msg:    fmt.Sprintf("%s%s", field, tag),
				Error:  fmt.Sprint(err),
			}
		}
	}
	var unmarshalTypeError *json.UnmarshalTypeError
	if !errors.As(err, &unmarshalTypeError) {
		return serializer.Response{
			Status: 40001,
			Msg:    "JSON类型不匹配",
			Error:  fmt.Sprint(err),
		}
	}
	return serializer.Response{
		Status: 40001,
		Msg:    "参数错误",
		Error:  fmt.Sprint(err),
	}
}
