package app

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
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
	// 上传表单
	// 缓冲的大小 - 4M
	r.ParseMultipartForm(1024 << 12)
	// 是上传表单域的名字 fileHeader
	upFile, upFileInfo, err := r.FormFile("userfile")
	if err != nil {
		w.Write(ResponseError(StatusForm))
		logrus.Printf("FormFile %v\n", err)
		return
	}
	defer upFile.Close()

	fileData, err := io.ReadAll(upFile)
	if err != nil {
		w.Write(ResponseError(StatusForm))
		logrus.Printf("FormFile %v\n", err)
		return
	}

	// 判断是否有这个图片类型
	imgtype, err := filetype.Image(fileData)
	if imgtype == types.Unknown {
		w.Write(ResponseError(StatusImgIsType))
		return
	}

	// 进行 MD5 算计，返回 16进制的 byte 数组
	fileMd5 := fmt.Sprintf("%x", md5.Sum(fileData))

	// 目录计算
	// 组合文件完整路径
	dirPath := c.ImgHand.JoinDir(fileMd5) // 目录
	filePath := dirPath + fileMd5         // 文件路径

	// 获取目录信息，并创建目录
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			w.Write(ResponseError(StatusMkdir))
			logrus.Printf("os.MkdirAll %v\n", err)
			return
		}
	} else {
		if !dirInfo.IsDir() {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				w.Write(ResponseError(StatusMkdir))
				logrus.Printf("os.MkdirAll %v\n", err)
				return
			}
		}
	}

	// 存入文件
	_, err = os.Stat(filePath)
	if err != nil {
		err = os.WriteFile(filePath, fileData, 0755)
		if err != nil {
			w.Write(ResponseError(StatusOpenFile))
			logrus.Printf("WriteFile %v\n", err)
			return
		}

	}
	// 响应返回
	res := new(UpdateResponse)
	res.Success = true
	res.Code = StatusOK
	res.Msg = StatusText(StatusOK)
	res.Data.ImgId = fileMd5
	res.Data.Mime = imgtype.Extension
	res.Data.Size = upFileInfo.Size
	w.Write(ResponseJson(res))
}

// Info 获取图片信息
func (c *Controller) Info(w http.ResponseWriter, r *http.Request) {
	// 获取要图片id
	imgid := r.FormValue("imgid")
	// 获取裁剪后图像的宽度、高度
	width := c.ImgHand.StringToInt(r.FormValue("w"))  // 宽度
	height := c.ImgHand.StringToInt(r.FormValue("h")) // 高度
	// 组合文件完整路径
	filePath := c.ImgHand.UrlParse(imgid)
	if filePath == "" {
		w.Write(ResponseError(StatusUrlNotFound))
		return
	}

	if width != 0 || height != 0 {
		filePath = fmt.Sprintf("%s_%d_%d", filePath, width, height)
	}

	fimg, err := os.Open(filePath)
	if err != nil {
		w.Write(ResponseError(StatusImgNotFound))
		logrus.Printf("os.Open %v\n", err)
		return
	}
	defer fimg.Close()

	bufimg := bufio.NewReader(fimg)
	_, imgtype, err := image.Decode(bufimg)
	if err != nil {
		w.Write(ResponseError(StatusImgNotFound))
		logrus.Printf("image.Decode %v\n", err)
		return
	}

	finfo, _ := fimg.Stat()

	// 响应返回
	res := new(UpdateResponse)
	res.Success = true
	res.Code = StatusOK
	res.Msg = StatusText(StatusOK)
	res.Data.ImgId = imgid
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
