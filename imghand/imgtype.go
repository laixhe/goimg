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

var imgType []string = []string{PNG, JPG, JPEG, GIF}

func GetImgType() []string {
	return imgType
}

// 判断是否有这个图片类型
func IsType(str string) bool {

	// 转小写
	str = strings.ToLower(str)

	for _, v := range imgType {
		if v == str {
			return true
		}
	}

	return false
}
