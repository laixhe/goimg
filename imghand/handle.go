package imghand

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

// 裁剪图像
func CutImage(w http.ResponseWriter, path string, width, height int) {

	// 没有宽高，就是在加载原图像
	if width == 0 && height == 0 {

		file, err := os.Open(path)
		if err != nil {
			NoImage(w)
			return
		}
		defer file.Close()

		io.Copy(w, file)

		return
	}

	// 裁剪图像 --------------------------------------

	// 裁剪图像的组合路径
	CutPath := fmt.Sprintf("%s_%d_%d", path, width, height)

	// 判断是否存在裁剪图像
	_, err := os.Stat(CutPath)
	if err == nil {

		file, err := os.Open(CutPath)
		if err != nil {
			NoImage(w)
			return
		}
		defer file.Close()

		io.Copy(w, file)

		return
	}

	// 判断是否存在原图像
	_, err = os.Stat(path)
	if err != nil {
		NoImage(w)
		return
	}

	// 原图像
	file, err := os.Open(path)
	if err != nil {
		NoImage(w)
		return
	}
	defer file.Close()

	// 图片解码 --------------------------------------

	bufFile := bufio.NewReader(file)
	img, imgtype, err := image.Decode(bufFile)
	if err != nil {
		NoImage(w)
		return
	}

	// gif 图就不处理了
	if imgtype == GIF {

		_, err = file.Seek(0, 0)
		if err != nil {
			NoImage(w)
			return
		}

		io.Copy(w, file)

		return
	}

	// 进行裁剪
	reimg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	// 裁剪的存储
	out, err := os.Create(CutPath)
	if err != nil {
		NoImage(w)
		return
	}
	defer out.Close()

	if imgtype == JPEG || imgtype == JPG {
		jpeg.Encode(out, reimg, nil)
	} else if imgtype == PNG {
		png.Encode(out, reimg)
	}

	_, err = out.Seek(0, 0)
	if err != nil {
		NoImage(w)
		return
	}

	io.Copy(w, out)

}

func NoImage(w http.ResponseWriter) {

	// 图片流方式输出
	w.Header().Set("Content-Type", "image/png")
	// 进行图片的编码
	png.Encode(w, noImg)

}
