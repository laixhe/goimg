## Goimg 轻量级的图片服务器

### 简介
> goImg是一个使用Golang语言编写的图片服务器
> 目前只实现单文件上传
> 支持 jpeg png gif 等图片上传

### 功能特点
> 文件存储目录采用md5算法生成

### 安装
> go get github.com/laixhe/goimg

### 获取图片
> GET /图片ID

### 上传图片
> POST /

> 表单参数: userfile

> 返回值: json 主要是 imgid

```
{
	"success": true,
	"code": 200,
	"msg": "OK",
	"version": "0.1",
	"data": {
		"size": 42445,
		"mime": "jpeg",
		"imgid": "9d32e3c40efb0b749270695d5f0afdfc"
	}
}
```