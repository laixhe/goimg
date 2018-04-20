package controller

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"net/http"
	"os"
	"fmt"
)

func Index(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.URL)

	fimg, _ := os.Open("./1.jpg")
	defer fimg.Close()

	// 进行图片的解码
	img, _, _ := image.Decode(fimg)

	// 创建 RGBA 画板大小 - 用于裁剪
	dst := image.NewRGBA(image.Rect(0, 0, 400, 400))

	// 进行图像合成 - 裁剪
	draw.Draw(dst, dst.Rect, img, image.Point{50, 50}, draw.Src)

	// 图片流方式输出
	w.Header().Set("Content-Type", "image/png")

	// 进行图片的编码
	png.Encode(w, dst)
}
