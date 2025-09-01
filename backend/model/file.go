package model

import (
	"time"

	"gorm.io/gorm"
)

// File 文件模型
type File struct {
	ID        uint64         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 基本信息
	OriginalName string `gorm:"type:varchar(255);not null" json:"original_name"` // 原始文件名
	StoredName   string `gorm:"type:varchar(255);not null" json:"stored_name"`   // 存储文件名
	FilePath     string `gorm:"type:varchar(500);not null" json:"file_path"`     // 文件路径
	FileSize     int64  `gorm:"not null" json:"file_size"`                       // 文件大小(字节)
	ContentType  string `gorm:"type:varchar(100)" json:"content_type"`           // 文件类型
	FileHash     string `gorm:"type:varchar(64);index" json:"file_hash"`         // 文件哈希值

	// 关联信息
	SubmissionID uint64 `gorm:"not null;index" json:"submission_id"`
	StudentID    uint64 `gorm:"not null;index" json:"student_id"`
	TaskID       uint64 `gorm:"not null;index" json:"task_id"`

	// 文件状态
	IsDeleted bool `gorm:"default:false" json:"is_deleted"` // 是否已删除
}

// TableName 设置表名
func (File) TableName() string {
	return "files"
}
