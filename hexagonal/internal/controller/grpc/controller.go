package grpc

import (
	"context"

	proto "user-service/api"
	"user-service/internal/core/entity/error_code"
	"user-service/internal/core/model/request"
	"user-service/internal/core/model/response"
	"user-service/internal/core/port/service"
)

var errorCodeMapper = map[error_code.ErrorCode]proto.ErrorCode{
	error_code.Success:        proto.ErrorCode_SUCCESS,
	error_code.InternalError:  proto.ErrorCode_EC_UNKNOWN,
	error_code.InvalidRequest: proto.ErrorCode_INVALID_REQUEST,
	error_code.DuplicateUser:  proto.ErrorCode_DUPLICATE_USER,
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) proto.UserServiceServer {
	return &userController{
		userService: userService,
	}
}

func (u userController) SignUp(
	ctx context.Context, request *proto.SignUpRequest,
) (*proto.SignUpResponse, error) {
	resp := u.userService.SignUp(u.newSignUpRequest(request))
	return u.newSignUpResponse(resp)
}

func (u userController) newSignUpRequest(protoRequest *proto.SignUpRequest) *request.SignUpRequest {
	return &request.SignUpRequest{
		Username: protoRequest.GetUserName(),
		Password: protoRequest.GetPassword(),
	}
}

func (u userController) newSignUpResponse(resp *response.Response) (
	*proto.SignUpResponse, error,
) {
	if !resp.Status {
		return &proto.SignUpResponse{
			Status:       resp.Status,
			ErrorCode:    u.mapErrorCode(resp.ErrorCode),
			ErrorMessage: resp.ErrorMessage,
		}, nil
	}

	data := resp.Data.(response.SignUpDataResponse)
	return &proto.SignUpResponse{
		Status:       resp.Status,
		ErrorCode:    u.mapErrorCode(resp.ErrorCode),
		ErrorMessage: resp.ErrorMessage,
		DisplayName:  data.DisplayName,
	}, nil
}

func (u userController) mapErrorCode(errCode error_code.ErrorCode) proto.ErrorCode {
	code, existed := errorCodeMapper[errCode]
	if existed {
		return code
	}

	return proto.ErrorCode_EC_UNKNOWN
}
