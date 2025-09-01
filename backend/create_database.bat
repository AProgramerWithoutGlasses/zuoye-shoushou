@echo off
echo ========================================
echo 创建数据库
echo ========================================
echo.

echo 正在创建数据库 zuoye-shoushou...
mysql -u root -p123456 -e "CREATE DATABASE IF NOT EXISTS zuoye-shoushou CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

if %errorlevel% equ 0 (
    echo ✅ 数据库创建成功！
    echo.
    echo 现在可以运行 run_sql_init.bat 来初始化数据
) else (
    echo ❌ 数据库创建失败！
    echo 请检查MySQL服务是否启动
)

echo.
echo 按任意键退出...
pause > nul
