package img

import (
	"fmt"
	"image"
	"math"
)

// 图片噪点值

//	原理如下：
//
//	首先将图像转换为亮度矩阵，专注于亮度通道的变化
//	使用 3x3 高斯模糊核创建图像的平滑版本，高斯模糊能有效去除高频噪声
//	计算原始图像与平滑图像之间的差异，这种差异主要来自于噪点
//	对所有像素的差异取平均值，作为整体噪点水平的度量

// 计算图片的噪点值
// 返回值越高，表示噪点越多
func calculateImageNoise(imgData image.Image) (float64, error) {
	// 获取图片边界
	bounds := imgData.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	if width < 3 || height < 3 {
		return 0.0, fmt.Errorf("图片尺寸过小，无法计算噪点")
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

	// 使用3x3高斯模糊核创建平滑图像
	gaussianKernel := [3][3]float64{
		{1, 2, 1},
		{2, 4, 2},
		{1, 2, 1},
	}
	kernelSum := 16.0 // 高斯核的总和

	// 创建平滑后的亮度矩阵
	smoothed := make([][]float64, height)
	for y := 0; y < height; y++ {
		smoothed[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			// 边界像素直接使用原始值
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				smoothed[y][x] = luminance[y][x]
				continue
			}

			// 应用高斯模糊
			var sum float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					sum += luminance[y+ky][x+kx] * gaussianKernel[ky+1][kx+1]
				}
			}
			smoothed[y][x] = sum / kernelSum
		}
	}

	// 计算原始图像与平滑图像的差异（噪点估计）
	var totalNoise float64
	noiseCount := 0

	// 忽略边缘像素
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// 计算亮度差异的绝对值
			diff := math.Abs(luminance[y][x] - smoothed[y][x])
			totalNoise += diff
			noiseCount++
		}
	}

	// 计算平均噪点值
	averageNoise := totalNoise / float64(noiseCount)
	return averageNoise, nil
}
