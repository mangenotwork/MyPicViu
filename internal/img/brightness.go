package img

import (
	"MyPicViu/common/logger"
	"MyPicViu/common/utils"
	"fmt"
	"image"
	"image/color"
)

// 图片亮度

// ITU-R BT.709 标准的亮度计算公式：
// Y = 0.2126*R + 0.7152*G + 0.0722*B
// 这个公式考虑了人眼对不同颜色的敏感度差异：
// 绿色对亮度感知贡献最大（71.52%）
// 红色次之（21.26%）
// 蓝色贡献最小（7.22%）

// 计算单个RGB像素的亮度
// 返回值范围: 0.0 (最暗) 到 1.0 (最亮)
func calculatePixelBrightness(r, g, b uint8) float64 {
	// 使用ITU-R BT.709标准的亮度计算公式
	// 人眼对绿色最敏感，红色次之，蓝色最不敏感
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0

	// 亮度公式：Y = 0.2126*R + 0.7152*G + 0.0722*B
	brightness := 0.2126*rNorm + 0.7152*gNorm + 0.0722*bNorm
	return brightness
}

// 计算图片的平均亮度
func calculateImageBrightness(imgData image.Image) (float64, error) {
	// 获取图片边界
	bounds := imgData.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, fmt.Errorf("图片尺寸为零")
	}

	// 计算所有像素的亮度并取平均值
	var totalBrightness float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := imgData.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 计算当前像素的亮度并累加
			totalBrightness += calculatePixelBrightness(r8, g8, b8)
		}
	}

	// 返回平均亮度
	averageBrightness := totalBrightness / float64(totalPixels)
	return averageBrightness, nil
}

func SetImageBrightness(imgData image.Image, value float64) image.Image {
	logger.Debug("调整亮度 value = ", value)
	bounds := imgData.Bounds()
	result := image.NewRGBA(bounds)

	// 转换百分比为绝对亮度增量（-255到255）
	delta := int(value * 255)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := imgData.At(x, y).RGBA()

			// 转换为0-255范围
			r8 := int(r >> 8)
			g8 := int(g >> 8)
			b8 := int(b >> 8)

			// 加法调整（保持色彩比例）
			newR := r8 + delta
			newG := g8 + delta
			newB := b8 + delta

			// 钳位到有效范围
			newR = utils.Clamp(newR, 0, 255)
			newG = utils.Clamp(newG, 0, 255)
			newB = utils.Clamp(newB, 0, 255)

			result.Set(x, y, color.RGBA{
				uint8(newR),
				uint8(newG),
				uint8(newB),
				uint8(a >> 8),
			})
		}
	}
	return result
}
