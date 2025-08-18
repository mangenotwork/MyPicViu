package img

import (
	"fmt"
	"image"
	"math"
)

// 图片色温值

// 计算色温的原理如下：
//	分析图像中所有非暗部像素的 RGB 值，忽略过暗像素（避免黑色影响判断）
//	计算红、绿、蓝三通道的平均值，并以绿色为基准进行归一化处理
//	通过红蓝光比例计算色温值，使用经验模型将比例转换为开尔文温度
//	根据色温值判断图片的色调偏向：
//	低于 5000K：暖色调（偏黄 / 红色）
//	5000K-7000K：中性色调
//	高于 7000K：冷色调（偏蓝色）

// 计算图片的色温（开尔文）
func calculateImageColorTemperature(imgData image.Image) (float64, error) {
	// 获取图片边界
	bounds := imgData.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, fmt.Errorf("图片尺寸为零")
	}

	// 统计所有像素的RGB值总和（忽略极暗像素，避免影响计算）
	var totalR, totalG, totalB float64
	validPixels := 0
	var darkThreshold uint8 = 30 // 忽略暗部像素（避免黑色影响色温判断）

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := imgData.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 忽略暗部像素
			if r8 < darkThreshold && g8 < darkThreshold && b8 < darkThreshold {
				continue
			}

			totalR += float64(r8)
			totalG += float64(g8)
			totalB += float64(b8)
			validPixels++
		}
	}

	if validPixels == 0 {
		return 0.0, fmt.Errorf("图片过暗，无法计算色温")
	}

	// 计算平均RGB值
	avgR := totalR / float64(validPixels)
	avgG := totalG / float64(validPixels)
	avgB := totalB / float64(validPixels)

	// 归一化处理（以绿色为基准）
	rNorm := avgR / avgG
	bNorm := avgB / avgG

	// 计算红蓝光比例
	rRatio := rNorm / (rNorm + bNorm)
	bRatio := bNorm / (rNorm + bNorm)

	// 色温转换公式（基于经验模型）
	// 色温范围大致在2000K（暖黄）到10000K（冷蓝）之间
	var temperature float64
	if rRatio > bRatio {
		// 暖色调
		ratio := rRatio / bRatio
		temperature = 6500 - 4500*math.Min(1.0, (ratio-1.0)/2.0)
	} else {
		// 冷色调
		ratio := bRatio / rRatio
		temperature = 6500 + 3500*math.Min(1.0, (ratio-1.0)/2.0)
	}

	// 确保色温在合理范围内
	if temperature < 2000 {
		temperature = 2000
	} else if temperature > 10000 {
		temperature = 10000
	}

	return temperature, nil
}
