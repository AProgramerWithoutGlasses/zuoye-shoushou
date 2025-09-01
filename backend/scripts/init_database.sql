-- 数据库初始化SQL脚本
-- 使用前请确保数据库已创建

-- 先删除可能存在的表（按依赖关系倒序）
DROP TABLE IF EXISTS `files`;
DROP TABLE IF EXISTS `submissions`;
DROP TABLE IF EXISTS `task_students`;
DROP TABLE IF EXISTS `tasks`;
DROP TABLE IF EXISTS `users`;

-- 1. 创建用户表
CREATE TABLE `users` (
  `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `username` varchar(50) NOT NULL UNIQUE,
  `password` varchar(255) NOT NULL,
  `name` varchar(50) NOT NULL,
  `role` enum('student','teacher') NOT NULL,
  `wx_open_id` varchar(100) UNIQUE,
  `student_id` varchar(20),
  `major` varchar(100),
  `grade` varchar(20),
  `class` varchar(50),
  `teacher_id` varchar(20),
  `phone` varchar(20),
  `department` varchar(100),
  `is_active` boolean DEFAULT true,
  INDEX `idx_users_deleted_at` (`deleted_at`),
  INDEX `idx_users_student_id` (`student_id`),
  INDEX `idx_users_teacher_id` (`teacher_id`),
  INDEX `idx_users_phone` (`phone`)
);

-- 2. 创建任务表
CREATE TABLE `tasks` (
  `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `title` varchar(200) NOT NULL,
  `description` text,
  `status` enum('draft','active','expired','completed') DEFAULT 'draft',
  `start_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `allowed_formats` json,
  `filename_template` varchar(200),
  `max_file_size` bigint DEFAULT 10485760,
  `teacher_id` bigint unsigned NOT NULL,
  `total_students` int DEFAULT 0,
  `submitted_count` int DEFAULT 0,
  `on_time_count` int DEFAULT 0,
  INDEX `idx_tasks_deleted_at` (`deleted_at`),
  INDEX `idx_tasks_teacher_id` (`teacher_id`)
);

-- 3. 创建任务学生关联表
CREATE TABLE `task_students` (
  `task_id` bigint unsigned NOT NULL,
  `student_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) NULL,
  PRIMARY KEY (`task_id`, `student_id`)
);

-- 4. 创建提交表
CREATE TABLE `submissions` (
  `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `task_id` bigint unsigned NOT NULL,
  `student_id` bigint unsigned NOT NULL,
  `status` enum('pending','submitted','late','reviewed') DEFAULT 'pending',
  `submitted_at` datetime(3) NULL,
  `is_on_time` boolean DEFAULT false,
  `score` decimal(5,2) NULL,
  `comment` text,
  `reviewed_at` datetime(3) NULL,
  `reviewed_by` bigint unsigned NULL,
  INDEX `idx_submissions_deleted_at` (`deleted_at`),
  INDEX `idx_submissions_task_id` (`task_id`),
  INDEX `idx_submissions_student_id` (`student_id`),
  UNIQUE KEY `idx_task_student` (`task_id`, `student_id`)
);

-- 5. 创建文件表
CREATE TABLE `files` (
  `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `original_name` varchar(255) NOT NULL,
  `stored_name` varchar(255) NOT NULL,
  `file_path` varchar(500) NOT NULL,
  `file_size` bigint NOT NULL,
  `content_type` varchar(100),
  `file_hash` varchar(64),
  `submission_id` bigint unsigned NULL,
  `student_id` bigint unsigned NOT NULL,
  `task_id` bigint unsigned NOT NULL,
  INDEX `idx_files_deleted_at` (`deleted_at`),
  INDEX `idx_files_submission_id` (`submission_id`),
  INDEX `idx_files_student_id` (`student_id`),
  INDEX `idx_files_task_id` (`task_id`)
);

-- 6. 插入教师用户数据
INSERT INTO `users` (`username`, `password`, `name`, `role`, `teacher_id`, `phone`, `department`, `is_active`, `wx_open_id`, `created_at`, `updated_at`) VALUES
('13800138001', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '张教授', 'teacher', 'T001', '13800138001', '计算机科学与技术学院', true, 'wx_teacher_001', NOW(), NOW()),
('13800138002', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '李老师', 'teacher', 'T002', '13800138002', '软件工程学院', true, 'wx_teacher_002', NOW(), NOW()),
('13800138003', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '王老师', 'teacher', 'T003', '13800138003', '计算机科学与技术学院', true, 'wx_teacher_003', NOW(), NOW());

-- 7. 插入学生用户数据
INSERT INTO `users` (`username`, `password`, `name`, `role`, `student_id`, `major`, `grade`, `class`, `is_active`, `wx_open_id`, `created_at`, `updated_at`) VALUES
('20210001', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '张三', 'student', '20210001', '计算机科学与技术', '2021级', '计科2101班', true, 'wx_student_001', NOW(), NOW()),
('20210002', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '李四', 'student', '20210002', '计算机科学与技术', '2021级', '计科2101班', true, 'wx_student_002', NOW(), NOW()),
('20210003', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '王五', 'student', '20210003', '计算机科学与技术', '2021级', '计科2101班', true, 'wx_student_003', NOW(), NOW()),
('20210004', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '赵六', 'student', '20210004', '软件工程', '2021级', '软工2101班', true, 'wx_student_004', NOW(), NOW()),
('20210005', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '钱七', 'student', '20210005', '软件工程', '2021级', '软工2101班', true, 'wx_student_005', NOW(), NOW()),
('20210006', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '孙八', 'student', '20210006', '计算机科学与技术', '2021级', '计科2102班', true, 'wx_student_006', NOW(), NOW()),
('20210007', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '周九', 'student', '20210007', '计算机科学与技术', '2021级', '计科2102班', true, 'wx_student_007', NOW(), NOW()),
('20210008', '$2a$10$6pq1lLvUJE9BHVw0WGnmTegvBASOq6JJGWA3dfVP3p5dx/naabdO6', '吴十', 'student', '20210008', '软件工程', '2021级', '软工2102班', true, 'wx_student_008', NOW(), NOW());

-- 8. 插入任务数据
INSERT INTO `tasks` (`title`, `description`, `status`, `start_time`, `end_time`, `allowed_formats`, `filename_template`, `max_file_size`, `teacher_id`, `total_students`, `created_at`, `updated_at`) VALUES
('期末论文提交', '请提交期末课程设计论文，要求原创，字数不少于5000字。论文格式按照学校统一要求，包含摘要、关键词、正文、参考文献等部分。', 'active', DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_ADD(NOW(), INTERVAL 5 DAY), '["pdf", "doc", "docx"]', '学号_姓名_期末论文.pdf', 10485760, 1, 4, NOW(), NOW()),
('数据分析报告', '完成第三章数据分析部分，包含数据预处理、统计分析、可视化图表等内容。', 'active', DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_ADD(NOW(), INTERVAL 10 DAY), '["pdf", "doc", "docx", "xlsx"]', '学号_姓名_数据分析报告', 20971520, 2, 4, NOW(), NOW()),
('实验照片提交', '提交实验室操作照片，要求清晰展示实验过程和结果。每个实验至少3张照片。', 'active', DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 15 DAY), '["jpg", "jpeg", "png"]', '学号_姓名_实验照片', 52428800, 3, 4, NOW(), NOW()),
('程序设计作业', '完成课程设计程序，包含源代码、可执行文件和说明文档。', 'active', DATE_SUB(NOW(), INTERVAL 14 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), '["zip", "rar", "7z"]', '学号_姓名_程序设计作业', 104857600, 1, 4, NOW(), NOW());

-- 9. 插入任务学生关联数据
INSERT INTO `task_students` (`task_id`, `student_id`, `created_at`) VALUES
(1, 4, NOW()), (1, 5, NOW()), (1, 6, NOW()), (1, 7, NOW()),
(2, 8, NOW()), (2, 9, NOW()), (2, 10, NOW()), (2, 11, NOW()),
(3, 4, NOW()), (3, 5, NOW()), (3, 6, NOW()), (3, 7, NOW()),
(4, 8, NOW()), (4, 9, NOW()), (4, 10, NOW()), (4, 11, NOW());

-- 10. 插入提交数据
INSERT INTO `submissions` (`task_id`, `student_id`, `status`, `submitted_at`, `is_on_time`, `score`, `comment`, `reviewed_at`, `reviewed_by`, `created_at`, `updated_at`) VALUES
(1, 4, 'submitted', DATE_SUB(NOW(), INTERVAL 1 DAY), true, NULL, NULL, NULL, NULL, NOW(), NOW()),
(1, 5, 'reviewed', DATE_SUB(NOW(), INTERVAL 2 DAY), true, 88.0, '作业完成质量良好，格式规范，内容充实。', DATE_SUB(NOW(), INTERVAL 1 DAY), 1, NOW(), NOW()),
(1, 6, 'pending', NULL, false, NULL, NULL, NULL, NULL, NOW(), NOW()),
(1, 7, 'submitted', DATE_SUB(NOW(), INTERVAL 1 DAY), true, NULL, NULL, NULL, NULL, NOW(), NOW()),
(2, 8, 'submitted', DATE_SUB(NOW(), INTERVAL 1 DAY), true, NULL, NULL, NULL, NULL, NOW(), NOW()),
(2, 9, 'reviewed', DATE_SUB(NOW(), INTERVAL 2 DAY), true, 92.0, '数据分析深入，图表清晰，结论合理。', DATE_SUB(NOW(), INTERVAL 1 DAY), 2, NOW(), NOW()),
(2, 10, 'pending', NULL, false, NULL, NULL, NULL, NULL, NOW(), NOW()),
(2, 11, 'submitted', DATE_SUB(NOW(), INTERVAL 1 DAY), true, NULL, NULL, NULL, NULL, NOW(), NOW());

-- 11. 更新任务统计
UPDATE `tasks` SET 
  `submitted_count` = (SELECT COUNT(*) FROM `submissions` WHERE `task_id` = `tasks`.`id` AND `status` IN ('submitted', 'reviewed')),
  `on_time_count` = (SELECT COUNT(*) FROM `submissions` WHERE `task_id` = `tasks`.`id` AND `is_on_time` = true)
WHERE `id` IN (1, 2, 3, 4);

-- 完成初始化
SELECT '数据库初始化完成！' as message;
