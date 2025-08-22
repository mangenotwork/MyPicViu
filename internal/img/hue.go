package img

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// 图片色调

// 计算图片的色调（Hue）需要将 RGB 颜色空间转换到 HSV（色相 - 饱和度 - 明度）颜色空间，
//其中 Hue（色相）代表了颜色的基本属性，如红色、绿色、蓝色等。色调通常用角度（0°-360°）表示不同的颜色
//
// 原理如下：
//	将每个像素的 RGB 值转换为 HSV 颜色空间，提取色调（Hue）值
//	色调用 0°-360° 的角度表示，不同角度对应不同颜色：
//	0°/360°：红色
//	60°：黄色
//	120°：绿色
//	180°：青色
//	240°：蓝色
//	300°：品红色
//	忽略过暗和低饱和度的像素（接近灰色的像素），避免影响色调判断
//	统计所有有效像素的色调分布，计算平均色调角度和主要色调类别
//
//	程序将色调分为 13 种常见类别，包括红色、橙色、黄色、绿色等，最终输出：
//	平均色调角度（0°-360°）
//	主要色调类别（图片中最占优势的颜色）

// 色调范围定义
type HueRange struct {
	Start    float64
	End      float64
	Category string
}

// 常见色调范围分类
var hueRanges = []HueRange{
	{345, 360, "红色"},
	{0, 15, "红色"},
	{15, 45, "橙色"},
	{45, 75, "黄色"},
	{75, 105, "绿色"},
	{105, 135, "青色"},
	{135, 165, "蓝色"},
	{165, 195, "靛蓝色"},
	{195, 225, "紫色"},
	{225, 255, "粉红色"},
	{255, 285, "品红色"},
	{285, 315, "紫红色"},
	{315, 345, "深红色"},
}

// 将RGB转换为HSV颜色空间，返回色调值（0-360）
func rgbToHue(r, g, b uint8) float64 {
	// 归一化到0.0-1.0范围
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0

	maxVal := math.Max(rNorm, math.Max(gNorm, bNorm))
	minVal := math.Min(rNorm, math.Min(gNorm, bNorm))
	delta := maxVal - minVal

	var hue float64

	// 如果delta为0，说明是灰色，没有色调
	if delta == 0 {
		return -1 // 用-1表示无色调（灰色）
	}

	// 计算色调
	switch {
	case maxVal == rNorm:
		hue = math.Mod((gNorm-bNorm)/delta, 6)
	case maxVal == gNorm:
		hue = (bNorm-rNorm)/delta + 2
	case maxVal == bNorm:
		hue = (rNorm-gNorm)/delta + 4
	}

	// 转换为0-360度
	hue *= 60
	if hue < 0 {
		hue += 360
	}

	return hue
}

// 确定色调所属类别
func getHueCategory(hue float64) string {
	if hue < 0 {
		return "灰色/无色调"
	}

	for _, r := range hueRanges {
		if (hue >= r.Start && hue <= r.End) ||
			(r.Start > r.End && (hue >= r.Start || hue <= r.End)) {
			return r.Category
		}
	}
	return "未知"
}

// 计算图片的主色调
func calculateImageHue(imgData image.Image) (float64, string, error) {
	// 获取图片边界
	bounds := imgData.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	totalPixels := width * height

	if totalPixels == 0 {
		return 0.0, "", fmt.Errorf("图片尺寸为零")
	}

	// 统计变量
	var totalHue float64
	hueCount := 0
	hueDistribution := make(map[string]int)
	nonGrayPixels := 0

	// 忽略过暗像素的阈值
	var darkThreshold uint8 = 30
	// 忽略低饱和度像素的阈值
	minSaturation := 0.15

	// 遍历所有像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, _ := imgData.At(x, y).RGBA()

			// 将16位RGBA值转换为8位RGB值
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// 忽略过暗像素
			if r8 < darkThreshold && g8 < darkThreshold && b8 < darkThreshold {
				continue
			}

			// 计算饱和度（用于过滤低饱和度像素）
			rNorm := float64(r8) / 255.0
			gNorm := float64(g8) / 255.0
			bNorm := float64(b8) / 255.0
			maxVal := math.Max(rNorm, math.Max(gNorm, bNorm))
			minVal := math.Min(rNorm, math.Min(gNorm, bNorm))
			var saturation float64
			if maxVal > 0 {
				saturation = (maxVal - minVal) / maxVal
			}

			// 忽略低饱和度像素（接近灰色）
			if saturation < minSaturation {
				continue
			}

			// 计算色调
			hue := rgbToHue(r8, g8, b8)
			if hue >= 0 {
				totalHue += hue
				hueCount++
				nonGrayPixels++

				// 统计色调分布
				category := getHueCategory(hue)
				hueDistribution[category]++
			}
		}
	}

	if nonGrayPixels == 0 {
		return 0.0, "灰色/无明显色调", nil
	}

	// 计算平均色调
	averageHue := totalHue / float64(hueCount)

	// 找到最主要的色调类别
	mainCategory := "未知"
	maxCount := 0
	for cat, count := range hueDistribution {
		if count > maxCount {
			maxCount = count
			mainCategory = cat
		}
	}

	return averageHue, mainCategory, nil
}

func SetImageHue(imgData image.Image, value float64) image.Image {
	bounds := imgData.Bounds()
	result := image.NewRGBA(bounds)

	// 将 -1.0~1.0 映射为色相角度偏移（-180~180度）
	// value=1.0 表示顺时针旋转180度，value=-1.0表示逆时针旋转180度
	hueOffset := value * 180

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := imgData.At(x, y).RGBA()

			// 正确转换为 0-255 范围（位运算更高效）
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8) // alpha通道单独处理

			// 转换为HSV色彩空间
			h, s, v := RGBToHSV(r8, g8, b8)

			// 调整色调（叠加偏移后取模，确保在0-360度范围）
			h = math.Mod(h+hueOffset, 360)
			if h < 0 {
				h += 360 // 处理负数情况（如hueOffset=-200，h=100时，结果应为160）
			}

			// 转换回RGB色彩空间
			r1, g1, b1 := HSVToRGB(h, s, v)

			result.Set(x, y, color.RGBA{r1, g1, b1, a8})
		}
	}
	return result
}
