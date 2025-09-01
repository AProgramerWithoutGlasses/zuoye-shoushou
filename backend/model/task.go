package model

import (
	"time"

	"gorm.io/gorm"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusDraft     TaskStatus = "draft"     // 草稿
	TaskStatusActive    TaskStatus = "active"    // 进行中
	TaskStatusExpired   TaskStatus = "expired"   // 已截止
	TaskStatusCompleted TaskStatus = "completed" // 已完成
)

// Task 任务模型
type Task struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 基本信息
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`                                         // 任务标题
	Description string     `gorm:"type:text" json:"description"`                                                    // 任务描述
	Status      TaskStatus `gorm:"type:enum('draft','active','expired','completed');default:'draft'" json:"status"` // 任务状态

	// 时间设置
	StartTime time.Time `gorm:"not null" json:"start_time"` // 开始时间
	EndTime   time.Time `gorm:"not null" json:"end_time"`   // 截止时间

	// 文件要求
	AllowedFormats   []string `gorm:"serializer:json" json:"allowed_formats"`     // 允许的文件格式
	FilenameTemplate string   `gorm:"type:varchar(200)" json:"filename_template"` // 文件名模板
	MaxFileSize      int64    `gorm:"default:10485760" json:"max_file_size"`      // 最大文件大小(字节)

	// 关联信息
	TeacherID uint64 `gorm:"not null;index" json:"teacher_id"`              // 发布教师ID
	Teacher   User   `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"` // 发布教师信息

	// 目标学生 (多对多关系)
	Students []User `gorm:"many2many:task_students;" json:"students,omitempty"`

	// 统计信息
	TotalStudents  int `gorm:"default:0" json:"total_students"`  // 总学生数
	SubmittedCount int `gorm:"default:0" json:"submitted_count"` // 已提交数
	OnTimeCount    int `gorm:"default:0" json:"on_time_count"`   // 按时提交数
}

// TableName 设置表名
func (Task) TableName() string {
	return "tasks"
}

// TaskStudent 任务学生关联表
type TaskStudent struct {
	TaskID    uint64    `gorm:"primaryKey;column:task_id" json:"task_id"`
	StudentID uint64    `gorm:"primaryKey;column:student_id" json:"student_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 设置表名
func (TaskStudent) TableName() string {
	return "task_students"
}

// BeforeCreate GORM钩子，在创建前设置时间
func (ts *TaskStudent) BeforeCreate(tx *gorm.DB) error {
	if ts.CreatedAt.IsZero() {
		ts.CreatedAt = time.Now()
	}
	return nil
}
