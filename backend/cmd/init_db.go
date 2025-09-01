package main

import (
	"fmt"
	"goweb_staging/dao"
	"goweb_staging/pkg/settings"
	"log"
)

func main() {
	fmt.Println("=== 数据库初始化程序 ===")

	// 1. 加载配置
	fmt.Println("1. 加载配置文件...")
	app, err := settings.Init("local")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v\n", err)
	}
	fmt.Println("✅ 配置加载成功")

	// 2. 初始化数据访问层
	fmt.Println("2. 初始化数据库连接...")
	daoInstance := dao.Init(app)
	fmt.Println("✅ 数据库连接成功")

	// 3. 初始化数据库数据
	fmt.Println("3. 开始初始化数据库数据...")
	if err := daoInstance.InitData(); err != nil {
		log.Fatalf("❌ 初始化数据库数据失败: %v\n", err)
	}

	fmt.Println("✅ 数据库初始化完成！")
	fmt.Println("")
	fmt.Println("=== 测试账号信息 ===")
	fmt.Println("")
	fmt.Println("📚 教师账号：")
	fmt.Println("   用户名: 13800138001  密码: 123456  姓名: 张教授")
	fmt.Println("   用户名: 13800138002  密码: 123456  姓名: 李老师")
	fmt.Println("   用户名: 13800138003  密码: 123456  姓名: 王老师")
	fmt.Println("")
	fmt.Println("🎓 学生账号：")
	fmt.Println("   用户名: 20210001  密码: 123456  姓名: 张三")
	fmt.Println("   用户名: 20210002  密码: 123456  姓名: 李四")
	fmt.Println("   用户名: 20210003  密码: 123456  姓名: 王五")
	fmt.Println("   用户名: 20210004  密码: 123456  姓名: 赵六")
	fmt.Println("   用户名: 20210005  密码: 123456  姓名: 钱七")
	fmt.Println("   用户名: 20210006  密码: 123456  姓名: 孙八")
	fmt.Println("   用户名: 20210007  密码: 123456  姓名: 周九")
	fmt.Println("   用户名: 20210008  密码: 123456  姓名: 吴十")
	fmt.Println("")
	fmt.Println("📝 已创建测试任务和部分提交记录")
	fmt.Println("")
	fmt.Println("=== 初始化完成 ===")
}
