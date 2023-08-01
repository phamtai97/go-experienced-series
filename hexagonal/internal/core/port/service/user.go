package service

import (
	"user-service/internal/core/model/request"
	"user-service/internal/core/model/response"
)

type UserService interface {
	SignUp(request *request.SignUpRequest) *response.Response
}
