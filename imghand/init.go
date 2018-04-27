package imghand

import (
	"log"
	"regexp"
	"image"
)

var regexpUrlParse *regexp.Regexp

var noImg *image.RGBA

func init() {

	var err error
	// 初始化正则表达式
	regexpUrlParse, err = regexp.Compile("[a-z0-9]{32}")
	if err != nil {
		log.Fatalln("regexpUrlParse:", err)
	}

	// 创建 RGBA 画板大小 - 用于找不到图片时用
	noImg = image.NewRGBA(image.Rect(0, 0, 400, 400))

}
