package imghand

import (
	"strings"
)

// 图像类型
const (
	PNG  = "png"
	JPG  = "jpg"
	JPEG = "jpeg"
	GIF  = "gif"
)

var imgType = []string{PNG, JPG, JPEG, GIF}

// GetImgType 图片类型
func GetImgType() []string {
	return imgType
}

// IsImgType 判断是否有这个图片类型
func IsImgType(str string) bool {
	// 转小写
	str = strings.ToLower(str)
	for _, v := range imgType {
		if v == str {
			return true
		}
	}
	return false
}
