package img

import (
	"golang.org/x/image/draw"
	"image"
	"math"
)

// 指纹字符串

/*
感知哈希算法是一类用于生成图像 “指纹”（即哈希值）的算法，这些 “指纹” 字符串能够以一种紧凑的形式概括图像的感知特征。通过比较这些 “指纹”，
可以快速判断两幅图像在视觉上是否相似，而无需像传统方法那样对整个图像像素进行逐点比较
这些算法常用于图像去重、图像检索等领域

均值哈希算法（Average Hash，AHash）
原理
图像缩放：将图像缩放到一个固定的尺寸（如 8x8 像素），忽略图像的细节和纵横比，这样可以确保所有图像在相同的尺度下进行比较。
灰度化：把缩放后的彩色图像转换为灰度图像，简化处理过程。
计算均值：计算灰度图像中所有像素的平均值。
生成哈希值：将每个像素的灰度值与平均值进行比较，如果像素的灰度值大于平均值，则该位置的哈希位设为 1，否则设为 0。最终得到一个由 0 和
1 组成的二进制字符串，这就是图像的 “指纹”。

感知哈希算法（Perceptual Hash，PHash）
原理
图像缩放：将图像缩放到一个固定的尺寸（如 32x32 像素），以减少计算量。
灰度化：将彩色图像转换为灰度图像。
离散余弦变换（DCT）：对灰度图像进行 DCT 变换，将图像从空间域转换到频率域。DCT 变换可以将图像的能量集中在低频部分，而高频部分则包含图像的细节信息。
取低频系数：选取 DCT 变换结果的左上角 8x8 区域，这些系数代表了图像的低频特征。
计算均值：计算选取的 8x8 区域的系数平均值。
生成哈希值：将每个系数与平均值进行比较，如果系数大于平均值，则该位置的哈希位设为 1，否则设为 0。最终得到一个 64 位的二进制字符串作为图像的 “指纹”。


差异哈希算法（Difference Hash，DHash）
原理
图像缩放：将图像缩放到一个固定的尺寸（如 9x8 像素）。
灰度化：将彩色图像转换为灰度图像。
计算差异：比较相邻像素的灰度值，如果右边像素的灰度值大于左边像素的灰度值，则该位置的哈希位设为 1，否则设为 0。这样可以得到一个 64 位的
二进制字符串作为图像的 “指纹”。

*/

// 计算图像的均值哈希值
func averageHash(imgData image.Image) string {
	// 缩放图像到 8x8
	resized := image.NewRGBA(image.Rect(0, 0, 8, 8))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), imgData, imgData.Bounds(), draw.Src, nil)

	// 灰度化并计算像素总和
	var total int
	grayPixels := make([]uint8, 64)
	index := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			// 将颜色值转换为 0 - 255 范围
			r = r >> 8
			g = g >> 8
			b = b >> 8
			gray := uint8((r + g + b) / 3)
			grayPixels[index] = gray
			total += int(gray)
			index++
		}
	}

	// 计算平均灰度值
	average := total / 64

	// 生成哈希值
	var hashStr string
	for _, pixel := range grayPixels {
		if int(pixel) > average {
			hashStr += "1"
		} else {
			hashStr += "0"
		}
	}

	return hashStr
}

// 二维离散余弦变换
func dct2d(data [][]float64) [][]float64 {
	N := len(data)
	result := make([][]float64, N)
	for i := range result {
		result[i] = make([]float64, N)
	}

	for u := 0; u < N; u++ {
		for v := 0; v < N; v++ {
			var sum float64
			Cu := 1.0
			Cv := 1.0
			for x := 0; x < N; x++ {
				for y := 0; y < N; y++ {
					if u == 0 {
						Cu = 1.0 / math.Sqrt(2)
					}
					if v == 0 {
						Cv = 1.0 / math.Sqrt(2)
					}
					sum += data[x][y] * math.Cos((2*float64(x)+1)*float64(u)*math.Pi/(2*float64(N))) *
						math.Cos((2*float64(y)+1)*float64(v)*math.Pi/(2*float64(N)))
				}
			}
			result[u][v] = 2.0 / float64(N) * Cu * Cv * sum
		}
	}
	return result
}

// 感知哈希算法
func pHash(imgData image.Image) string {
	// 缩放图像到 32x32
	resized := image.NewRGBA(image.Rect(0, 0, 32, 32))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), imgData, imgData.Bounds(), draw.Src, nil)

	// 灰度化
	grayPixels := make([][]float64, 32)
	for i := range grayPixels {
		grayPixels[i] = make([]float64, 32)
	}
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			gray := float64(r+g+b) / 3
			grayPixels[y][x] = gray
		}
	}

	// 二维离散余弦变换
	dctResult := dct2d(grayPixels)

	// 取左上角 8x8 低频分量
	lowFreq := make([][]float64, 8)
	for i := range lowFreq {
		lowFreq[i] = make([]float64, 8)
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			lowFreq[y][x] = dctResult[y][x]
		}
	}

	// 计算低频分量的平均值
	var sum float64
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			sum += lowFreq[y][x]
		}
	}
	average := sum / 64

	// 生成哈希值
	var hashStr string
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if lowFreq[y][x] > average {
				hashStr += "1"
			} else {
				hashStr += "0"
			}
		}
	}

	return hashStr
}

// 计算图像的差异哈希值
func differenceHash(imgData image.Image) string {
	// 缩放图像到 9x8
	resized := image.NewRGBA(image.Rect(0, 0, 9, 8))
	draw.NearestNeighbor.Scale(resized, resized.Bounds(), imgData, imgData.Bounds(), draw.Src, nil)

	// 灰度化
	grayPixels := make([]uint8, 9*8)
	index := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 9; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			gray := uint8((r + g + b) / 3)
			grayPixels[index] = gray
			index++
		}
	}

	// 生成哈希值
	var hashStr string
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			currentIndex := y*9 + x
			nextIndex := currentIndex + 1
			if grayPixels[nextIndex] > grayPixels[currentIndex] {
				hashStr += "1"
			} else {
				hashStr += "0"
			}
		}
	}

	return hashStr
}
