package model

import (
	"time"

	"gorm.io/gorm"
)

// UserRole 用户角色类型
type UserRole string

const (
	RoleStudent UserRole = "student" // 学生
	RoleTeacher UserRole = "teacher" // 教师
)

// User 用户模型
type User struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 基本信息
	Username string   `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"` // 学号或手机号
	Password string   `gorm:"type:varchar(255);not null" json:"-"`                   // 密码
	Name     string   `gorm:"type:varchar(50);not null" json:"name"`                 // 真实姓名
	Role     UserRole `gorm:"type:enum('student','teacher');not null" json:"role"`   // 用户角色
	WxOpenID string   `gorm:"type:varchar(100);uniqueIndex" json:"wx_open_id"`       // 微信openid

	// 学生特有字段
	StudentID string `gorm:"type:varchar(20);index" json:"student_id"` // 学号
	Major     string `gorm:"type:varchar(100)" json:"major"`           // 专业
	Grade     string `gorm:"type:varchar(20)" json:"grade"`            // 年级
	Class     string `gorm:"type:varchar(50)" json:"class"`            // 班级

	// 教师特有字段
	TeacherID  string `gorm:"type:varchar(20);index" json:"teacher_id"` // 教师工号
	Phone      string `gorm:"type:varchar(20);index" json:"phone"`      // 手机号
	Department string `gorm:"type:varchar(100)" json:"department"`      // 部门

	// 状态字段
	IsActive bool `gorm:"default:true" json:"is_active"` // 是否激活
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}
