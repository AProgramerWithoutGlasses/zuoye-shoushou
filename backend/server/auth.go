package server

import (
	"goweb_staging/pkg/response"
	"goweb_staging/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// wxLogin 微信授权登录
func wxLogin(c *gin.Context) {
	var req service.WxLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	data, err := svc.WxLogin(&req)
	if err != nil {
		zap.L().Error("wx login failed", zap.Error(err))
		response.FailWithMsg(c, response.LoginErrCode, err.Error())
		return
	}

	response.Success(c, data)
}

// login 账号密码登录
func login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	data, err := svc.Login(&req)
	if err != nil {
		zap.L().Error("login failed", zap.Error(err))
		response.FailWithMsg(c, response.LoginErrCode, err.Error())
		return
	}

	response.Success(c, data)
}

// bindWxAccount 绑定微信账号
func bindWxAccount(c *gin.Context) {
	var req struct {
		WxCode string `json:"wx_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	userID := getCurrentUserID(c)
	err := svc.BindWxAccount(userID, req.WxCode)
	if err != nil {
		zap.L().Error("bind wx account failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, nil)
}

// getUserInfo 获取用户信息
func getUserInfo(c *gin.Context) {
	userID := getCurrentUserID(c)
	user, err := svc.GetUserInfo(userID)
	if err != nil {
		zap.L().Error("get user info failed", zap.Error(err))
		response.Fail(c, response.ServerErrCode)
		return
	}

	response.Success(c, user)
}

// updateUserInfo 更新用户信息
func updateUserInfo(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	userID := getCurrentUserID(c)
	err := svc.UpdateUserInfo(userID, req)
	if err != nil {
		zap.L().Error("update user info failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, nil)
}

// getCurrentUserID 从context中获取当前用户ID
func getCurrentUserID(c *gin.Context) uint64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	if id, ok := userID.(uint64); ok {
		return id
	}
	return 0
}
