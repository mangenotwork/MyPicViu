package img

import (
	"golang.org/x/image/draw"
	"image"
)

// 缩放图像

// DrawCatmullRom 缩放图像 三次卷积插值
func DrawCatmullRom(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// DrawNearestNeighbor 缩放图像 最近邻插值
func DrawNearestNeighbor(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// DrawApproxBiLinear 缩放图像 近似双线性插值
func DrawApproxBiLinear(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// DrawBiLinear 缩放图像 双线性插值
func DrawBiLinear(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}
