package model

import (
	"time"

	"gorm.io/gorm"
)

// SubmissionStatus 提交状态
type SubmissionStatus string

const (
	SubmissionStatusPending   SubmissionStatus = "pending"   // 未提交
	SubmissionStatusSubmitted SubmissionStatus = "submitted" // 已提交
	SubmissionStatusLate      SubmissionStatus = "late"      // 迟交
	SubmissionStatusReviewed  SubmissionStatus = "reviewed"  // 已批阅
)

// Submission 提交记录模型
type Submission struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联信息
	TaskID    uint64 `gorm:"not null;index" json:"task_id"`
	StudentID uint64 `gorm:"not null;index" json:"student_id"`

	// 提交信息
	Status      SubmissionStatus `gorm:"type:enum('pending','submitted','late','reviewed');default:'pending'" json:"status"`
	SubmittedAt *time.Time       `json:"submitted_at"`                    // 提交时间
	IsOnTime    bool             `gorm:"default:false" json:"is_on_time"` // 是否按时提交

	// 文件信息
	Files []File `gorm:"foreignKey:SubmissionID" json:"files,omitempty"`

	// 批阅信息
	Score      *float64   `json:"score"`                    // 分数
	Comment    string     `gorm:"type:text" json:"comment"` // 批阅评语
	ReviewedAt *time.Time `json:"reviewed_at"`              // 批阅时间
	ReviewedBy *uint64    `json:"reviewed_by"`              // 批阅教师ID
}

// TableName 设置表名
func (Submission) TableName() string {
	return "submissions"
}

// 创建唯一索引
func (Submission) Index() []string {
	return []string{
		"idx_task_student:task_id,student_id", // 确保一个学生在一个任务中只能有一条提交记录
	}
}
