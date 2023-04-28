package imghand

import (
	"bufio"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/laixhe/goimg/config"
)

type ImgHand struct {
	regexpUrlParse *regexp.Regexp // 匹配md5的长度
	noImg          *image.RGBA    // 创建 RGBA 画板大小 - 用于找不到图片时用
	stringPool     *sync.Pool     // 字符串对象池
}

func NewImgHand() *ImgHand {
	regexpUrlParse, _ := regexp.Compile("[a-z0-9]{32}")
	return &ImgHand{
		regexpUrlParse: regexpUrlParse,
		noImg:          image.NewRGBA(image.Rect(0, 0, 400, 400)),
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

// IsMD5Path 匹配是否是 md5 的长度
func (h *ImgHand) IsMD5Path(str string) bool {
	return h.regexpUrlParse.MatchString(str)
}

// SortPath 路径部分排序做目录
func (h *ImgHand) SortPath(str []byte) string {

	// 对 byte 进行排序
	strLen := len(str)
	for i := 0; i < strLen; i++ {
		for j := 1 + i; j < strLen; j++ {
			if str[i] > str[j] {
				str[i], str[j] = str[j], str[i]
			}
		}
	}

	// 对 byte 依次组成数字符串
	var ret = strings.Builder{}

	for i := 0; i < strLen; i++ {
		ret.WriteString(strconv.Itoa(int(str[i])))
	}

	return ret.String()
}

// JoinPath 组合文件目录路径
func (h *ImgHand) JoinPath(md5Str string) string {

	// 路径部分排序做目录
	sortPath := h.SortPath([]byte(md5Str[:5]))

	var strBuilder = h.GetStringPool()

	strBuilder.WriteString(config.ImgDir())
	strBuilder.WriteString(sortPath)
	strBuilder.WriteString("/")
	strBuilder.WriteString(md5Str)
	str := strBuilder.String()

	h.SetStringPool(strBuilder)
	return str
}

// UrlParse 进行 url 部分解析 - md5，并组合文件完整路径
func (h *ImgHand) UrlParse(md5Url string) string {
	if md5Url == "" {
		return ""
	}
	if len(md5Url) < 32 {
		return ""
	}
	// 进行 url 解析
	parse, err := url.Parse(md5Url)
	if err != nil {
		return ""
	}
	if len(parse.Path) != 32 {
		return ""
	}
	// 匹配是否是 md5 的长度
	if !h.IsMD5Path(parse.Path) {
		return ""
	}
	// 组合文件完整路径
	return h.JoinPath(parse.Path) + "/" + parse.Path
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
		h.NoImage(w)

		log.Println("file, err = os.Open(path)", err)
		return
	}
	defer file.Close()

	// 图片解码 --------------------------------------

	bufFile := bufio.NewReader(file)
	img, imgtype, err := image.Decode(bufFile)
	if err != nil {
		h.NoImage(w)

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
		h.NoImage(w)

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

// NoImage 用于找不到图片时用
func (h *ImgHand) NoImage(w http.ResponseWriter) {

	// 图片流方式输出
	w.Header().Set("Content-Type", "image/png")
	// 进行图片的编码
	png.Encode(w, h.noImg)

}
