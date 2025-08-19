package utils

import (
	"MyPicViu/common/logger"
	"fmt"
	"net/http"
	"os"
	"sync"
)

// CopyFileHandle 创建一个新的*os.File句柄，指向与源文件相同的路径
func CopyFileHandle(src *os.File) (*os.File, error) {
	// 获取源文件的路径
	path := src.Name()

	// 获取源文件的打开模式（只读、读写等）
	// 注意：需要根据实际需求调整模式，这里假设源文件是可读的
	mode := os.O_RDONLY
	if src.Fd() >= 0 {
		// 简单判断是否可写（实际场景可能需要更复杂的权限检查）
		if _, err := src.WriteString(""); err == nil {
			mode = os.O_RDWR // 源文件可写，则新句柄也以读写模式打开
		}
	}

	// 再次打开文件，创建新的句柄
	newFile, err := os.OpenFile(path, mode, 0644)
	if err != nil {
		return nil, fmt.Errorf("无法创建新文件句柄: %w", err)
	}

	return newFile, nil
}

// 缓冲区对象池：缓存512字节的字节切片
var bufferPool = sync.Pool{
	// New函数：当池为空时，创建新的缓冲区
	New: func() interface{} {
		return make([]byte, 512) // 固定512字节，满足http.DetectContentType需求
	},
}

// DetectByStdLib 使用标准库（基础类型识别）
func DetectByStdLib(filePath string) (string, error) {
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	n, err := file.Read(buf)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buf[:n]), nil
}

func IsImgFile(filePath string) bool {
	ext, err := DetectByStdLib(filePath)
	if err != nil {
		logger.Error("获取文件类型失败: ", err)
		return false
	}
	for _, v := range []string{"image/png", "image/jpeg", "image/gif", "image/webp", "image/bmp", "image/tiff", "image/svg+xml"} {
		if ext == v {
			return true
		}
	}
	return false
}
