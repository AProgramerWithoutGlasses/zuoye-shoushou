package dao

import (
	"goweb_staging/model"
	"time"
)

// CreateSubmission 创建提交记录
func (dao *Dao) CreateSubmission(submission *model.Submission) error {
	return dao.db.Create(submission).Error
}

// GetSubmissionByTaskAndStudent 根据任务和学生获取提交记录
func (dao *Dao) GetSubmissionByTaskAndStudent(taskID, studentID uint64) (*model.Submission, error) {
	var submission model.Submission
	err := dao.db.Preload("Files").
		Where("task_id = ? AND student_id = ?", taskID, studentID).
		First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

// GetSubmissionByID 根据ID获取提交记录
func (dao *Dao) GetSubmissionByID(id uint64) (*model.Submission, error) {
	var submission model.Submission
	err := dao.db.Preload("Task").Preload("Student").Preload("Files").
		First(&submission, id).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

// UpdateSubmission 更新提交记录
func (dao *Dao) UpdateSubmission(submission *model.Submission) error {
	return dao.db.Save(submission).Error
}

// GetSubmissionsByTask 获取任务的所有提交记录
func (dao *Dao) GetSubmissionsByTask(taskID uint64, status string, limit, offset int) ([]model.Submission, int64, error) {
	var submissions []model.Submission
	var total int64

	query := dao.db.Model(&model.Submission{}).Where("task_id = ?", taskID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Student").Preload("Files").
		Order("submitted_at DESC, student_id ASC").
		Limit(limit).Offset(offset).Find(&submissions).Error

	return submissions, total, err
}

// GetSubmissionsByStudent 获取学生的提交记录
func (dao *Dao) GetSubmissionsByStudent(studentID uint64, limit, offset int) ([]model.Submission, int64, error) {
	var submissions []model.Submission
	var total int64

	query := dao.db.Model(&model.Submission{}).Where("student_id = ?", studentID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Task").Preload("Files").
		Order("submitted_at DESC").
		Limit(limit).Offset(offset).Find(&submissions).Error

	return submissions, total, err
}

// SubmitTask 提交任务
func (dao *Dao) SubmitTask(submission *model.Submission) error {
	tx := dao.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新提交状态和时间
	now := time.Now()
	submission.SubmittedAt = &now
	submission.Status = model.SubmissionStatusSubmitted

	// 检查是否按时提交
	var task model.Task
	if err := tx.First(&task, submission.TaskID).Error; err != nil {
		tx.Rollback()
		return err
	}

	submission.IsOnTime = now.Before(task.EndTime) || now.Equal(task.EndTime)
	if !submission.IsOnTime {
		submission.Status = model.SubmissionStatusLate
	}

	// 保存提交记录
	if err := tx.Save(submission).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetTaskSubmissionStatistics 获取任务提交统计
func (dao *Dao) GetTaskSubmissionStatistics(taskID uint64) (map[string]interface{}, error) {
	var result struct {
		TotalStudents  int64 `json:"total_students"`
		SubmittedCount int64 `json:"submitted_count"`
		OnTimeCount    int64 `json:"on_time_count"`
		LateCount      int64 `json:"late_count"`
		PendingCount   int64 `json:"pending_count"`
	}

	// 获取总学生数
	dao.db.Model(&model.TaskStudent{}).Where("task_id = ?", taskID).Count(&result.TotalStudents)

	// 获取已提交数
	dao.db.Model(&model.Submission{}).
		Where("task_id = ? AND status IN ?", taskID, []model.SubmissionStatus{
			model.SubmissionStatusSubmitted,
			model.SubmissionStatusLate,
			model.SubmissionStatusReviewed,
		}).Count(&result.SubmittedCount)

	// 获取按时提交数
	dao.db.Model(&model.Submission{}).
		Where("task_id = ? AND is_on_time = true", taskID).Count(&result.OnTimeCount)

	// 获取迟交数
	dao.db.Model(&model.Submission{}).
		Where("task_id = ? AND status = ?", taskID, model.SubmissionStatusLate).Count(&result.LateCount)

	// 计算未提交数
	result.PendingCount = result.TotalStudents - result.SubmittedCount

	// 计算提交率
	var submitRate float64
	if result.TotalStudents > 0 {
		submitRate = float64(result.SubmittedCount) / float64(result.TotalStudents) * 100
	}

	// 计算按时率
	var onTimeRate float64
	if result.SubmittedCount > 0 {
		onTimeRate = float64(result.OnTimeCount) / float64(result.SubmittedCount) * 100
	}

	return map[string]interface{}{
		"total_students":  result.TotalStudents,
		"submitted_count": result.SubmittedCount,
		"on_time_count":   result.OnTimeCount,
		"late_count":      result.LateCount,
		"pending_count":   result.PendingCount,
		"submit_rate":     submitRate,
		"on_time_rate":    onTimeRate,
	}, nil
}

// GetSubmissionsByTaskID 获取任务的所有提交记录
func (dao *Dao) GetSubmissionsByTaskID(taskID uint64) ([]model.Submission, error) {
	var submissions []model.Submission
	err := dao.db.Where("task_id = ?", taskID).
		Preload("Files").
		Order("submitted_at ASC").
		Find(&submissions).Error
	return submissions, err
}
