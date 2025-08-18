package utils

import (
	"fmt"
	"os"
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
