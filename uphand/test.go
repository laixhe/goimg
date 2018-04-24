package uphand

import (
	"encoding/json"
	"net/http"
)

// 测试上传
func Test(w http.ResponseWriter, r *http.Request) {

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

	w.Write([]byte(html))

}

// 状态码
func StatusCode(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(GetStatusText())
	w.Write(data)
}
