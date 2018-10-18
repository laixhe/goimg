## Goimg 轻量级的图片服务器

### 简介
> goImg是一个使用Golang语言编写的图片服务器

> 目前只实现单文件上传

> 支持 jpeg png gif 等图片上传

> 请使用 Go1.10 以上

### 功能特点
> 文件存储目录采用md5算法生成

> 目前支持 jpeg png 的缩略图，gif 暂不支持

### 安装
> go get github.com/laixhe/goimg

### 获取图片
> GET /9d32e3c40efb0b749270695d5f0afdfc


### 获取图片 并缩略，宽度=300 高度=100
> GET /9d32e3c40efb0b749270695d5f0afdfc?w=300&h=100

### 上传图片
> POST  /

> 表单参数: userfile

> 返回值: json 主要是 imgid

```
{
	"success": true,
	"code": 200,
	"msg": "OK",
	"version": "v0.1.1",
	"data": {
		"size": 42445,
		"mime": "jpeg",
		"imgid": "9d32e3c40efb0b749270695d5f0afdfc"
	}
}
```

### 获取图片信息
> GET /info?imgid=9d32e3c40efb0b749270695d5f0afdfc

> GET /info?imgid=9d32e3c40efb0b749270695d5f0afdfc&w=300&h=100

> 返回值: json

> 与 上传图片 的返回一致

### 获取状态码
> GET /statuscode

> 返回值: json