package server

import (
	"crypto/md5"
	"fmt"
	"goweb_staging/pkg/response"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// uploadFile 文件上传
func uploadFile(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMsg(c, response.ParamErrCode, "获取文件失败")
		return
	}
	defer file.Close()

	// 验证文件大小 (最大10MB)
	maxSize := int64(10 << 20) // 10MB
	if header.Size > maxSize {
		response.FailWithMsg(c, response.ParamErrCode, "文件大小超过限制")
		return
	}

	// 验证文件格式
	allowedExt := []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".jpg", ".jpeg", ".png"}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	isValidExt := false
	for _, validExt := range allowedExt {
		if ext == validExt {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		response.FailWithMsg(c, response.ParamErrCode, "不支持的文件格式")
		return
	}

	// 生成文件哈希
	fileHash, err := calculateFileHash(file)
	if err != nil {
		response.FailWithMsg(c, response.ServerErrCode, "计算文件哈希失败")
		return
	}

	// 重置文件指针
	file.Seek(0, 0)

	// 生成存储文件名
	storedName := fmt.Sprintf("%d_%s%s", time.Now().Unix(), fileHash[:8], ext)

	// 创建上传目录
	uploadDir := "uploads/" + time.Now().Format("2006/01/02")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.FailWithMsg(c, response.ServerErrCode, "创建上传目录失败")
		return
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, storedName)
	dst, err := os.Create(filePath)
	if err != nil {
		response.FailWithMsg(c, response.ServerErrCode, "创建文件失败")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		response.FailWithMsg(c, response.ServerErrCode, "保存文件失败")
		return
	}

	// 返回文件信息
	fileInfo := map[string]interface{}{
		"original_name": header.Filename,
		"stored_name":   storedName,
		"file_path":     filePath,
		"file_size":     header.Size,
		"content_type":  header.Header.Get("Content-Type"),
		"file_hash":     fileHash,
	}

	response.Success(c, fileInfo)
}

// downloadFile 文件下载
func downloadFile(c *gin.Context) {
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		response.Fail(c, response.ParamErrCode)
		return
	}

	// 获取文件信息
	fileInfo, err := svc.GetFileByID(fileID)
	if err != nil {
		zap.L().Error("get file failed", zap.Error(err))
		response.Fail(c, response.ServerErrCode)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(fileInfo.FilePath); os.IsNotExist(err) {
		response.FailWithMsg(c, response.ServerErrCode, "文件不存在")
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.OriginalName)
	c.Header("Content-Type", "application/octet-stream")

	// 返回文件
	c.File(fileInfo.FilePath)
}

// calculateFileHash 计算文件MD5哈希
func calculateFileHash(file multipart.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
