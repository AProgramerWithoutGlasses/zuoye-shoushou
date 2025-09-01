package server

import (
	"goweb_staging/pkg/response"
	"goweb_staging/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// createTask 创建任务
func createTask(c *gin.Context) {
	var req service.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	teacherID := getCurrentUserID(c)
	task, err := svc.CreateTask(teacherID, &req)
	if err != nil {
		zap.L().Error("create task failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, task)
}

// updateTask 更新任务
func updateTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	var req service.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	teacherID := getCurrentUserID(c)
	task, err := svc.UpdateTask(teacherID, taskID, &req)
	if err != nil {
		zap.L().Error("update task failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, task)
}

// publishTask 发布任务
func publishTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	teacherID := getCurrentUserID(c)
	err = svc.PublishTask(teacherID, taskID)
	if err != nil {
		zap.L().Error("publish task failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, nil)
}

// deleteTask 删除任务
func deleteTask(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	teacherID := getCurrentUserID(c)
	err = svc.DeleteTask(teacherID, taskID)
	if err != nil {
		zap.L().Error("delete task failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, nil)
}

// getTaskDetail 获取任务详情
func getTaskDetail(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	userID := getCurrentUserID(c)
	task, err := svc.GetTaskDetail(userID, taskID)
	if err != nil {
		zap.L().Error("get task detail failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, task)
}

// getTeacherTasks 获取教师任务列表
func getTeacherTasks(c *gin.Context) {
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
	data, err := svc.GetTeacherTasks(teacherID, status, page, size)
	if err != nil {
		zap.L().Error("get teacher tasks failed", zap.Error(err))
		response.Fail(c, response.ServerErrCode)
		return
	}

	response.Success(c, data)
}

// getStudentTasks 获取学生任务列表
func getStudentTasks(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	studentID := getCurrentUserID(c)
	data, err := svc.GetStudentTasks(studentID, status, page, size)
	if err != nil {
		zap.L().Error("get student tasks failed", zap.Error(err))
		response.Fail(c, response.ServerErrCode)
		return
	}

	response.Success(c, data)
}

// getTaskStatistics 获取任务统计
func getTaskStatistics(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	teacherID := getCurrentUserID(c)
	data, err := svc.GetTaskStatistics(teacherID, taskID)
	if err != nil {
		zap.L().Error("get task statistics failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, data)
}

// getTaskAllStudentsStatus 获取任务所有学生的提交状态
func getTaskAllStudentsStatus(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	data, err := svc.GetTaskAllStudentsStatus(taskID)
	if err != nil {
		zap.L().Error("get task all students status failed", zap.Error(err))
		response.FailWithMsg(c, response.ServerErrCode, err.Error())
		return
	}

	response.Success(c, data)
}
