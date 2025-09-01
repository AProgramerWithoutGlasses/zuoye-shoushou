package dao

import (
	"goweb_staging/model"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// InitData 初始化数据库数据
func (dao *Dao) InitData() error {
	// 0. 先清理可能存在的表（可选）
	if err := dao.cleanDatabase(); err != nil {
		return err
	}

	// 1. 先创建用户表（无外键依赖）
	if err := dao.db.AutoMigrate(&model.User{}); err != nil {
		return err
	}

	// 2. 创建任务表
	if err := dao.db.AutoMigrate(&model.Task{}); err != nil {
		return err
	}

	// 3. 手动创建关联表，避免GORM的复合主键问题
	if err := dao.createTaskStudentTable(); err != nil {
		return err
	}

	// 4. 创建提交表和文件表
	if err := dao.db.AutoMigrate(&model.Submission{}, &model.File{}); err != nil {
		return err
	}

	// 5. 创建初始用户数据
	if err := dao.createInitialUsers(); err != nil {
		return err
	}

	// 6. 创建初始任务数据
	if err := dao.createInitialTasks(); err != nil {
		return err
	}

	// 7. 创建初始提交数据
	if err := dao.createInitialSubmissions(); err != nil {
		return err
	}

	return nil
}

// createTaskStudentTable 手动创建任务学生关联表
func (dao *Dao) createTaskStudentTable() error {
	// 先删除可能存在的表
	dao.db.Exec("DROP TABLE IF EXISTS task_students")

	// 手动创建表
	sql := `
	CREATE TABLE task_students (
		task_id bigint unsigned NOT NULL,
		student_id bigint unsigned NOT NULL,
		created_at datetime(3) NULL,
		PRIMARY KEY (task_id, student_id)
	)`

	return dao.db.Exec(sql).Error
}

// cleanDatabase 清理数据库表
func (dao *Dao) cleanDatabase() error {
	// 按依赖关系倒序删除表
	tables := []string{"files", "submissions", "task_students", "tasks", "users"}

	for _, table := range tables {
		// 检查表是否存在
		var count int64
		dao.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", table).Count(&count)

		if count > 0 {
			// 删除表
			if err := dao.db.Exec("DROP TABLE IF EXISTS " + table).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// createInitialUsers 创建初始用户数据
func (dao *Dao) createInitialUsers() error {
	// 检查是否已有用户数据
	var count int64
	dao.db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，跳过初始化
	}

	// 密码加密
	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	passwordStr := string(password)

	// 创建教师用户
	teachers := []model.User{
		{
			Username:   "13800138001",
			Password:   passwordStr,
			Name:       "张教授",
			Role:       model.RoleTeacher,
			TeacherID:  "T001",
			Phone:      "13800138001",
			Department: "计算机科学与技术学院",
			IsActive:   true,
			WxOpenID:   "wx_teacher_001",
		},
		{
			Username:   "13800138002",
			Password:   passwordStr,
			Name:       "李老师",
			Role:       model.RoleTeacher,
			TeacherID:  "T002",
			Phone:      "13800138002",
			Department: "软件工程学院",
			IsActive:   true,
			WxOpenID:   "wx_teacher_002",
		},
		{
			Username:   "13800138003",
			Password:   passwordStr,
			Name:       "王老师",
			Role:       model.RoleTeacher,
			TeacherID:  "T003",
			Phone:      "13800138003",
			Department: "计算机科学与技术学院",
			IsActive:   true,
			WxOpenID:   "wx_teacher_003",
		},
	}

	for _, teacher := range teachers {
		if err := dao.db.Create(&teacher).Error; err != nil {
			return err
		}
	}

	// 创建学生用户
	students := []model.User{
		{
			Username:  "20210001",
			Password:  passwordStr,
			Name:      "张三",
			Role:      model.RoleStudent,
			StudentID: "20210001",
			Major:     "计算机科学与技术",
			Grade:     "2021级",
			Class:     "计科2101班",
			IsActive:  true,
			WxOpenID:  "wx_student_001",
		},
		{
			Username:  "20210002",
			Password:  passwordStr,
			Name:      "李四",
			Role:      model.RoleStudent,
			StudentID: "20210002",
			Major:     "计算机科学与技术",
			Grade:     "2021级",
			Class:     "计科2101班",
			IsActive:  true,
			WxOpenID:  "wx_student_002",
		},
		{
			Username:  "20210003",
			Password:  passwordStr,
			Name:      "王五",
			Role:      model.RoleStudent,
			StudentID: "20210003",
			Major:     "计算机科学与技术",
			Grade:     "2021级",
			Class:     "计科2101班",
			IsActive:  true,
			WxOpenID:  "wx_student_003",
		},
		{
			Username:  "20210004",
			Password:  passwordStr,
			Name:      "赵六",
			Role:      model.RoleStudent,
			StudentID: "20210004",
			Major:     "软件工程",
			Grade:     "2021级",
			Class:     "软工2101班",
			IsActive:  true,
			WxOpenID:  "wx_student_004",
		},
		{
			Username:  "20210005",
			Password:  passwordStr,
			Name:      "钱七",
			Role:      model.RoleStudent,
			StudentID: "20210005",
			Major:     "软件工程",
			Grade:     "2021级",
			Class:     "软工2101班",
			IsActive:  true,
			WxOpenID:  "wx_student_005",
		},
		{
			Username:  "20210006",
			Password:  passwordStr,
			Name:      "孙八",
			Role:      model.RoleStudent,
			StudentID: "20210006",
			Major:     "计算机科学与技术",
			Grade:     "2021级",
			Class:     "计科2102班",
			IsActive:  true,
			WxOpenID:  "wx_student_006",
		},
		{
			Username:  "20210007",
			Password:  passwordStr,
			Name:      "周九",
			Role:      model.RoleStudent,
			StudentID: "20210007",
			Major:     "计算机科学与技术",
			Grade:     "2021级",
			Class:     "计科2102班",
			IsActive:  true,
			WxOpenID:  "wx_student_007",
		},
		{
			Username:  "20210008",
			Password:  passwordStr,
			Name:      "吴十",
			Role:      model.RoleStudent,
			StudentID: "20210008",
			Major:     "软件工程",
			Grade:     "2021级",
			Class:     "软工2102班",
			IsActive:  true,
			WxOpenID:  "wx_student_008",
		},
	}

	for _, student := range students {
		if err := dao.db.Create(&student).Error; err != nil {
			return err
		}
	}

	return nil
}

// createInitialTasks 创建初始任务数据
func (dao *Dao) createInitialTasks() error {
	// 检查是否已有任务数据
	var count int64
	dao.db.Model(&model.Task{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，跳过初始化
	}

	// 获取教师ID
	var teacher1, teacher2, teacher3 model.User
	dao.db.Where("teacher_id = ?", "T001").First(&teacher1)
	dao.db.Where("teacher_id = ?", "T002").First(&teacher2)
	dao.db.Where("teacher_id = ?", "T003").First(&teacher3)

	now := time.Now()

	// 创建测试任务
	tasks := []model.Task{
		{
			Title:            "期末论文提交",
			Description:      "请提交期末课程设计论文，要求原创，字数不少于5000字。论文格式按照学校统一要求，包含摘要、关键词、正文、参考文献等部分。",
			Status:           model.TaskStatusActive,
			StartTime:        now.AddDate(0, 0, -7), // 7天前开始
			EndTime:          now.AddDate(0, 0, 5),  // 5天后截止
			AllowedFormats:   []string{".pdf", ".doc", ".docx"},
			FilenameTemplate: "学号_姓名_期末论文.pdf",
			MaxFileSize:      10485760, // 10MB
			TeacherID:        teacher1.ID,
			TotalStudents:    4,
		},
		{
			Title:            "数据分析报告",
			Description:      "完成第三章数据分析部分，包含数据预处理、统计分析、可视化图表等内容。",
			Status:           model.TaskStatusActive,
			StartTime:        now.AddDate(0, 0, -3), // 3天前开始
			EndTime:          now.AddDate(0, 0, 10), // 10天后截止
			AllowedFormats:   []string{".pdf", ".doc", ".docx", ".xlsx"},
			FilenameTemplate: "学号_姓名_数据分析报告",
			MaxFileSize:      20971520, // 20MB
			TeacherID:        teacher2.ID,
			TotalStudents:    4,
		},
		{
			Title:            "实验照片提交",
			Description:      "提交实验室操作照片，要求清晰展示实验过程和结果。每个实验至少3张照片。",
			Status:           model.TaskStatusActive,
			StartTime:        now.AddDate(0, 0, -1), // 1天前开始
			EndTime:          now.AddDate(0, 0, 15), // 15天后截止
			AllowedFormats:   []string{".jpg", ".jpeg", ".png"},
			FilenameTemplate: "学号_姓名_实验照片",
			MaxFileSize:      52428800, // 50MB
			TeacherID:        teacher3.ID,
			TotalStudents:    4,
		},
		{
			Title:            "程序设计作业",
			Description:      "完成课程设计程序，包含源代码、可执行文件和说明文档。",
			Status:           model.TaskStatusActive,
			StartTime:        now.AddDate(0, 0, -14), // 14天前开始
			EndTime:          now.AddDate(0, 0, -2),  // 2天前截止（已过期）
			AllowedFormats:   []string{".zip", ".rar", ".7z"},
			FilenameTemplate: "学号_姓名_程序设计作业",
			MaxFileSize:      104857600, // 100MB
			TeacherID:        teacher1.ID,
			TotalStudents:    4,
		},
	}

	for i, task := range tasks {
		if err := dao.db.Create(&task).Error; err != nil {
			return err
		}

		// 为每个任务分配学生
		var students []model.User
		dao.db.Where("role = ?", model.RoleStudent).Limit(4).Offset(i).Find(&students)

		// 创建任务-学生关联
		for _, student := range students {
			taskStudent := model.TaskStudent{
				TaskID:    task.ID,
				StudentID: student.ID,
				CreatedAt: time.Now(),
			}
			dao.db.Create(&taskStudent)
		}
	}

	return nil
}

// createInitialSubmissions 创建初始提交数据
func (dao *Dao) createInitialSubmissions() error {
	// 检查是否已有提交数据
	var count int64
	dao.db.Model(&model.Submission{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，跳过初始化
	}

	// 获取所有任务和学生
	var tasks []model.Task
	var students []model.User
	dao.db.Find(&tasks)
	dao.db.Where("role = ?", model.RoleStudent).Find(&students)

	now := time.Now()

	// 为每个任务创建一些提交记录
	for _, task := range tasks {
		// 获取该任务的学生
		var taskStudents []model.User
		dao.db.Joins("JOIN task_students ON users.id = task_students.student_id").
			Where("task_students.task_id = ?", task.ID).Find(&taskStudents)

		for i, student := range taskStudents {
			submission := model.Submission{
				TaskID:    task.ID,
				StudentID: student.ID,
			}

			// 模拟不同的提交状态
			if i%3 == 0 {
				// 已提交
				submitTime := now.AddDate(0, 0, -1)
				submission.Status = model.SubmissionStatusSubmitted
				submission.SubmittedAt = &submitTime
				submission.IsOnTime = submitTime.Before(task.EndTime)
				if !submission.IsOnTime {
					submission.Status = model.SubmissionStatusLate
				}
			} else if i%3 == 1 {
				// 已批阅
				submitTime := now.AddDate(0, 0, -2)
				reviewTime := now.AddDate(0, 0, -1)
				score := float64(85 + i*3) // 模拟分数
				submission.Status = model.SubmissionStatusReviewed
				submission.SubmittedAt = &submitTime
				submission.IsOnTime = submitTime.Before(task.EndTime)
				submission.Score = &score
				submission.Comment = "作业完成质量良好，格式规范，内容充实。"
				submission.ReviewedAt = &reviewTime
				submission.ReviewedBy = &task.TeacherID
			} else {
				// 未提交
				submission.Status = model.SubmissionStatusPending
			}

			if err := dao.db.Create(&submission).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
