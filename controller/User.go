package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"muxi_auditor/api/request"
	"muxi_auditor/api/response"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/repository/model"
	"muxi_auditor/service"
)

type UserController struct {
	service UserService
}
type UserService interface {
	GetUsers(ctx context.Context, id uint) ([]model.UserResponse, error)
	UpdateUserRole(ctx context.Context, userId uint, projectPermit []model.ProjectPermit, role int) error
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}
func (uc *UserController) GetUsers(ctx *gin.Context, req request.GetUserReq) (response.Response, error) {
	if req.Role != 2 {
		return response.Response{
			Msg:  "无权限",
			Code: 40001,
			Data: nil,
		}, nil

	}
	userResponse, err := uc.service.GetUsers(ctx, req.Project_id)
	if err != nil {
		return response.Response{}, err
	}

	return response.Response{
		Msg:  "获取成功",
		Code: 200,
		Data: userResponse,
	}, nil
}
func (uc *UserController) UpdateUsers(ctx *gin.Context, req request.UpdateUserRoleReq, cla jwt.UserClaims) (response.Response, error) {

	if cla.UserRule != 2 {
		return response.Response{
			Msg:  "无权限",
			Code: 40001,
			Data: nil,
		}, nil
	}
	err := uc.service.UpdateUserRole(ctx, req.UserId, req.ProjectPermit, req.Role)
	if err != nil {
		return response.Response{}, err
	}
	return response.Response{
		Code: 200,
		Msg:  "修改成功",
		Data: nil,
	}, nil
}
