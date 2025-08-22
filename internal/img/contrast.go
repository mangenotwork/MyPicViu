package img

import (
	"MyPicViu/common/utils"
	"fmt"
	"image"
	"image/color"
	"math"
)

// 图片对比度

// 两种计算对比度的方法：
//
//1. Michelson 对比度（默认使用）：
//	公式：(最大亮度 - 最小亮度) / (最大亮度 + 最小亮度)
//	这种方法简单直观，通过计算图像中最亮和最暗像素的差异来衡量对比度
//	值范围：0.0（无对比度，完全灰阶）到 1.0（最大对比度）
//2. 标准差对比度（注释中提供）：
//	基于亮度值的统计分布，使用标准差除以平均值（变异系数）
//	这种方法能更好地反映整体图像的对比度分布情况
//	适用于需要更精确对比度评估的场景

// 计算单个RGB像素的亮度（使用ITU-R BT.709标准）
func getLuminance(r, g, b uint8) float64 {
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0
	return 0.2126*rNorm + 0.7152*gNorm + 0.0722*bNorm
}

// 计算图片的对比度
// 对比度公式: (最大亮度 - 最小亮度) / (最大亮度 + 最小亮度)
// 返回值范围: 0.0 (无对比度) 到 1.0 (最大对比度)
func calculateImageContrast(imgData image.Image) (float64, error) {

	// 获取图片边界
	bounds := imgData.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, fmt.Errorf("图片尺寸为零")
	}

	// 初始化最大和最小亮度
	maxLuminance := 0.0
	minLuminance := 1.0
	var totalLuminance float64
	var totalLuminanceSquared float64

	// 遍历所有像素计算亮度统计值
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := imgData.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算当前像素的亮度
			luminance := getLuminance(r8, g8, b8)

			// 更新最大和最小亮度
			if luminance > maxLuminance {
				maxLuminance = luminance
			}
			if luminance < minLuminance {
				minLuminance = luminance
			}

			// 累加亮度值用于计算标准差
			totalLuminance += luminance
			totalLuminanceSquared += luminance * luminance
		}
	}

	//// 方法1: 使用最大最小亮度计算对比度 (Michelson对比度)
	//// 这种方法适合简单场景
	//contrast := (maxLuminance - minLuminance) / (maxLuminance + minLuminance)

	// 方法2: 使用标准差计算对比度 (更精确但计算稍复杂)
	mean := totalLuminance / float64(totalPixels)
	variance := (totalLuminanceSquared / float64(totalPixels)) - (mean * mean)
	stdDev := math.Sqrt(variance)
	contrast := stdDev / mean // 变异系数作为对比度度量

	return contrast, nil
}

func SetImageContrast(imgData image.Image, value float64) image.Image {
	bounds := imgData.Bounds()
	result := image.NewRGBA(bounds)

	// 将 -1.0~1.0 转换为对比度因子（映射到 0~255 等效范围）
	// value=0 时保持原图；value=1 时最大对比度；value=-1 时最小对比度（灰度）
	adjustedValue := value * 255 // 转换为 -255~255 范围
	factor := (259 * (adjustedValue + 255)) / (255 * (259 - adjustedValue))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := imgData.At(x, y).RGBA()

			// 正确转换为 0-255 范围（位运算更高效）
			r8 := int(r >> 8)
			g8 := int(g >> 8)
			b8 := int(b >> 8)
			a8 := uint8(a >> 8) // alpha 通道单独处理

			// 对比度调整公式：factor*(色值-128) + 128
			newR := int(factor*(float64(r8)-128) + 128)
			newG := int(factor*(float64(g8)-128) + 128)
			newB := int(factor*(float64(b8)-128) + 128)

			// 钳位到 0-255 范围（使用工具函数更简洁）
			newR = utils.Clamp(newR, 0, 255)
			newG = utils.Clamp(newG, 0, 255)
			newB = utils.Clamp(newB, 0, 255)

			result.Set(x, y, color.RGBA{
				uint8(newR),
				uint8(newG),
				uint8(newB),
				a8,
			})
		}
	}
	return result
}
