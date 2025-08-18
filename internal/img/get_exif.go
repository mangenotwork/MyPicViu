package img

import (
	"MyPicViu/common/logger"
	"github.com/rwcarlsen/goexif/exif"
	"io"
	"math/big"
	"time"
)

// todo 这里需要优化，提取exif还是有很多问题
func (info *ImgInfo) GetExif(f io.Reader) {

	// 解析EXIF数据
	x, err := exif.Decode(f)
	if err != nil {
		logger.ErrorF("解析EXIF失败: %w", err)
		return
	}

	// 打印基本信息
	logger.Debug("===== 基本EXIF信息 =====")

	// 1. 拍摄时间
	if tag, err := x.Get(exif.DateTimeOriginal); err == nil {
		if val, err := tag.StringVal(); err == nil {
			logger.DebugF("拍摄时间: %s\n", parseExifTime(val))
		}
	}

	// 2. 设备信息
	if tag, err := x.Get(exif.Model); err == nil {
		if val, err := tag.StringVal(); err == nil {
			logger.DebugF("设备型号: %s\n", val)
		}
	}
	if tag, err := x.Get(exif.Make); err == nil {
		if val, err := tag.StringVal(); err == nil {
			logger.DebugF("设备厂商: %s\n", val)
		}
	}

	// 3. 光圈值（使用Rat()替代Rat64()）
	if tag, err := x.Get(exif.FNumber); err == nil {
		// 关键修正：用Rat()获取分数值，返回*big.Rat
		rat, err := tag.Rat(0)
		if err == nil {
			// 转换为浮点值（分子/分母）
			fNumber := new(big.Float).Quo(
				new(big.Float).SetInt(rat.Num()),
				new(big.Float).SetInt(rat.Denom()),
			)
			fVal, _ := fNumber.Float64()
			logger.DebugF("光圈值: f/%.1f\n", fVal)
		}
	}

	// 4. 曝光时间
	if tag, err := x.Get(exif.ExposureTime); err == nil {
		rat, err := tag.Rat(0)
		if err == nil {
			logger.DebugF("曝光时间: %d/%ds\n", rat.Num(), rat.Denom())
		}
	}

	// 5. 焦距
	if tag, err := x.Get(exif.FocalLength); err == nil {
		rat, err := tag.Rat(0)
		if err == nil {
			focal := new(big.Float).Quo(
				new(big.Float).SetInt(rat.Num()),
				new(big.Float).SetInt(rat.Denom()),
			)
			fVal, _ := focal.Float64()
			logger.DebugF("焦距: %.1fmm\n", fVal)
		}
	}

	// 6. 水平DPI
	if tag, err := x.Get(exif.XResolution); err == nil {
		rat, err := tag.Rat(0)
		if err == nil && rat.Denom().Cmp(big.NewInt(0)) != 0 {
			dpi := new(big.Float).Quo(
				new(big.Float).SetInt(rat.Num()),
				new(big.Float).SetInt(rat.Denom()),
			)
			dpiVal, _ := dpi.Float64()
			logger.DebugF("水平DPI: %.1f\n", dpiVal)
		}
	}

}

// 解析EXIF时间格式（"2006:01:02 15:04:05"）为标准时间
func parseExifTime(exifTime string) string {
	t, err := time.Parse("2006:01:02 15:04:05", exifTime)
	if err != nil {
		return exifTime // 解析失败则返回原始字符串
	}
	return t.Format("2006-01-02 15:04:05")
}
