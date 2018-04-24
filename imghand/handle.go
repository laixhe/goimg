package imghand

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// 裁剪图像
func CutImage() {

	fimg, _ := os.Open("./1.jpg")
	defer fimg.Close()

	// 进行图片的解码
	img, _, _ := image.Decode(fimg)

	// 创建 RGBA 画板大小 - 用于裁剪
	dst := image.NewRGBA(image.Rect(0, 0, 400, 400))

	// 进行图像合成 - 裁剪
	draw.Draw(dst, dst.Rect, img, image.Point{50, 50}, draw.Src)

	// 图片流方式输出
	//w.Header().Set("Content-Type", "image/png")

	// 进行图片的编码
	//png.Encode(w, dst)

}

// 翻转图像
func FlipImage() {
}

// 旋转图像
func RotateImage() {
}
