package imghand

import (
	"bufio"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

// 裁剪图像
func CutImage(w http.ResponseWriter, path string, width, height int) {

	// 没有宽高，就是在加载原图像
	if width == 0 && height == 0 {

		file, err := os.Open(path)
		if err != nil {
			NoImage(w)

			log.Println("file, err = os.Open(path)", err)
			return
		}

		io.Copy(w, file)
		file.Close()

		return
	}

	// 裁剪图像 --------------------------------------

	// 裁剪图像的组合路径
	var str = strings.Builder{}
	str.WriteString(path)
	str.WriteString("_")
	str.WriteString(strconv.Itoa(width))
	str.WriteString("_")
	str.WriteString(strconv.Itoa(height))
	CutPath := str.String()

	// 判断是否存在裁剪图像
	file, err := os.Open(CutPath)
	if err == nil {

		io.Copy(w, file)
		file.Close()

		return
	}

	// 原图像
	file, err = os.Open(path)
	if err != nil {
		NoImage(w)

		log.Println("file, err = os.Open(path)", err)
		return
	}
	defer file.Close()

	// 图片解码 --------------------------------------

	bufFile := bufio.NewReader(file)
	img, imgtype, err := image.Decode(bufFile)
	if err != nil {
		NoImage(w)

		log.Println("img, imgtype, err := image.Decode(bufFile)", err)
		return
	}

	// 要裁剪的宽高不能大于自身的宽高
	Rwidth := img.Bounds().Max.X
	if width > Rwidth {
		width = Rwidth
	}

	Rheight := img.Bounds().Max.Y
	if height > Rheight {
		height = Rheight
	}

	// gif 图就不处理了
	if imgtype == GIF || (width == Rwidth && height == Rheight) {

		// 设置文件的偏移量 - 因为文件被 image.Decode 后文件的偏移量到尾部
		file.Seek(0, 0)

		// 向浏览器输出
		io.Copy(w, file)

		return
	}

	// 进行裁剪
	reimg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	// 裁剪的存储
	out, err := os.Create(CutPath)
	if err != nil {
		NoImage(w)

		log.Println("out, err := os.Create(CutPath)", err)
		return
	}
	defer out.Close()

	if imgtype == JPEG || imgtype == JPG {

		// 保存裁剪的图片
		jpeg.Encode(out, reimg, nil)

		// 向浏览器输出
		jpeg.Encode(w, reimg, nil)

	} else if imgtype == PNG {

		// 保存裁剪的图片
		png.Encode(out, reimg)

		// 向浏览器输出
		png.Encode(w, reimg)
	}

}

// 用于找不到图片时用
func NoImage(w http.ResponseWriter) {

	// 图片流方式输出
	w.Header().Set("Content-Type", "image/png")
	// 进行图片的编码
	png.Encode(w, noImg)

}
