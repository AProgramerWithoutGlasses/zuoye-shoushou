package dao

import (
	"goweb_staging/model"
	"time"
)

// CreateTask 创建任务
func (dao *Dao) CreateTask(task *model.Task) error {
	return dao.db.Create(task).Error
}

// GetTaskByID 根据ID获取任务详情
func (dao *Dao) GetTaskByID(id uint64) (*model.Task, error) {
	var task model.Task
	err := dao.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask 更新任务
func (dao *Dao) UpdateTask(task *model.Task) error {
	return dao.db.Save(task).Error
}

// DeleteTask 删除任务
func (dao *Dao) DeleteTask(id uint64) error {
	return dao.db.Delete(&model.Task{}, id).Error
}

// GetTasksByTeacher 获取教师发布的任务列表
func (dao *Dao) GetTasksByTeacher(teacherID uint64, status string, limit, offset int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := dao.db.Model(&model.Task{}).Where("teacher_id = ?", teacherID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&tasks).Error

	return tasks, total, err
}

// GetTasksByStudent 获取学生的任务列表
func (dao *Dao) GetTasksByStudent(studentID uint64, status string, limit, offset int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := dao.db.Model(&model.Task{}).
		Joins("JOIN task_students ON tasks.id = task_students.task_id").
		Where("task_students.student_id = ?", studentID)

	if status != "" {
		query = query.Where("tasks.status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Teacher").
		Order("tasks.end_time ASC").
		Limit(limit).Offset(offset).Find(&tasks).Error

	return tasks, total, err
}

// AssignTaskToStudents 将任务分配给学生
func (dao *Dao) AssignTaskToStudents(taskID uint64, studentIDs []uint64) error {
	tx := dao.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 清除原有分配
	if err := tx.Where("task_id = ?", taskID).Delete(&model.TaskStudent{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新分配
	for _, studentID := range studentIDs {
		taskStudent := model.TaskStudent{
			TaskID:    taskID,
			StudentID: studentID,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&taskStudent).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新任务的学生总数
	if err := tx.Model(&model.Task{}).Where("id = ?", taskID).
		Update("total_students", len(studentIDs)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetTaskStudents 获取任务的学生列表
func (dao *Dao) GetTaskStudents(taskID uint64) ([]model.User, error) {
	var students []model.User
	err := dao.db.Joins("JOIN task_students ON users.id = task_students.student_id").
		Where("task_students.task_id = ?", taskID).
		Find(&students).Error
	return students, err
}

// UpdateTaskStatistics 更新任务统计信息
func (dao *Dao) UpdateTaskStatistics(taskID uint64) error {
	var submittedCount, onTimeCount int64

	// 统计已提交数量
	dao.db.Model(&model.Submission{}).
		Where("task_id = ? AND status IN ?", taskID, []model.SubmissionStatus{
			model.SubmissionStatusSubmitted,
			model.SubmissionStatusLate,
			model.SubmissionStatusReviewed,
		}).Count(&submittedCount)

	// 统计按时提交数量
	dao.db.Model(&model.Submission{}).
		Where("task_id = ? AND is_on_time = true", taskID).Count(&onTimeCount)

	// 更新任务统计
	return dao.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"submitted_count": submittedCount,
		"on_time_count":   onTimeCount,
	}).Error
}
