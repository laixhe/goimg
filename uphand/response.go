package uphand

import "encoding/json"

type ResponseInterface interface {
	SetVersion(string)
}

// 响应请求的公共模型
type ResponseModel struct {
	Success bool   `json:"success"` // 是否成功
	Code    int    `json:"code"`    // 响应码
	Msg     string `json:"msg"`     // 响应信息
	Version string `json:"version"` // 版本号
	Data    string `json:"data"`    // 数据
}

func (this *ResponseModel) SetVersion(str string) {
	this.Version = str
}

// 上传响应数据
type UpdateDate struct {
	Size  int64  `json:"size"`  // 大小
	Mime  string `json:"mime"`  // 图片类型
	Imgid string `json:"imgid"` // 图片id
}

// 上传响应模型
type UpdateResponse struct {
	ResponseModel
	Data UpdateDate `json:"data"`
}

// 响应 json 打包
func ResponseJson(res ResponseInterface) []byte {

	res.SetVersion("v0.1.1")

	data, err := json.Marshal(res)
	if err != err {

		// 打包失败
		data, _ = json.Marshal(ResponseModel{false,
			StatusJson,
			StatusText(StatusJson),
			"", ""})
	}

	return data
}
