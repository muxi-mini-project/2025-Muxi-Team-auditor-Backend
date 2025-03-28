package controller

import (
	"context"
	"fmt"
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
	UpdateUserRole(ctx context.Context, userId uint, projectPermit []model.ProjectPermit, role int) error
	GetMyInfo(ctx context.Context, id uint) (*model.User, error)
	UpdateMyInfo(ctx context.Context, req request.UpdateUserReq, id uint) error
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// UpdateUsers 更新用户权限等信息
// @Summary 更新用户角色
// @Description 修改指定用户的角色，根据项目权限设置角色信息
// @Tags User
// @Accept json
// @Produce json
// @Param UpdateUserRoleReq body request.UpdateUserRoleReq true "更新用户角色请求体"
// @Success 200 {object} response.Response "成功更新用户角色"
// @Failure 40001 {object} response.Response "无权限"
// @Failure 400 {object} response.Response "修改失败"
// @Security ApiKeyAuth
// @Router /api/v1/user/updateUser [post]
func (c *UserController) UpdateUsers(ctx *gin.Context, req request.UpdateUserRoleReq, cla jwt.UserClaims) (response.Response, error) {

	if cla.UserRule != 2 {
		return response.Response{
			Msg:  "无权限",
			Code: 40001,
			Data: nil,
		}, nil
	}
	fmt.Println(1)
	err := c.service.UpdateUserRole(ctx, req.UserId, req.ProjectPermit, req.Role)
	if err != nil {
		return response.Response{}, err
	}
	return response.Response{
		Code: 200,
		Msg:  "修改成功",
		Data: nil,
	}, nil
}

// GetMyInfo 获取自己信息
// @Summary 获取自己信息
// @Description 获取用户名，邮箱，权限
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.UserInfo} "获取信息成功"
// @Failure 400 {object} response.Response{data=nil} "失败"
// @Security ApiKeyAuth
// @Router /api/v1/user/getMyInfo [get]
func (c *UserController) GetMyInfo(ctx *gin.Context, cla jwt.UserClaims) (response.Response, error) {

	user, err := c.service.GetMyInfo(ctx, cla.Uid)
	if err != nil {
		return response.Response{

			Code: 400,
			Data: nil,
			Msg:  "获取失败",
		}, err
	}

	return response.Response{
		Code: 200,
		Data: response.UserInfo{
			Name:   user.Name,
			Role:   user.UserRole,
			Email:  user.Email,
			Avatar: user.Avatar,
		},
		Msg: "获取成功",
	}, nil
}

// UpdateMyInfo 更新自己信息
// @Summary 更新用户信息
// @Description 更新当前用户的信息，如邮箱、名称和头像
// @Tags User
// @Accept json
// @Produce json
// @Param update body request.UpdateUserReq true "更新用户信息请求体"
// @Success 200 {object} response.Response "成功更新用户信息"
// @Failure 400 {object} response.Response "Invalid or expired token"
// @Security ApiKeyAuth
// @Router /api/v1/user/updateMyInfo [post]
func (c *UserController) UpdateMyInfo(ctx *gin.Context, req request.UpdateUserReq, cla jwt.UserClaims) (response.Response, error) {

	err := c.service.UpdateMyInfo(ctx, req, cla.Uid)
	if err != nil {
		return response.Response{}, err
	}
	return response.Response{
		Msg:  "更新用户信息成功",
		Code: 200,
		Data: nil,
	}, nil
}
