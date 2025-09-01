package service

import (
	"errors"
	"goweb_staging/model"
	"time"
)

// SubmitTaskRequest 提交任务请求
type SubmitTaskRequest struct {
	Files []FileInfo `json:"files" binding:"required"`
}

// FileInfo 文件信息
type FileInfo struct {
	OriginalName string `json:"original_name" binding:"required"`
	StoredName   string `json:"stored_name" binding:"required"`
	FilePath     string `json:"file_path" binding:"required"`
	FileSize     int64  `json:"file_size" binding:"required"`
	ContentType  string `json:"content_type"`
	FileHash     string `json:"file_hash"`
}

// SubmissionListResponse 提交列表响应
type SubmissionListResponse struct {
	Submissions []model.Submission `json:"submissions"`
	Total       int64              `json:"total"`
	Page        int                `json:"page"`
	Size        int                `json:"size"`
}

// SubmitTask 提交任务
func (s *Service) SubmitTask(studentID, taskID uint64, req *SubmitTaskRequest) (*model.Submission, error) {
	// 获取任务信息
	task, err := s.dao.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	// 检查任务状态
	if task.Status != model.TaskStatusActive {
		return nil, errors.New("任务未开放提交")
	}

	// 检查提交时间
	now := time.Now()
	if now.Before(task.StartTime) {
		return nil, errors.New("任务尚未开始")
	}

	// 检查文件数量
	if len(req.Files) == 0 {
		return nil, errors.New("请至少上传一个文件")
	}

	// 验证文件格式和大小
	for _, file := range req.Files {
		if file.FileSize > task.MaxFileSize {
			return nil, errors.New("文件大小超过限制")
		}

		// TODO: 验证文件格式
	}

	// 查找或创建提交记录
	submission, err := s.dao.GetSubmissionByTaskAndStudent(taskID, studentID)
	if err != nil {
		// 创建新的提交记录
		submission = &model.Submission{
			TaskID:    taskID,
			StudentID: studentID,
			Status:    model.SubmissionStatusPending,
		}
		err = s.dao.CreateSubmission(submission)
		if err != nil {
			return nil, err
		}
	}

	// 如果已经提交过，检查是否允许重新提交
	if submission.Status != model.SubmissionStatusPending && now.After(task.EndTime) {
		return nil, errors.New("任务已截止，不能重新提交")
	}

	// 创建文件记录
	var files []model.File
	for _, fileInfo := range req.Files {
		file := model.File{
			OriginalName: fileInfo.OriginalName,
			StoredName:   fileInfo.StoredName,
			FilePath:     fileInfo.FilePath,
			FileSize:     fileInfo.FileSize,
			ContentType:  fileInfo.ContentType,
			FileHash:     fileInfo.FileHash,
			SubmissionID: submission.ID,
			StudentID:    studentID,
			TaskID:       taskID,
		}
		files = append(files, file)
	}

	err = s.dao.BatchCreateFiles(files)
	if err != nil {
		return nil, err
	}

	// 更新提交状态
	err = s.dao.SubmitTask(submission)
	if err != nil {
		return nil, err
	}

	// 更新任务统计
	s.dao.UpdateTaskStatistics(taskID)

	return submission, nil
}

// GetStudentSubmission 获取学生的提交记录
func (s *Service) GetStudentSubmission(studentID, taskID uint64) (*model.Submission, error) {
	return s.dao.GetSubmissionByTaskAndStudent(taskID, studentID)
}

// GetSubmissionDetail 获取提交详情
func (s *Service) GetSubmissionDetail(userID, submissionID uint64) (*model.Submission, error) {
	submission, err := s.dao.GetSubmissionByID(submissionID)
	if err != nil {
		return nil, err
	}

	// 权限检查
	user, err := s.dao.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 学生只能看自己的提交，教师只能看自己发布任务的提交
	if user.Role == model.RoleStudent && submission.StudentID != userID {
		return nil, errors.New("无权限查看此提交")
	}
	if user.Role == model.RoleTeacher {
		// 查询任务信息来获取教师ID
		task, err := s.dao.GetTaskByID(submission.TaskID)
		if err != nil {
			return nil, err
		}
		if task.TeacherID != userID {
			return nil, errors.New("无权限查看此提交")
		}
	}

	return submission, nil
}

// GetTaskSubmissions 获取任务的所有提交记录
func (s *Service) GetTaskSubmissions(teacherID, taskID uint64, status string, page, size int) (*SubmissionListResponse, error) {
	// 验证权限
	task, err := s.dao.GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}

	if task.TeacherID != teacherID {
		return nil, errors.New("无权限查看此任务的提交")
	}

	offset := (page - 1) * size
	submissions, total, err := s.dao.GetSubmissionsByTask(taskID, status, size, offset)
	if err != nil {
		return nil, err
	}

	return &SubmissionListResponse{
		Submissions: submissions,
		Total:       total,
		Page:        page,
		Size:        size,
	}, nil
}

// GetStudentSubmissions 获取学生的提交历史
func (s *Service) GetStudentSubmissions(studentID uint64, page, size int) (*SubmissionListResponse, error) {
	offset := (page - 1) * size
	submissions, total, err := s.dao.GetSubmissionsByStudent(studentID, size, offset)
	if err != nil {
		return nil, err
	}

	return &SubmissionListResponse{
		Submissions: submissions,
		Total:       total,
		Page:        page,
		Size:        size,
	}, nil
}

// ReviewSubmission 批阅提交
func (s *Service) ReviewSubmission(teacherID, submissionID uint64, score *float64, comment string) error {
	submission, err := s.dao.GetSubmissionByID(submissionID)
	if err != nil {
		return err
	}

	// 验证权限
	task, err := s.dao.GetTaskByID(submission.TaskID)
	if err != nil {
		return err
	}
	if task.TeacherID != teacherID {
		return errors.New("无权限批阅此提交")
	}

	// 更新批阅信息
	now := time.Now()
	submission.Score = score
	submission.Comment = comment
	submission.ReviewedAt = &now
	submission.ReviewedBy = &teacherID
	submission.Status = model.SubmissionStatusReviewed

	return s.dao.UpdateSubmission(submission)
}
