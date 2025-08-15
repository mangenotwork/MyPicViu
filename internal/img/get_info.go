package img

// 获取图片信息
type ImgInfo struct {
	FileName          string  // 图片文件名
	FileSize          int64   // 图片文件大小
	FileMd5           string  // 图片文件md5
	FilePath          string  // 图片文件路径
	FileMode          string  // 权限
	FileModTime       string  // 最后修改时间
	FileSys           string  // 系统原生信息
	ImgWidth          int     // 图片宽
	ImgHeight         int     // 图片高
	ImgDPI            int     // 图片dpi
	ImgBit            int     // 图片是多少位的
	ImgMimeType       string  // 图片MIME类型
	ImgOrientation    int     // 图片方向
	ImgCameraMake     string  // 图片相机制造商
	ImgCameraModel    string  // 相机型号
	ImgDateTime       string  // 拍摄时间
	ImgFocalLength    string  // 焦距
	ImgISO            string  // iso
	ImgAperture       string  // 光圈
	ImgShutterSpeed   string  // 快门
	ImgSaturation     float64 // 图片饱和度值
	ImgBrightness     float64 // 图片亮度值
	ImgContrast       float64 // 图片对比度值
	ImgSharpness      float64 // 图片锐度值
	ImgExposure       float64 // 图片曝光度值
	ImgTemperature    float64 // 图片色温值
	ImgHue            float64 // 图片色调值
	ImgNoise          float64 // 图片噪点值
	ImgDifferenceHash string  // 图片的差异哈希值获取图像的指纹字符串
	ImgPHash          string  // 感知哈希算法获取图像的指纹字符串
	ImgAverageHash    string  // 均值哈希值获取图像的指纹字符串
}

// todo ...
