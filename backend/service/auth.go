package service

import (
	"errors"
	"goweb_staging/model"
	"goweb_staging/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// WxLoginRequest 微信登录请求
type WxLoginRequest struct {
	Code string `json:"code" binding:"required"` // 微信小程序code
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名（学号或手机号）
	Password string `json:"password" binding:"required"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// WxLogin 微信授权登录
func (s *Service) WxLogin(req *WxLoginRequest) (*LoginResponse, error) {
	// TODO: 调用微信API获取openid
	// 这里暂时模拟
	openID := "mock_openid_" + req.Code

	// 查找用户
	user, err := s.dao.GetUserByWxOpenID(openID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户未绑定，请先进行账号验证")
		}
		return nil, err
	}

	// 检查用户状态
	if !user.IsActive {
		return nil, errors.New("用户已被禁用")
	}

	// 生成token
	token, err := jwt.GenToken(user.Username, user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login 账号密码登录
func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	user, err := s.dao.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if !user.IsActive {
		return nil, errors.New("用户已被禁用")
	}

	// 生成token
	token, err := jwt.GenToken(user.Username, user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// BindWxAccount 绑定微信账号
func (s *Service) BindWxAccount(userID uint64, wxCode string) error {
	// TODO: 调用微信API获取openid
	openID := "mock_openid_" + wxCode

	// 检查openid是否已被绑定
	_, err := s.dao.GetUserByWxOpenID(openID)
	if err == nil {
		return errors.New("该微信账号已被绑定")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 获取用户
	user, err := s.dao.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 绑定微信
	user.WxOpenID = openID
	return s.dao.UpdateUser(user)
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(userID uint64) (*model.User, error) {
	return s.dao.GetUserByID(userID)
}

// UpdateUserInfo 更新用户信息
func (s *Service) UpdateUserInfo(userID uint64, updates map[string]interface{}) error {
	user, err := s.dao.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 更新允许的字段
	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}
	if phone, ok := updates["phone"].(string); ok {
		user.Phone = phone
	}
	if department, ok := updates["department"].(string); ok {
		user.Department = department
	}

	return s.dao.UpdateUser(user)
}
