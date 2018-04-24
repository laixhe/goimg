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

	// 响应返回
	res := UpdateResponse{}

	// 上传表单 --------------------------------------

	// 缓冲的大小 - 4M
	r.ParseMultipartForm(1024 << 12)
	//是上传表单域的名字fileHeader
	upfile, upFileInfo, err := r.FormFile("userfile")
	if err != nil {

		res.Code = StatusForm
		res.Msg = StatusText(StatusForm)
		w.Write(ResponseJson(res))

		return
	}
	defer upfile.Close()

	// 图片解码 --------------------------------------

	// 读入缓存
	bufUpFile := bufio.NewReader(upfile)
	// 进行图片的解码
	img, imgtype, err := image.Decode(bufUpFile)
	if err != nil {

		res.Code = StatusImgDecode
		res.Msg = StatusText(StatusImgDecode)
		w.Write(ResponseJson(res))

		return
	}

	// 判断是否有这个图片类型
	if !imghand.IsType(imgtype) {

		res.Code = StatusImgIsType
		res.Msg = StatusText(StatusImgIsType)
		w.Write(ResponseJson(res))

		return
	}

	// 设置文件读写下标 --------------------------------

	// 设置下次读写位置（移动文件指针位置）
	_, err = upfile.Seek(0, 0)
	if err != nil {

		res.Code = StatusFileSeek
		res.Msg = StatusText(StatusFileSeek)
		w.Write(ResponseJson(res))

		return
	}

	// 计算文件的 MD5 值 -----------------------------

	// 初始化 MD5 实例
	md5Hash := md5.New()
	// 读入缓存
	bufFile := bufio.NewReader(upfile)
	_, err = io.Copy(md5Hash, bufFile)
	if err != nil {

		res.Code = StatusFileMd5
		res.Msg = StatusText(StatusFileMd5)
		w.Write(ResponseJson(res))

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
		err = os.MkdirAll(dirPath, 0666)
		if err != nil {

			res.Code = StatusMkdir
			res.Msg = StatusText(StatusMkdir)
			w.Write(ResponseJson(res))

		}
	} else {
		if !dirInfo.IsDir() {
			err = os.MkdirAll(dirPath, 0666)
			if err != nil {

				res.Code = StatusMkdir
				res.Msg = StatusText(StatusMkdir)
				w.Write(ResponseJson(res))

			}
		}
	}

	// 存入文件 --------------------------------------

	// 打开一个文件,文件不存在就会创建
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {

		res.Code = StatusOpenFile
		res.Msg = StatusText(StatusOpenFile)
		w.Write(ResponseJson(res))

		return
	}
	defer file.Close()

	if imgtype == imghand.PNG {
		err = png.Encode(file, img)
		if err != nil {

			res.Code = StatusImgEncode
			res.Msg = StatusText(StatusImgEncode)
			w.Write(ResponseJson(res))

			return
		}

	} else if imgtype == imghand.JPG {
		err = jpeg.Encode(file, img, nil)
		if err != nil {

			res.Code = StatusImgEncode
			res.Msg = StatusText(StatusImgEncode)
			w.Write(ResponseJson(res))

			return
		}

	} else if imgtype == imghand.JPEG {
		err = jpeg.Encode(file, img, nil)
		if err != nil {

			res.Code = StatusImgEncode
			res.Msg = StatusText(StatusImgEncode)
			w.Write(ResponseJson(res))

			return
		}

	}

	res.Success = true
	res.Code = StatusOK
	res.Msg = StatusText(StatusOK)
	res.Data.Imgid = fileMd5
	res.Data.Mime = imgtype
	res.Data.Size = upFileInfo.Size

	w.Write(ResponseJson(res))

}
