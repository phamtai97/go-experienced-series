package service

import (
	"user-service/internal/core/common/utils"
	"user-service/internal/core/dto"
	"user-service/internal/core/entity/error_code"
	"user-service/internal/core/model/request"
	"user-service/internal/core/model/response"
	"user-service/internal/core/port/repository"
	"user-service/internal/core/port/service"
)

const (
	invalidUserNameErrMsg = "invalid username"
	invalidPasswordErrMsg = "invalid password"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u userService) SignUp(request *request.SignUpRequest) *response.Response {
	// validate request
	if len(request.Username) == 0 {
		return u.createFailedResponse(error_code.InvalidRequest, invalidUserNameErrMsg)
	}

	if len(request.Password) == 0 {
		return u.createFailedResponse(error_code.InvalidRequest, invalidPasswordErrMsg)
	}

	currentTime := utils.GetUTCCurrentMillis()
	userDTO := dto.UserDTO{
		UserName:    request.Username,
		Password:    request.Password,
		DisplayName: u.getRandomDisplayName(request.Username),
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
	}

	// save a new user
	err := u.userRepo.Insert(userDTO)
	if err != nil {
		if err == repository.DuplicateUser {
			return u.createFailedResponse(error_code.DuplicateUser, err.Error())
		}

		return u.createFailedResponse(error_code.InternalError, error_code.InternalErrMsg)
	}

	// create data response
	signUpData := response.SignUpDataResponse{
		DisplayName: userDTO.DisplayName,
	}

	return u.createSuccessResponse(signUpData)
}

func (u userService) getRandomDisplayName(username string) string {
	return username + utils.ConvertUInt64ToString(utils.GetUTCCurrentMillis())
}

func (u userService) createFailedResponse(
	code error_code.ErrorCode, message string,
) *response.Response {
	return &response.Response{
		ErrorCode:    code,
		ErrorMessage: message,
		Status:       false,
	}
}

func (u userService) createSuccessResponse(data response.SignUpDataResponse) *response.Response {
	return &response.Response{
		Data:         data,
		ErrorCode:    error_code.Success,
		ErrorMessage: error_code.SuccessErrMsg,
		Status:       true,
	}
}
