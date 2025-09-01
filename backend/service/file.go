package service

import (
	"goweb_staging/model"
)

// GetFileByID 获取文件信息
func (s *Service) GetFileByID(fileID uint64) (*model.File, error) {
	return s.dao.GetFileByID(fileID)
}

// GetFilesByTask 获取任务的所有文件
func (s *Service) GetFilesByTask(taskID uint64) ([]model.File, error) {
	return s.dao.GetFilesByTask(taskID)
}

// GetFileStatistics 获取文件统计信息
func (s *Service) GetFileStatistics(taskID uint64) (map[string]interface{}, error) {
	return s.dao.GetFileStatistics(taskID)
}
