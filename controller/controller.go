package controller

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"os"
	"io"
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
	if err !=nil {
		w.Write(show404("404"))
		return
	}
	if len(parse.Path) != 33 {
		w.Write(show404(parse.Path))
		return
	}

	// 匹配是否是 md5 的长度
	if !regexpUrlParse.MatchString(parse.Path[1:]) {
		log.Println("RegexpUrlParse:", err)

		w.Write(show404(parse.Path))
		return
	}
	log.Println(parse.Path[1:])

	// 打开文件
	file, err := os.Open("img/" + parse.Path[1:])
	if err != nil {
		log.Println("Open File:", err)

		w.Write(show404(parse.Path))
		return
	}
	defer file.Close()

	io.Copy(w, file)

	a := []byte("45bc")
	//sort.Ints(a)
	log.Println(a)
	log.Printf("%d%d%d%d\n", a[0], a[1], a[2], a[3])

	a = []byte("b45c")
	log.Println(a)
	log.Printf("%d%d%d%d\n", a[0], a[1], a[2], a[3])

}

// 上传图片
func (this Controller) Post(w http.ResponseWriter, r *http.Request) {

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
