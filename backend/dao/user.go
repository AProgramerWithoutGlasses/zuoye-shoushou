package dao

import (
	"goweb_staging/model"
)

// CreateUser 创建用户
func (dao *Dao) CreateUser(user *model.User) error {
	return dao.db.Create(user).Error
}

// GetUserByUsername 根据用户名获取用户
func (dao *Dao) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := dao.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByWxOpenID 根据微信OpenID获取用户
func (dao *Dao) GetUserByWxOpenID(openID string) (*model.User, error) {
	var user model.User
	err := dao.db.Where("wx_open_id = ?", openID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func (dao *Dao) GetUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := dao.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (dao *Dao) UpdateUser(user *model.User) error {
	return dao.db.Save(user).Error
}

// GetStudentsByClass 根据班级获取学生列表
func (dao *Dao) GetStudentsByClass(major, grade, class string) ([]model.User, error) {
	var users []model.User
	query := dao.db.Where("role = ?", model.RoleStudent)

	if major != "" {
		query = query.Where("major = ?", major)
	}
	if grade != "" {
		query = query.Where("grade = ?", grade)
	}
	if class != "" {
		query = query.Where("class = ?", class)
	}

	err := query.Find(&users).Error
	return users, err
}

// GetStudentsByIDs 根据ID列表获取学生
func (dao *Dao) GetStudentsByIDs(ids []uint64) ([]model.User, error) {
	var users []model.User
	err := dao.db.Where("id IN ? AND role = ?", ids, model.RoleStudent).Find(&users).Error
	return users, err
}

// GetTeacherList 获取教师列表
func (dao *Dao) GetTeacherList(limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	err := dao.db.Model(&model.User{}).Where("role = ?", model.RoleTeacher).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = dao.db.Where("role = ?", model.RoleTeacher).
		Limit(limit).Offset(offset).Find(&users).Error

	return users, total, err
}
