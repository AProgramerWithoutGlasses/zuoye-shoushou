package service

import (
	"errors"
	"goweb_staging/model"
	"time"
)

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Title            string    `json:"title" binding:"required"`
	Description      string    `json:"description"`
	StartTime        time.Time `json:"start_time" binding:"required"`
	EndTime          time.Time `json:"end_time" binding:"required"`
	AllowedFormats   []string  `json:"allowed_formats"`
	FilenameTemplate string    `json:"filename_template"`
	MaxFileSize      int64     `json:"max_file_size"`
	StudentIDs       []uint64  `json:"student_ids" binding:"required"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	AllowedFormats   []string  `json:"allowed_formats"`
	FilenameTemplate string    `json:"filename_template"`
	MaxFileSize      int64     `json:"max_file_size"`
	StudentIDs       []uint64  `json:"student_ids"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Tasks []model.Task `json:"tasks"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

// CreateTask 创建任务
func (s *Service) CreateTask(teacherID uint64, req *CreateTaskRequest) (*model.Task, error) {
	// 验证时间
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("截止时间不能早于开始时间")
	}

	// 创建任务
	task := &model.Task{
		Title:            req.Title,
		Description:      req.Description,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		AllowedFormats:   req.AllowedFormats,
		FilenameTemplate: req.FilenameTemplate,
		MaxFileSize:      req.MaxFileSize,
		TeacherID:        teacherID,
		Status:           model.TaskStatusDraft,
	}

	// 设置默认值
	if task.MaxFileSize == 0 {
		task.MaxFileSize = 10485760 // 10MB
	}

	err := s.dao.CreateTask(task)
	if err != nil {
		return nil, err
	}

	// 分配给学生
	if len(req.StudentIDs) > 0 {
		err = s.dao.AssignTaskToStudents(task.ID, req.StudentIDs)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}

// UpdateTask 更新任务
func (s *Service) UpdateTask(teacherID, taskID uint64, req *UpdateTaskRequest) (*model.Task, error) {
	// 获取任务
	task, err := s.dao.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	// 验证权限
	if task.TeacherID != teacherID {
		return nil, errors.New("无权限操作此任务")
	}

	// 验证任务状态
	if task.Status == model.TaskStatusCompleted {
		return nil, errors.New("已完成的任务不能修改")
	}

	// 更新字段
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if !req.StartTime.IsZero() {
		task.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		task.EndTime = req.EndTime
	}
	if req.AllowedFormats != nil {
		task.AllowedFormats = req.AllowedFormats
	}
	if req.FilenameTemplate != "" {
		task.FilenameTemplate = req.FilenameTemplate
	}
	if req.MaxFileSize > 0 {
		task.MaxFileSize = req.MaxFileSize
	}

	// 验证时间
	if task.EndTime.Before(task.StartTime) {
		return nil, errors.New("截止时间不能早于开始时间")
	}

	err = s.dao.UpdateTask(task)
	if err != nil {
		return nil, err
	}

	// 更新学生分配
	if req.StudentIDs != nil {
		err = s.dao.AssignTaskToStudents(task.ID, req.StudentIDs)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}

// PublishTask 发布任务
func (s *Service) PublishTask(teacherID, taskID uint64) error {
	task, err := s.dao.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	if task.TeacherID != teacherID {
		return errors.New("无权限操作此任务")
	}

	if task.Status != model.TaskStatusDraft {
		return errors.New("只能发布草稿状态的任务")
	}

	task.Status = model.TaskStatusActive
	return s.dao.UpdateTask(task)
}

// DeleteTask 删除任务
func (s *Service) DeleteTask(teacherID, taskID uint64) error {
	task, err := s.dao.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	if task.TeacherID != teacherID {
		return errors.New("无权限操作此任务")
	}

	if task.Status == model.TaskStatusActive {
		return errors.New("进行中的任务不能删除")
	}

	return s.dao.DeleteTask(taskID)
}

// GetTaskDetail 获取任务详情
func (s *Service) GetTaskDetail(userID, taskID uint64) (*model.Task, error) {
	task, err := s.dao.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	// 更新任务统计
	s.dao.UpdateTaskStatistics(taskID)

	return task, nil
}

// GetTeacherTasks 获取教师任务列表
func (s *Service) GetTeacherTasks(teacherID uint64, status string, page, size int) (*TaskListResponse, error) {
	offset := (page - 1) * size
	tasks, total, err := s.dao.GetTasksByTeacher(teacherID, status, size, offset)
	if err != nil {
		return nil, err
	}

	return &TaskListResponse{
		Tasks: tasks,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// GetStudentTasks 获取学生任务列表
func (s *Service) GetStudentTasks(studentID uint64, status string, page, size int) (*TaskListResponse, error) {
	offset := (page - 1) * size
	tasks, total, err := s.dao.GetTasksByStudent(studentID, status, size, offset)
	if err != nil {
		return nil, err
	}

	// 为每个任务更新统计
	for i := range tasks {
		s.dao.UpdateTaskStatistics(tasks[i].ID)
	}

	return &TaskListResponse{
		Tasks: tasks,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// GetTaskStatistics 获取任务统计
func (s *Service) GetTaskStatistics(teacherID, taskID uint64) (map[string]interface{}, error) {
	task, err := s.dao.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	if task.TeacherID != teacherID {
		return nil, errors.New("无权限查看此任务")
	}

	return s.dao.GetTaskSubmissionStatistics(taskID)
}

// GetTaskAllStudentsStatus 获取任务所有学生的提交状态
func (s *Service) GetTaskAllStudentsStatus(taskID uint64) (map[string]interface{}, error) {
	// 获取任务的所有学生
	students, err := s.dao.GetTaskStudents(taskID)
	if err != nil {
		return nil, err
	}

	// 获取任务的提交记录
	submissions, err := s.dao.GetSubmissionsByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	// 创建学生ID到提交记录的映射
	submissionMap := make(map[uint64]*model.Submission)
	for i := range submissions {
		submissionMap[submissions[i].StudentID] = &submissions[i]
	}

	// 构建结果
	var result []map[string]interface{}
	for _, student := range students {
		studentData := map[string]interface{}{
			"student":       student,
			"has_submitted": false,
			"submission":    nil,
		}

		if submission, exists := submissionMap[student.ID]; exists {
			studentData["has_submitted"] = true
			studentData["submission"] = submission
		}

		result = append(result, studentData)
	}

	return map[string]interface{}{
		"students":        result,
		"total_students":  len(students),
		"submitted_count": len(submissions),
	}, nil
}
