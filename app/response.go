package app

import (
	"encoding/json"

	"github.com/laixhe/goimg/config"
)

type ResponseInterface interface {
	SetVersion(string)
}

// ResponseModel 响应请求的公共模型
type ResponseModel struct {
	Success bool   `json:"success"` // 是否成功
	Code    int    `json:"code"`    // 响应码
	Msg     string `json:"msg"`     // 响应信息
	Version string `json:"version"` // 版本号
	Data    string `json:"data"`    // 数据
}

func (r *ResponseModel) SetVersion(str string) {
	r.Version = str
}

// UpdateDate 上传响应数据
type UpdateDate struct {
	Size  int64  `json:"size"`  // 大小
	Mime  string `json:"mime"`  // 图片类型
	Imgid string `json:"imgid"` // 图片id
}

// UpdateResponse 上传响应模型
type UpdateResponse struct {
	ResponseModel
	Data UpdateDate `json:"data"`
}

// ResponseJson 响应 json 打包
func ResponseJson(res ResponseInterface) []byte {
	res.SetVersion(config.Get().App.Version)
	data, err := json.Marshal(res)
	if err != err {
		// 打包失败
		data, _ = json.Marshal(ResponseModel{false,
			StatusJson,
			StatusText(StatusJson),
			config.Get().App.Version,
			""})
	}
	return data
}
