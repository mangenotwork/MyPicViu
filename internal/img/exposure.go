package img

import (
	"fmt"
	"image"
)

// 图片曝光度

// 曝光度分析结果
type ExposureResult struct {
	OverexposedRatio  float64 // 过曝像素比例 (0.0-1.0)
	UnderexposedRatio float64 // 欠曝像素比例 (0.0-1.0)
	AverageLuminance  float64 // 平均亮度 (0.0-1.0)
	ExposureRating    string  // 曝光评级: "过曝", "欠曝", "正常"
}

// 计算图片的曝光度
func calculateImageExposure(imgData image.Image) (ExposureResult, error) {
	// 获取图片边界
	bounds := imgData.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return ExposureResult{}, fmt.Errorf("图片尺寸为零")
	}

	// 阈值定义
	overexposureThreshold := 0.9  // 亮度超过此值视为过曝
	underexposureThreshold := 0.1 // 亮度低于此值视为欠曝

	// 统计变量
	var overexposedCount, underexposedCount int
	var totalLuminance float64

	// 遍历所有像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := imgData.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算亮度
			luminance := getLuminance(r8, g8, b8)
			totalLuminance += luminance

			// 判断是否过曝或欠曝
			if luminance >= overexposureThreshold {
				overexposedCount++
			} else if luminance <= underexposureThreshold {
				underexposedCount++
			}
		}
	}

	// 计算比例
	overexposedRatio := float64(overexposedCount) / float64(totalPixels)
	underexposedRatio := float64(underexposedCount) / float64(totalPixels)
	averageLuminance := totalLuminance / float64(totalPixels)

	// 确定曝光评级
	var exposureRating string
	if overexposedRatio > 0.2 { // 超过20%像素过曝
		exposureRating = "过曝"
	} else if underexposedRatio > 0.4 { // 超过40%像素欠曝
		exposureRating = "欠曝"
	} else {
		exposureRating = "正常"
	}

	return ExposureResult{
		OverexposedRatio:  overexposedRatio,
		UnderexposedRatio: underexposedRatio,
		AverageLuminance:  averageLuminance,
		ExposureRating:    exposureRating,
	}, nil
}
