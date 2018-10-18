package uphand

// 状态码
const (
	StatusJson        = 1
	StatusForm        = 2
	StatusImgDecode   = 3
	StatusImgIsType   = 4
	StatusFileSeek    = 5
	StatusFileMd5     = 6
	StatusMkdir       = 7
	StatusOpenFile    = 8
	StatusImgEncode   = 9
	StatusImgNotFound = 10
	StatusUrlNotFound = 11
	StatusOK          = 200
	StatusNotFound    = 404
)

var statusText = map[int]string{
	StatusJson:        "json打包失败",
	StatusForm:        "表单字段 userfile 缺少",
	StatusImgDecode:   "图片解码不符合",
	StatusImgIsType:   "图片类型不符合",
	StatusFileSeek:    "设置文件读写位置失败",
	StatusFileMd5:     "计算文件MD5失败",
	StatusMkdir:       "创建目录失败",
	StatusOpenFile:    "文件创建失败",
	StatusImgEncode:   "图片生成失败",
	StatusImgNotFound: "没有找到图片",
	StatusUrlNotFound: "Url Not Found",
	StatusOK:          "OK",
	StatusNotFound:    "Not Found",
}

func StatusText(code int) string {
	return statusText[code]
}

func GetStatusText() map[int]string {
	return statusText
}

//------------------------------------------------------------------

func showMain() []byte {

	show := `<html>
  <div>Goimg 轻量级的图片服务器</div>
  <div><a href="/test">开始吧</a></div>
</html>`

	return []byte(show)
}

