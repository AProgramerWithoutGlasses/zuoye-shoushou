package server

import (
	"goweb_staging/logger"
	"goweb_staging/middleware"
	"goweb_staging/pkg/settings"
	"goweb_staging/service"

	"github.com/gin-gonic/gin"
)

var svc *service.Service

func initRouter() *gin.Engine {
	// 创建一个默认的路由引擎
	r := gin.Default()

	// 注册使用的中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middleware.CORSMiddleware())

	// 公开路由（不需要认证）
	public := r.Group("/api")
	{
		// 认证相关
		public.POST("/auth/wx-login", wxLogin) // 微信授权登录
		public.POST("/auth/login", login)      // 账号密码登录
	}

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		// 用户相关
		auth.GET("/user/info", getUserInfo)       // 获取用户信息
		auth.PUT("/user/info", updateUserInfo)    // 更新用户信息
		auth.POST("/user/bind-wx", bindWxAccount) // 绑定微信账号

		// 任务相关
		auth.POST("/tasks", createTask)                                  // 创建任务（教师）
		auth.GET("/tasks", getTeacherTasks)                              // 获取教师任务列表
		auth.GET("/tasks/student", getStudentTasks)                      // 获取学生任务列表
		auth.GET("/tasks/:id", getTaskDetail)                            // 获取任务详情
		auth.GET("/tasks/:id/students-status", getTaskAllStudentsStatus) // 获取任务所有学生提交状态
		auth.PUT("/tasks/:id", updateTask)                               // 更新任务（教师）
		auth.POST("/tasks/:id/publish", publishTask)                     // 发布任务（教师）
		auth.DELETE("/tasks/:id", deleteTask)                            // 删除任务（教师）
		auth.GET("/tasks/:id/statistics", getTaskStatistics)             // 获取任务统计（教师）

		// 提交相关
		auth.POST("/tasks/:id/submit", submitTask)              // 提交任务（学生）
		auth.GET("/tasks/:id/submission", getStudentSubmission) // 获取学生提交记录
		auth.GET("/tasks/:id/submissions", getTaskSubmissions)  // 获取任务的所有提交记录（教师）
		auth.GET("/submissions", getStudentSubmissions)         // 获取学生提交历史
		auth.GET("/submissions/:id", getSubmissionDetail)       // 获取提交详情
		auth.POST("/submissions/:id/review", reviewSubmission)  // 批阅提交（教师）

		// 文件相关
		auth.POST("/files/upload", uploadFile)        // 文件上传
		auth.GET("/files/:id/download", downloadFile) // 文件下载

		// 测试接口
		auth.POST("/test", test)
	}

	return r
}

func Init(app *settings.AppConfig) *gin.Engine {
	svc = service.InitService(app)
	return initRouter()
}
