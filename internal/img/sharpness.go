package img

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// 图片锐度值

// 使用 Sobel 算子来计算图像的锐度，工作原理如下：
//	首先将图像转换为亮度矩阵，忽略色彩信息，只关注明暗变化
//	使用 Sobel 边缘检测算子，该算子包含两个卷积核：
//	水平方向卷积核（检测垂直边缘）
//	垂直方向卷积核（检测水平边缘）
//	对每个像素应用这两个卷积核，计算梯度幅度（边缘强度）
//	所有像素的平均边缘强度作为图像的锐度值

// 使用Sobel算子计算锐度
// 锐度值越高，表示图像越清晰
func calculateImageSharpness(imgData image.Image) (float64, error) {

	// 获取图片边界
	bounds := imgData.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	if width < 3 || height < 3 {
		return 0.0, fmt.Errorf("图片尺寸过小，无法计算锐度")
	}

	// 创建亮度矩阵
	luminance := make([][]float64, height)
	for y := 0; y < height; y++ {
		luminance[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := imgData.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			luminance[y][x] = getLuminance(r8, g8, b8)
		}
	}

	// Sobel算子 - 水平和垂直方向的卷积核
	sobelX := [3][3]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	sobelY := [3][3]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	// 计算边缘强度
	var totalEdgeStrength float64
	edgeCount := 0

	// 遍历每个像素（跳过边界像素）
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// 应用Sobel算子
			var gx, gy float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					gx += luminance[y+ky][x+kx] * sobelX[ky+1][kx+1]
					gy += luminance[y+ky][x+kx] * sobelY[ky+1][kx+1]
				}
			}

			// 计算梯度幅度（边缘强度）
			edgeStrength := math.Sqrt(gx*gx + gy*gy)
			totalEdgeStrength += edgeStrength
			edgeCount++
		}
	}

	// 计算平均边缘强度作为锐度度量
	averageSharpness := totalEdgeStrength / float64(edgeCount)
	return averageSharpness, nil
}

func SetImageSharpness(imgData image.Image, value float64) image.Image {
	// 拉普拉斯算子
	var laplacianKernel = [][]int{
		{0, -1, 0},
		{-1, 4, -1},
		{0, -1, 0},
	}

	bounds := imgData.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var rSum, gSum, bSum int

			// 应用拉普拉斯算子
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					r, g, b, _ := imgData.At(x+kx, y+ky).RGBA()
					r = r / 256
					g = g / 256
					b = b / 256

					kernelValue := laplacianKernel[ky+1][kx+1]
					rSum += int(r) * kernelValue
					gSum += int(g) * kernelValue
					bSum += int(b) * kernelValue
				}
			}

			// 获取当前像素的原始值
			r, g, b, a := imgData.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256

			// 调整锐度
			newR := int(r) + int(float64(rSum)*value)
			newG := int(g) + int(float64(gSum)*value)
			newB := int(b) + int(float64(bSum)*value)

			// 确保颜色值在 0 到 255 范围内
			if newR < 0 {
				newR = 0
			} else if newR > 255 {
				newR = 255
			}
			if newG < 0 {
				newG = 0
			} else if newG > 255 {
				newG = 255
			}
			if newB < 0 {
				newB = 0
			} else if newB > 255 {
				newB = 255
			}

			result.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), uint8(a)})
		}
	}

	// 处理边缘像素，简单复制原始像素值
	for y := 0; y < height; y++ {
		result.Set(0, y, imgData.At(0, y))
		result.Set(width-1, y, imgData.At(width-1, y))
	}
	for x := 0; x < width; x++ {
		result.Set(x, 0, imgData.At(x, 0))
		result.Set(x, height-1, imgData.At(x, height-1))
	}

	return result
}
