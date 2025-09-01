package server

import (
	"goweb_staging/pkg/response"
	"goweb_staging/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// submitTask 提交任务
func submitTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	var req service.SubmitTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	studentID := getCurrentUserID(c)
	submission, err := svc.SubmitTask(studentID, taskID, &req)
	if err != nil {
		zap.L().Error("submit task failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, submission)
}

// getStudentSubmission 获取学生的提交记录
func getStudentSubmission(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	studentID := getCurrentUserID(c)
	submission, err := svc.GetStudentSubmission(studentID, taskID)
	if err != nil {
		zap.L().Error("get student submission failed", zap.Error(err))
		response.Fail(c, response.ServerErrCode)
		return
	}

	response.Success(c, submission)
}

// getSubmissionDetail 获取提交详情
func getSubmissionDetail(c *gin.Context) {
	submissionIDStr := c.Param("id")
	submissionID, err := strconv.ParseUint(submissionIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	userID := getCurrentUserID(c)
	submission, err := svc.GetSubmissionDetail(userID, submissionID)
	if err != nil {
		zap.L().Error("get submission detail failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, submission)
}

// getTaskSubmissions 获取任务的所有提交记录
func getTaskSubmissions(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	teacherID := getCurrentUserID(c)
	data, err := svc.GetTaskSubmissions(teacherID, taskID, status, page, size)
	if err != nil {
		zap.L().Error("get task submissions failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, data)
}

// getStudentSubmissions 获取学生的提交历史
func getStudentSubmissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	studentID := getCurrentUserID(c)
	data, err := svc.GetStudentSubmissions(studentID, page, size)
	if err != nil {
		zap.L().Error("get student submissions failed", zap.Error(err))
		response.Fail(c, response.ServerErrCode)
		return
	}

	response.Success(c, data)
}

// reviewSubmission 批阅提交
func reviewSubmission(c *gin.Context) {
	submissionIDStr := c.Param("id")
	submissionID, err := strconv.ParseUint(submissionIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	var req struct {
		Score   *float64 `json:"score"`
		Comment string   `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	teacherID := getCurrentUserID(c)
	err = svc.ReviewSubmission(teacherID, submissionID, req.Score, req.Comment)
	if err != nil {
		zap.L().Error("review submission failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, nil)
}
