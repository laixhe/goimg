package imghand

import (
	"bufio"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/nfnt/resize"
	"github.com/sirupsen/logrus"

	"github.com/laixhe/goimg/config"
)

type ImgHand struct {
	noImg      *image.RGBA // 创建 RGBA 画板大小 - 用于找不到图片时用
	stringPool *sync.Pool  // 字符串对象池
}

func NewImgHand() *ImgHand {
	return &ImgHand{
		noImg: image.NewRGBA(image.Rect(0, 0, 400, 400)),
		stringPool: &sync.Pool{
			New: func() any {
				return &strings.Builder{}
			},
		},
	}
}

func (h *ImgHand) GetStringPool() *strings.Builder {
	return h.stringPool.Get().(*strings.Builder)
}

func (h *ImgHand) SetStringPool(str *strings.Builder) {
	str.Reset()
	h.stringPool.Put(str)
}

// UrlParse 进行 url 部分解析 - md5，并组合文件目录路径
func (h *ImgHand) UrlParse(path string) string {
	if len(path) != 32 {
		return ""
	}
	// 组合文件目录路径
	return h.JoinDir(path)
}

// JoinDir 组合文件目录路径
func (h *ImgHand) JoinDir(path string) string {
	var strBuilder = h.GetStringPool()

	strBuilder.WriteString(config.ImgDir())
	strBuilder.WriteString("/")
	strBuilder.WriteString(path[:5])
	strBuilder.WriteString("/")
	strBuilder.WriteString(path)
	strBuilder.WriteString("/")
	str := strBuilder.String()

	h.SetStringPool(strBuilder)
	return str
}

// StringToInt 字符串的数字转int
func (h *ImgHand) StringToInt(str string) int {
	if str == "" {
		return 0
	}
	toint, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	if toint < 0 {
		return 0
	}
	return toint
}

// CutImage 裁剪图像
func (h *ImgHand) CutImage(w http.ResponseWriter, path string, width, height int) {
	// 没有宽高，就是在加载原图像
	if width == 0 && height == 0 {
		file, err := os.Open(path)
		if err != nil {
			h.NoImage(w)

			logrus.Println("file, err = os.Open(path)", err)
			return
		}
		io.Copy(w, file)
		file.Close()
		return
	}

	// 裁剪图像 --------------------------------------

	// 裁剪图像的组合路径
	var strBuilder = h.GetStringPool()
	strBuilder.WriteString(path)
	strBuilder.WriteString("_")
	strBuilder.WriteString(strconv.Itoa(width))
	strBuilder.WriteString("_")
	strBuilder.WriteString(strconv.Itoa(height))
	CutPath := strBuilder.String()
	h.SetStringPool(strBuilder)

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
		h.NoImage(w)
		logrus.Println("file, err = os.Open(path)", err)
		return
	}
	defer file.Close()

	// 图片解码 --------------------------------------

	bufFile := bufio.NewReader(file)
	img, imgtype, err := image.Decode(bufFile)
	if err != nil {
		h.NoImage(w)
		logrus.Println("img, imgtype, err := image.Decode(bufFile)", err)
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
		h.NoImage(w)
		logrus.Println("out, err := os.Create(CutPath)", err)
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

// NoImage 用于找不到图片时用
func (h *ImgHand) NoImage(w http.ResponseWriter) {
	// 图片流方式输出
	w.Header().Set("Content-Type", "image/png")
	// 进行图片的编码
	png.Encode(w, h.noImg)
}
