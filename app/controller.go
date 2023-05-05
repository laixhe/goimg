package app

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/laixhe/goimg/imghand"
)

type Controller struct {
	ImgHand *imghand.ImgHand
}

func NewController() *Controller {
	return &Controller{
		ImgHand: imghand.NewImgHand(),
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/favicon.ico" {
		return
	}
	if r.Method == http.MethodGet {
		c.Get(w, r)
		return
	}
	if r.Method == http.MethodPost {
		c.Post(w, r)
		return
	}
}

// Get 输出图片
func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	// 组合文件目录路径
	pathDir := c.ImgHand.UrlParse(path)
	if pathDir == "" {
		_, _ = w.Write(showTestHtml())
		return
	}

	// 获取要裁剪图像的宽度、高度
	width := c.ImgHand.StringToInt(r.FormValue("w"))  // 宽度
	height := c.ImgHand.StringToInt(r.FormValue("h")) // 高度

	// 加载图片
	c.ImgHand.CutImage(w, pathDir+path, width, height)
}

// Post 上传图片
func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {

	// 响应返回
	res := new(UpdateResponse)

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
		logrus.Printf("res is %v\n", ResponseJson(res))
		return
	}

	// 判断是否有这个图片类型
	if !imghand.IsImgType(imgtype) {
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
	fileMd5 := hex.EncodeToString(md5Hash.Sum(nil))

	// 目录计算 --------------------------------------

	// 组合文件完整路径
	dirPath := c.ImgHand.JoinDir(fileMd5) // 目录
	filePath := dirPath + fileMd5         // 文件路径

	// 获取目录信息，并创建目录
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			logrus.Printf("res is %v\n", res)
			res.Code = StatusMkdir
			res.Msg = StatusText(StatusMkdir)
			w.Write(ResponseJson(res))
			return
		}
	} else {
		if !dirInfo.IsDir() {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				res.Code = StatusMkdir
				res.Msg = StatusText(StatusMkdir)
				w.Write(ResponseJson(res))
				return
			}
		}
	}

	// 存入文件 --------------------------------------

	_, err = os.Stat(filePath)
	if err != nil {
		// 打开一个文件,文件不存在就会创建
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			res.Code = StatusOpenFile
			res.Msg = StatusText(StatusOpenFile)
			w.Write(ResponseJson(res))
			return
		}
		defer file.Close()

		if imgtype == imghand.PNG {
			err = png.Encode(file, img)

		} else if imgtype == imghand.JPG || imgtype == imghand.JPEG {
			err = jpeg.Encode(file, img, nil)

		} else if imgtype == imghand.GIF {

			// 重新对 gif 格式进行解码
			// image.Decode 只能读取 gif 的第一帧

			// 设置下次读写位置（移动文件指针位置）
			_, err = upfile.Seek(0, 0)
			if err != nil {
				res.Code = StatusFileSeek
				res.Msg = StatusText(StatusFileSeek)
				w.Write(ResponseJson(res))
				return
			}

			gifimg, giferr := gif.DecodeAll(upfile)
			if giferr != nil {
				res.Code = StatusImgDecode
				res.Msg = StatusText(StatusImgDecode)
				w.Write(ResponseJson(res))
				return
			}
			err = gif.EncodeAll(file, gifimg)
		}

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

// Info 获取图片信息
func (c *Controller) Info(w http.ResponseWriter, r *http.Request) {
	// 响应返回
	res := new(UpdateResponse)
	// 获取要图片id
	imgid := r.FormValue("imgid")
	// 获取裁剪后图像的宽度、高度
	width := c.ImgHand.StringToInt(r.FormValue("w"))  // 宽度
	height := c.ImgHand.StringToInt(r.FormValue("h")) // 高度
	// 组合文件完整路径
	filePath := c.ImgHand.UrlParse(imgid)
	if filePath == "" {
		res.Code = StatusUrlNotFound
		res.Msg = StatusText(StatusUrlNotFound)
		w.Write(ResponseJson(res))
		return
	}

	if width != 0 || height != 0 {
		filePath = fmt.Sprintf("%s_%d_%d", filePath, width, height)
	}

	fimg, err := os.Open(filePath)
	if err != nil {
		res.Code = StatusImgNotFound
		res.Msg = StatusText(StatusImgNotFound)
		w.Write(ResponseJson(res))
		return
	}
	defer fimg.Close()

	bufimg := bufio.NewReader(fimg)
	_, imgtype, err := image.Decode(bufimg)
	if err != nil {
		res.Code = StatusImgNotFound
		res.Msg = StatusText(StatusImgNotFound)
		w.Write(ResponseJson(res))
		return
	}

	finfo, _ := fimg.Stat()

	res.Success = true
	res.Code = StatusOK
	res.Msg = StatusText(StatusOK)
	res.Data.Imgid = imgid
	res.Data.Mime = imgtype
	res.Data.Size = finfo.Size()

	_, _ = w.Write(ResponseJson(res))
}

// StatusCode 状态码
func (c *Controller) StatusCode(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(GetStatusText())
	_, _ = w.Write(data)
}

// Test 测试上传
func (c *Controller) Test(w http.ResponseWriter, r *http.Request) {
	html := `
<html>
    <head>
        <title>测试上传</title>
    </head>
    <body>
        <div style="padding: 20px;">
            <p>测试上传</p>
            <form action="/" method="post" enctype="multipart/form-data">
                <p><input type="file" name="userfile" /></p>
                <p><input type="submit" name="submit" value="上传" /></p>
            </form>
        </div>
    </body>
</html>`
	_, _ = w.Write([]byte(html))
}
