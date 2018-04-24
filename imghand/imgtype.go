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

var ImgType []string = []string{PNG, JPG, JPEG, GIF}

// 判断是否有这个图片类型
func IsType(str string) bool {

	// 转小写
	str = strings.ToLower(str)

	for _, v := range ImgType {
		if v == str {
			return true
		}
	}

	return false
}
