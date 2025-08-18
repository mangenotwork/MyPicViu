package img

import (
	"fmt"
	"image"
)

// 饱和度

// 计算单个RGB像素的饱和度
// 返回值范围: 0.0 (灰度) 到 1.0 (最大饱和度)
func calculatePixelSaturation(r, g, b uint8) float64 {
	// 将uint8转换为0.0-1.0范围的浮点数
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0

	// 找到RGB中的最大值和最小值
	maxVal := max3(rNorm, gNorm, bNorm)
	minVal := min3(rNorm, gNorm, bNorm)

	// 如果是灰度（max == min），饱和度为0
	if maxVal == minVal {
		return 0.0
	}

	// 计算亮度
	luminance := (maxVal + minVal) / 2.0

	var saturation float64
	if luminance <= 0.5 {
		saturation = (maxVal - minVal) / (maxVal + minVal)
	} else {
		saturation = (maxVal - minVal) / (2.0 - maxVal - minVal)
	}

	return saturation
}

// 辅助函数：返回三个数中的最大值
func max3(a, b, c float64) float64 {
	max := a
	if b > max {
		max = b
	}
	if c > max {
		max = c
	}
	return max
}

// 辅助函数：返回三个数中的最小值
func min3(a, b, c float64) float64 {
	min := a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return min
}

// 计算图片的平均饱和度
func calculateImageSaturation(imgData image.Image) (float64, error) {

	// 获取图片边界
	bounds := imgData.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, fmt.Errorf("图片尺寸为零")
	}

	// 计算所有像素的饱和度并取平均值
	var totalSaturation float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := imgData.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算当前像素的饱和度并累加
			totalSaturation += calculatePixelSaturation(r8, g8, b8)
		}
	}

	// 返回平均饱和度
	averageSaturation := totalSaturation / float64(totalPixels)
	return averageSaturation, nil
}
