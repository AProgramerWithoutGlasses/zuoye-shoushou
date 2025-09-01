package dao

import (
	"goweb_staging/model"
)

// CreateFile 创建文件记录
func (dao *Dao) CreateFile(file *model.File) error {
	return dao.db.Create(file).Error
}

// GetFileByID 根据ID获取文件
func (dao *Dao) GetFileByID(id uint64) (*model.File, error) {
	var file model.File
	err := dao.db.First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// GetFilesBySubmission 根据提交记录获取文件列表
func (dao *Dao) GetFilesBySubmission(submissionID uint64) ([]model.File, error) {
	var files []model.File
	err := dao.db.Where("submission_id = ? AND is_deleted = false", submissionID).Find(&files).Error
	return files, err
}

// GetFilesByTask 根据任务获取所有文件
func (dao *Dao) GetFilesByTask(taskID uint64) ([]model.File, error) {
	var files []model.File
	err := dao.db.Where("task_id = ? AND is_deleted = false", taskID).Find(&files).Error
	return files, err
}

// UpdateFile 更新文件信息
func (dao *Dao) UpdateFile(file *model.File) error {
	return dao.db.Save(file).Error
}

// DeleteFile 软删除文件
func (dao *Dao) DeleteFile(id uint64) error {
	return dao.db.Model(&model.File{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetFileByHash 根据文件哈希获取文件（用于去重）
func (dao *Dao) GetFileByHash(hash string) (*model.File, error) {
	var file model.File
	err := dao.db.Where("file_hash = ? AND is_deleted = false", hash).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// BatchCreateFiles 批量创建文件记录
func (dao *Dao) BatchCreateFiles(files []model.File) error {
	if len(files) == 0 {
		return nil
	}
	return dao.db.Create(&files).Error
}

// GetFileStatistics 获取文件统计信息
func (dao *Dao) GetFileStatistics(taskID uint64) (map[string]interface{}, error) {
	var result struct {
		TotalFiles int64 `json:"total_files"`
		TotalSize  int64 `json:"total_size"`
	}

	// 统计文件数量和总大小
	err := dao.db.Model(&model.File{}).
		Select("COUNT(*) as total_files, COALESCE(SUM(file_size), 0) as total_size").
		Where("task_id = ? AND is_deleted = false", taskID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_files": result.TotalFiles,
		"total_size":  result.TotalSize,
	}, nil
}
