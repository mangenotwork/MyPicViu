package img

import "image"

// Rotate90Clockwise 顺时针旋转 90 度
func Rotate90Clockwise(src image.Image) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, srcHeight, srcWidth))

	for y := 0; y < srcHeight; y++ {
		for x := 0; x < srcWidth; x++ {
			newX := y
			newY := srcWidth - 1 - x
			dst.Set(newX, newY, src.At(x, y))
		}
	}

	return dst
}

// Rotate180 旋转 180 度
func Rotate180Clockwise(src image.Image) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, srcWidth, srcHeight))

	for y := 0; y < srcHeight; y++ {
		for x := 0; x < srcWidth; x++ {
			newX := srcWidth - 1 - x
			newY := srcHeight - 1 - y
			dst.Set(newX, newY, src.At(x, y))
		}
	}

	return dst
}

// Rotate270Clockwise 顺时针旋转 270 度
func Rotate270Clockwise(src image.Image) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, srcHeight, srcWidth))

	for y := 0; y < srcHeight; y++ {
		for x := 0; x < srcWidth; x++ {
			newX := srcHeight - 1 - y
			newY := x
			dst.Set(newX, newY, src.At(x, y))
		}
	}

	return dst
}

// HorizontalMirror 图像水平镜像
func HorizontalMirror(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(width-1-x, y, src.At(x, y))
		}
	}
	return dst
}

// VerticalMirror 图像垂直镜像
func VerticalMirror(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(x, height-1-y, src.At(x, y))
		}
	}
	return dst
}
