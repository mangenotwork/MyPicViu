package img

import (
	"MyPicViu/common/logger"
	"MyPicViu/common/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

// ImgInfo 获取图片信息
type ImgInfo struct {
	FileName       string         // 图片文件名
	FileSize       int64          // 图片文件大小
	FileMd5        string         // 图片文件md5
	FilePath       string         // 图片文件路径
	FileMode       string         // 权限
	FileModTime    string         // 最后修改时间
	FileSys        string         // 系统原生信息
	Width          int            // 图片宽
	Height         int            // 图片高
	Format         string         // 图片格式
	DPI            int            // 图片dpi
	Bit            int            // 图片是多少位的
	MimeType       string         // 图片MIME类型
	Orientation    int            // 图片方向
	CameraMake     string         // 图片相机制造商
	CameraModel    string         // 相机型号
	DateTime       string         // 拍摄时间
	FocalLength    string         // 焦距
	ISO            string         // iso
	Aperture       string         // 光圈
	ShutterSpeed   string         // 快门
	Saturation     float64        // 图片饱和度值
	Brightness     float64        // 图片亮度值
	Contrast       float64        // 图片对比度值
	Sharpness      float64        // 图片锐度值
	Exposure       ExposureResult // 图片曝光度值
	Temperature    float64        // 图片色温值
	Hue            float64        // 图片色调值
	Noise          float64        // 图片噪点值
	DifferenceHash string         // 图片的差异哈希值获取图像的指纹字符串
	PHash          string         // 感知哈希算法获取图像的指纹字符串
	AverageHash    string         // 均值哈希值获取图像的指纹字符串
}

func (info *ImgInfo) GetFileInfo() {
	fileInfo, err := os.Stat(info.FilePath)
	if err != nil {
		// 处理错误（如文件不存在）
		if os.IsNotExist(err) {
			logger.ErrorF("文件不存在: %s\n", info.FilePath)
		} else if os.IsPermission(err) {
			logger.ErrorF("没有权限访问文件: %s\n", info.FilePath)
		} else {
			logger.ErrorF("获取文件信息失败: %v\n", err)
		}
		return
	}

	info.FileName = fileInfo.Name()
	info.FileSize = fileInfo.Size()
	info.FileMode = fileInfo.Mode().String()
	info.FileModTime = fileInfo.ModTime().Format(utils.TimeTemplateChinese)
	info.FileSys = fmt.Sprintf("%+v", fileInfo.Sys())

}

func (info *ImgInfo) GetImgInfo() {
	// 打开文件
	file, err := os.Open(info.FilePath)
	if err != nil {
		logger.Error("无法打开文件:", err)
		return
	}
	defer file.Close()

	imgConfig, format, err := image.DecodeConfig(file)
	if err != nil {
		logger.ErrorF("无法解析图片: %w", err)
		return
	}
	info.Width = imgConfig.Width
	info.Height = imgConfig.Height
	info.Format = format

	// 创建MD5哈希计算器
	hash := md5.New()

	// 以流的方式读取文件并更新哈希（适合大文件）
	hashFile, err := utils.CopyFileHandle(file)
	defer func() {
		_ = hashFile.Close()
	}()
	if err != nil {
		logger.Error("copy文件失败")
		return
	}
	if _, err := io.Copy(hash, hashFile); err != nil {
		logger.Error("读取文件失败:", err)
	}

	// 计算最终哈希值并转换为十六进制字符串
	md5Sum := hex.EncodeToString(hash.Sum(nil))
	info.FileMd5 = md5Sum

	logger.Debug("图片宽高格式 : ", info.Width, "|", info.Height, "|", info.Format)

	exifFile, err := utils.CopyFileHandle(file)
	defer func() {
		_ = exifFile.Close()
	}()
	if err != nil {
		logger.Error("copy文件失败")
		return
	}
	info.GetExif(exifFile)

	imgFile, err := utils.CopyFileHandle(file)
	if err != nil {
		logger.Error("copy文件失败")
		return
	}
	defer func() {
		_ = imgFile.Close()
	}()
	imgData, _, err := image.Decode(imgFile)
	if err != nil {
		logger.Error("读取图片文件失败", err)
		return
	}

	info.Saturation, err = calculateImageSaturation(imgData)
	if err != nil {
		logger.Error("获取饱和度失败：", err)
	}
	logger.Debug("饱和度 ： ", info.Saturation)

	info.Brightness, err = calculateImageBrightness(imgData)
	if err != nil {
		logger.Error("获取亮度失败：", err)
	}
	logger.Debug("亮度 ： ", info.Saturation)

	info.Contrast, err = calculateImageContrast(imgData)
	if err != nil {
		logger.Error("获取对比度失败：", err)
	}
	logger.Debug("对比度 ： ", info.Saturation)

	info.Sharpness, err = calculateImageSharpness(imgData)
	if err != nil {
		logger.Error("获取锐度失败：", err)
	}
	logger.Debug("锐度 ： ", info.Saturation)

	info.Exposure, err = calculateImageExposure(imgData)
	if err != nil {
		logger.Error("获取曝光度失败：", err)
	}
	logger.Debug("曝光度 ： ", info.Exposure.ExposureRating)

	info.Temperature, err = calculateImageColorTemperature(imgData)
	if err != nil {
		logger.Error("获取色温失败：", err)
	}
	logger.Debug("色温 ： ", info.Temperature)

	info.Hue, _, err = calculateImageHue(imgData)
	if err != nil {
		logger.Error("获取色调失败：", err)
	}
	logger.Debug("色调 ： ", info.Temperature)

	info.Noise, err = calculateImageNoise(imgData)
	if err != nil {
		logger.Error("获取噪点值失败：", err)
	}
	logger.Debug("噪点值 ： ", info.Temperature)

	info.DifferenceHash = differenceHash(imgData)
	info.PHash = pHash(imgData)
	info.AverageHash = averageHash(imgData)

	logger.Debug("指纹1 ", info.DifferenceHash)
	logger.Debug("指纹2 ", info.PHash)
	logger.Debug("指纹3 ", info.AverageHash)
}
