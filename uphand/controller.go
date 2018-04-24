package uphand

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/laixhe/goimg/imghand"
)

type Controller struct {
}

func (this Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/favicon.ico" {
		// 设置 http请求状态
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "GET" {
		this.Get(w, r)
		return
	}

	if r.Method == "POST" {
		this.Post(w, r)
		return
	}
}

// 输出图片
func (this Controller) Get(w http.ResponseWriter, r *http.Request) {

	urlParse := r.URL.String()
	if len(urlParse) < 32 {
		w.Write(showMain())
		return
	}

	// 进行 url 解析
	parse, err := url.Parse(urlParse)
	if err != nil {
		w.Write(show404("404 not found"))
		return
	}
	if len(parse.Path) != 33 {
		w.Write(show404(parse.Path))
		return
	}

	parsePath := parse.Path[1:]

	// 匹配是否是 md5 的长度
	if !regexpUrlParse.MatchString(parsePath) {
		w.Write(show404(parse.Path))
		return
	}

	// 组合文件完整路径
	sortPath := SortPath([]byte(parsePath[:5]))
	filePath := "img/" + sortPath + "/" + parsePath

	log.Println(parsePath)
	log.Println(filePath)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		w.Write(show404(filePath + " : " + err.Error()))
		return
	}
	defer file.Close()

	io.Copy(w, file)

}

// 上传图片
func (this Controller) Post(w http.ResponseWriter, r *http.Request) {

	// 上传表单 --------------------------------------

	// 缓冲的大小 - 4M
	r.ParseMultipartForm(1024 << 12)
	//是上传表单域的名字fileHeader
	upfile, _, err := r.FormFile("userfile")
	if err != nil {
		w.Write(show404("表单字段: [ userfile ]" + err.Error()))
		return
	}
	defer upfile.Close()

	// 图片解码 --------------------------------------

	// 读入缓存
	bufUpFile := bufio.NewReader(upfile)
	// 进行图片的解码
	img, imgtype, err := image.Decode(bufUpFile)
	if err != nil {
		w.Write(show404("图片解码错误:" + err.Error()))
		return
	}

	// 判断是否有这个图片类型
	if !imghand.IsType(imgtype) {
		w.Write(show404("没有这个图片类型:" + imgtype))
		return
	}

	// 设置文件读写下标 --------------------------------

	// 设置下次读写位置（移动文件指针位置）
	_, err = upfile.Seek(0, 0)
	if err != nil {
		w.Write(show404("设置下次读写位置失败"))
		return
	}

	// 计算文件的 MD5 值 -----------------------------

	// 初始化 MD5 实例
	md5Hash := md5.New()
	// 读入缓存
	bufFile := bufio.NewReader(upfile)
	_, err = io.Copy(md5Hash, bufFile)
	if err != nil {
		w.Write(show404("MD5计算文件失败:" + err.Error()))
		return
	}
	// 进行 MD5 算计，返回 16进制的 byte 数组
	fileMd5FX := md5Hash.Sum(nil)
	fileMd5 := fmt.Sprintf("%x", fileMd5FX)

	// 目录计算 --------------------------------------

	// 组合文件完整路径
	sortPath := SortPath([]byte(fileMd5[:5]))
	dirPath := "img/" + sortPath + "/" // 目录
	filePath := dirPath + fileMd5      // 文件路径

	// 获取目录信息，并创建目录
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		os.MkdirAll(dirPath, 0666)
	} else {
		if !dirInfo.IsDir() {
			os.MkdirAll(dirPath, 0666)
		}
	}

	// 存入文件 --------------------------------------

	// 打开一个文件,文件不存在就会创建
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		w.Write(show404("上传创建失败:" + err.Error()))
		return
	}
	defer file.Close()

	if imgtype == imghand.PNG {
		png.Encode(file, img)
	} else if imgtype == imghand.JPG {
		jpeg.Encode(file, img, nil)
	} else if imgtype == imghand.JPEG {
		jpeg.Encode(file, img, nil)
	}

	w.Write([]byte(fileMd5))
}
